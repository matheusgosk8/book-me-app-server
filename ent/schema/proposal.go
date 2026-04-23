package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Proposal struct {
	ent.Schema
}

func (Proposal) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Float("price"),
		field.Text("message"),
		field.Text("status").Default("pending"),
	}
}

func (Proposal) Edges() []ent.Edge {
	return []ent.Edge{
		// A proposta pertence a uma Ordem
		edge.From("order", Order.Type).
			Ref("proposals").
			Unique().
			Required(),
		// A proposta foi enviada por um Prestador (User)
		edge.From("provider", User.Type).
			Ref("proposals_sent").
			Unique().
			Required(),
	}
}