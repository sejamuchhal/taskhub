# TaskHub

![image](https://github.com/user-attachments/assets/3321f286-70f8-48e6-9c92-68a56ff916e3)


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
   OR
   ```bash
   make build
   ```

   This will start the following services:

   - Gateway REST API service at port 3000
      [http://localhost:3000/health](http://localhost:3000/health)
   - Auth gRPC service at port 4040
   - Task gRPC service at port 8080
   - Notification background worker
   - PostgreSQL database server on port 5432
   - RabbitMQ message broker on ports 5672 (with UI exposed at 15672)
      [http://localhost:15672](http://localhost:15672)
   - Prometheus basic monitoring server on port 9090, with metrics exposed via the gateway endpoint 
      [http://localhost:3000/matrics](http://localhost:3000/matrics)
