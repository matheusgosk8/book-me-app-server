package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),

		// CAMPOS OBRIGATÓRIOS
		field.String("nome"),
		field.String("email").Unique(),
		field.String("senha"),
		field.String("user_type"), // CUSTOMER, PROVIDER, COLLABORATOR

		// CAMPOS OPCIONAIS
		field.String("telefone").Optional(),
		field.String("cpf").Optional().Unique(),
		field.String("cnpj").Optional().Unique(),
		field.String("cep").Optional(),
		field.String("estado").Optional(),
		field.String("cidade").Optional(),
		field.String("logradouro").Optional(),
		field.String("rua").Optional(),
		field.String("confirma_senha").Optional(),

		//Refressh Token
		field.String("refresh_token").
			Optional().
			Unique(),

		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("provider_profile", ProviderProfile.Type).Unique(),
		edge.To("addresses", Address.Type),
		edge.To("orders_created", Order.Type),
		edge.To("orders_fulfilled", Order.Type),
		edge.To("proposals_sent", Proposal.Type),
		edge.To("reviews_written", Review.Type),
		edge.To("reviews_received", Review.Type),
		edge.To("services", Service.Type),
	}
}
