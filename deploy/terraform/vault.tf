locals {
  hour_in_seconds = 3600
  day_in_seconds  = 86400
  services = ["booking", "catalog", "pricing", "tee_time"]
}

resource "vault_mount" "this" {
  path = "db"
  type = "database"
}

resource "vault_database_secret_backend_connection" "this" {
  backend = vault_mount.this.path
  name    = "local"
  allowed_roles = ["booking", "catalog", "pricing", "tee_time"]

  mysql {
    connection_url = "root:root@tcp(database:3306)/"
  }
}

resource "vault_policy" "policy" {
  name   = "service"
  policy = <<EOF
path "db/creds/{{identity.entity.aliases.kubernetes.metadata.service_account_name}}/*" {
  capabilities = ["read", "list"]
}
EOF
}

resource "vault_database_secret_backend_role" "role" {
  count = length(local.services)

  backend = vault_mount.this.path
  name    = local.services[count.index]
  db_name = vault_database_secret_backend_connection.this.name
  creation_statements = [
    "CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT, INSERT, UPDATE ON ${local.services[count.index]}.* TO '{{name}}'@'%';"
  ]
  default_ttl = local.hour_in_seconds
  max_ttl     = local.day_in_seconds
}