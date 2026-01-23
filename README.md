# URL-Shortner
A high-performance, modular URL Shortener service built with Go (Golang), featuring PostgreSQL, Redis caching, PASETO authentication, and CI/CD integration.




## Api list:
| Method        | route  | description |
| ------------- |:-------------|:-------------|
|POST| /api/v1/auth/register | User can register, all users created in this route are simple users|
|POST| /api/v1/auth/login | Users can login, they get (token, refreshToken)|
|POST| /api/v1/change-password | User can change their password|
|POST| /api/v1/urls | Submit a new url for the user (url, expiration date)|
|GET| /api/v1/urls | Return the list of all user's urls (it has pagination)|
|GET| /api/v1/urls/{url_id} | Return the information of the url by id (should it have authorization) |
|DELETE| /api/v1/urls/{url_id} | Remove a url by id (Only owner can do this) |
|PUT| /api/v1/urls/{url_id}/renew | Change the url's expiration time, or change the destination (Only owner can do this) |
|GET| /api/v1/urls/{url_id}/stats| Return the usage statistics of a url by id, it has date filter and pagination (Admin or owner) |
|GET| /api/v1/admin/urls | Return list of urls in the system, it has date filter and pagination |
|DELETE| /api/v1/admin/urls/{url_id} | Remove a url by id |
|GET| /api/v1/admin/users/{user_name} | Return user info by username |
|GET | /api/v1/admin/users/{user_id}/urls | Return an specific user's urls |

### Optional Api list:
| Method        | route  | description |
| ------------- |:-------------|:-------------|
|DELETE| /api/v1/admin/users/{user_id} | Deletes all user's links and ban his username |
|GET| /api/v1/admin/urls/stats | Shows stats of the system (number of urls and their count per day) |
|PUT| /api/v1/admin/users/{user_id}/role | promote/demote a user to admin, this can only be called by supper admin (Can admin user still do the user tasks?)


## Todo:
1. create base fiber server
    1. create json error (errors should be in json formats)
1. add zap logger in json format
    1. log request and response
    1. log errors
1. create config and inject in fiber
    1. deferent config should be applied in (prod, dev, and test environments)
1. inject database and redis to fiber
    1. created base on config (pay extra attention to testability)
1. create user table and repositories, and routes
    1. create super admin user on startup (ignore if already exists)
1. create token for users
1. implement refresh token
1. implement authorization and authentication, add middleware to fiber
1. create integration test client
    1. tests for authorization and authentication
1. create url to short code table and repositories
1. create short code generator service
    1. use worker, channel to to generate and put them in redis list
        1. worker should always fill the redis list (on heavy load list might be empty should we use channel or redis list in api?)
1. implement url registration
1. implement url resolve and submit statistics (just log it for now)
1. create url `usage statistics` table and repositories
1. refactor submit statistics to work with db
1. create usage log aggregation worker
1. add usage log routes
1. add rate limiter middleware
    1. Can we load rate limit rules from config?
    1. Key can be `METHOD:ROUTE:<USER_ID or IP>` (We need to see what does go fiber provides in the middleware)
1. add migration scripts to match indexes and db design
    1. goose with ent?
1. create multistage docker image for the app
    1. create run migration docker image and it main app should depend on exiting it


------------

###
internal/infrastructure/database/     # DB connections
internal/interfaces/http/handlers/    # HTTP handlers
internal/interfaces/grpc/             # gRPC handlers
pkg/shared/                           # Shared utilities


# Developer's documentations:
- [Database schema](./docs/database-schema.md)
- [How to set up a new db entity using ent framework in hexagonal architecture with DDD](./docs/how-to-add-a-new-entity.md)
