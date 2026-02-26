package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Urls holds the schema definition for the Urls entity.
type Urls struct {
	ent.Schema
}

func (Urls) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "urls",
		},
	}
}

func (Urls) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Urls.
func (Urls) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			NotEmpty().
			MaxLen(32).
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(32)",
			}),
		field.String("destination").
			NotEmpty().
			SchemaType(map[string]string{
				dialect.Postgres: "TEXT",
			}),
		field.Int64("user_id").
			Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "BIGINT",
			}),
		field.Time("expires_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "TIMESTAMP",
			}),
	}
}

func (Urls) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "updated_at").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")).
			StorageKey("ix_urls_user_active"),
	}
}

// Edges of the Urls.
func (Urls) Edges() []ent.Edge {
	return nil
}
