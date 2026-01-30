package orm

import (
        "testing"
)

func TestHasOneStructure(t *testing.T) {
        hasOne := &HasOne{
                related:    "profiles",
                foreignKey: "user_id",
                localKey:   "id",
                localValue: 1,
        }

        if hasOne.ForeignKeyName() != "user_id" {
                t.Errorf("Expected 'user_id', got '%s'", hasOne.ForeignKeyName())
        }
        if hasOne.LocalKeyName() != "id" {
                t.Errorf("Expected 'id', got '%s'", hasOne.LocalKeyName())
        }
}

func TestHasManyStructure(t *testing.T) {
        hasMany := &HasMany{
                related:    "posts",
                foreignKey: "user_id",
                localKey:   "id",
                localValue: 1,
        }

        if hasMany.ForeignKeyName() != "user_id" {
                t.Errorf("Expected 'user_id', got '%s'", hasMany.ForeignKeyName())
        }
        if hasMany.LocalKeyName() != "id" {
                t.Errorf("Expected 'id', got '%s'", hasMany.LocalKeyName())
        }
}

func TestBelongsToStructure(t *testing.T) {
        belongsTo := &BelongsTo{
                related:      "users",
                foreignKey:   "user_id",
                ownerKey:     "id",
                foreignValue: 1,
        }

        if belongsTo.ForeignKeyName() != "user_id" {
                t.Errorf("Expected 'user_id', got '%s'", belongsTo.ForeignKeyName())
        }
        if belongsTo.OwnerKeyName() != "id" {
                t.Errorf("Expected 'id', got '%s'", belongsTo.OwnerKeyName())
        }
}

func TestBelongsToManyStructure(t *testing.T) {
        belongsToMany := &BelongsToMany{
                related:         "roles",
                pivotTable:      "role_user",
                foreignPivotKey: "user_id",
                relatedPivotKey: "role_id",
                parentKey:       "id",
                relatedKey:      "id",
                parentValue:     1,
        }

        if belongsToMany.related != "roles" {
                t.Errorf("Expected 'roles', got '%s'", belongsToMany.related)
        }
        if belongsToMany.pivotTable != "role_user" {
                t.Errorf("Expected 'role_user', got '%s'", belongsToMany.pivotTable)
        }
}

func TestNewHasOne(t *testing.T) {
        hasOne := NewHasOne(nil, "profiles", "user_id", "id", 1)

        if hasOne.related != "profiles" {
                t.Errorf("Expected 'profiles', got '%s'", hasOne.related)
        }
        if hasOne.foreignKey != "user_id" {
                t.Errorf("Expected 'user_id', got '%s'", hasOne.foreignKey)
        }
        if hasOne.localKey != "id" {
                t.Errorf("Expected 'id', got '%s'", hasOne.localKey)
        }
        if hasOne.localValue != 1 {
                t.Errorf("Expected 1, got %v", hasOne.localValue)
        }
}

func TestNewHasMany(t *testing.T) {
        hasMany := NewHasMany(nil, "posts", "user_id", "id", 1)

        if hasMany.related != "posts" {
                t.Errorf("Expected 'posts', got '%s'", hasMany.related)
        }
        if hasMany.foreignKey != "user_id" {
                t.Errorf("Expected 'user_id', got '%s'", hasMany.foreignKey)
        }
}

func TestNewBelongsTo(t *testing.T) {
        belongsTo := NewBelongsTo(nil, "users", "user_id", "id", 1)

        if belongsTo.related != "users" {
                t.Errorf("Expected 'users', got '%s'", belongsTo.related)
        }
        if belongsTo.foreignKey != "user_id" {
                t.Errorf("Expected 'user_id', got '%s'", belongsTo.foreignKey)
        }
        if belongsTo.ownerKey != "id" {
                t.Errorf("Expected 'id', got '%s'", belongsTo.ownerKey)
        }
}

func TestNewBelongsToMany(t *testing.T) {
        belongsToMany := NewBelongsToMany(nil, "roles", "role_user", "user_id", "role_id", "id", "id", 1)

        if belongsToMany.related != "roles" {
                t.Errorf("Expected 'roles', got '%s'", belongsToMany.related)
        }
        if belongsToMany.pivotTable != "role_user" {
                t.Errorf("Expected 'role_user', got '%s'", belongsToMany.pivotTable)
        }
        if belongsToMany.foreignPivotKey != "user_id" {
                t.Errorf("Expected 'user_id', got '%s'", belongsToMany.foreignPivotKey)
        }
        if belongsToMany.relatedPivotKey != "role_id" {
                t.Errorf("Expected 'role_id', got '%s'", belongsToMany.relatedPivotKey)
        }
}

func TestRelationshipBuilder(t *testing.T) {
        rb := NewRelationshipBuilder(nil)

        hasOne := rb.HasOne("profiles", "user_id", 1)
        if hasOne == nil {
                t.Error("Expected HasOne to be created")
        }

        hasMany := rb.HasMany("posts", "user_id", 1)
        if hasMany == nil {
                t.Error("Expected HasMany to be created")
        }

        belongsTo := rb.BelongsTo("users", "user_id", 1)
        if belongsTo == nil {
                t.Error("Expected BelongsTo to be created")
        }

        belongsToMany := rb.BelongsToMany("roles", "role_user", 1)
        if belongsToMany == nil {
                t.Error("Expected BelongsToMany to be created")
        }
}
