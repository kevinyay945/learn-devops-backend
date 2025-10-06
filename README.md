# Learn DevOps Backend

This is a simple Go application built with the Gin web framework, demonstrating basic DevOps practices including containerization with Docker and health checks.

## Features

- Go application using Gin framework.
- Environment variable loading with `godotenv/autoload`.
- Liveness and Readiness health check endpoints.
- Graceful shutdown endpoint for controlled server termination.
- Multi-stage Dockerfile for efficient image building.
- Makefile for common development and deployment tasks.

## Prerequisites

- Go (version 1.22 or higher recommended)
- Docker
- Make

## Getting Started

1.  **Clone the repository:**

    ```bash
    git clone <repository-url>
    cd learn-devops/backend
    ```

2.  **Environment Variables:**

    Create a `.env` file in the root of the project based on `.env.example`:

    ```bash
    cp .env.example .env
    ```

    You can modify the `PORT` variable in `.env` if you want to run the application on a different port.

## Running the Application

### Locally

To run the application directly using Go:

```bash
make run
```

The application will be accessible at `http://localhost:5000` (or your specified port).

### With Docker

1.  **Build the Docker image:**

    ```bash
    make build-docker
    ```

2.  **Run the Docker container:**

    ```bash
    docker run -d -p 5000:5000 --name learn-devops-backend learn-devops-backend
    ```

    This will run the container in detached mode and map port 5000 of the host to port 5000 of the container.

## Endpoints

-   **Root:** `GET /`
    Returns a simple JSON message.

-   **Liveness Probe:** `GET /health/liveness`
    Returns `{"status": "UP"}` if the application is running. Can be toggled to `{"status": "DOWN"}` via the `/health/liveness/toggle` endpoint.

-   **Toggle Liveness:** `POST /health/liveness/toggle`
    Toggles the liveness state between `UP` and `DOWN`.

-   **Readiness Probe:** `GET /health/readiness`
    Returns `{"status": "READY"}`. In a real-world scenario, this would check external dependencies (e.g., database connections).

-   **Environment Variable:** `GET /env`
    Returns the value of the `CURR_ENV` environment variable.

-   **Graceful Shutdown:** `POST /shutdown`
    Triggers a graceful shutdown of the server. Returns `{"message": "Server is shutting down gracefully..."}` and initiates the shutdown process with a 5-second timeout for ongoing requests.

### Health Check Command

You can also run a health check from the command line:

```bash
go run main.go health
```

Or, if you have the compiled binary:

```bash
./backend health
```

This command sends a `GET` request to the `/health/liveness` endpoint and will exit with a non-zero status code if the health check fails.

## Building

### Build Go Executable

To compile the Go application into an executable:

```bash
make build-app
```

This will create an executable named `backend` in the project root.

### Build Docker Image

To build the Docker image:

```bash
make build-docker
```

## Deployment

### Automated Versioning with GitHub Actions

This repository includes automated versioning workflows that follow [Semantic Versioning](https://semver.org/):

-   **Bump Minor Version:** Navigate to the "Actions" tab on GitHub, select "Bump Minor Version" workflow, and click "Run workflow". This increments the minor version (e.g., v1.0.0 → v1.1.0).
-   **Bump Patch Version:** Navigate to the "Actions" tab on GitHub, select "Bump Patch Version" workflow, and click "Run workflow". This increments the patch version (e.g., v1.0.0 → v1.0.1).

After a new tag is created, the "Build and Push Docker Image on Tag" workflow automatically builds and pushes a Docker image to GitHub Container Registry (ghcr.io).

#### Required Setup

To enable the automated workflows to trigger each other, you need to create a Personal Access Token (PAT):

1.  **Create a Personal Access Token:**
    - Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
    - Click "Generate new token (classic)"
    - Give it a descriptive name (e.g., "Learn DevOps Backend Workflows")
    - Select scopes: `repo` (Full control of private repositories), `workflow` (Update GitHub Action workflows), and `write:packages` (Upload packages to GitHub Package Registry)
    - Click "Generate token" and copy the token

2.  **Add the token as a repository secret:**
    - Navigate to your repository's Settings → Secrets and variables → Actions
    - Click "New repository secret"
    - Name: `PAT`
    - Value: Paste your Personal Access Token
    - Click "Add secret"

This PAT allows the bump version workflows to push tags that trigger the "Build and Push Docker Image on Tag" workflow, and enables pushing Docker images to GitHub Container Registry.

### Deploy to Docker Hub

To push the Docker image to Docker Hub:

1.  **Log in to Docker Hub:**

    ```bash
    docker login
    ```

2.  **Update `Makefile`:**

    Edit the `Makefile` and replace `<your-dockerhub-username>` with your actual Docker Hub username.

3.  **Push the image:**

    ```bash
    make deploy-dockerhub
    ```

## Project Structure

```
.
├── Dockerfile          # Defines the Docker image for the application.
├── Makefile            # Contains common commands for building, running, and deploying.
├── README.md           # This file.
├── go.mod              # Manages dependencies for the Go project.
├── go.sum              # Contains the checksums of the direct and indirect dependencies.
└── main.go             # The main entry point for the application.
```
