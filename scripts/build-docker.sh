docker container stop postgres
docker container rm postgres
docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 23729:5432 -v postgres_data:/var/lib/postgresql/data --restart unless-stopped postgres

# psql "postgresql://postgres:password@localhost:23729/postgres?sslmode=disable" -f internal/db/schema.sql