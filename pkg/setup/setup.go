package setup

import (
	"context"
	"errors"
	"fmt"
	"github.com/jamiewhitney/fairways-core/pkg/cache"
	"github.com/jamiewhitney/fairways-core/pkg/database"
	databasesql "github.com/jamiewhitney/fairways-core/pkg/database/mysql"
	"github.com/jamiewhitney/fairways-core/pkg/environment"
	internal_services "github.com/jamiewhitney/fairways-core/pkg/internal-services"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"github.com/jamiewhitney/fairways-core/pkg/pubsub"
	"github.com/jamiewhitney/fairways-core/pkg/secretstore"
	"github.com/jamiewhitney/fairways-core/pkg/transport"
	booking_pb "github.com/jamiewhitney/fairways-core/protobufs/booking"
	pricing_pb "github.com/jamiewhitney/fairways-core/protobufs/pricing"
	"github.com/sethvargo/go-envconfig"
	"os"
	"time"
)

type SecretStoreProvider interface {
	SecretStoreConfig() *secretstore.Config
}

type DatabaseConfigProvider interface {
	DatabaseConfig() *database.Config
}

type CacheConfigProvider interface {
	CacheConfig() *cache.Config
}

type PricingServiceProvider interface {
	PricingServiceConfig() internal_services.ServiceConfig
}

type BookingServiceProvider interface {
	BookingServiceConfig() internal_services.ServiceConfig
}

type PubsubConfigProvider interface {
	PubsubConfig() *pubsub.Config
}

func Setup(ctx context.Context, config interface{}) (*environment.Environment, error) {
	logger := logging.FromContext(ctx)

	var env environment.Environment

	if err := envconfig.Process(ctx, &env); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	var mutatorFuncs []envconfig.Mutator

	var serverEnvOpts []environment.Option

	if provider, ok := config.(SecretStoreProvider); ok {
		logger.Debug("configuring secret manager")

		ssConfig := provider.SecretStoreConfig()
		if err := envconfig.Process(ctx, ssConfig); err != nil {
			return nil, err
		}

		ss, err := secretstore.GetSecretStore(ctx, ssConfig)
		if err != nil {
			return nil, err
		}

		mutatorFuncs = append(mutatorFuncs, secretstore.Resolver(ss))

		serverEnvOpts = append(serverEnvOpts, environment.WithSecretStore(ss))

		logger.Debugf("secret manager configured %+v", *ssConfig)
	}

	if err := envconfig.Process(ctx, config, mutatorFuncs...); err != nil {
		return nil, err
	}

	if _, ok := config.(PricingServiceProvider); ok {
		logger.Debug("configuring pricing service")

		addr, found := os.LookupEnv("PRICING_ADDR")
		if !found {
			return nil, errors.New("service address not found for pricing service")
		}

		conn, err := transport.NewClient(ctx, addr, nil)
		if err != nil {
			return nil, err
		}

		serverEnvOpts = append(serverEnvOpts, environment.WithPricingServiceClient(pricing_pb.NewPricingServiceClient(conn)))
	}

	if _, ok := config.(BookingServiceProvider); ok {
		logger.Debug("configuring booking service")

		addr, found := os.LookupEnv("BOOKING_ADDR")
		if !found {
			return nil, errors.New("service address not found for booking service")
		}

		conn, err := transport.NewClient(ctx, addr, nil)
		if err != nil {
			return nil, err
		}

		serverEnvOpts = append(serverEnvOpts, environment.WithBookingServiceClient(booking_pb.NewBookingServiceClient(conn)))
	}

	if provider, ok := config.(PubsubConfigProvider); ok {
		logger.Debug("configuring pubsub")

		pubsubConfig := provider.PubsubConfig()

		if err := envconfig.Process(ctx, pubsubConfig); err != nil {
			return nil, err
		}

		ps, err := pubsub.New(ctx)
		if err != nil {
			return nil, err
		}

		serverEnvOpts = append(serverEnvOpts, environment.WithPubsubConfig(ps))

		logger.Debugf("pubsub configured %+v", pubsubConfig)
	}

	if provider, ok := config.(CacheConfigProvider); ok {
		logger.Debug("configuring cache")
		cacheConfig := provider.CacheConfig()

		if err := envconfig.Process(ctx, cacheConfig); err != nil {
			return nil, err
		}

		cc, err := cache.New(fmt.Sprintf("%s:%s", provider.CacheConfig().RedisHost, provider.CacheConfig().RedisPort), provider.CacheConfig().RedisUsername, provider.CacheConfig().RedisPassword, 1*time.Minute)
		if err != nil {
			return nil, err
		}

		serverEnvOpts = append(serverEnvOpts, environment.WithCache(*cc))

		logger.Debugf("cache configured %+v", cacheConfig)
	}

	if provider, ok := config.(DatabaseConfigProvider); ok {
		logger.Debug("configuring database connection")
		dbConfig := provider.DatabaseConfig()

		if err := envconfig.Process(ctx, dbConfig, mutatorFuncs...); err != nil {
			return nil, err
		}
		db, err := databasesql.NewFromConfig(dbConfig)
		if err != nil {
			return nil, err
		}

		serverEnvOpts = append(serverEnvOpts, environment.WithDatabase(db))
	}

	return environment.New(&env, serverEnvOpts...), nil
}
