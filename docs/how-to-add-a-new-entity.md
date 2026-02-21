
### How to create new Entity with ent in `hexagonal` architecture as infrastructure/persistence

To add a `user` entity to the project, you can follow this pattern
| Folders to create | Role  | Description |
| ----------------- |:------|:------------|
| internal/domain/user/ | User domain | Define business domain classes, functions and services that compose the domain model. All business rules should be declared in this layer so the application layer can use it to compose your use cases |
| internal/application/user/ | User Application | Responsible to mediate between input interfaces and business domain |
| infrastructure/persistence/ent | --- | -- , Since ent create its own folders, we put them here |
| internal/adapters/repository/ | | --, repository is meant to retrieve and modify aggregate roots |

#### internal/domain/user/ - Domain Layer
**Purpose**: The innermost layer containing business logic, domain entities, and business rules.

What to put here:
- user.go - Core domain entity (aggregate root)
- value_objects.go - Value objects (Email, Password, UserID, etc.)
- errors.go - Domain-specific errors
- events.go - Domain events (if using event-driven approach)
- specifications.go - Business rules/policies
- repository.go - Repository interface (port)

```go
// Aggregate Root
type User struct {
    ID        UserID
    Email     Email
    Password  PasswordHash
    Status    UserStatus
    CreatedAt time.Time
}

func NewUser(email, password string) (*User, error) {
    // Business logic & validation
}

func (u *User) ChangeEmail(newEmail string) error {
    // Business rules
}
```

#### internal/application/user/ - Application Layer
**Purpose**: Orchestrates use cases, coordinates domain objects, and implements business transactions.

**What to put here:**
- service.go or usecases.go - Application services/use cases
- commands.go - Command objects (for CQRS)
- queries.go - Query objects (for CQRS)
- dtos.go - Data Transfer Objects

```go
type UserService struct {
    repo domain.UserRepository
    // other dependencies
}

func (s *UserService) RegisterUser(cmd RegisterUserCommand) error {
    // Orchestrate domain logic
    user := domain.NewUser(cmd.Email, cmd.Password)
    return s.repo.Save(user)
}
```

#### internal/adapters/repository/ - Infrastructure Layer (Persistence)
**Purpose**: Implements domain interfaces for data persistence (adapters).

**What to put here:**
- user_repository.go - Concrete implementation of domain.UserRepository
- postgres_user_repository.go or mysql_user_repository.go
- inmemory_user_repository.go - For testing

```go
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(user *domain.User) error {
    // Database-specific implementation
}

func (r *PostgresUserRepository) FindByID(id domain.UserID) (*domain.User, error) {
    // Database query logic
}
```

#### infrastructure/persistence/ent - Infrastructure Layer (Framework/ORM)
**Purpose**: Adapters for external tools/libraries (like Ent ORM).

**What to put here:**
- ent_client.go - Ent client setup/configuration
- ent_schema.go - Entity definitions for Ent
- mappers.go - Convert between domain entities and ORM models
- migrations.go - Database migrations

`mappers.go`:
```go
func ToDomainUser(entUser *ent.User) (*domain.User, error) {
    // Map Ent model to domain entity
}

func ToEntUser(domainUser *domain.User) *ent.User {
    // Map domain entity to Ent model
}
```

### Dependency Flow

| Level | Path  | Description |
| ----------------- |:------|:------------|
|1| external/http/handlers.go | Controllers/HTTP |
|2| internal/application/user/service.go | Use Cases |
|3| internal/domain/user/ (user.go, repository.go (interface)) | Business Logic |
|4| internal/adapters/repository/ | Persistence Impl |
|5| infrastructure/persistence/ent/ | ORM/Framework generated files |



User the following commands to create necessary folders:

```bash
mkdir -p internal/domain/user/
mkdir -p internal/application/user/
mkdir -p infrastructure/persistence/ent/
mkdir -p internal/adapters/repository/
```

To create a base schema with `ent` run the following command. Pay attention that you have to set `--target`, so `ent` creates its own folder structure in the `infrastructure/persistence/ent` folder.
```bash
go run -mod=mod entgo.io/ent/cmd/ent new --target infrastructure/persistence/ent/schema User
```

Now edit the generated file in `infrastructure/persistence/ent/user.go` and add `Fields` and `Indexes` inside that. After that you can run the following command and `ent` will create all its files and folders. You need to specify the location of `schema` folder here:
```shell
go run -mod=mod entgo.io/ent/cmd/ent generate ./infrastructure/persistence/ent/schema
```

#### Add repository
In `internal/domain/user/` create `repository.go` file and add user repository interface, something like this:
```go
type Repository interface {
	Create(ctx context.Context, u *User) (*User, error)
	GetByID(ctx context.Context, id ID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}
```
Now we need to implement the `domain/user/repository` for `ent` adaptor. In `internal/adapters/repository/` create `user_ent_repository.go` folder. Here you need to implement all the repository functions with `ent.client`
```go
type UserEntRepository struct {
	Client *ent.Client
}

func NewUserEntRepository(client *ent.Client) *UserEntRepository {
	return &UserEntRepository{Client: client}
}
```

then in `internal/application/user/` create `input` and `output` conversions and domain conversions, and services.

In next step you need to create a `ent user repo` then use it to create a `user app`
```go
clientEnt := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
userRepo := &repository.UserEntRepository{Client: clientEnt}
createUserApp := userapp.CreateUser{Repo: userRepo}
```