# go-auth-service — K3s Deployment Guide

Complete step-by-step guide to deploy the go-auth-service on a K3s cluster, from cluster creation to debugging and teardown.

> **Windows users**: K3s is Linux-only. You must use **WSL2** (Windows Subsystem for Linux). All `kubectl` and K3s commands run inside WSL2. Docker Desktop with WSL2 backend is used for building images.

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Prerequisites](#2-prerequisites)
3. [K3s Cluster Setup](#3-k3s-cluster-setup)
4. [Pre-Deployment Configuration](#4-pre-deployment-configuration)
5. [Build & Import Docker Image](#5-build--import-docker-image)
6. [Deploy to K3s](#6-deploy-to-k3s)
7. [Verify Deployment](#7-verify-deployment)
8. [Access the Service](#8-access-the-service)
9. [Debugging & Troubleshooting](#9-debugging--troubleshooting)
10. [Scaling & Updates](#10-scaling--updates)
11. [Teardown & Cluster Destruction](#11-teardown--cluster-destruction)

---

## 1. Architecture Overview

### Docker Setup (2 variants)

| File | Builder | Runtime | Notes |
|------|---------|---------|-------|
| `Dockerfile` (root) | `golang:1.25-alpine` | `alpine:3.21` | Uses `.env` file via docker-compose |
| `docker-deploy/Dockerfile` | `golang:1.25-alpine` | `alpine:3.21` | Production-oriented, `.env` file |

Both use multi-stage builds producing a static binary with `CGO_ENABLED=0`. The service runs as non-root user (UID 1000) and exposes port **8080**.

### Docker Compose Services

- **auth-service** — the Go binary, port 8080, health check via `/health`
- **mysql** — MySQL 8.0, port 3306, persistent volume, auto-runs migration SQL scripts from `migrations/`

### Kubernetes Manifests (`k8s/`)

| Manifest | Kind | Purpose |
|----------|------|---------|
| `deployment.yaml` | Deployment | 3 replicas, resource requests/limits, liveness + readiness probes, security context (non-root) |
| `service.yaml` | Service (ClusterIP) | Maps port 80 → 8080 |
| `configmap.yaml` | ConfigMap | Non-sensitive config: mysql host/port, database name, CORS origins |
| `secrets.yaml` | Secret | mysql-user, mysql-password, jwt-secret, super-admin-code |
| `mysql.yaml` | StatefulSet + Service + ConfigMap | MySQL 8.0 with 5Gi PVC, init SQL scripts embedded in ConfigMap |
| `ingress.yaml` | Ingress | Traefik ingress, host `auth.<NODE_IP>.nip.io` |
| `hpa.yaml` | HorizontalPodAutoscaler | 3–10 replicas, scale on CPU 70% / Memory 80% |
| `pdb.yaml` | PodDisruptionBudget | minAvailable: 2 |

### Environment Variables

| Variable | Source | Description |
|----------|--------|-------------|
| `APP_PORT` | Hardcoded `8080` | Application listen port |
| `APP_ENV` | Hardcoded `production` | Environment mode |
| `APP_VERSION` | Hardcoded `1.0.0` | Version string |
| `MYSQL_HOST` | ConfigMap `auth-config` | MySQL service hostname |
| `MYSQL_PORT` | ConfigMap `auth-config` | MySQL port |
| `MYSQL_USER` | Secret `auth-secrets` | Database user |
| `MYSQL_PASSWORD` | Secret `auth-secrets` | Database password |
| `MYSQL_DB` | ConfigMap `auth-config` | Database name |
| `JWT_SECRET` | Secret `auth-secrets` | JWT signing key (min 32 chars) |
| `SUPER_ADMIN_CODE` | Secret `auth-secrets` | Admin registration code |
| `CORS_ALLOWED_ORIGINS` | ConfigMap `auth-config` | Comma-separated CORS origins |

### Health Check Endpoints

| Endpoint | Used By | Description |
|----------|---------|-------------|
| `/health` | Docker HEALTHCHECK | General health |
| `/health/live` | K8s livenessProbe | Liveness check |
| `/health/ready` | K8s readinessProbe | Readiness check (DB connectivity) |

---

## 2. Prerequisites

### Windows (WSL2)

```powershell
# 1. Enable WSL2 (PowerShell as Admin)
wsl --install
# Reboot your machine after this

# 2. Install Ubuntu (default distro)
wsl --install -d Ubuntu

# 3. Open WSL2 terminal
wsl

# Inside WSL2, update packages
sudo apt update && sudo apt upgrade -y
```

Install Docker Desktop with WSL2 backend:

```powershell
# Download and install Docker Desktop from https://www.docker.com/products/docker-desktop/
# In Docker Desktop settings:
#   → General → check "Use the WSL 2 based engine"
#   → Resources → WSL Integration → enable your Ubuntu distro
```

### Linux (native)

```bash
# Required tools
# - Ubuntu 22.04+ or similar
# - Docker
# - curl, sudo/root access
```

### Install kubectl (inside WSL2 or Linux)

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Verify
kubectl version --client
```

---

## 3. K3s Cluster Setup

> **All commands in this section run inside WSL2 (Windows) or terminal (Linux).**

### 3.1 Install K3s (Single Node)

```bash
# Install K3s with Traefik enabled (default)
curl -sfL https://get.k3s.io | sh -

# Wait ~30 seconds, then check
sudo systemctl status k3s

# Set up kubectl config so your user can use kubectl
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $(id -u):$(id -g) ~/.kube/config
echo 'export KUBECONFIG=~/.kube/config' >> ~/.bashrc
source ~/.bashrc
```

### 3.2 Verify Cluster

```bash
# Check node status
kubectl get nodes

# Expected output:
# NAME        STATUS   ROLES                  AGE   VERSION
# hostname    Ready    control-plane,master   30s   v1.28.x+k3s1

# Check system pods
kubectl get pods -A

# Verify Traefik is running (K3s includes it by default)
kubectl get pods -n kube-system -l app.kubernetes.io/name=traefik
```

### 3.3 Get Node IP (for Ingress)

```bash
# Inside WSL2 — get the WSL2 VM IP
hostname -I | awk '{print $1}'

# Store it
NODE_IP=$(hostname -I | awk '{print $1}')
echo "Node IP: $NODE_IP"

# This IP will be used in the ingress host: auth.${NODE_IP}.nip.io
```

**Windows note**: If accessing from Windows browser, use `localhost` or the WSL2 IP directly. See [Section 8](#8-access-the-service).

---

## 4. Pre-Deployment Configuration

### 4.1 Update Secrets

Edit `k8s/secrets.yaml` — **change all placeholder values**:

```bash
cd go-auth-service

cat > k8s/secrets.yaml << 'EOF'
apiVersion: v1
kind: Secret
metadata:
  name: auth-secrets
type: Opaque
stringData:
  mysql-user: "root"
  mysql-password: "YourStrongPassword123!"
  jwt-secret: "my-super-secret-jwt-key-at-least-32-characters-long"
  super-admin-code: "MY_SUPER_ADMIN_SECRET_2024"
EOF
```

### 4.2 Update ConfigMap and Ingress with Node IP

The `configmap.yaml` and `ingress.yaml` contain `<NODE_IP>` placeholders. Replace them:

```bash
NODE_IP=$(hostname -I | awk '{print $1}')

# Update configmap
sed -i "s|<NODE_IP>|${NODE_IP}|g" k8s/configmap.yaml

# Update ingress
sed -i "s|<NODE_IP>|${NODE_IP}|g" k8s/ingress.yaml

# Verify
grep "${NODE_IP}" k8s/configmap.yaml k8s/ingress.yaml
```

**Windows (PowerShell alternative)**:

```powershell
# Get WSL2 IP from PowerShell
wsl hostname -I
# Example output: 172.25.160.1

# Then manually edit the files in VS Code:
code go-auth-service/k8s/configmap.yaml
code go-auth-service/k8s/ingress.yaml
# Replace <NODE_IP> with the IP (e.g., 172.25.160.1)
```

### 4.3 Create Namespace (Optional)

```bash
kubectl create namespace auth-service
# Then add: -n auth-service to all kubectl apply commands
```

---

## 5. Build & Import Docker Image

K3s bundles containerd, not Docker. You build the image with Docker, then import it into K3s.

### 5.1 Build the Image

**WSL2 / Linux:**

```bash
# From the go-auth-service directory
docker build -t auth-service:latest -f docker-deploy/Dockerfile .

# Or using the root Dockerfile
# docker build -t auth-service:latest .
```

**Windows (PowerShell):**

```powershell
# Build from PowerShell using Docker Desktop
docker build -t auth-service:latest -f go-auth-service\docker-deploy\Dockerfile go-auth-service\
```

### 5.2 Import into K3s

**WSL2 (if Docker runs inside WSL2):**

```bash
# Save image to tarball
docker save auth-service:latest -o /tmp/auth-service.tar

# Import into K3s containerd
sudo k3s ctr images import /tmp/auth-service.tar

# Verify
sudo k3s crictl images | grep auth-service
```

**Windows (Docker Desktop, build from PowerShell, import from WSL2):**

```powershell
# 1. Save image from PowerShell
docker save auth-service:latest -o $env:TEMP\auth-service.tar

# 2. Copy to WSL2
wsl cp /mnt/c/Users/$env:USERNAME/AppData/Local/Temp/auth-service.tar /tmp/auth-service.tar

# 3. Import into K3s (from WSL2)
wsl sudo k3s ctr images import /tmp/auth-service.tar
```

**Alternative**: If you have a private registry, update `deployment.yaml` `imagePullPolicy` and image name accordingly.

---

## 6. Deploy to K3s

### 6.1 Deploy All Resources

```bash
cd go-auth-service

# Apply in order (dependencies first)
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/mysql.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/pdb.yaml
kubectl apply -f k8s/ingress.yaml
kubectl apply -f k8s/hpa.yaml
```

Or use the Makefile (run inside WSL2):

```bash
make k8s-deploy
```

### 6.2 Watch Deployment Progress

```bash
# Watch pods coming up
kubectl get pods -w

# Watch MySQL first (auth-service depends on it)
kubectl get pods -l app=mysql -w
```

---

## 7. Verify Deployment

### 7.1 Check All Resources

```bash
kubectl get pods -o wide
kubectl get svc
kubectl get ingress
kubectl get hpa
kubectl get pdb
kubectl get statefulset
```

### 7.2 Check Pod Status

```bash
kubectl describe pods -l app=auth-service
kubectl describe pods -l app=mysql
```

### 7.3 Verify Health

```bash
# Port-forward to test
kubectl port-forward svc/auth-service 8080:80 &

# Test health endpoint
curl http://localhost:8080/health
curl http://localhost:8080/health/ready
curl http://localhost:8080/health/live

# Kill port-forward
kill %1
```

**Windows (test from browser)**: Open `http://localhost:8080/health` in your browser after running `kubectl port-forward` from WSL2.

---

## 8. Access the Service

### 8.1 Via Ingress (Recommended for Linux)

```bash
NODE_IP=$(hostname -I | awk '{print $1}')
curl http://auth.${NODE_IP}.nip.io/health
```

### 8.2 Via Port-Forward (Best for Windows)

This is the easiest method for Windows users:

```bash
# In WSL2 terminal
kubectl port-forward svc/auth-service 8080:80
```

Then open in Windows browser: `http://localhost:8080/health`

### 8.3 Via NodePort (Alternative)

```bash
kubectl patch svc auth-service -p '{"spec":{"type":"NodePort"}}'
NODE_PORT=$(kubectl get svc auth-service -o jsonpath='{.spec.ports[0].nodePort}')
NODE_IP=$(hostname -I | awk '{print $1}')
echo "Access at: http://${NODE_IP}:${NODE_PORT}/health"
```

### 8.4 Test API Endpoints

```bash
# Health check
curl -s http://localhost:8080/health | jq

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPassword123!",
    "super_admin_code": "MY_SUPER_ADMIN_SECRET_2024"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPassword123!"
  }'
```

---

## 9. Debugging & Troubleshooting

### 9.1 Pod Not Starting

```bash
kubectl get pods
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl logs <pod-name> --previous   # previous container if crashlooping
```

### 9.2 MySQL Connection Issues

```bash
# Check MySQL pod
kubectl get pods -l app=mysql
kubectl logs <mysql-pod-name>

# Exec into auth-service pod and test connectivity
kubectl exec -it <auth-service-pod> -- sh
wget -qO- http://localhost:8080/health/ready

# Check if MySQL service DNS resolves
kubectl exec -it <auth-service-pod> -- nslookup mysql-service

# Port-forward MySQL to test from Windows
kubectl port-forward svc/mysql-service 3306:3306
# Then from another terminal:
mysql -h 127.0.0.1 -u root -p -D auth_service
```

### 9.3 Check Application Logs

```bash
kubectl logs -f -l app=auth-service --all-containers=true
kubectl logs -f <pod-name>
kubectl logs <pod-name> --timestamps
kubectl logs <pod-name> --tail=100
```

### 9.4 Exec into Running Pod

```bash
kubectl exec -it <pod-name> -- sh

# Inside the pod:
env | grep MYSQL
env | grep JWT
wget -qO- http://localhost:8080/health
nslookup mysql-service
```

### 9.5 Check Resource Usage

```bash
kubectl top nodes
kubectl top pods
kubectl top pods -l app=auth-service
```

### 9.6 Ingress Not Working

```bash
kubectl describe ingress auth-service-ingress

# Check Traefik pods
kubectl get pods -n kube-system -l app.kubernetes.io/name=traefik
kubectl logs -n kube-system -l app.kubernetes.io/name=traefik
kubectl get svc -n kube-system traefik
```

### 9.7 Secrets/ConfigMap Issues

```bash
kubectl get configmap auth-config -o yaml
kubectl get secret auth-secrets -o yaml
kubectl get secret auth-secrets -o jsonpath='{.data.mysql-password}' | base64 -d
```

### 9.8 Events & Diagnostics

```bash
kubectl get events --sort-by='.lastTimestamp'
kubectl describe pod <pod-name> | grep -A 10 Events
kubectl get events --field-selector type=Warning
```

### 9.9 HPA Not Scaling

```bash
kubectl describe hpa auth-service-hpa
kubectl get pods -n kube-system -l k8s-app=metrics-server
kubectl top pods -l app=auth-service
```

### 9.10 PVC / Storage Issues (MySQL)

```bash
kubectl get pvc
kubectl describe pvc mysql-data-mysql-0
kubectl get storageclass
```

### 9.11 Common Debug Commands Reference

```bash
kubectl get all
kubectl cluster-info
kubectl cluster-info dump > cluster-dump.txt
kubectl api-resources

# K3s specific
sudo systemctl status k3s
sudo journalctl -u k3s --since "1 hour ago"
```

---

## 10. Scaling & Updates

### 10.1 Manual Scaling

```bash
kubectl scale deployment auth-service --replicas=5
kubectl get pods -l app=auth-service
```

**Note**: When HPA is active, manual scaling will be overridden by HPA based on CPU/memory metrics.

### 10.2 HPA Auto-Scaling

The HPA is configured to scale between 3–10 replicas based on:
- CPU utilization target: 70%
- Memory utilization target: 80%

```bash
kubectl get hpa auth-service-hpa -w
```

### 10.3 Rolling Update

```bash
# 1. Rebuild
docker build -t auth-service:v2 -f docker-deploy/Dockerfile .

# 2. Save and import into K3s
docker save auth-service:v2 -o /tmp/auth-service-v2.tar
sudo k3s ctr images import /tmp/auth-service-v2.tar

# 3. Update deployment image
kubectl set image deployment/auth-service auth-service=auth-service:v2

# 4. Watch rollout
kubectl rollout status deployment/auth-service

# Rollback if needed
kubectl rollout undo deployment/auth-service
```

### 10.4 Update ConfigMap/Secrets

```bash
kubectl edit configmap auth-config
kubectl edit secret auth-secrets

# Restart pods to pick up changes
kubectl rollout restart deployment/auth-service
kubectl rollout restart statefulset/mysql
```

---

## 11. Teardown & Cluster Destruction

### 11.1 Delete Application Resources

```bash
cd go-auth-service

# Delete in reverse order
kubectl delete -f k8s/hpa.yaml
kubectl delete -f k8s/ingress.yaml
kubectl delete -f k8s/pdb.yaml
kubectl delete -f k8s/service.yaml
kubectl delete -f k8s/deployment.yaml
kubectl delete -f k8s/mysql.yaml
kubectl delete -f k8s/secrets.yaml
kubectl delete -f k8s/configmap.yaml

# Or use Makefile
make k8s-delete
```

### 11.2 Delete PVCs (Persistent Data)

```bash
kubectl get pvc
kubectl delete pvc mysql-data-mysql-0
# WARNING: This deletes all MySQL data permanently
```

### 11.3 Delete Namespace (if used)

```bash
kubectl delete namespace auth-service
```

### 11.4 Remove Docker Image from K3s

```bash
sudo k3s ctr images remove docker.io/library/auth-service:latest
```

### 11.5 Uninstall K3s

```bash
# K3s provides an uninstall script
/usr/local/bin/k3s-uninstall.sh

# This removes:
# - K3s binary
# - All containerd data
# - All cluster data
# - systemd service
# - kubeconfig

# Clean up leftover files (if any)
sudo rm -rf /etc/rancher/k3s
sudo rm -rf /var/lib/rancher/k3s
sudo rm -rf /var/lib/kubelet
```

### 11.6 Full Cleanup Checklist

```bash
# 1. Delete app resources
make k8s-delete

# 2. Delete PVCs
kubectl delete pvc --all

# 3. Remove images
sudo k3s ctr images rm auth-service:latest
sudo k3s ctr images rm mysql:8.0

# 4. Uninstall K3s
/usr/local/bin/k3s-uninstall.sh

# 5. Remove saved tarballs
rm -f /tmp/auth-service.tar /tmp/auth-service-v2.tar

# 6. Remove kubeconfig
rm -rf ~/.kube/config

echo "Cleanup complete."
```

---

## Quick Reference Card

```bash
# --- Cluster ---
curl -sfL https://get.k3s.io | sh -                  # Install K3s
sudo systemctl status k3s                             # Check K3s status
kubectl get nodes                                     # List nodes
kubectl get all                                       # All resources

# --- Deploy ---
make k8s-deploy                                       # Deploy everything
make k8s-delete                                       # Delete everything

# --- Debug ---
kubectl get pods -w                                   # Watch pods
kubectl logs -f -l app=auth-service                   # Stream logs
kubectl describe pod <name>                           # Pod details
kubectl exec -it <pod> -- sh                          # Shell into pod
kubectl top pods                                      # Resource usage

# --- Access (Windows) ---
kubectl port-forward svc/auth-service 8080:80         # Then open localhost:8080

# --- Access (Linux) ---
NODE_IP=$(hostname -I | awk '{print $1}')
curl http://auth.${NODE_IP}.nip.io/health

# --- Cleanup ---
/usr/local/bin/k3s-uninstall.sh                       # Remove K3s
```
