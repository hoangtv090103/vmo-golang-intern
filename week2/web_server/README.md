# Project Structure - Clean Architecture

```plaintext
week2/web_server
│
├── cmd
│   └── main.go                   --> Application entry point
│
├── internal
│
│   module_user
│   │
│   ├── domain
│   │   └── user.go               --> Core entities (User struct, business logic)
│   │
│   ├── handlers
│   │   └── user_handler.go       --> Interface/Controller layer (Handles HTTP requests, forwards to use case)
│   │
│   ├── infra
│   │   └── memory
│   │       └── user.go           --> Infrastructure layer (In-memory storage, repository implementation)
│   │
│   ├── repositories
│   │   └── user_repository.go    --> Repository interfaces
│   │
│   └── usecases
│       └── user_usecase.go       --> Use case layer (Application logic, coordinates between domain and repositories)
│
└── tmp                           --> (Empty or used for temporary files)
```

## Explanation

- **cmd/main.go**: The application entry point where the HTTP server is set up, and dependencies are wired together.
- **internal/module_user/domain**: This folder contains domain logic and entities. user.go defines the core User struct, which may also include business rules.
- **internal/module_user/handlers**: This folder contains the controller logic that handles HTTP requests. user_handler.go manages the HTTP routes for user operations and forwards them to the use case layer.
- **internal/module_user/infra/memory**: This contains infrastructure-specific logic, in this case, an in-memory implementation of the user repository (user.go), storing user data in-memory for simplicity.
- **internal/module_user/repositories**: This folder defines repository interfaces. user_repository.go contains methods that abstract data access logic. The interface is implemented by the infrastructure layer.
- **internal/module_user/usecases**: This is the use case layer that coordinates between user requests (from the handlers) and domain logic (business rules) or repositories. user_usecase.go handles specific user operations (e.g., create, get users).
