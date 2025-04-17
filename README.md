# Simple Bank 

## Migration

1. Create Migration File
```shell
migrate create -ext sql -dir db/migration -seq <nama schema>
```
2. Execute migration
```shell
migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up <step>
```
3. Rollback migration
```shell
migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down <step>
```

## SQLC
1. Initialize
```shell
sqlc init
```

2. Configuration 
```yaml
version: "1"
packages:
    - name: "db"
      path: "./db/sqlc"
      queries: "./db/query/"
      schema: "./db/migration/"
      engine: "postgresql"
      emit_json_tags: true
      emit_prepared_queries: true
      emit_interface: false
      emit_exact_table_names: false
```
3. For generate
```shell
sqlc generate
```

## Testing Coverage
for test all go test cover
```sh
go test -v -cover ./...
```

## Initialize Mock
```sh
mockgen -package mockdb -destination db/mock/store github.com/Cell6969/go_bank/db/sqlc Store
```
## Run HTTP Server
```sh
go run main.go
```

## Build Image
```sh
docker build -t gobank:latest .
```

## Run Container
for run in development mode:
```sh
docker run --name gobank -p 8080:8080 gobank:latest
```

for run in release mode with same network
```sh
docker run --name gobank --network bank-network -e GIN_MODE=release -p 8080:8080 -e DB_SOURCE=postgres://root:root@po
stgres:5432/simple_bank?sslmode=disable gobank:latest
```