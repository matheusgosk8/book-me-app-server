package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/google/uuid"
)

type Category struct {
    ent.Schema
}

func (Category) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New),
        field.Text("name").Unique(),
        field.Bool("is_active").Default(true), 
    }
}

func (Category) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("children", Category.Type).
            From("parent").
            Unique(),
        edge.To("services", Service.Type),
        edge.To("orders", Order.Type),
    }
}