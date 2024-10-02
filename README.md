# Garantex-proxy

### Запуск приложения через docker-compose

```
docker-compose up
```

### Запуск приложения c флагами конфигурации
```
go run cmd/*.go -pg-user=postgres -pg-pass=postgres -pg-db=backend -pg-host=localhost -pg-port=5432
```
####

### Запуск приложения c переменными окружения
```
export POSTGRES_USER=postgres && export POSTGRES_PASSWORD=postgres && export POSTGRES_DATABASE=backend && export POSTGRES_HOST=localhost && export POSTGRES_PORT=5432 &&  go run cmd/*.go
```

### Запуск postgres в docker

```
make run-db
```