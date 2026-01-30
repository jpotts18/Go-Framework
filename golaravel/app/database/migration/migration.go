package migration

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type Migration struct {
	Name string
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

type Migrator struct {
	db         *sql.DB
	migrations []Migration
	tableName  string
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
		tableName:  "migrations",
	}
}

func (m *Migrator) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) createMigrationsTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			migration VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`, m.tableName)

	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) getRanMigrations() ([]string, error) {
	query := fmt.Sprintf("SELECT migration FROM %s ORDER BY batch, migration", m.tableName)
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations = append(migrations, name)
	}

	return migrations, rows.Err()
}

func (m *Migrator) getLastBatch() (int, error) {
	query := fmt.Sprintf("SELECT COALESCE(MAX(batch), 0) FROM %s", m.tableName)
	var batch int
	err := m.db.QueryRow(query).Scan(&batch)
	return batch, err
}

func (m *Migrator) Run() error {
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	ranMigrations, err := m.getRanMigrations()
	if err != nil {
		return fmt.Errorf("failed to get ran migrations: %w", err)
	}

	ranMap := make(map[string]bool)
	for _, name := range ranMigrations {
		ranMap[name] = true
	}

	batch, err := m.getLastBatch()
	if err != nil {
		return fmt.Errorf("failed to get last batch: %w", err)
	}
	batch++

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Name < m.migrations[j].Name
	})

	for _, migration := range m.migrations {
		if ranMap[migration.Name] {
			continue
		}

		log.Printf("Migrating: %s", migration.Name)

		if err := migration.Up(m.db); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.Name, err)
		}

		insertQuery := fmt.Sprintf("INSERT INTO %s (migration, batch) VALUES ($1, $2)", m.tableName)
		if _, err := m.db.Exec(insertQuery, migration.Name, batch); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}

		log.Printf("Migrated: %s", migration.Name)
	}

	return nil
}

func (m *Migrator) Rollback() error {
	batch, err := m.getLastBatch()
	if err != nil {
		return fmt.Errorf("failed to get last batch: %w", err)
	}

	if batch == 0 {
		log.Println("Nothing to rollback")
		return nil
	}

	query := fmt.Sprintf("SELECT migration FROM %s WHERE batch = $1 ORDER BY migration DESC", m.tableName)
	rows, err := m.db.Query(query, batch)
	if err != nil {
		return fmt.Errorf("failed to get migrations for batch: %w", err)
	}
	defer rows.Close()

	var migrationsToRollback []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		migrationsToRollback = append(migrationsToRollback, name)
	}

	migrationMap := make(map[string]Migration)
	for _, migration := range m.migrations {
		migrationMap[migration.Name] = migration
	}

	for _, name := range migrationsToRollback {
		migration, exists := migrationMap[name]
		if !exists {
			log.Printf("Warning: Migration %s not found in registered migrations", name)
			continue
		}

		log.Printf("Rolling back: %s", name)

		if err := migration.Down(m.db); err != nil {
			return fmt.Errorf("rollback of %s failed: %w", name, err)
		}

		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE migration = $1", m.tableName)
		if _, err := m.db.Exec(deleteQuery, name); err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", name, err)
		}

		log.Printf("Rolled back: %s", name)
	}

	return nil
}

func (m *Migrator) Reset() error {
	for {
		batch, err := m.getLastBatch()
		if err != nil {
			return err
		}
		if batch == 0 {
			break
		}
		if err := m.Rollback(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrator) Refresh() error {
	if err := m.Reset(); err != nil {
		return err
	}
	return m.Run()
}

func (m *Migrator) Status() ([]MigrationStatus, error) {
	if err := m.createMigrationsTable(); err != nil {
		return nil, err
	}

	ranMigrations, err := m.getRanMigrations()
	if err != nil {
		return nil, err
	}

	ranMap := make(map[string]bool)
	for _, name := range ranMigrations {
		ranMap[name] = true
	}

	var statuses []MigrationStatus
	for _, migration := range m.migrations {
		statuses = append(statuses, MigrationStatus{
			Name: migration.Name,
			Ran:  ranMap[migration.Name],
		})
	}

	return statuses, nil
}

type MigrationStatus struct {
	Name string
	Ran  bool
}

type Schema struct {
	columns []ColumnDefinition
}

type ColumnDefinition struct {
	Name       string
	Type       string
	Nullable   bool
	Default    interface{}
	Primary    bool
	Unique     bool
	Index      bool
	References *ForeignKey
}

type ForeignKey struct {
	Table      string
	Column     string
	OnDelete   string
	OnUpdate   string
}

func CreateTable(name string, callback func(*Blueprint)) string {
	bp := &Blueprint{tableName: name}
	callback(bp)
	return bp.toSQL()
}

type Blueprint struct {
	tableName string
	columns   []string
	indexes   []string
	foreigns  []string
}

func (b *Blueprint) ID() {
	b.columns = append(b.columns, "id SERIAL PRIMARY KEY")
}

func (b *Blueprint) String(name string, length int) *ColumnBuilder {
	col := fmt.Sprintf("%s VARCHAR(%d)", name, length)
	b.columns = append(b.columns, col)
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) Text(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" TEXT")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) Integer(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" INTEGER")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) BigInteger(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" BIGINT")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) Float(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" REAL")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) Decimal(name string, precision, scale int) *ColumnBuilder {
	col := fmt.Sprintf("%s DECIMAL(%d, %d)", name, precision, scale)
	b.columns = append(b.columns, col)
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) Boolean(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" BOOLEAN")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) Date(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" DATE")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) DateTime(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" TIMESTAMP")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) Timestamps() {
	b.columns = append(b.columns, "created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	b.columns = append(b.columns, "updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
}

func (b *Blueprint) SoftDeletes() {
	b.columns = append(b.columns, "deleted_at TIMESTAMP NULL")
}

func (b *Blueprint) JSON(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" JSONB")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) UUID(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" UUID")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1}
}

func (b *Blueprint) ForeignID(name string) *ColumnBuilder {
	b.columns = append(b.columns, name+" BIGINT")
	return &ColumnBuilder{blueprint: b, index: len(b.columns) - 1, columnName: name}
}

func (b *Blueprint) toSQL() string {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", b.tableName)
	sql += strings.Join(b.columns, ",\n")

	if len(b.foreigns) > 0 {
		sql += ",\n" + strings.Join(b.foreigns, ",\n")
	}

	sql += "\n)"
	return sql
}

type ColumnBuilder struct {
	blueprint  *Blueprint
	index      int
	columnName string
}

func (cb *ColumnBuilder) Nullable() *ColumnBuilder {
	cb.blueprint.columns[cb.index] += " NULL"
	return cb
}

func (cb *ColumnBuilder) NotNull() *ColumnBuilder {
	cb.blueprint.columns[cb.index] += " NOT NULL"
	return cb
}

func (cb *ColumnBuilder) Default(value interface{}) *ColumnBuilder {
	switch v := value.(type) {
	case string:
		cb.blueprint.columns[cb.index] += fmt.Sprintf(" DEFAULT '%s'", v)
	default:
		cb.blueprint.columns[cb.index] += fmt.Sprintf(" DEFAULT %v", v)
	}
	return cb
}

func (cb *ColumnBuilder) Unique() *ColumnBuilder {
	cb.blueprint.columns[cb.index] += " UNIQUE"
	return cb
}

func (cb *ColumnBuilder) Primary() *ColumnBuilder {
	cb.blueprint.columns[cb.index] += " PRIMARY KEY"
	return cb
}

func (cb *ColumnBuilder) References(table, column string) *ColumnBuilder {
	foreign := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", cb.columnName, table, column)
	cb.blueprint.foreigns = append(cb.blueprint.foreigns, foreign)
	return cb
}

func (cb *ColumnBuilder) OnDelete(action string) *ColumnBuilder {
	if len(cb.blueprint.foreigns) > 0 {
		cb.blueprint.foreigns[len(cb.blueprint.foreigns)-1] += " ON DELETE " + action
	}
	return cb
}

func DropTable(name string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
}

func DropTableCascade(name string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", name)
}

func MigrationName(name string) string {
	return fmt.Sprintf("%s_%s", time.Now().Format("2006_01_02_150405"), name)
}
