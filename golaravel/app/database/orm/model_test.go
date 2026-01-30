package orm

import (
	"testing"
)

func TestQueryBuilderSelect(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.Select("id", "name", "email")

	sql := qb.ToSQL()
	expected := "SELECT id, name, email FROM users"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderWhere(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.Where("id", 1)

	sql := qb.ToSQL()
	expected := "SELECT * FROM users WHERE id = ?"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}

	if len(qb.bindings) != 1 || qb.bindings[0] != 1 {
		t.Error("Expected binding to be added")
	}
}

func TestQueryBuilderWhereWithOperator(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.Where("age", ">", 18)

	sql := qb.ToSQL()
	expected := "SELECT * FROM users WHERE age > ?"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderMultipleWhere(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.Where("status", "active").Where("role", "admin")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users WHERE status = ? AND role = ?"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}

	if len(qb.bindings) != 2 {
		t.Errorf("Expected 2 bindings, got %d", len(qb.bindings))
	}
}

func TestQueryBuilderOrWhere(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.Where("status", "active").OrWhere("role", "admin")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users WHERE status = ? OR role = ?"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderWhereNull(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.WhereNull("deleted_at")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users WHERE deleted_at IS NULL"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderWhereNotNull(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.WhereNotNull("email_verified_at")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users WHERE email_verified_at IS NOT NULL"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderWhereIn(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.WhereIn("id", []interface{}{1, 2, 3})

	sql := qb.ToSQL()
	expected := "SELECT * FROM users WHERE id IN (?, ?, ?)"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}

	if len(qb.bindings) != 3 {
		t.Errorf("Expected 3 bindings, got %d", len(qb.bindings))
	}
}

func TestQueryBuilderOrderBy(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.OrderBy("created_at", "DESC")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users ORDER BY created_at DESC"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderMultipleOrderBy(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.OrderBy("role", "ASC").OrderBy("created_at", "DESC")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users ORDER BY role ASC, created_at DESC"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderLimit(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.Limit(10)

	sql := qb.ToSQL()
	expected := "SELECT * FROM users LIMIT 10"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderOffset(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.Limit(10).Offset(20)

	sql := qb.ToSQL()
	expected := "SELECT * FROM users LIMIT 10 OFFSET 20"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderJoin(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
		joins:    make([]joinClause, 0),
	}

	qb.Join("posts", "users.id", "=", "posts.user_id")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users INNER JOIN posts ON users.id = posts.user_id"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderLeftJoin(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
		joins:    make([]joinClause, 0),
	}

	qb.LeftJoin("posts", "users.id", "=", "posts.user_id")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users LEFT JOIN posts ON users.id = posts.user_id"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderRightJoin(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
		joins:    make([]joinClause, 0),
	}

	qb.RightJoin("posts", "users.id", "=", "posts.user_id")

	sql := qb.ToSQL()
	expected := "SELECT * FROM users RIGHT JOIN posts ON users.id = posts.user_id"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderGroupBy(t *testing.T) {
	qb := &QueryBuilder{
		table:    "orders",
		columns:  []string{"user_id", "COUNT(*) as total"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.GroupByColumn("user_id")

	sql := qb.ToSQL()
	expected := "SELECT user_id, COUNT(*) as total FROM orders GROUP BY user_id"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderHaving(t *testing.T) {
	qb := &QueryBuilder{
		table:    "orders",
		columns:  []string{"user_id", "COUNT(*) as total"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
	}

	qb.GroupByColumn("user_id").Having("COUNT(*) > 5")

	sql := qb.ToSQL()
	expected := "SELECT user_id, COUNT(*) as total FROM orders GROUP BY user_id HAVING COUNT(*) > 5"
	if sql != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sql)
	}
}

func TestQueryBuilderComplexQuery(t *testing.T) {
	qb := &QueryBuilder{
		table:    "users",
		columns:  []string{"*"},
		wheres:   make([]whereClause, 0),
		orderBys: make([]orderByClause, 0),
		bindings: make([]interface{}, 0),
		joins:    make([]joinClause, 0),
	}

	qb.Select("users.name", "posts.title").
		LeftJoin("posts", "users.id", "=", "posts.user_id").
		Where("users.status", "active").
		Where("posts.published", true).
		OrderBy("posts.created_at", "DESC").
		Limit(10).
		Offset(5)

	sql := qb.ToSQL()

	if len(sql) == 0 {
		t.Error("Expected SQL to be generated")
	}

	if len(qb.bindings) != 2 {
		t.Errorf("Expected 2 bindings, got %d", len(qb.bindings))
	}
}

func TestPaginatorHasMorePages(t *testing.T) {
	p := &Paginator{
		Total:       100,
		PerPage:     10,
		CurrentPage: 5,
		LastPage:    10,
	}

	if !p.HasMorePages() {
		t.Error("Expected HasMorePages to return true")
	}

	p.CurrentPage = 10
	if p.HasMorePages() {
		t.Error("Expected HasMorePages to return false on last page")
	}
}

func TestPaginatorHasPages(t *testing.T) {
	p := &Paginator{
		Total:       100,
		PerPage:     10,
		CurrentPage: 1,
		LastPage:    10,
	}

	if !p.HasPages() {
		t.Error("Expected HasPages to return true")
	}

	p.LastPage = 1
	if p.HasPages() {
		t.Error("Expected HasPages to return false for single page")
	}
}

func TestModel(t *testing.T) {
	m := Model{
		ID: 1,
	}

	if m.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", m.ID)
	}
}
