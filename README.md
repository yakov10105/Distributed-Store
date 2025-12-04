# E-Commerce Distributed System

A learning project for building a distributed e-commerce system with Go microservices, React, Kafka, and Kubernetes.

## Prerequisites

- Go 1.22+
- Node.js 20+ & npm
- Docker
- Minikube (Kubernetes)
- Helm (Kubernetes Package Manager)
- Skaffold (Dev Workflow Tool)
- Protoc (Protocol Buffers Compiler)

## Getting Started (Recommended)

We use **Skaffold** for a seamless "Build, Deploy, Watch" workflow.

### 1. Start Minikube
```bash
minikube start
```

### 2. Run with Skaffold
This will build all images, deploy the Helm chart to Minikube, and tail the logs.
```bash
skaffold dev
```
*   **Skaffold** automatically detects changes in your code and re-deploys the services.
*   **Helm** manages the Kubernetes manifests.

## Manual Deployment

If you prefer to do it manually without Skaffold:

### 1. Configure Docker Environment
```bash
# PowerShell
minikube -p minikube docker-env --shell powershell | Invoke-Expression

# Bash
eval $(minikube docker-env)
```

### 2. Install via Helm
```bash
helm install my-store ./deploy/helm/my-store
```

### 3. Access Services
- **Frontend**: `minikube service frontend`
- **BFF**: `minikube service bff-service`

## Development

### Regenerate Protobufs
If you have `protoc` installed:
```bash
# Windows (PowerShell)
protoc --go_out=./pkg --go_opt=paths=source_relative --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative proto/*.proto
```
