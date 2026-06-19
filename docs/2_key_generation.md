#  Generate PKCS#12 (`.p12`) Keystore

This guide explains how to generate a private key, self-signed certificate, and bundle them into a `.p12` (PKCS#12) keystore. This is commonly used for enabling HTTPS in applications (e.g., Go servers).

---

##  Overview

A `.p12` file (PKCS#12 keystore) contains:

* Private Key 
* Certificate 
* (Optional) Certificate chain

It is widely used in Java, Go, and other systems requiring TLS/SSL configuration.

---

##  Step 1: Generate Private Key & Certificate

### Using RSA 2048 (standard)

```bash
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
```

###  Explanation

* `req -x509` → Generate self-signed certificate
* `-newkey rsa:2048` → Create new RSA key (2048-bit)
* `-keyout key.pem` → Output private key
* `-out cert.pem` → Output certificate
* `-days 365` → Valid for 1 year
* `-nodes` → No password on private key

---

###  Stronger option (RSA 4096)

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

👉 Use this for stronger security (at the cost of slightly higher CPU usage).

---

##  Step 2: Create PKCS#12 Keystore

```bash
openssl pkcs12 -export -out keystore.p12 -inkey key.pem -in cert.pem
```

### 🔍 Explanation

* `pkcs12 -export` → Create `.p12` bundle
* `-out keystore.p12` → Output file
* `-inkey key.pem` → Input private key
* `-in cert.pem` → Input certificate

 You will be prompted to set a **password** for the keystore.

---

##  Output Files

After running the commands, you will have:

```text
key.pem         # Private key
cert.pem        # Self-signed certificate
keystore.p12    # PKCS#12 keystore
```

---

##  Verify the Keystore

```bash
openssl pkcs12 -info -in keystore.p12
```

 This lets you inspect the contents of the `.p12` file.

---

##  Notes

* Self-signed certificates are suitable for **local development only**
* For production, use a trusted CA (e.g., Let's Encrypt)
* Keep your private key secure (`key.pem` should not be committed to Git)

---


##  Reference

* https://trenchesdeveloper.medium.com/how-to-add-https-and-http-2-to-your-go-server-b7dc7834b0a1

