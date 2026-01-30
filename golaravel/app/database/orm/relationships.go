package orm

import (
        "database/sql"
)

type Relationship interface {
        ForeignKeyName() string
}

type HasOne struct {
        conn       *Connection
        related    string
        foreignKey string
        localKey   string
        localValue interface{}
}

func NewHasOne(conn *Connection, related, foreignKey, localKey string, localValue interface{}) *HasOne {
        return &HasOne{
                conn:       conn,
                related:    related,
                foreignKey: foreignKey,
                localKey:   localKey,
                localValue: localValue,
        }
}

func (r *HasOne) Get() *sql.Row {
        return r.conn.Query(r.related).Where(r.foreignKey, r.localValue).First()
}

func (r *HasOne) ForeignKeyName() string {
        return r.foreignKey
}

func (r *HasOne) LocalKeyName() string {
        return r.localKey
}

type HasMany struct {
        conn       *Connection
        related    string
        foreignKey string
        localKey   string
        localValue interface{}
}

func NewHasMany(conn *Connection, related, foreignKey, localKey string, localValue interface{}) *HasMany {
        return &HasMany{
                conn:       conn,
                related:    related,
                foreignKey: foreignKey,
                localKey:   localKey,
                localValue: localValue,
        }
}

func (r *HasMany) Get() (*sql.Rows, error) {
        return r.conn.Query(r.related).Where(r.foreignKey, r.localValue).Get()
}

func (r *HasMany) Query() *QueryBuilder {
        return r.conn.Query(r.related).Where(r.foreignKey, r.localValue)
}

func (r *HasMany) ForeignKeyName() string {
        return r.foreignKey
}

func (r *HasMany) LocalKeyName() string {
        return r.localKey
}

type BelongsTo struct {
        conn       *Connection
        related    string
        foreignKey string
        ownerKey   string
        foreignValue interface{}
}

func NewBelongsTo(conn *Connection, related, foreignKey, ownerKey string, foreignValue interface{}) *BelongsTo {
        return &BelongsTo{
                conn:         conn,
                related:      related,
                foreignKey:   foreignKey,
                ownerKey:     ownerKey,
                foreignValue: foreignValue,
        }
}

func (r *BelongsTo) Get() *sql.Row {
        return r.conn.Query(r.related).Where(r.ownerKey, r.foreignValue).First()
}

func (r *BelongsTo) ForeignKeyName() string {
        return r.foreignKey
}

func (r *BelongsTo) OwnerKeyName() string {
        return r.ownerKey
}

type BelongsToMany struct {
        conn           *Connection
        related        string
        pivotTable     string
        foreignPivotKey string
        relatedPivotKey string
        parentKey      string
        relatedKey     string
        parentValue    interface{}
}

func NewBelongsToMany(conn *Connection, related, pivotTable, foreignPivotKey, relatedPivotKey, parentKey, relatedKey string, parentValue interface{}) *BelongsToMany {
        return &BelongsToMany{
                conn:            conn,
                related:         related,
                pivotTable:      pivotTable,
                foreignPivotKey: foreignPivotKey,
                relatedPivotKey: relatedPivotKey,
                parentKey:       parentKey,
                relatedKey:      relatedKey,
                parentValue:     parentValue,
        }
}

func (r *BelongsToMany) Get() (*sql.Rows, error) {
        return r.Query().Get()
}

func (r *BelongsToMany) Query() *QueryBuilder {
        return r.conn.Query(r.related).
                Join(r.pivotTable, r.related+"."+r.relatedKey, "=", r.pivotTable+"."+r.relatedPivotKey).
                Where(r.pivotTable+"."+r.foreignPivotKey, r.parentValue)
}

func (r *BelongsToMany) Attach(relatedID interface{}) (sql.Result, error) {
        return r.conn.Query(r.pivotTable).Insert(map[string]interface{}{
                r.foreignPivotKey: r.parentValue,
                r.relatedPivotKey: relatedID,
        })
}

func (r *BelongsToMany) Detach(relatedID interface{}) (sql.Result, error) {
        return r.conn.Query(r.pivotTable).
                Where(r.foreignPivotKey, r.parentValue).
                Where(r.relatedPivotKey, relatedID).
                Delete()
}

func (r *BelongsToMany) ForeignKeyName() string {
        return r.foreignPivotKey
}

func (r *BelongsToMany) Sync(relatedIDs []interface{}) error {
        _, err := r.conn.Query(r.pivotTable).
                Where(r.foreignPivotKey, r.parentValue).
                Delete()
        if err != nil {
                return err
        }

        for _, id := range relatedIDs {
                _, err := r.Attach(id)
                if err != nil {
                        return err
                }
        }

        return nil
}

type RelationshipBuilder struct {
        conn *Connection
}

func NewRelationshipBuilder(conn *Connection) *RelationshipBuilder {
        return &RelationshipBuilder{conn: conn}
}

func (rb *RelationshipBuilder) HasOne(related, foreignKey string, localValue interface{}) *HasOne {
        return NewHasOne(rb.conn, related, foreignKey, "id", localValue)
}

func (rb *RelationshipBuilder) HasMany(related, foreignKey string, localValue interface{}) *HasMany {
        return NewHasMany(rb.conn, related, foreignKey, "id", localValue)
}

func (rb *RelationshipBuilder) BelongsTo(related, foreignKey string, foreignValue interface{}) *BelongsTo {
        return NewBelongsTo(rb.conn, related, foreignKey, "id", foreignValue)
}

func (rb *RelationshipBuilder) BelongsToMany(related, pivotTable string, parentValue interface{}) *BelongsToMany {
        return NewBelongsToMany(
                rb.conn,
                related,
                pivotTable,
                "parent_id",
                "related_id",
                "id",
                "id",
                parentValue,
        )
}
