# Go micro-service

create migration:
```
migrate create -ext sql -dir migrations mg_name
```

migrate:
```
migrate -path migrations -database "postgres://localhost:5432/db_name?sslmode=disable" up
```
