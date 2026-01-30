package migration

import (
        "database/sql"
        "strings"
        "testing"
)

func TestCreateTable(t *testing.T) {
        sql := CreateTable("users", func(bp *Blueprint) {
                bp.ID()
                bp.String("name", 255)
                bp.String("email", 255).Unique()
                bp.Timestamps()
        })

        if !strings.Contains(sql, "CREATE TABLE IF NOT EXISTS users") {
                t.Error("Expected CREATE TABLE statement")
        }
        if !strings.Contains(sql, "id SERIAL PRIMARY KEY") {
                t.Error("Expected ID column")
        }
        if !strings.Contains(sql, "name VARCHAR(255)") {
                t.Error("Expected name column")
        }
        if !strings.Contains(sql, "email VARCHAR(255) UNIQUE") {
                t.Error("Expected email column with UNIQUE")
        }
        if !strings.Contains(sql, "created_at TIMESTAMP") {
                t.Error("Expected created_at column")
        }
        if !strings.Contains(sql, "updated_at TIMESTAMP") {
                t.Error("Expected updated_at column")
        }
}

func TestBlueprintID(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.ID()

        if len(bp.columns) != 1 {
                t.Errorf("Expected 1 column, got %d", len(bp.columns))
        }
        if bp.columns[0] != "id SERIAL PRIMARY KEY" {
                t.Errorf("Expected 'id SERIAL PRIMARY KEY', got '%s'", bp.columns[0])
        }
}

func TestBlueprintString(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.String("name", 100)

        if len(bp.columns) != 1 {
                t.Errorf("Expected 1 column, got %d", len(bp.columns))
        }
        if bp.columns[0] != "name VARCHAR(100)" {
                t.Errorf("Expected 'name VARCHAR(100)', got '%s'", bp.columns[0])
        }
}

func TestBlueprintText(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.Text("description")

        if bp.columns[0] != "description TEXT" {
                t.Errorf("Expected 'description TEXT', got '%s'", bp.columns[0])
        }
}

func TestBlueprintInteger(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.Integer("count")

        if bp.columns[0] != "count INTEGER" {
                t.Errorf("Expected 'count INTEGER', got '%s'", bp.columns[0])
        }
}

func TestBlueprintBigInteger(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.BigInteger("big_count")

        if bp.columns[0] != "big_count BIGINT" {
                t.Errorf("Expected 'big_count BIGINT', got '%s'", bp.columns[0])
        }
}

func TestBlueprintFloat(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.Float("price")

        if bp.columns[0] != "price REAL" {
                t.Errorf("Expected 'price REAL', got '%s'", bp.columns[0])
        }
}

func TestBlueprintDecimal(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.Decimal("amount", 10, 2)

        if bp.columns[0] != "amount DECIMAL(10, 2)" {
                t.Errorf("Expected 'amount DECIMAL(10, 2)', got '%s'", bp.columns[0])
        }
}

func TestBlueprintBoolean(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.Boolean("active")

        if bp.columns[0] != "active BOOLEAN" {
                t.Errorf("Expected 'active BOOLEAN', got '%s'", bp.columns[0])
        }
}

func TestBlueprintDate(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.Date("birth_date")

        if bp.columns[0] != "birth_date DATE" {
                t.Errorf("Expected 'birth_date DATE', got '%s'", bp.columns[0])
        }
}

func TestBlueprintDateTime(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.DateTime("event_time")

        if bp.columns[0] != "event_time TIMESTAMP" {
                t.Errorf("Expected 'event_time TIMESTAMP', got '%s'", bp.columns[0])
        }
}

func TestBlueprintTimestamps(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.Timestamps()

        if len(bp.columns) != 2 {
                t.Errorf("Expected 2 columns, got %d", len(bp.columns))
        }
        if !strings.Contains(bp.columns[0], "created_at") {
                t.Error("Expected created_at column")
        }
        if !strings.Contains(bp.columns[1], "updated_at") {
                t.Error("Expected updated_at column")
        }
}

func TestBlueprintSoftDeletes(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.SoftDeletes()

        if bp.columns[0] != "deleted_at TIMESTAMP NULL" {
                t.Errorf("Expected 'deleted_at TIMESTAMP NULL', got '%s'", bp.columns[0])
        }
}

func TestBlueprintJSON(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.JSON("metadata")

        if bp.columns[0] != "metadata JSONB" {
                t.Errorf("Expected 'metadata JSONB', got '%s'", bp.columns[0])
        }
}

func TestBlueprintUUID(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.UUID("uuid")

        if bp.columns[0] != "uuid UUID" {
                t.Errorf("Expected 'uuid UUID', got '%s'", bp.columns[0])
        }
}

func TestBlueprintForeignID(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.ForeignID("user_id")

        if bp.columns[0] != "user_id BIGINT" {
                t.Errorf("Expected 'user_id BIGINT', got '%s'", bp.columns[0])
        }
}

func TestColumnBuilderNullable(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.String("bio", 500).Nullable()

        if !strings.Contains(bp.columns[0], "NULL") {
                t.Error("Expected column to be nullable")
        }
}

