# Simple Bank 

## Migration

1. Create Migration File
```shell
migrate create -ext sql -dir db/migration -seq <nama schema>
```
2. Execute migration
```shell
migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up
```
3. Rollback migration
```shell
migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down
```