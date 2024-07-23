provider "kubernetes" {
  config_path = "~/.kube/config"

}

resource "kubernetes_service_account_v1" "vault-auth" {
  metadata {
    name = "vault-auth"
  }
  automount_service_account_token = true
}

resource "kubernetes_cluster_role_binding_v1" "example" {
  metadata {
    name = "role-tokenreview-binding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "system:auth-delegator"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.vault-auth.metadata[0].name
    namespace = kubernetes_service_account_v1.vault-auth.metadata[0].namespace
  }
}


resource "kubernetes_secret_v1" "sa_token" {
  metadata {
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.vault-auth.metadata[0].name
    }

    name = "vault-auth-secret"
  }

  type                           = "kubernetes.io/service-account-token"
  wait_for_service_account_token = true
}

resource "vault_auth_backend" "k8s" {
  type = "kubernetes"
  path = "kubernetes"
}

resource "vault_kubernetes_auth_backend_config" "k8s_auth_backend" {
  backend = vault_auth_backend.k8s.path

  kubernetes_host    = "https://host.docker.internal:62593"
  token_reviewer_jwt = kubernetes_secret_v1.sa_token.data["token"]
  kubernetes_ca_cert = <<EOF
-----BEGIN CERTIFICATE-----
MIIC/jCCAeagAwIBAgIBADANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwprdWJl
cm5ldGVzMB4XDTI0MDYyMDA4NTQyOVoXDTM0MDYxODA4NTQyOVowFTETMBEGA1UE
AxMKa3ViZXJuZXRlczCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALHF
jh3+H4d90NgIppsQ+3r8cj94k479m1vz6+Q58XhzhdKLuC1rzxPFDeimI/MEo7nI
xZuvXt2QqJ1vodEI6JrLSuvpzaM3V5J6KGIfGr4PgGKEKPaPo9TMHE86iMIVzbtI
VhG7oY8zLt92orNYeGIDCda1xI9zFFQPJ1xZ5endn1HnLSKibY55Q9mumcyyzhnk
wHPopPpqeAWEt/MtfoMTYOSuY7ZNNhyCfvJm18pf0cjppYg54ssWNmh0KvHCbqnZ
IFbV6/Rl3Sr1fVWtShD87AI1JczsIajHitpH8OdTl+JDrSw1j2XPCl5G5kV/jfdI
+aCTn6Ya2n/flhruePECAwEAAaNZMFcwDgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB
/wQFMAMBAf8wHQYDVR0OBBYEFF5Y8mXQNU4Oeo2bMTch8pqX3W4wMBUGA1UdEQQO
MAyCCmt1YmVybmV0ZXMwDQYJKoZIhvcNAQELBQADggEBAGZ+2twgZ02cvxo6XZiD
PUqCiBeFEiCWKAnSylGezNjTcsmVP72ZdCy/bl0kA4BoUIvttC2stiZmCC6fErWs
ZhYFxeOGOKwvRiLo2N2K5Eg4yK3+gbjWCMaYPGRUVmdaoKlsVOBl+ROBhLf5Pwod
PQGU7Er6wbhiW7fFGI5AIwFiuK1iMzumIzKVzI8cFZfzTqdirhR9z9xSzNHUMW5o
mrwpIvql0jLXZhDAyaQbLhopDNY1kJ3zh6Xpe/NiL9adAMi9xMeH+Gbx4yf3Yn2J
0/3RSccGJzaqbxrTtKr+1pZR4VXNr9MGURywOQc5S+uNhpN2DAYNTcMejPGaPIv7
6gQ=
-----END CERTIFICATE-----
  EOF
}

resource "vault_kubernetes_auth_backend_role" "k8s_auth_backend_role" {
  bound_service_account_names = ["vault-auth", "pricing", "booking", "catalog", "tee_time"]
  bound_service_account_namespaces = ["default", "pricing", "booking", "catalog", "tee_time"]
  token_policies = [vault_policy.policy.id]
  role_name = "default"
}