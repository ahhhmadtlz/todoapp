package migrator

import (
	"database/sql"
	"fmt"
	"log"
	"todoapp/internal/repository/mysql"

	migrate "github.com/rubenv/sql-migrate"
)

const (
	defaultMigrationTable = "schema_migrations"
	defaultMigrationsDir  = "./internal/repository/mysql/migrations"
)

type Migrator struct {
	db             *sql.DB
	dialect        string
	migrations     *migrate.FileMigrationSource
	migrationTable string
}

type Config struct {
	MigrationsDir  string
	MigrationTable string
}

// New creates a new Migrator instance with an existing DB connection
func New(mysqlDB *mysql.MySQLDB, cfg Config) Migrator {
	migrationsDir := cfg.MigrationsDir
	if migrationsDir == "" {
		migrationsDir = defaultMigrationsDir
	}

	migrationTable := cfg.MigrationTable
	if migrationTable == "" {
		migrationTable = defaultMigrationTable
	}

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	return Migrator{
		db:             mysqlDB.Conn(),
		dialect:        "mysql",
		migrations:     migrations,
		migrationTable: migrationTable,
	}
}

// setMigrationTable configures the migration table name
func (m Migrator) setMigrationTable() {
	migrate.SetTable(m.migrationTable)
}

// Up applies all pending migrations or up to limit if specified
func (m Migrator) Up(limit int) error {
	m.setMigrationTable()

	var n int
	var err error

	if limit > 0 {
		n, err = migrate.ExecMax(m.db, m.dialect, m.migrations, migrate.Up, limit)
	} else {
		n, err = migrate.Exec(m.db, m.dialect, m.migrations, migrate.Up)
	}

	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if n == 0 {
		log.Println("No new migrations to apply")
	} else {
		log.Printf("Successfully applied %d migration(s)\n", n)
	}

	return nil
}

// Down rolls back migrations
func (m Migrator) Down(limit int) error {
	m.setMigrationTable()

	var n int
	var err error

	if limit > 0 {
		n, err = migrate.ExecMax(m.db, m.dialect, m.migrations, migrate.Down, limit)
	} else {
		// Rollback all
		n, err = migrate.Exec(m.db, m.dialect, m.migrations, migrate.Down)
	}

	if err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	if n == 0 {
		log.Println("No migrations to rollback")
	} else {
		log.Printf("Successfully rolled back %d migration(s)\n", n)
	}

	return nil
}

// Status returns the current migration status
func (m Migrator) Status() error {
	m.setMigrationTable()

	records, err := migrate.GetMigrationRecords(m.db, m.dialect)
	if err != nil {
		return fmt.Errorf("failed to get migration records: %w", err)
	}

	migrations, err := m.migrations.FindMigrations()
	if err != nil {
		return fmt.Errorf("failed to find migrations: %w", err)
	}

	// Create a map of applied migrations
	appliedMap := make(map[string]*migrate.MigrationRecord)
	for _, record := range records {
		appliedMap[record.Id] = record
	}

	fmt.Println("\nMigration Status:")
	fmt.Println("==================")

	if len(migrations) == 0 {
		fmt.Println("No migrations found")
		return nil
	}

	pendingCount := 0
	appliedCount := 0

	for _, migration := range migrations {
		if record, applied := appliedMap[migration.Id]; applied {
			fmt.Printf("✓ [Applied] %s (applied at: %s)\n", 
				migration.Id, 
				record.AppliedAt.Format("2006-01-02 15:04:05"))
			appliedCount++
		} else {
			fmt.Printf("✗ [Pending] %s\n", migration.Id)
			pendingCount++
		}
	}

	fmt.Println("==================")
	fmt.Printf("Total: %d | Applied: %d | Pending: %d\n\n", 
		len(migrations), appliedCount, pendingCount)

	return nil
}

// Redo rolls back and re-applies the last migration
func (m Migrator) Redo() error {
	log.Println("Rolling back last migration...")
	if err := m.Down(1); err != nil {
		return fmt.Errorf("failed to rollback: %w", err)
	}

	log.Println("Re-applying last migration...")
	if err := m.Up(1); err != nil {
		return fmt.Errorf("failed to re-apply: %w", err)
	}

	log.Println("Redo completed successfully")
	return nil
}

// Reset rolls back all migrations
func (m Migrator) Reset() error {
	log.Println("Rolling back all migrations...")
	if err := m.Down(0); err != nil {
		return fmt.Errorf("failed to reset: %w", err)
	}

	log.Println("Reset completed successfully")
	return nil
}

// Fresh rolls back all migrations and re-applies them
func (m Migrator) Fresh() error {
	if err := m.Reset(); err != nil {
		return err
	}

	log.Println("Re-applying all migrations...")
	if err := m.Up(0); err != nil {
		return fmt.Errorf("failed to re-apply migrations: %w", err)
	}

	log.Println("Fresh completed successfully")
	return nil
}