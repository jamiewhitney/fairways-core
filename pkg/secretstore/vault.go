package secretstore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/sethvargo/go-envconfig"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type HashiCorpVault struct {
	client *api.Client
}

type HashiCorpVaultConfig struct {
	Role      string `env:"VAULT_LOGIN_ROLE" json:"role"`
	LoginPath string `env:"VAULT_LOGIN_PATH" json:"login_path"`
	Addr      string `env:"VAULT_ADDR" json:"addr"`

	Token string `env:"VAULT_TOKEN"`
}

func NewHashiCorpVault(ctx context.Context) (*HashiCorpVault, error) {
	hcv := HashiCorpVaultConfig{}
	if err := envconfig.Process(context.Background(), &hcv); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	hv := &HashiCorpVault{
		client: client,
	}

	if IsRunningInKubernetes() {
		loginInfo, err := hv.kubernetesAuth(hcv.Role, hcv.LoginPath)
		if err != nil {
			return nil, err
		}

		if loginInfo == nil {
			return nil, fmt.Errorf("login secret not exists")
		}

		go func() {
			err := manageTokenLifecycle(ctx, hv.client, loginInfo)
			if err != nil {
				log.Printf("error managing token lifecycle: %v", err)
			}
		}()
	}

	hv.client.SetToken(hcv.Token)

	return hv, nil
}

func IsRunningInKubernetes() bool {
	_, hostExists := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	_, portExists := os.LookupEnv("KUBERNETES_SERVICE_PORT")
	return hostExists && portExists
}

func (hv *HashiCorpVault) GetSecret(ctx context.Context, name string) (string, error) {
	u, err := url.Parse(name)
	if err != nil {
		return "", fmt.Errorf("failed to parse name: %w", err)
	}

	name, version := u.Path, u.Query().Get("version")
	if version == "" {
		version = "1"
	}

	secret, err := hv.client.Logical().ReadWithData(name, map[string][]string{
		"version": {version},
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("secret data is nil")
	}

	dataRaw, ok := secret.Data["data"]
	if !ok {
		if username, ok := secret.Data["username"].(string); ok {
			if password, ok := secret.Data["password"].(string); ok {
				return fmt.Sprintf("%s:%s", username, password), nil
			}
		}

		return "", fmt.Errorf("missing 'data' key")

	}

	rewable, err := secret.TokenIsRenewable()
	if err != nil {
		return "", fmt.Errorf("failed to check token renewability: %w", err)
	}

	if rewable {
		go manageTokenLifecycle(ctx, hv.client, secret)
	}

	data, ok := dataRaw.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("data is not a map")
	}

	valueRaw, ok := data["value"]
	if !ok {
		return "", fmt.Errorf("missing 'value' key")
	}

	switch typ := valueRaw.(type) {
	case string:
		return typ, nil
	case []byte:
		return string(typ), nil
	case bool:
		return strconv.FormatBool(typ), nil
	case json.Number:
		return typ.String(), nil
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", typ), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", typ), nil
	default:
		return "", fmt.Errorf("found secret %v, but is of unknown type %T", name, typ)
	}
}

func (hv *HashiCorpVault) GetDatabaseCredentials(ctx context.Context, path string) (map[string]string, error) {
	secret, err := hv.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Data == nil {
		fmt.Println(secret)
		fmt.Println(err)
		return nil, fmt.Errorf("secret data is nil")
	}

	return map[string]string{
		"username": secret.Data["username"].(string),
		"password": secret.Data["password"].(string),
	}, nil
}

func (hv *HashiCorpVault) kubernetesAuth(role string, loginPath string) (*api.Secret, error) {
	k8sAuth, err := auth.NewKubernetesAuth(role, auth.WithMountPath(loginPath))
	if err != nil {
		return nil, err
	}

	authInfo, err := hv.client.Auth().Login(context.Background(), k8sAuth)
	if err != nil {
		return authInfo, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}

	if authInfo == nil {
		return authInfo, fmt.Errorf("no auth info was returned after login")
	}

	return authInfo, nil
}

func manageTokenLifecycle(ctx context.Context, client *api.Client, token *api.Secret) error {
	renew := token.Auth.Renewable
	if !renew {
		log.Printf("Token is not configured to be renewable. Re-attempting login.")
		return nil
	}

	watcher, err := client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
		Secret:    token,
		Increment: 3600,
	})
	if err != nil {
		return fmt.Errorf("unable to initialize new lifetime watcher for renewing auth token: %w", err)
	}

	go watcher.Start()
	defer watcher.Stop()

	for {
		select {
		case err := <-watcher.DoneCh():
			if err != nil {
				log.Printf("Failed to renew token: %v. Re-attempting login.", err)
				return nil
			}
			log.Printf("Token can no longer be renewed. Re-attempting login.")
			return nil
		case <-ctx.Done():
			log.Printf("Context cancelled. Re-attempting login.")
			cancelCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			err := client.Sys().RevokeWithContext(cancelCtx, token.LeaseID)
			if err != nil {
				log.Printf("Failed to revoke token: %v", err)
				return err
			}
		case renewal := <-watcher.RenewCh():
			log.Printf("Successfully renewed: %#v", renewal)
		}
	}
}
