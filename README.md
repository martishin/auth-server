# Auth Server
Authentication/Authorization microservice written in Go, provides gRPC endpoints for managing JWT tokens and user's data.  

## Running Locally
* To run the server and PostgreSQL locally you can execute:   
  `make run`
* To stop them, run:  
 `make down`
* **Alternatively** you can start the database and apply migrations by running:  
  `make run-postgres DETACHED=true && make migrate`
* And then start a server:   
`go run cmd/sso/main.go --config=./configs/local.yaml`
* API will be available at http://localhost:8080/

## Testing
* Run tests: `make test`

## How to Connect From Another Service  
* Fetch schemas from the [auth-server-schemas](https://github.com/tty-monkey/auth-server-schemas) repo:  
  `go get github.com/tty-monkey/auth-server-schemas`
* Create a client connection:
```go
cc, err := grpc.DialContext(context.Background(),
	net.JoinHostPort(grpcHost, grpcPort),
	grpc.WithTransportCredentials(insecure.NewCredentials()),
)
```
* Initialize a client:
```go
authClient := ssov1.NewAuthClient(cc)
```
* Make requests using the client:
```go
resp, err := authClient.Register(ctx, &ssov1.RegisterRequest{
  Email:    email,
  Password: password,
})
```
* Usage examples can be found in [e2e tests](https://github.com/tty-monkey/auth-server/blob/d6b9a6ddf5d998fc11b75273124f8597fc4bc1ae/tests/auth_register_login_test.go#L24-L24)

## Endpoints
gRPC protobuf schemas can be found [here](https://github.com/tty-monkey/auth-server-schemas).
You can import the schema file into [Postman](https://blog.postman.com/postman-now-supports-grpc/) and send requests from it.

### Auth / Register
Registers a new user.

- **Request: `RegisterRequest`**
  - `email` (string): User's email address.
  - `password` (string): User's password.

- **Response: `RegisterResponse`**
  - `user_id` (int64): Unique identifier of the registered user.

### Auth / Login
Authenticates a user and provides a token.

- **Request: `LoginRequest`**
  - `email` (string): User's email address.
  - `password` (string): User's password.
  - `app_id` (int32): Application identifier.

- **Response: `LoginResponse`**
  - `token` (string): Authentication token for the user.

### Auth / IsAdmin
Checks if the user is an admin.

- **Request**: `IsAdminRequest`
  - `user_id` (int64): Unique identifier of the user.
 
- **Response**: `IsAdminResponse`
  - `is_admin` (bool): Indicates whether the user is an admin.

## Technologies/Packages Used
* [Go](https://go.dev/)
* [PostgreSQL](https://www.postgresql.org/)
* [Protocol Buffers](https://protobuf.dev/getting-started/gotutorial/)
* [Docker](https://www.docker.com/)
* [pgx](https://pkg.go.dev/github.com/jackc/pgx/v5)
* [jwt](https://pkg.go.dev/github.com/golang-jwt/jwt)
* [golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)
* [crypto](https://pkg.go.dev/golang.org/x/crypto)
* [validator](https://pkg.go.dev/github.com/go-playground/validator/v10)
* [testcontainers-go](https://pkg.go.dev/github.com/testcontainers/testcontainers-go)
* [testify](https://pkg.go.dev/github.com/stretchr/testify)
* [gofakeit](https://pkg.go.dev/github.com/brianvoe/gofakeit)
* [cleanenv](https://pkg.go.dev/github.com/ilyakaznacheev/cleanenv)
