package orm

import (
        "database/sql"
        "fmt"
        "reflect"
        "strings"
        "time"
)

type Model struct {
        ID        int64      `db:"id" json:"id"`
        CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
        UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type ModelInterface interface {
        TableName() string
        PrimaryKey() string
}

type Connection struct {
        db     *sql.DB
        driver string
}

var defaultConnection *Connection

func Connect(driver, dsn string) (*Connection, error) {
        db, err := sql.Open(driver, dsn)
        if err != nil {
                return nil, err
        }

        if err := db.Ping(); err != nil {
                return nil, err
        }

        conn := &Connection{
                db:     db,
                driver: driver,
        }

        if defaultConnection == nil {
                defaultConnection = conn
        }

        return conn, nil
}

func (c *Connection) DB() *sql.DB {
        return c.db
}

func (c *Connection) Close() error {
        return c.db.Close()
}

type Transaction struct {
        tx *sql.Tx
}

func (c *Connection) Begin() (*Transaction, error) {
        tx, err := c.db.Begin()
        if err != nil {
                return nil, err
        }
        return &Transaction{tx: tx}, nil
}

func (c *Connection) Transaction(fn func(tx *Transaction) error) error {
        tx, err := c.Begin()
        if err != nil {
                return err
        }
        
        if err := fn(tx); err != nil {
                tx.Rollback()
                return err
        }
        
        return tx.Commit()
}

func (t *Transaction) Commit() error {
        return t.tx.Commit()
}

func (t *Transaction) Rollback() error {
        return t.tx.Rollback()
}

func (t *Transaction) Query(tableName string) *QueryBuilder {
        return &QueryBuilder{
                tx:        t.tx,
                table:     tableName,
                columns:   []string{"*"},
                wheres:    make([]whereClause, 0),
                orderBys:  make([]orderByClause, 0),
                bindings:  make([]interface{}, 0),
        }
}

func (t *Transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
        return t.tx.Exec(query, args...)
}

func (t *Transaction) QueryRow(query string, args ...interface{}) *sql.Row {
        return t.tx.QueryRow(query, args...)
}

func DB() *Connection {
        return defaultConnection
}

func (c *Connection) Query(tableName string) *QueryBuilder {
        return &QueryBuilder{
                conn:      c,
                table:     tableName,
                columns:   []string{"*"},
                wheres:    make([]whereClause, 0),
                orderBys:  make([]orderByClause, 0),
                bindings:  make([]interface{}, 0),
        }
}

type QueryBuilder struct {
        conn      *Connection
        tx        *sql.Tx
        table     string
        columns   []string
        wheres    []whereClause
        orderBys  []orderByClause
        groupBy   []string
        having    string
        limit     int
        offset    int
        bindings  []interface{}
        joins     []joinClause
}

func (qb *QueryBuilder) executor() interface {
        Query(string, ...interface{}) (*sql.Rows, error)
        QueryRow(string, ...interface{}) *sql.Row
        Exec(string, ...interface{}) (sql.Result, error)
} {
        if qb.tx != nil {
                return qb.tx
        }
        return qb.conn.db
}

type whereClause struct {
        column   string
        operator string
        value    interface{}
        boolean  string
        isRaw    bool
}

type orderByClause struct {
        column    string
        direction string
}

type joinClause struct {
        joinType string
        table    string
        first    string
        operator string
        second   string
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
        qb.columns = columns
        return qb
}

func (qb *QueryBuilder) Where(column string, args ...interface{}) *QueryBuilder {
        operator := "="
        var value interface{}

        if len(args) == 1 {
                value = args[0]
        } else if len(args) >= 2 {
                operator = args[0].(string)
                value = args[1]
        }

        qb.wheres = append(qb.wheres, whereClause{
                column:   column,
                operator: operator,
                value:    value,
                boolean:  "AND",
        })
        qb.bindings = append(qb.bindings, value)

        return qb
}

func (qb *QueryBuilder) OrWhere(column string, args ...interface{}) *QueryBuilder {
        operator := "="
        var value interface{}

        if len(args) == 1 {
                value = args[0]
        } else if len(args) >= 2 {
                operator = args[0].(string)
                value = args[1]
        }

        qb.wheres = append(qb.wheres, whereClause{
                column:   column,
                operator: operator,
                value:    value,
                boolean:  "OR",
        })
        qb.bindings = append(qb.bindings, value)

        return qb
}

func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
        placeholders := make([]string, len(values))
        for i := range values {
                placeholders[i] = "?"
                qb.bindings = append(qb.bindings, values[i])
        }

        qb.wheres = append(qb.wheres, whereClause{
                column:   column,
                operator: "IN",
                value:    "(" + strings.Join(placeholders, ", ") + ")",
                boolean:  "AND",
                isRaw:    true,
        })

        return qb
}

