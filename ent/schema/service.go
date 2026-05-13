package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	ent.Schema
}

func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Text("title").NotEmpty(),
		field.Text("description").Optional(),

		// Valor base do serviço
		field.Float("price_base").Min(0),

		// Define a lógica de cobrança: Fixo ou Por Hora
		field.Enum("price_type").
			Values("fixed", "hourly").
			Default("fixed").
			Comment("fixed: valor único | hourly: valor multiplicado pelas horas"),

		// Tempo estimado (importante para o agendamento não encavalar)
		field.Int("duration_minutes").
			Positive().
			Comment("Duração média em minutos do serviço"),

		field.Bool("is_active").Default(true),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),

		field.Bool("is_in_place").Default(true),
		
		field.UUID("address_id", uuid.UUID{}).Optional().Nillable(),
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
