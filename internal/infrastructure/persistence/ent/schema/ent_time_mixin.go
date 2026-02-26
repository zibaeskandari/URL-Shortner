package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now().UTC).
			Annotations(
				entsql.Annotation{
					Default: "NOW()",
				},
			).
			SchemaType(map[string]string{
				dialect.Postgres: "TIMESTAMP",
			}),
		field.Time("updated_at").
			Default(time.Now().UTC).
			UpdateDefault(time.Now().UTC).
			Annotations(
				entsql.Annotation{
					Default: "NOW()",
				},
			).
			SchemaType(map[string]string{
				dialect.Postgres: "TIMESTAMP",
			}),
		field.Time("deleted_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "TIMESTAMP",
			}),
	}
}