func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
        qb.wheres = append(qb.wheres, whereClause{
                column:   column,
                operator: "IS",
                value:    "NULL",
                boolean:  "AND",
                isRaw:    true,
        })
        return qb
}

func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
        qb.wheres = append(qb.wheres, whereClause{
                column:   column,
                operator: "IS NOT",
                value:    "NULL",
                boolean:  "AND",
                isRaw:    true,
        })
        return qb
}

func (qb *QueryBuilder) OrderBy(column, direction string) *QueryBuilder {
        qb.orderBys = append(qb.orderBys, orderByClause{
                column:    column,
                direction: strings.ToUpper(direction),
        })
        return qb
}

func (qb *QueryBuilder) GroupByColumn(columns ...string) *QueryBuilder {
        qb.groupBy = append(qb.groupBy, columns...)
        return qb
}

func (qb *QueryBuilder) Having(condition string) *QueryBuilder {
        qb.having = condition
        return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
        qb.limit = n
        return qb
}

func (qb *QueryBuilder) Offset(n int) *QueryBuilder {
        qb.offset = n
        return qb
}

func (qb *QueryBuilder) Join(table, first, operator, second string) *QueryBuilder {
        qb.joins = append(qb.joins, joinClause{
                joinType: "INNER",
                table:    table,
                first:    first,
                operator: operator,
                second:   second,
        })
        return qb
}

func (qb *QueryBuilder) LeftJoin(table, first, operator, second string) *QueryBuilder {
        qb.joins = append(qb.joins, joinClause{
                joinType: "LEFT",
                table:    table,
                first:    first,
                operator: operator,
                second:   second,
        })
        return qb
}

func (qb *QueryBuilder) RightJoin(table, first, operator, second string) *QueryBuilder {
        qb.joins = append(qb.joins, joinClause{
                joinType: "RIGHT",
                table:    table,
                first:    first,
                operator: operator,
                second:   second,
        })
        return qb
}

func (qb *QueryBuilder) buildSelectQuery() string {
        query := "SELECT " + strings.Join(qb.columns, ", ") + " FROM " + qb.table

        for _, join := range qb.joins {
                query += fmt.Sprintf(" %s JOIN %s ON %s %s %s",
                        join.joinType, join.table, join.first, join.operator, join.second)
        }

        if len(qb.wheres) > 0 {
                query += " WHERE "
                for i, where := range qb.wheres {
                        if i > 0 {
                                query += " " + where.boolean + " "
                        }
                        if where.isRaw {
                                query += where.column + " " + where.operator + " " + where.value.(string)
                        } else {
                                query += where.column + " " + where.operator + " ?"
                        }
                }
        }

        if len(qb.groupBy) > 0 {
                query += " GROUP BY " + strings.Join(qb.groupBy, ", ")
        }

        if qb.having != "" {
                query += " HAVING " + qb.having
        }

        if len(qb.orderBys) > 0 {
                orderParts := make([]string, len(qb.orderBys))
                for i, ob := range qb.orderBys {
                        orderParts[i] = ob.column + " " + ob.direction
                }
                query += " ORDER BY " + strings.Join(orderParts, ", ")
        }

        if qb.limit > 0 {
                query += fmt.Sprintf(" LIMIT %d", qb.limit)
        }

        if qb.offset > 0 {
                query += fmt.Sprintf(" OFFSET %d", qb.offset)
        }

        return query
}

func (qb *QueryBuilder) ToSQL() string {
        return qb.buildSelectQuery()
}

func (qb *QueryBuilder) Get() (*sql.Rows, error) {
        query := qb.buildSelectQuery()
        return qb.executor().Query(query, qb.bindings...)
}

