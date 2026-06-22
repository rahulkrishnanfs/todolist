#  OpenSSL: Generate CA and Sign Server Certificate

This guide walks through creating a **Certificate Authority (CA)** and using it to sign a **server certificate (`todolist.ai`)**.

---

##  Overview

We will generate:

### Certificate Authority (CA)

* `ca.key` → Private key (**keep secure**)
* `ca.crt` → Public certificate

### Server Certificate

* `server.key` → Private key
* `server.csr` → Certificate Signing Request
* `server.crt` → Signed certificate

---

##  Step 1: Generate CA Private Key

```bash
openssl genrsa -out ca.key 4096
```

* Generates a 4096-bit RSA private key
* This is your **root of trust**

---

##  Step 2: Generate Self-Signed CA Certificate

```bash
openssl req -x509 -newkey rsa:4096 -sha512 -days 3653 -nodes \
-keyout ca.key -out ca.crt \
-subj "/C=UK/ST=Northamptonshire/L=Northampton/O=zaagpro/CN=zaagpro.com" \
-addext "basicConstraints=critical,CA:TRUE" \
-addext "keyUsage=critical,keyCertSign,cRLSign"
```

###  Explanation

* `-x509` → Create self-signed certificate
* `-days 3653` → Valid for ~10 years
* `-nodes` → No password on private key
* `CA:TRUE` → Marks certificate as a CA
* `keyCertSign` → Allows signing other certificates

---

##  Step 3: Verify CA Certificate

```bash
openssl x509 -in ca.crt -text -noout
```

Check for:

* `CA:TRUE`
* Key usage includes `Certificate Sign`

---

##  Step 4: Generate Server Private Key

```bash
openssl genrsa -out server.key 4096
```

---

##  Step 5: Generate CSR (Certificate Signing Request)

```bash
openssl req -new -key server.key -out server.csr \
-subj "/C=UK/ST=Northamptonshire/L=Northampton/O=zaagpro/CN=todolist.ai" \
-addext "subjectAltName=DNS:todolist.ai"
```

###  Important

* `CN = todolist.ai`
* `subjectAltName` is **mandatory** (modern TLS ignores CN)

---

##  Step 6: Sign Server Certificate with CA

```bash
openssl x509 -req -in server.csr \
-CA ca.crt -CAkey ca.key -CAcreateserial \
-out server.crt -days 3653 -sha256 \
-copy_extensions copyall
```

###  Explanation

* `-CAcreateserial` → Generates serial file
* `-copy_extensions copyall` → Copies SAN from CSR

---

##  Final Output

```
ca.key       # CA private key
ca.crt       # CA certificate

server.key   # server private key
server.csr   # certificate request
server.crt   # signed server certificate
```

---

##  Step 7: Verify Server Certificate

```bash
openssl verify -CAfile ca.crt server.crt
```

Expected output:

```
server.crt: OK
```

---

##  Key Concepts

###  Certificate Authority (CA)

Trusted entity that signs certificates.

###  CSR (Certificate Signing Request)

Contains identity + public key.

###  SAN (Subject Alternative Name)

Required for hostname validation.


---

##  Example: Multiple SANs

```bash
-addext "subjectAltName=DNS:todolist.ai,DNS:www.todolist.ai,IP:127.0.0.1"
```

---

##  Usage

* Use `server.crt` + `server.key` in your HTTPS server
* Use `ca.crt` to establish trust

---

##  Summary

1. Create CA → `ca.key`, `ca.crt`
2. Create server key → `server.key`
3. Generate CSR → `server.csr`
4. Sign with CA → `server.crt`
5. Verify → `openssl verify`

---

##  Suitable For

* Local development
* Internal services
* Kubernetes clusters
* mTLS setups

---
