# GmsTemp

### Env variables:

```
DEBUG=true
LOG_LEVEL=debug
HTTP_LISTEN=:80
HTTP_CORS=false
REDIS_URL=redis:6379
REDIS_PSW=psw
REDIS_DB=0
PG_DSN=postgres://user:pass@host:5432/dbname?sslmode=disable
MS_SMS_URL=http://url/
MS_WS_URL=http://url/
MS_WS_API_KEY=api_key
MS_WS_CHANNEL_NS=channel_ns
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
