# TaskHub

TaskHub is a microservices-based task management system.

## Prerequisites

- Docker and Docker Compose

## Getting Started

To try this project:

1. Clone the repository:

   ```bash
   git clone github.com/sejamuchhal/taskhub
   cd taskhub
   ```

2. Create an environment file by copying and modifying `.env.example` to `.env`. Set the following environment variables:

   ```plaintext
   MAILERSEND_API_KEY
   MAILERSEND_SENDER_EMAIL
   ```

3. Start the services using Docker Compose:

   ```bash
   docker-compose up --build
   ```

   This will start the following services:

   - Gateway REST API service at port 3000
   - Auth gRPC service at port 4040
   - Task gRPC service at port 8080
   - Notification background worker
   - PostgreSQL database server on port 5432
   - RabbitMQ message broker on ports 5672 (with UI exposed at 15672)
   - Prometheus basic monitoring server on port 9090, with metrics exposed via the gateway endpoint [matrics](http://localhost:3000/matrics)

## TODO

- [ ] **Authentication**: Implement refresh tokens and session functionality to manage access tokens securely.
- [ ] **Authorization**: Utilize user roles stored in user records for JWT token claims.
- [ ] **Testing**: Implement unit tests to ensure code reliability.
- [ ] **Documenation**: Provide proper documentation for the codebase.
