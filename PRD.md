# Project Requirements Document (PRD) & Deliverables

## Project Overview

A distributed e-commerce system for learning advanced microservices concepts. The system uses React (TypeScript) for the frontend, Go for backend services, Kafka for asynchronous communication, and runs on Kubernetes (Minikube).

## Tech Stack

- **Frontend**: React, TypeScript, Vite, Tailwind CSS
- **BFF (Backend for Frontend)**: Go, gRPC Client, REST API
- **Microservices**: Go, gRPC Server
- **Messaging**: Apache Kafka
- **Database**: PostgreSQL (Auth, Order, Shipping)
- **Infrastructure**: Kubernetes (Minikube), Docker, Helm, Tilt, ArgoCD

## Deliverable Phases

### Phase 1: Scaffolding & Infrastructure Setup

- [x] **Task 1.1**: Initialize Monorepo Structure
  - Set up root directory, git, and go workspace.
  - Create `pkg` directory for shared libraries.
- [x] **Task 1.2**: Define gRPC Protocol Buffers
  - Create `/proto` directory.
  - Define `auth.proto`, `order.proto`, `shipping.proto`.
  - configure `protoc` generation.
- [x] **Task 1.3**: Scaffold Go Services
  - Initialize `bff`, `auth`, `order`, `shipping`, `notification`, `analytics` services.
  - Implement basic gRPC server scaffolding.
- [x] **Task 1.4**: Initialize Frontend
  - Setup React + TypeScript + Vite in `/frontend`.
- [x] **Task 1.5**: Containerization
  - Create `Dockerfile` for all services and frontend.
- [x] **Task 1.6**: Kubernetes Manifests
  - Create deployments for Kafka, Zookeeper, Postgres.
  - Create deployments/services for all microservices.

### Phase 2: Core Logic Implementation (Current)

#### Task 2.1: Authentication Service Logic

- [x] **2.1.1 Implement User Store (PostgreSQL)**
  - Update Helm chart to provision `auth_db`.
  - Create `users` table with email/password columns.
  - Implement `pgx` database connection and queries in `auth` service.
- [x] **2.1.2 Implement Registration**
  - Implement `Register` gRPC method.
  - Validate email uniqueness.
  - Hash password using `bcrypt`.
- [x] **2.1.3 Implement Login & JWT Generation**
  - Implement `Login` gRPC method.
  - Verify password hash.
  - Generate JWT token with claims (UserId, Expiry).
- [x] **2.1.4 Implement Validation**
  - Implement `Validate` gRPC method to parse/verify JWT tokens.

#### Task 2.2: Order Service Logic

- [x] **2.2.1 Implement Order Store (PostgreSQL)**
  - Update Helm chart to provision `order_db`.
  - Create `orders` table (stores items as JSONB for simplicity initially).
  - Implement `pgx` database connection and queries in `order` service.
- [x] **2.2.2 Implement Create Order**
  - Implement `CreateOrder` gRPC method.
  - Generate unique Order ID.
  - Set initial status to `PENDING`.

#### Task 2.3: Shipping Service Logic (PostgreSQL)

- [ ] **2.3.1 Database Connection**
  - Add `lib/pq` or `pgx` driver to Shipping service.
  - Implement DB connection logic using environment variables.
  - Create migration script for `shipments` table.
- [ ] **2.3.2 Implement Create Shipment**
  - Implement `CreateShipment` gRPC method.
  - Insert shipment record into PostgreSQL.
  - Return generated Tracking ID.

#### Task 2.4: BFF & Inter-Service Communication

- [x] **2.4.1 Setup gRPC Clients**
  - Initialize gRPC clients (Auth, Order, Shipping) in BFF main.
  - Implement connection pooling/management.
- [x] **2.4.2 Implement HTTP Handlers**
  - `POST /api/auth/register` -> Calls Auth.Register
  - `POST /api/auth/login` -> Calls Auth.Login
  - `POST /api/orders` -> Validates Token (Auth) -> Calls Order.CreateOrder
- [x] **2.4.3 Setup CORS & Router**
  - Install `rs/cors`.
  - Configure Middleware for Auth.

### Phase 3: Async Communication & Analytics (Upcoming)

#### Task 3.1: Kafka Integration

- [ ] **3.1.1 Setup Kafka Producer (Order Service)**
  - Configure Sarama or Confluent Go client.
  - Publish `OrderCreated` event to `orders` topic when order is created.
- [ ] **3.1.2 Setup Kafka Consumers**
  - **Shipping Service**: Listen to `orders` -> Create Shipment automatically.
  - **Notification Service**: Listen to `orders` -> Log "Email Sent".
  - **Analytics Service**: Listen to `orders` -> Update stats.

#### Task 3.2: Frontend Integration

- [x] **3.2.1 Build API Client** (Partially done for Auth)
  - Generate/Create TypeScript interfaces for BFF endpoints.
- [x] **3.2.2 Implement UI Components**
  - [x] Registration Form.
  - [ ] Login Form.
  - [ ] Product List (Mock).
  - [ ] Order History View.
- [x] **3.2.3 UI Architecture**
  - Setup React Router.
  - Configure Layouts and Pages.
  - Configure Tailwind CSS.

### Phase 4: DevEx & Deployment (In Progress)

#### Task 4.1: Helm Chart Migration

- [x] **4.1.1 Create Base Helm Chart**
  - Initialize `deploy/helm/my-store` chart.
  - Create templates for Deployment, Service, Ingress.
- [x] **4.1.2 Migrate Infrastructure**
  - Move Postgres and Kafka/Zookeeper manifests to Helm templates (or use sub-charts).
  - Add pgAdmin for database management.
- [x] **4.1.3 Migrate Microservices**
  - Move all service deployments to use the Helm chart with `values.yaml`.

#### Task 4.2: Local Development Setup (Tilt)

- [x] **4.2.1 Configure Tilt**
  - Create `Tiltfile`.
  - Configure build artifacts for all services (Auth, Order, Shipping, Notification, Analytics, BFF, Frontend).
  - Configure deploy to use Helm Chart (`deploy/helm/my-store`).
  - Setup port forwarding and resource grouping.

#### Task 4.3: CI/CD Pipeline

- [ ] **4.3.1 GitHub Actions**
  - Create workflow for Build & Push.
- [ ] **4.3.2 ArgoCD Setup**
  - Install ArgoCD in Minikube.
  - Configure Application to point to the Helm chart.
