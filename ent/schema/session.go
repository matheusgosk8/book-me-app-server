package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Session struct {
	ent.Schema
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("refresh_token").Unique(),
		field.Time("last_login_at").Default(time.Now),
		field.Time("expires_at"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		// A sessão volta para o usuário dono dela
		edge.From("user", User.Type).
			Ref("sessions").
			Unique().
			Required(),
	}
}