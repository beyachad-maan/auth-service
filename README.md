# auth-service
Authentication and Authorization service for beyachad-maan application

## Overview
This is a Go-based authentication service that provides JWT-based authentication and user management for the beyachad-maan application. The service runs on HTTPS using TLS certificates and connects to a PostgreSQL database.

## Prerequisites
- [CodeReady Containers (CRC)](https://developers.redhat.com/products/openshift-local/overview) running
- Docker installed and configured
- `kubectl` or `oc` CLI tools installed
- Go 1.19+ (for local development)

## Quick Start with CRC

### 1. Start CRC and Setup Environment
```bash
# Start CRC (if not already running)
crc start

# Login to the cluster
oc login -u kubeadmin https://api.crc.testing:6443

# Create project namespace
oc new-project beyachad-maan
```

### 2. Build and Deploy

Use the provided Makefile targets for easy deployment:

```bash
# Build the Go application locally
make build

# Build Docker image
make image

# Push image to CRC's internal registry
make image-push

# Deploy to Kubernetes
make deploy
```

### 3. Verify Deployment
```bash
# Check pod status
kubectl get pods -l app=auth-service

# Check service
kubectl get svc auth-service

# View logs
kubectl logs -l app=auth-service
```

## Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Compile the Go application locally |
| `make image` | Build Docker image with tag `auth-service:latest` |
| `make image-push` | Push Docker image to CRC's internal OpenShift registry |
| `make deploy` | Deploy the application and PostgreSQL to Kubernetes using `deployment.yaml` |
| `make clean-deploy` | Remove all deployed resources from Kubernetes |

## Manual Deployment Steps

If you prefer to run commands manually:

### Build and Push Image
```bash
# Build the application
go build .

# Build Docker image
docker build -f ./Dockerfile . -t auth-service:latest

# Login to CRC registry
docker login -u kubeadmin -p $(oc whoami -t) default-route-openshift-image-registry.apps-crc.testing

# Tag and push image
docker tag auth-service:latest default-route-openshift-image-registry.apps-crc.testing/beyachad-maan/auth-service:latest
docker push default-route-openshift-image-registry.apps-crc.testing/beyachad-maan/auth-service:latest
```

### Deploy to Kubernetes
```bash
# Apply the deployment configuration
kubectl apply -f ./deployment.yaml
```

## Service Configuration

The service includes:
- **Auth Service**: Runs on port 8443 (mapped to external port 443)
- **PostgreSQL Database**: Internal database for user storage
- **TLS Configuration**: Uses mounted certificates from Kubernetes secrets
- **JWT Keys**: RSA key pair for token signing and verification

## Certificates and Secrets

The deployment expects the following Kubernetes secrets:
- `auth-service-tls-secret`: TLS certificate and private key
- `jwt-key-pair`: JWT signing keys (public and private)

These are included in the `deployment.yaml` file.

## Local Development

For local development without CRC:
```bash
# Build the application
make build

# Run locally (requires local PostgreSQL)
./auth-service run
```

## Cleanup

To remove all deployed resources:
```bash
make clean-deploy
```
