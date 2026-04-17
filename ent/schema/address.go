package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Address holds the schema definition for the Address entity.
type Address struct {
	ent.Schema
}

// Fields of the Address.
func (Address) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(), field.String("street"),
		field.String("city"),
		field.String("state"),
		field.String("postal_code"),
		field.String("country"),
		field.Time("creation_date").Default(time.Now),
	}
}

// Edges of the Address.
func (Address) Edges() []ent.Edge {
	return nil
}
