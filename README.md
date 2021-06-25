# GmsTemp

### Env variables:

```
DEBUG: true
LOG_LEVEL: "debug"
HTTP_LISTEN: ":80"
METRICS: false
MS_SMS_URL: https://url/
MS_WS_URL: https://url/
```

### DB dump:

```
pg_dump --no-owner -Fc -U postgres gms_temp -f ./gms_temp.custom
```

### DB restore:

```
dropdb -U postgres gms_temp
createdb -U postgres gms_temp
pg_restore --no-owner -d gms_temp -U postgres ./gms_temp.custom
```

### Install `migrate` command-tool:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

### Create new migration:

```
migrate create -ext sql -dir migrations mg_name
```

### Apply migration:

```
migrate -path migrations -database "postgres://localhost:5432/db_name?sslmode=disable" up
```

<br/>

### Install `swagger-cli`:

```
dir=$(mktemp -d) 
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
go install ./cmd/swagger
rm -rf "$dir"
```
