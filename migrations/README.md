# Migrations

SQL migration files live in this directory.

Files must follow the naming convention expected by golang-migrate:

```
{version}_{title}.up.sql
{version}_{title}.down.sql
```

Example:

```
000001_create_users_table.up.sql
000001_create_users_table.down.sql
```

## Regenerate

Run migrations with:

```bash
make migrate DATABASE_URL="postgres://user:pass@localhost:5432/arena?sslmode=disable"
```

Or directly:

```bash
go run ./cmd migrate --database-url "postgres://user:pass@localhost:5432/arena?sslmode=disable"
```
