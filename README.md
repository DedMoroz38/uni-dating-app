# uni-dating-app

This is a dating app for university students.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

### Running the application

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/DedMoroz38/uni-dating-app.git
    cd uni-dating-app
    ```

2.  **Run the application with Docker Compose:**

    ```sh
    docker-compose up -d
    ```

This will start the API server and a PostgreSQL database. The API will be available at `http://localhost:3000`.

The database schema will be automatically applied on the first startup.

### Environment Variables

The application uses the following environment variable:

*   `DB_SOURCE`: The connection string for the PostgreSQL database. This is already configured in the `docker-compose.yaml` for the development environment.

Example: `postgresql://user:password@db:5432/uni_dating_app_db?sslmode=disable`


3. To apply the migration 
a. create the migration:
`migrate create -ext sql -dir internal/db/migrations -seq <migration name>`

b. apply the migration:
`migrate -path internal/db/migrations -database "postgres://user:password@localhost:5432/postgres?sslmode=disable" up`