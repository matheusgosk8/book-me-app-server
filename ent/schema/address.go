package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
        field.UUID("id", uuid.UUID{}).Default(uuid.New),
        field.Float("latitude"),
        field.Float("longitude"),
        field.Text("label"), // Ex: "Casa", "Trabalho"
        field.Bool("is_primary").Default(false),
    }
}

func (Address) Edges() []ent.Edge {
    return []ent.Edge{
        // FK: O endereço aponta para o dono dele (User)
        edge.From("user", User.Type).
            Ref("addresses").
            Unique().
            Required(),
    }
}
