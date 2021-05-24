# Go micro-service

create migration:
```
migrate create -ext sql -dir migrations mg_name
```

migrate:
```
migrate -path migrations -database "postgres://localhost:5432/db_name?sslmode=disable" up
```

To generate documentation, `go-swagger` need to be installed:
```
dir=$(mktemp -d) 
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
go install ./cmd/swagger
rm -rf "$dir"
```
