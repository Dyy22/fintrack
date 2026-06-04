#!/usr/bin/env sh
set -eu

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.dev.yml}"
DB_NAME="${DB_NAME:-fintrack}"
DB_USER="${DB_USER:-fintrack}"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-migrations}"

if ! command -v docker >/dev/null 2>&1; then
  echo "Missing required command: docker" >&2
  exit 1
fi

if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "Migrations directory not found: $MIGRATIONS_DIR" >&2
  exit 1
fi

echo "Waiting for PostgreSQL..."
docker compose -f "$COMPOSE_FILE" exec -T postgres sh -c "until pg_isready -U $DB_USER -d $DB_NAME; do sleep 1; done"

docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -c "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())" >/dev/null

existing_users_table="$(docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT to_regclass('public.users')")"
existing_initial_migration="$(docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT version FROM schema_migrations WHERE version='000001_init'")"
if [ "$existing_users_table" = "users" ] && [ "$existing_initial_migration" != "000001_init" ]; then
  echo "Marking existing initial schema as applied: 000001_init"
  docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -c "INSERT INTO schema_migrations (version) VALUES ('000001_init') ON CONFLICT DO NOTHING" >/dev/null
fi

applied_any=0
for migration_file in "$MIGRATIONS_DIR"/*.up.sql; do
  if [ ! -f "$migration_file" ]; then
    continue
  fi

  migration_name="$(basename "$migration_file")"
  version="${migration_name%.up.sql}"
  existing_version="$(docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" -tAc "SELECT version FROM schema_migrations WHERE version='$version'")"
  if [ "$existing_version" = "$version" ]; then
    echo "Skipping already applied migration: $version"
    continue
  fi

  echo "Applying migration: $migration_file"
  docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 < "$migration_file"
  docker compose -f "$COMPOSE_FILE" exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -c "INSERT INTO schema_migrations (version) VALUES ('$version')" >/dev/null
  applied_any=1
done

if [ "$applied_any" = "0" ]; then
  echo "No pending migrations."
else
  echo "Migrations applied."
fi
