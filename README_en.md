# User Service

A learning project in Go: a user registration service. It shows how to split a small application into layers and packages so that the business logic depends on neither the storage nor the way notifications are sent, and how to test such code.

## Features

- Registering a user by name and email with validation (empty name, email already taken)
- A welcome notification after registration
- Looking up a user by email and listing all users (in registration order)
- Several ways to send notifications: console, email, Telegram, and a test double
- An in-memory storage that is safe to use from multiple goroutines

## Requirements

Go 1.25 or newer (`main.go` uses `sync.WaitGroup.Go`). Standard library only.

## Structure

```
.
├── main.go                    # wiring: creates implementations and connects them, plus a demo
├── go.mod
└── internal/
    ├── user/                  # core: types, rules, interfaces
    │   ├── user.go            # User, ErrNotFound, ErrMailTaken, Repository interface
    │   ├── service.go         # Notifier interface, Service
    │   └── service_test.go
    ├── storage/
    │   ├── memory.go          # MemoryRepo: in-memory repository
    │   └── memory_test.go
    └── notify/
        └── notify.go          # Console, Email, Telegram, Fake
```

Dependency direction:

```
main ──► user ◄── storage
  │
  └────► notify
```

`user` imports neither `storage` nor `notify`. Both interfaces, `Repository` and `Notifier`, are declared in `user`, that is, on the side of the code that uses them. The implementations in `storage` and `notify` satisfy these interfaces implicitly: matching methods is enough. So replacing the storage or the notification channel never requires changing `user.Service`.

## How it works

- **Service.** Registration rules: an empty name (after trimming spaces) is rejected, then the user is saved through `Repository`, then a notification is sent through `Notifier`. If the notification fails, the user stays saved and the error is returned to the caller.
- **Errors.** `ErrNotFound` and `ErrMailTaken` are declared as named errors, errors are wrapped with `fmt.Errorf` and `%w`, and checked with `errors.Is`.
- **MemoryRepo.** A map keyed by email guarded by `sync.RWMutex`. `List` sorts its result by ID because map iteration order is random.
- **Fake.** A notification test double: it records sent messages, can simulate a send failure, and is goroutine-safe.

## Usage

```bash
go run .

# with the race detector
go run -race .
```

Besides a usage example, `main.go` contains a concurrency check: 100 goroutines register users at the same time, then it verifies that there are exactly 100 users and that the IDs have no gaps or duplicates.

## Tests

```bash
go test ./...                 # all tests
go test -race ./...           # with the race detector
go test -cover ./...          # coverage
```

Coverage: `storage` 100%, `user` 100%. The repository tests check ID assignment, a duplicate email, lookup, and list order (the list is requested 20 times in a row to catch random map ordering). The service tests use `Fake` and check that the welcome message is sent once, is not sent when registration fails, and that the user stays saved when the notification fails. The service is tested from an external test package, `user_test`, because `storage` imports `user` and otherwise there would be an import cycle.

## Limitations and ideas for improvement

- Data lives in memory only and is lost on restart
- No HTTP layer: the next step is a REST API on top of this service
- Replacing `MemoryRepo` with a PostgreSQL-backed repository requires only a new implementation of `user.Repository` and one line in `main`

## What I learned

Interfaces and dependency injection through constructors, laying code out in packages and `internal/`, dependency direction and import cycles, named errors and `errors.Is`, the repository pattern, table-driven tests and subtests, test doubles instead of real dependencies, coverage, and the race detector.