func TestColumnBuilderNotNull(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.String("name", 100).NotNull()

        if !strings.Contains(bp.columns[0], "NOT NULL") {
                t.Error("Expected column to be NOT NULL")
        }
}

func TestColumnBuilderDefault(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.Boolean("active").Default(true)

        if !strings.Contains(bp.columns[0], "DEFAULT true") {
                t.Errorf("Expected DEFAULT true, got: %s", bp.columns[0])
        }
}

func TestColumnBuilderDefaultString(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.String("status", 50).Default("pending")

        if !strings.Contains(bp.columns[0], "DEFAULT 'pending'") {
                t.Errorf("Expected DEFAULT 'pending', got: %s", bp.columns[0])
        }
}

func TestColumnBuilderUnique(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.String("email", 255).Unique()

        if !strings.Contains(bp.columns[0], "UNIQUE") {
                t.Error("Expected column to be UNIQUE")
        }
}

func TestColumnBuilderPrimary(t *testing.T) {
        bp := &Blueprint{tableName: "test"}
        bp.String("code", 10).Primary()

        if !strings.Contains(bp.columns[0], "PRIMARY KEY") {
                t.Error("Expected column to be PRIMARY KEY")
        }
}

func TestColumnBuilderReferences(t *testing.T) {
        bp := &Blueprint{tableName: "posts", foreigns: make([]string, 0)}
        bp.ForeignID("user_id").References("users", "id")

        if len(bp.foreigns) != 1 {
                t.Errorf("Expected 1 foreign key, got %d", len(bp.foreigns))
        }
        if !strings.Contains(bp.foreigns[0], "FOREIGN KEY (user_id)") {
                t.Error("Expected FOREIGN KEY clause")
        }
        if !strings.Contains(bp.foreigns[0], "REFERENCES users(id)") {
                t.Error("Expected REFERENCES clause")
        }
}

func TestColumnBuilderOnDelete(t *testing.T) {
        bp := &Blueprint{tableName: "posts", foreigns: make([]string, 0)}
        bp.ForeignID("user_id").References("users", "id").OnDelete("CASCADE")

        if !strings.Contains(bp.foreigns[0], "ON DELETE CASCADE") {
                t.Errorf("Expected ON DELETE CASCADE, got: %s", bp.foreigns[0])
        }
}

func TestDropTable(t *testing.T) {
        sql := DropTable("users")
        expected := "DROP TABLE IF EXISTS users"
        if sql != expected {
                t.Errorf("Expected '%s', got '%s'", expected, sql)
        }
}

func TestDropTableCascade(t *testing.T) {
        sql := DropTableCascade("users")
        expected := "DROP TABLE IF EXISTS users CASCADE"
        if sql != expected {
                t.Errorf("Expected '%s', got '%s'", expected, sql)
        }
}

func TestMigrationName(t *testing.T) {
        name := MigrationName("create_users_table")
        if !strings.Contains(name, "create_users_table") {
                t.Error("Expected migration name to contain 'create_users_table'")
        }
        if !strings.Contains(name, "_") {
                t.Error("Expected migration name to have timestamp prefix")
        }
}

func TestMigration(t *testing.T) {
        m := Migration{
                Name: "create_test_table",
                Up: func(db *sql.DB) error {
                        return nil
                },
                Down: func(db *sql.DB) error {
                        return nil
                },
        }

        if m.Name != "create_test_table" {
                t.Errorf("Expected name to be 'create_test_table', got '%s'", m.Name)
        }
}

func TestMigrationStatus(t *testing.T) {
        status := MigrationStatus{
                Name: "create_users_table",
                Ran:  true,
        }

        if !status.Ran {
                t.Error("Expected Ran to be true")
        }
}

func TestComplexTableCreation(t *testing.T) {
        sql := CreateTable("posts", func(bp *Blueprint) {
                bp.ID()
                bp.ForeignID("user_id").References("users", "id").OnDelete("CASCADE")
                bp.String("title", 255).NotNull()
                bp.Text("content").Nullable()
                bp.Boolean("published").Default(false)
                bp.JSON("metadata").Nullable()
                bp.Timestamps()
                bp.SoftDeletes()
        })

        expectations := []string{
                "CREATE TABLE IF NOT EXISTS posts",
                "id SERIAL PRIMARY KEY",
                "user_id BIGINT",
                "title VARCHAR(255) NOT NULL",
                "content TEXT NULL",
                "published BOOLEAN DEFAULT false",
                "metadata JSONB NULL",
                "created_at TIMESTAMP",
                "updated_at TIMESTAMP",
                "deleted_at TIMESTAMP NULL",
                "FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
        }

        for _, exp := range expectations {
                if !strings.Contains(sql, exp) {
                        t.Errorf("Expected SQL to contain '%s'", exp)
                }
        }
}

func TestColumnBuilderChaining(t *testing.T) {
        bp := &Blueprint{tableName: "test"}

        bp.String("email", 255).NotNull().Unique()

        column := bp.columns[0]
        if !strings.Contains(column, "NOT NULL") {
                t.Error("Expected NOT NULL")
        }
        if !strings.Contains(column, "UNIQUE") {
                t.Error("Expected UNIQUE")
        }
}
