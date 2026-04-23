package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Service struct {
	ent.Schema
}

func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Text("title"),
		field.Text("description"),
		field.Float("price_base"),
		field.Text("service_type"),
		field.Bool("is_active").Default(true),
	}
}

func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("provider", User.Type).
			Ref("services").
			Unique().
			Required(),
		edge.From("category", Category.Type).
			Ref("services").
			Unique().
			Required(),
	}
}