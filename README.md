# ja-eb
## 🔐 HD Wallet Service  
### Go Backend + Hugo Frontend

A production-oriented hierarchical deterministic (HD) wallet service built in Go, implementing:

- **BIP-39** → Mnemonic phrase generation  
- **BIP-32** → Deterministic hierarchical key derivation  
- RESTful API design using `net/http`  
- Thread-safe shared state using `sync.Mutex`  
- Static frontend powered by Hugo  

This project focuses on backend systems engineering principles including concurrency safety, deterministic cryptography, API design, and service observability.

---

# 📖 Table of Contents

- [Overview](#-overview)
- [System Architecture](#-system-architecture)
- [Backend Design](#-backend-design)
- [API Endpoints](#-api-endpoints)
- [Concurrency Model](#-concurrency-model)
- [Frontend (Hugo)](#-frontend-hugo)
- [Running the Project](#-running-the-project)
- [Example Usage](#-example-usage)
- [Security Considerations](#-security-considerations)
- [Engineering Goals](#-engineering-goals)

---

# 🚀 Overview

This service allows users to:

- Generate secure 12-word mnemonic phrases
- Derive multiple deterministic wallets from a single seed
- Interact through clean REST APIs
- Access a lightweight frontend dashboard
- Run the system locally with minimal setup

The wallet derivation is fully deterministic — meaning the same mnemonic + index always generates the same wallet.

The backend is intentionally designed to reflect infrastructure-ready service patterns.

---

# 🏗 System Architecture

```
Hugo Frontend (Static UI)
        ↓
HTTP Requests (Fetch API)
        ↓
Go HTTP Server (net/http)
        ↓
Validation Layer
        ↓
BIP-39 Seed Generation
        ↓
BIP-32 Child Key Derivation
        ↓
JSON Response
```

### Architectural Characteristics

- Stateless wallet derivation
- Deterministic key hierarchy
- Explicit validation layer
- Thread-safe shared memory
- Health endpoint for observability
- Separation of frontend and backend concerns

---

# ⚙ Backend Design

The backend is implemented using:

- `net/http` → HTTP server
- `encoding/json` → JSON responses
- `sync.Mutex` → Concurrency safety
- `go-bip39` → Mnemonic generation
- `go-bip32` → HD wallet derivation

### Key Design Decisions

✔ Explicit error handling (Go-style error returns)  
✔ Defensive validation (mnemonic + index checks)  
✔ Minimal dependencies  
✔ Clean modular project structure  
✔ No global state mutations without locking  

---

# 📡 API Endpoints

## 1️⃣ Generate Mnemonic

**POST** `/mnemonic`

Response:

```json
{
  "mnemonic": "believe tonight wool act artist unit grunt exact voyage dwarf one day"
}
```

---

## 2️⃣ Derive Wallet

**GET** `/derive?mnemonic=<mnemonic>&index=<number>`

Example:

```
/derive?mnemonic=<12-word-phrase>&index=0
```

Response:

```json
{
  "address": "derived_public_key"
}
```

Each incremented index produces a new deterministic wallet.

---

## 3️⃣ Health Check

**GET** `/health`

Response:

```
OK
```

### Why a Health Endpoint?

The health endpoint enables:

- Load balancer health checks
- Container liveness probes
- Service monitoring
- Infrastructure readiness validation

This mirrors real-world backend service patterns.

---

# 🔒 Concurrency Model

Go’s HTTP server handles requests concurrently.

To protect shared in-memory state:

```go
sync.Mutex
```

is used to guard access to the wallet map.

This prevents:

- Race conditions
- Data corruption
- Concurrent write conflicts

The design is ready to be upgraded to `sync.RWMutex` for read-heavy optimization.

---

# 🖥 Frontend (Hugo)

The frontend is built using **Hugo**, a static site generator.

Responsibilities:

- Trigger mnemonic generation
- Request wallet derivation
- Render deterministic wallet outputs
- Provide a minimal dashboard interface

The frontend communicates with the backend via Fetch API.

No client-side cryptography is performed.

---

# 📂 Project Structure

```
.
├── wallet-api/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
│
└── wallet-frontend/
    ├── layouts/
    ├── hugo.toml
    └── ...
```

---

# 🛠 Running the Project

## 1️⃣ Clone Repository

```bash
git clone https://github.com/<your-username>/<repo-name>.git
cd <repo-name>
```
Install Go if not already installed:

```bash
brew install go
```
---

## 2️⃣ Start Backend

```bash
cd wallet-api
go mod tidy
go run main.go
```

Backend runs at:

```
http://localhost:8080
```

---

## 3️⃣ Start Frontend

Install Hugo if not already installed:

```bash
brew install hugo
```

Run:

```bash
cd ../wallet-frontend
hugo server
```

Frontend runs at:

```
http://localhost:1313
```

⚠ Ensure backend is running before starting the frontend.

---

# 🧪 Example Usage

Generate mnemonic:

```bash
curl -X POST http://localhost:8080/mnemonic
```

Derive wallet:

```bash
curl "http://localhost:8080/derive?mnemonic=<mnemonic>&index=0"
```

Health check:

```bash
curl http://localhost:8080/health
```

---

# 🔐 Security Considerations

- Mnemonic validation using `bip39.IsMnemonicValid`
- Deterministic derivation without persistent storage
- No private key storage on server
- Query parameters used for simplicity (production systems should use JSON body)
- No logging of sensitive values

---


# 🎯 Engineering Goals

This project explores:

- Deterministic cryptographic wallet generation
- Backend API design patterns
- Concurrency-safe service development
- Infrastructure-aware system design
- Clean separation between service and UI

