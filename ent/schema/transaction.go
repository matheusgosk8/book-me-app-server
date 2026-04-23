package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Transaction struct {
	ent.Schema
}

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Float("amount"),
		field.Text("type"),   // Ex: "pix", "credit_card"
		field.Text("status"), // Ex: "pending", "paid"
		field.Time("created_at").Default(time.Now),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("transactions").Unique().Required(),
	}
}