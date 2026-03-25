package schema

import (
	"time"

	"entgo.io/ent"
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
		field.String("cep"),
		field.String("cidade"),
		field.String("cnpj"),
		field.String("confirma_senha"),
		field.String("cpf"),
		field.String("email"),
		field.String("estado"),
		field.String("logradouro"),
		field.String("nome"),
		field.String("rua"),
		field.String("senha"),
		field.String("telefone"),
		field.String("user_type"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
