
# Setting Kind cluster for Service Deployment

This guide explains how to create a local Kubernetes cluster using `kind`, install NGINX Ingress, enable SSL passthrough, and deploy your application using Helm.

---

## 1. Create Kind Cluster

```bash
kind create cluster --config kind-config.yaml
```

###  Explanation

* Creates a local Kubernetes cluster using **kind (Kubernetes in Docker)**.
* The `kind-config.yaml` should expose ports **80 and 443** so ingress works properly.

Example:

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 80
        hostPort: 80
      - containerPort: 443
        hostPort: 443
```

---

##  2. Install NGINX Ingress Controller

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

###  Explanation

* Installs the **NGINX Ingress Controller** into your cluster.
* This component handles incoming HTTP/HTTPS traffic and routes it to your services.

---

##  3. Enable SSL Passthrough

```bash
kubectl edit deployment ingress-nginx-controller -n ingress-nginx
```

### Add this argument under `args:`:

```yaml
- --enable-ssl-passthrough
```

###  Explanation

* Enables **TLS passthrough**, meaning:

  * NGINX does NOT terminate TLS
  * Traffic is forwarded directly to your backend service
* Required when your **application itself handles HTTPS**

---

##  4. Restart Ingress Controller

```bash
kubectl rollout restart deployment ingress-nginx-controller -n ingress-nginx
```

###  Explanation

* Applies the configuration change by restarting the controller.
* Without this, the new flag won’t take effect.

---

##  5. Deploy Application using Helm

```bash
helm install todolist charts/todolist
```

###  Explanation

* Deploys your application using a Helm chart.
* The chart should include:

  * Deployment
  * Service (exposing port 443 → targetPort 8080 if HTTPS)
  * Ingress configuration

---

##  6. Uninstall Application

```bash
helm uninstall todolist
```

###  Explanation

* Removes the deployed application and all associated Kubernetes resources.

---

##  Testing

Add domain mapping:

```bash
sudo nano /etc/hosts
```

```text
127.0.0.1 todolist.ai
```

Test:

```bash
curl -k https://todolist.ai
```

---


##  Architecture Overview

```
Browser → https://todolist.ai
        ↓
NGINX Ingress (SSL Passthrough)
        ↓
Service (port 443)
        ↓
Pod (HTTPS on port 8080)
```