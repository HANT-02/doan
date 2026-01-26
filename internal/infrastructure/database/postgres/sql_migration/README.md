# Migration Files - SQL Database Migrations

## 📁 Structure

```
sql_migration/
├── 000001_enable_uuid.up.sql          # Enable UUID extension
├── 000001_enable_uuid.down.sql        # Disable UUID extension
├── 000002_create_users_table.up.sql   # Create users table
├── 000002_create_users_table.down.sql # Drop users table
├── 000003_seed_users.up.sql           # Insert sample data
└── 000003_seed_users.down.sql         # Remove sample data
```

## 🔢 Naming Convention

Files must follow the pattern: `{version}_{description}.{direction}.sql`

Example:
- `000001_enable_uuid.up.sql` ✅
- `1_create_table.up.sql` ❌ (wrong - needs zero padding)

## 🚀 Running Migrations

### Run all migrations:
```bash
make migrate
```

### Manual migration:
```bash
cd cmd/cli/migration
go run main.go wire_gen.go
```

### Check migration status:
Connect to database and check:
```bash
psql -U postgres -d doan -c "SELECT * FROM schema_migrations;"
```

## 📊 Sample Data

After running migrations, you will have:

**10 sample users:**
- `admin` - Active admin user
- `user1`, `user2` - Test users
- `testuser` - Testing account
- `john_doe`, `jane_smith`, `alice_jones`, `charlie_brown` - Active users
- `bob_wilson` - Inactive user
- `david_miller` - Suspended user

**Default password:** `password123` (bcrypt hashed)

## 🔐 Password Hashing

To generate a new bcrypt password hash:

```go
import "golang.org/x/crypto/bcrypt"

hash, _ := bcrypt.GenerateFromPassword([]byte("yourpassword"), bcrypt.DefaultCost)
fmt.Println(string(hash))
```

## 🔄 Rollback Migrations

To rollback last migration:
```bash
# You need to implement this in migration.go
# Or manually:
psql -U postgres -d doan < internal/infrastructure/database/postgres/sql_migration/000003_seed_users.down.sql
```

## ✅ Verify Data

```bash
# Check users table
psql -U postgres -d doan -c "SELECT id, username, status, created_at FROM users;"

# Count users
psql -U postgres -d doan -c "SELECT COUNT(*) FROM users;"

# Check UUID extension
psql -U postgres -d doan -c "SELECT * FROM pg_extension WHERE extname = 'uuid-ossp';"
```

## 📝 Adding New Migrations

1. Create new migration files with next version number:
   ```bash
   touch internal/infrastructure/database/postgres/sql_migration/000004_add_email_column.up.sql
   touch internal/infrastructure/database/postgres/sql_migration/000004_add_email_column.down.sql
   ```

2. Write SQL in `.up.sql` file:
   ```sql
   ALTER TABLE users ADD COLUMN email VARCHAR(255) UNIQUE;
   ```

3. Write rollback in `.down.sql` file:
   ```sql
   ALTER TABLE users DROP COLUMN email;
   ```

4. Run migration:
   ```bash
   make migrate
   ```

## 🛠️ Troubleshooting

### Migration says "no change" but table doesn't exist

This means migration thinks it's already run. Check:
```bash
psql -U postgres -d doan -c "SELECT * FROM schema_migrations;"
```

To force re-run, you may need to:
1. Drop the database
2. Recreate it
3. Run migrations again

Or manually reset migration version in `schema_migrations` table.

### UUID generation error

Make sure UUID extension is enabled:
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### Permission denied

Check PostgreSQL user permissions:
```sql
GRANT ALL PRIVILEGES ON DATABASE doan TO postgres;
```

## 📚 Migration Library

This project uses `golang-migrate/migrate`:
- GitHub: https://github.com/golang-migrate/migrate
- Docs: https://github.com/golang-migrate/migrate/tree/master/database/postgres

