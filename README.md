1. To start the app
a. add DB_SOURCE:
`export DB_SOURCE=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable`

b. run the postgres in docker
`docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 -v postgres_data:/var/lib/postgresql/data --restart unless-stopped postgres`

2. To apply the migration 
a. create the migration:
`migrate create -ext sql -dir internal/db/migrations -seq <migration name>`

b. apply the migration:
`migrate -path internal/db/migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" up`