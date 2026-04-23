package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Review struct {
	ent.Schema
}

func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Int("rating"),
		field.Text("comment").Optional(),
		field.Time("created_at").Default(time.Now),
	}
}

func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		
		edge.From("order", Order.Type).Ref("reviews").Unique().Required(),
		edge.From("author", User.Type).Ref("reviews_written").Unique().Required(),
		edge.From("subject", User.Type).Ref("reviews_received").Unique().Required(),
	}
}