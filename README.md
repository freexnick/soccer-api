
## 📑 Table of Contents

- [📘 User Story](#user-story-soccer-manager-api-development)
- [🛠️ Project Setup & Installation](#running-the-soccer-api-project)

# Soccer Manager API

### User Story: Soccer Manager API Development

**Title**: Develop RESTful API for Soccer Manager Application

**As a** football fan,  
**I want** to create fantasy teams and buy or sell players via a RESTful API,  
**So that** I can engage in managing my own football team.

---

### Acceptance Criteria

#### Platform
- [x] Implement using the **latest version of the Go Project**

#### Authorization
- [x] Develop an authorization system

#### Database
- [x] Create a database architecture

#### User Account Management
- [x] Users must be able to create an account and log in using the API  
- [x] Each user can have **only one team**

#### Team Creation
- [x] Upon user signup, automatically assign a team of **20 players**, consisting of:
  - [x] 3 goalkeepers
  - [x] 6 defenders
  - [x] 6 midfielders
  - [x] 5 attackers
- [x] Each player has an **initial value of $1,000,000**
- [x] Each team has an additional budget of **$5,000,000** to purchase players

#### Team Information
- [x] Users can view their team and player information
- [x] Team details must include:
  - [x] Team name and country (editable)
  - [x] Total team value (sum of player values)

#### Player Information
- [x] Players must have the following attributes:
  - [x] First name, last name, country (editable by team owner)
  - [x] Age (randomly generated between 18 and 40)
  - [x] Market value

#### Transfer List Functionality
- [x] Team owners can set players on a transfer list with an asking price
- [x] Asking price must be displayed on a public market list
- [x] Players must be purchased at the listed price
- [x] Users can view all players on the transfer list
- [x] Team budgets must be updated with each transfer
- [x] Transferred players must have their value increased by a random **10% to 100%**

#### API Functionality
- [x] Ensure **all user actions** can be performed via RESTful API
- [x] Include user authentication

---

### Additional Notes
- [x] Provide Postman collection
- [x] Upload project to **GitHub as a public repository**
- [ ] Unit tests are **not required** but are recommended
- [x] Application must support **localization in English and Georgian**


# Running the Soccer API Project

This guide provides instructions for running and managing the Soccer API project using Docker and Docker Compose, based on the provided configuration. The project includes a PostgreSQL database, an Adminer interface, and a server application, all orchestrated via Docker Compose.

## Prerequisites

- **Docker** and **Docker Compose** installed on your system.
- A `.env` file in the `configs/` directory with the following variables:
  - `POSTGRES_USER`: PostgreSQL username.
  - `POSTGRES_PASSWORD`: PostgreSQL password.
  - `POSTGRES_DB`: PostgreSQL database name.
  - `DB_URL`: Database connection string (e.g., `postgresql://user:pass@localhost:5432/dbname?sslmode=disable&search_path=public`).
  - Optional:
    - `POSTGRES_HOST_PORT`: Host port for PostgreSQL (default: 5432).
    - `ADMINER_HOST_PORT`: Host port for Adminer (default: 8080).
    - `SERVER_HOST_PORT`: Host port for the server (default: 7777).
    - `HTTP_SERVER_ADDRESS`: Server port inside the container (default: 7777).
- Git installed to retrieve `GIT_COMMIT_SHA` and `APP_VERSION`.

## Project Structure

The project uses the following services defined in the `docker-compose.yaml`:

1. **PostgreSQL Database (`soccer-api_postgres_db`)**
   - Image: `postgres:16.3-alpine3.19`
   - Port: Maps `${POSTGRES_HOST_PORT:-5432}` to `5432` in the container.
   - Volume: Persistent storage at `soccer-api_postgres_db`.
   - Healthcheck: Ensures the database is ready using `pg_isready`.
   - Environment: Loaded from `configs/.env`.

2. **Adminer (`soccer-api_adminer`)**
   - Image: `adminer:4.8.1-standalone`
   - Port: Maps `${ADMINER_HOST_PORT:-8080}` to `8080` in the container.
   - Connects to the database via `ADMINER_DEFAULT_SERVER=soccer-api_postgres_db`.

3. **Server (`soccer-api_server`)**
   - Built from `deployments/server.dockerfile` in the project root.
   - Depends on a healthy database service.
   - Port: Maps `${SERVER_HOST_PORT:-7777}` to `${HTTP_SERVER_ADDRESS:-7777}` in the container.
   - Environment: Loaded from `configs/.env`.
   - Build arguments:
     - `APP_VERSION`: Application version (defaults to git tag or `dev`).
     - `GIT_COMMIT_SHA`: Git commit hash (defaults to `unknown`).

4. **Network**
   - A bridge network named `soccer-api` connects all services.

5. **Volume**
   - `soccer-api_postgres_db`: Persistent storage for PostgreSQL data.

## Makefile Commands

The provided `Makefile` includes commands to manage the project. Run these from the project root where the `Makefile` is located.

### Help
Display available commands:
```bash
make help
```

### Create a New Migration
Create a new SQL migration file in `internal/infrastructure/database/migrations`:
```bash
make create-migration NAME=<migration_name>
```
- Example: `make create-migration NAME=add_users_table`
- Creates `up.sql` and `down.sql` files for the migration. Edit these files to define the schema changes.

### Apply Migrations
Apply all pending migrations to the database:
```bash
make apply-migrations
```
- Requires `DB_URL` to be set in the `.env` file or environment.
- Uses the `migrate/migrate:v4.18.3` Docker image.

### Rollback Migration
Revert the last applied migration:
```bash
make rollback-migration
```
- Requires `DB_URL` to be set.
- Rolls back only the most recent migration.

### Lint Code
Run code quality checks:
```bash
make lint
```
- Executes `go vet` and `staticcheck` with predefined checks (`all,-ST1000,-ST1003`).

### Start Services
Start all services (database, Adminer, server) in detached mode:
```bash
make run
```
- Builds the server image with `APP_VERSION` and `GIT_COMMIT_SHA`.
- Uses the `docker-compose.yaml` file in `deployments/`.

### Start Services with Migrations
Start services and apply migrations automatically:
```bash
make run-auto
```
- Runs `make run` followed by `make apply-migrations`.

### Build Images
Build Docker images without cache:
```bash
make build
```
- Builds the server image with `APP_VERSION` and `GIT_COMMIT_SHA`.

### Stop Services
Stop and remove all services, networks, and volumes:
```bash
make stop
```

### View Logs
Display logs for all services (last 50 lines, following in real-time):
```bash
make logs
```

## Example Workflow

1. **Set Up Environment**
   Create a `configs/.env` file:
   ```env
   POSTGRES_USER=admin
   POSTGRES_PASSWORD=secret
   POSTGRES_DB=soccer_api
   DB_URL=postgresql://admin:secret@localhost:5432/soccer_api?sslmode=disable&search_path=public
   POSTGRES_HOST_PORT=5432
   ADMINER_HOST_PORT=8080
   SERVER_HOST_PORT=7777
   HTTP_SERVER_ADDRESS=7777
   ```

2. **Start the Project**
   ```bash
   make run-auto
   ```
   - Starts PostgreSQL, Adminer, and the server.
   - Applies any pending migrations.

3. **Access Services**
   - **Database**: Connect to PostgreSQL at `localhost:5432` (or the port defined in `POSTGRES_HOST_PORT`).
   - **Adminer**: Open `http://localhost:8080` (or the port defined in `ADMINER_HOST_PORT`) in a browser to manage the database.
   - **Server**: Access the API at `http://localhost:7777` (or the port defined in `SERVER_HOST_PORT`).

4. **Create and Apply a Migration**
   ```bash
   make create-migration NAME=add_teams_table
   make apply-migrations
   ```

5. **Stop the Project**
   ```bash
   make stop
   ```

## Notes
- Ensure the `DB_URL` is correctly formatted and accessible before running migration commands.
- The `GIT_COMMIT_SHA` and `APP_VERSION` are automatically set using Git commands if available.
- If you encounter issues, check logs with `make logs` for debugging.
- The project uses a bridge network (`soccer-api`) for service communication, ensuring isolation from other Docker networks.

For further details or troubleshooting, refer to the logs or the Docker Compose configuration.