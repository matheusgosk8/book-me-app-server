package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New),
        field.Text("name").Unique(),
    }
}

func (Category) Edges() []ent.Edge {
    return []ent.Edge{
        // Relacionamento Pai/Filho (Subcategorias)
        edge.To("children", Category.Type).
            From("parent").
            Unique(),
        edge.To("services", Service.Type),
		edge.To("orders", Order.Type),
    }
}