func (qb *QueryBuilder) First() *sql.Row {
        qb.limit = 1
        query := qb.buildSelectQuery()
        return qb.executor().QueryRow(query, qb.bindings...)
}

func (qb *QueryBuilder) Find(id interface{}) *sql.Row {
        return qb.Where("id", id).First()
}

func (qb *QueryBuilder) Count() (int64, error) {
        qb.columns = []string{"COUNT(*) as count"}
        query := qb.buildSelectQuery()
        
        var count int64
        err := qb.executor().QueryRow(query, qb.bindings...).Scan(&count)
        return count, err
}

func (qb *QueryBuilder) Exists() (bool, error) {
        count, err := qb.Count()
        return count > 0, err
}

func (qb *QueryBuilder) Insert(data map[string]interface{}) (sql.Result, error) {
        columns := make([]string, 0, len(data))
        placeholders := make([]string, 0, len(data))
        values := make([]interface{}, 0, len(data))

        for col, val := range data {
                columns = append(columns, col)
                placeholders = append(placeholders, "?")
                values = append(values, val)
        }

        query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
                qb.table,
                strings.Join(columns, ", "),
                strings.Join(placeholders, ", "))

        return qb.executor().Exec(query, values...)
}

func (qb *QueryBuilder) Update(data map[string]interface{}) (sql.Result, error) {
        sets := make([]string, 0, len(data))
        values := make([]interface{}, 0, len(data))

        for col, val := range data {
                sets = append(sets, col+" = ?")
                values = append(values, val)
        }

        values = append(values, qb.bindings...)

        query := fmt.Sprintf("UPDATE %s SET %s", qb.table, strings.Join(sets, ", "))

        if len(qb.wheres) > 0 {
                query += " WHERE "
                for i, where := range qb.wheres {
                        if i > 0 {
                                query += " " + where.boolean + " "
                        }
                        query += where.column + " " + where.operator + " ?"
                }
        }

        return qb.executor().Exec(query, values...)
}

func (qb *QueryBuilder) Delete() (sql.Result, error) {
        query := "DELETE FROM " + qb.table

        if len(qb.wheres) > 0 {
                query += " WHERE "
                for i, where := range qb.wheres {
                        if i > 0 {
                                query += " " + where.boolean + " "
                        }
                        query += where.column + " " + where.operator + " ?"
                }
        }

        return qb.executor().Exec(query, qb.bindings...)
}

func (qb *QueryBuilder) Paginate(page, perPage int) (*Paginator, error) {
        total, err := qb.Count()
        if err != nil {
                return nil, err
        }

        qb.columns = []string{"*"}
        qb.limit = perPage
        qb.offset = (page - 1) * perPage

        rows, err := qb.Get()
        if err != nil {
                return nil, err
        }

        return &Paginator{
                Rows:        rows,
                Total:       total,
                PerPage:     perPage,
                CurrentPage: page,
                LastPage:    int((total + int64(perPage) - 1) / int64(perPage)),
        }, nil
}

type Paginator struct {
        Rows        *sql.Rows
        Total       int64
        PerPage     int
        CurrentPage int
        LastPage    int
}

func (p *Paginator) HasMorePages() bool {
        return p.CurrentPage < p.LastPage
}

func (p *Paginator) HasPages() bool {
        return p.LastPage > 1
}

func ScanStruct(row *sql.Row, dest interface{}) error {
        v := reflect.ValueOf(dest).Elem()
        t := v.Type()

        fields := make([]interface{}, 0)
        for i := 0; i < t.NumField(); i++ {
                field := v.Field(i)
                if field.CanSet() {
                        fields = append(fields, field.Addr().Interface())
                }
        }

        return row.Scan(fields...)
}

func ScanStructRows(rows *sql.Rows, destSlice interface{}) error {
        sliceVal := reflect.ValueOf(destSlice).Elem()
        elemType := sliceVal.Type().Elem()

        for rows.Next() {
                elem := reflect.New(elemType).Elem()
                fields := make([]interface{}, 0)

                for i := 0; i < elem.NumField(); i++ {
                        field := elem.Field(i)
                        if field.CanSet() {
                                fields = append(fields, field.Addr().Interface())
                        }
                }

                if err := rows.Scan(fields...); err != nil {
                        return err
                }

                sliceVal.Set(reflect.Append(sliceVal, elem))
        }

        return rows.Err()
}
