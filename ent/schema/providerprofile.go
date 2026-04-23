package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ProviderProfile holds the schema definition for the ProviderProfile entity.
type ProviderProfile struct {
	ent.Schema
}

// Fields of the ProviderProfile.
func (ProviderProfile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id"),
		field.Text("bio").
			Optional(), // No diagrama está como text
		field.Float("rating_avg").
			Default(0), // numeric rating_avg
		field.Bool("is_active").
			Default(true), // boolean is_active
	}
}

// Edges of the ProviderProfile.
func (ProviderProfile) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("user", User.Type).
            Ref("provider_profile").
            Unique().
            Required(), 
            
        // Comente estas linhas abaixo temporariamente:
        // edge.To("services", Service.Type), 
    }
}