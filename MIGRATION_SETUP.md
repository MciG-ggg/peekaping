# Database Migration Setup

This project includes a separate migration container that gives you full control over when database migrations run.

## How it works

1. **Migration Container**: The project uses a dedicated migration container built from `Dockerfile.migrate`
2. **Migration Tool**: Uses the bun migration tool located in `apps/server/cmd/bun/`
3. **Separate Control**: Migrations run in their own container, separate from the application server
4. **Server Dependency**: The server waits for migrations to complete successfully before starting

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
DB_TYPE=postgres
DB_HOST=postgres
DB_PORT=5432
DB_NAME=peekaping
DB_USER=postgres
DB_PASS=password

# Server Configuration
SERVER_PORT=8033
```

## Running with Docker Compose

### Development (PostgreSQL)
```bash
docker-compose -f docker-compose.dev.postgres.yml up --build
```

### Production (PostgreSQL)
```bash
docker-compose -f docker-compose.prod.postgres.yml up --build
```

### MongoDB (No migrations needed)
```bash
# Development
docker-compose -f docker-compose.dev.mongo.yml up --build

# Production
docker-compose -f docker-compose.prod.mongo.yml up --build
```

## Controlling Migration Execution

The migration container runs automatically when you start the full stack, but you have several options for control:

### 1. Run migrations only
```bash
# Run just the database and migration services
docker-compose -f docker-compose.dev.postgres.yml up postgres migrate
```

### 2. Skip migrations and run everything else
```bash
# Start database first
docker-compose -f docker-compose.dev.postgres.yml up -d postgres

# Run migrations manually when ready
docker-compose -f docker-compose.dev.postgres.yml up migrate

# Start the rest of the stack
docker-compose -f docker-compose.dev.postgres.yml up server web
```

### 3. Re-run migrations
```bash
# Remove the old migration container
docker-compose -f docker-compose.dev.postgres.yml rm -f migrate

# Run migrations again
docker-compose -f docker-compose.dev.postgres.yml up migrate
```

## Manual Migration Commands

If you need to run migrations manually outside of Docker:

1. **Build the bun migration tool**:
   ```bash
   cd apps/server
   go build -o bun ./cmd/bun
   ```

2. **Run migrations**:
   ```bash
   ./bun db migrate
   ```

3. **Check migration status**:
   ```bash
   ./bun db status
   ```

4. **Rollback last migration**:
   ```bash
   ./bun db rollback
   ```

### Creating New Migrations

For creating new migrations, use the bun tool:

1. **Build the bun tool**:
   ```bash
   cd apps/server
   go build -o bun ./cmd/bun
   ```

2. **Create new migration**:
   ```bash
   ./bun db create_tx_sql migration_name
   ```

## Migration Files

Migrations are stored in `apps/server/cmd/bun/migrations/` and follow the naming convention:
- `YYYYMMDDHHMMSS_description.tx.up.sql` - Migration to apply
- `YYYYMMDDHHMMSS_description.tx.down.sql` - Migration to rollback

## Troubleshooting

1. **Migration container fails**:
   - Check that the database service is healthy and accessible
   - Verify your `.env` file has correct database credentials
   - Check migration container logs: `docker-compose logs migrate`

2. **Server won't start**:
   - Ensure the migration container completed successfully
   - Check if migration service is in "exited" state: `docker-compose ps`

3. **Connection refused**:
   - Verify database service health: `docker-compose ps postgres`
   - Check if the database port is correctly exposed

4. **Need to re-run migrations**:
   - Remove the migration container: `docker-compose rm -f migrate`
   - Run migrations again: `docker-compose up migrate`

## Migration Container Details

- **Image**: Built from `apps/server/Dockerfile.migrate`
- **Tool**: Uses bun migration tool (`./bun db migrate`)
- **Restart Policy**: `no` (runs once and exits)
- **Dependencies**: Waits for database to be healthy before running
- **Server Dependency**: Server waits for migration container to complete successfully

## Database Support

The migration system supports:
- PostgreSQL (recommended)
- MySQL
- SQLite

The startup script automatically detects the database type and waits appropriately.
