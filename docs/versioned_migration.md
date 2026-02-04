## Start:
Create a `migrations` folder in `./internal/infrastructure/persistence/ent/`
```bash
mkdir ./internal/infrastructure/persistence/ent/migrations
```

Create a Schema first using `ent`:
```bash
go run -mod=mod entgo.io/ent/cmd/ent new --target internal/infrastructure/persistence/ent/schema User
```

Edit the `internal/adapters/ent/schema/user.go` file to implement what you want in db for example:
```go
// User holds the schema definition for the User entity.
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "users",
		},
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			NotEmpty().
			MaxLen(32).
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(32)",
			}),
		field.String("password_hash").
			NotEmpty().
			MaxLen(128).
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(128)",
			}),
		field.String("role").
			NotEmpty().
			Default("user").
			MaxLen(32).
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(32)",
			}),
		field.Time("deleted_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}),
		field.Time("created_on").
			Optional().
			Nillable(),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").
			Unique().
			StorageKey("ux_users_username_active").
			Annotations(
				entsql.IndexAnnotation{
					Where: "deleted_at IS NULL",
				},
			),
	}
}

// TimeMixin for time-based fields.
type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("id").
			Unique().
			Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "serial",
			}),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(
				entsql.Annotation{
					Default: "NOW()", // Database default
				},
			).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}),
	}
}
```

After that you can run generate for ent to create all necessary files for your schema using this command:
```bash
go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./internal/infrastructure/persistence/ent/schema
```

Pay attention that you need to use `--feature sql/versioned-migration` option. Otherwise `ent` will not generate `Diff` and `NamedDiff` that are necessary for `atlas`. 

#### Important:
all files in `./internal/infrastructure/persistence/ent` are generated each time you run the generate command. You should just edit your schema files that are in `./internal/infrastructure/persistence/ent/schema`. So, If you changed anything in schema you do not need to modify or remove other files, just run `generate`

we also have put our `migrations` folder in `ent` folder. 

## Atlas needs a real db!!
To run atlas to generate migrations, it needs a real db! The example in the [documentation](https://entgo.io/docs/versioned-migrations/#option-2-create-a-migration-generation-script) is like this:
```bash
atlas migrate diff migration_name \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://mysql/8/ent"
```

The `--dev-url` creates a container form `mysql:8` and works with `/ent` db in it.

### Generate migration files
To make it work with our project we modify it to this:

- we change `migration_name` to something that describes what this changes do
- we change the `--dir` path to where we want to save our migrations files
- we change `--to` path to where our schema files are.
- we change `--dev-url` to a docker image that is same as ours, we just need to replace the `:` in `postgres:18.1-alpine` to `/`

```bash
atlas migrate diff create_user \
  --dir "file://internal/infrastructure/persistence/ent/migrations" \
  --to "ent://internal/infrastructure/persistence/ent/schema" \
  --dev-url "docker://postgres/18.1-alpine/ent"
```

Since `atlas` tries to spin a container it needs to be able to run `docker` commands. You can add your user to docker group (in linux `sudo usermod -aG docker $USER` or run it with `sudo`) or run the command with `sudo`. Pay attention if you run it with `sudo` the generated migration files will belong to `root` user and you need to modify that.


### Inspect generated migration file:
running above command will create a file in `internal/infrastructure/persistence/ent/migrations/20260204215221_create_user.sql`

Like this:
```sql
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" serial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "username" character varying(32) NOT NULL,
  "password_hash" character varying(128) NOT NULL,
  "role" character varying(32) NOT NULL DEFAULT 'user',
  "deleted_at" timestamptz NULL,
  "created_on" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "ux_users_username_active" to table: "users"
CREATE UNIQUE INDEX "ux_users_username_active" ON "public"."users" ("username") WHERE (deleted_at IS NULL);
```

You need to read it and verify it. If this is not what you want, just remove it.
run this command for recalculating the hashes
```bash
atlas migrate hash --dir "file://internal/infrastructure/persistence/ent/migrations"
```

run the following command before you start fixing your schema.
```bash
atlas migrate validate --dir "file://internal/infrastructure/persistence/ent/migrations"
```


Edit your schema and run `generate` with ent and then create migration with atlas.


### lint
Lint is a paid option so we ignore it

### apply the migrations:
Run this command after you are sure about changes:

In `--url` put your database connection string:
```bash
atlas migrate apply \
  --dir "file://internal/infrastructure/persistence/ent/migrations" \
  --url "postgres://postgres:pgpassword@localhost:5432/url_shortner?search_path=public&sslmode=disable"
```


### Example:
In above schema in `internal/infrastructure/persistence/ent/schema/user.go` the following field is not needed:
```go
		field.Time("created_on").
			Optional().
			Nillable(),
```

- remove it from `internal/infrastructure/persistence/ent/schema/user.go`

- run generate command from ent:
```bash
go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./internal/infrastructure/persistence/ent/schema
```

- run atlas to generate migration:
```
atlas migrate diff remove_create_on_field_from_users \
  --dir "file://internal/infrastructure/persistence/ent/migrations" \
  --to "ent://internal/infrastructure/persistence/ent/schema" \
  --dev-url "docker://postgres/18.1-alpine/ent"
```

- verify the sql in the newly generated file `internal/infrastructure/persistence/ent/migrations/20260204225037_remove_create_on_field_from_users.sql`
```sql
ALTER TABLE "public"."users" DROP COLUMN "created_on";
```

- apply the migration:
```bash
atlas migrate apply \
  --dir "file://internal/infrastructure/persistence/ent/migrations" \
  --url "postgres://postgres:pgpassword@localhost:5432/url_shortner?search_path=public&sslmode=disable"
```