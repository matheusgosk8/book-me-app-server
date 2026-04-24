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
        field.String("country"),    
        field.String("street"),
        field.String("city"),       
        field.String("state"),       
        field.String("postal_code"), 
        field.Float("latitude").Optional(),
        field.Float("longitude").Optional(),
        field.Text("label").Optional(), 
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
