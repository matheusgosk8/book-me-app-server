package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Text("description"),
		field.Text("status").Default("pending"),
		field.Float("latitude"),
		field.Float("longitude"),
		field.Time("scheduled_at").Optional(),
		field.Time("created_at").Default(time.Now),
	}
}

func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("client", User.Type).Ref("orders_created").Unique().Required(),
		edge.From("provider", User.Type).Ref("orders_fulfilled").Unique(),
		edge.From("category", Category.Type).Ref("orders").Unique().Required(),
		edge.To("proposals", Proposal.Type),
		edge.To("reviews", Review.Type),
		edge.To("transactions", Transaction.Type),
	}
}
