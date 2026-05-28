# Openway

A production-grade, non-custodial cross-border B2B settlement protocol. Built to demonstrate DevSecOps practices in Web3 infrastructure, from secure smart contract development to hardened Kubernetes deployment. Every layer is secured, every decision is explained, and every tool is real.

[CI Pipeline] [Go Report Card] [License: MIT]

## What This Is

This is a reference implementation for decentralized financial infrastructure, not a tutorial. It addresses the systemic failures of centralized African fintech clearinghouses (which collapsed due to pre-funded fiat treasury risks and opaque ledgers) by replacing them with an asset-light, trustless routing layer.

It demonstrates:

- **Secure Web3 Development:** Foundry fuzzing, Slither static analysis, EIP-712 signature verification.
- **Non-Custodial Architecture:** Smart contract escrows eliminating counterparty and treasury risk.
- **High-Throughput Ingestion:** Go-based webhook listeners handling Telco API callbacks (Daraja, MTN) asynchronously.
- **Infrastructure:** Hardened Distroless Docker images, Kubernetes deployment with Vault secrets injection.
- **Observability:** Structured JSON logging, Prometheus metrics, RPC node health monitoring.
- **CI/CD:** SAST, SCA, smart contract auditing, container scanning, artifact signing.

Built for the Infrastructure & Tooling / DeFi & Finance track.

## Architecture

```bash
┌─────────────┐     ┌─────────────────┐     ┌─────────────┐
│  Telco APIs │────▶│  Relayer / API  │────▶│  PostgreSQL │
│ (Daraja/MTN)│     │   (Go :8080)    │     │             │
└─────────────┘     └─────────────────┘     └─────────────┘
                                 │
                                 │ (EIP-712 Signatures)
                                 ▼
                    ┌─────────────────────────┐
                    │   Settlement Engine     │
                    │ (Celo / Smart Contracts)│
                    └─────────────────────────┘
                                 │
                                 ▼
┌─────────────┐     ┌─────────────────┐
│  Admin Web  │◀───▶│   Web3 RPCs     │
│  (Next.js)  │     │ (Alchemy/Ankr)  │
└─────────────┘     └─────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 20+ & pnpm
- Foundry (Forge, Cast, Anvil)
- Docker & Docker Compose
- kubectl & a local Kubernetes cluster (k3d, minikube, or kind)
- HashiCorp Vault CLI

## Local Development

```bash
# Clone and enter
git clone https://github.com/philaturo/openway.git
cd openway

# Start dependencies (Postgres, Redis, Local Celo Node)
docker-compose up -d

# Deploy local smart contracts
cd contracts && forge script script/Deploy.s.sol --broadcast --rpc-url http://localhost:8545
cd ..

# Run backend tests
cd backend && go test -race ./...
cd ..

# Start backend services
cd backend
go run ./cmd/relayer
go run ./cmd/api
cd ..

# Start frontend
cd frontend && pnpm install && pnpm dev
```

## System Endpoints

| Method   | Endpoint                  | Auth | Description                                     |
| :------- | :------------------------ | :--- | :---------------------------------------------- |
| **POST** | `/webhooks/daraja`        | HMAC | Ingest M-Pesa B2C/B2B transaction callbacks     |
| **POST** | `/webhooks/mtn`           | HMAC | Ingest MTN MoMo collection callbacks            |
| **POST** | `/api/v1/escrow/initiate` | JWT  | Lock stablecoins and generate settlement intent |
| **GET**  | `/api/v1/settlements`     | JWT  | List cross-border routing transactions          |
| **GET**  | `/api/v1/health/rpc`      | No   | Check Web3 RPC fallback array latency           |

## Security Features

| Layer               | Control                | Implementation                                                              |
| :------------------ | :--------------------- | :-------------------------------------------------------------------------- |
| **Smart Contracts** | Reentrancy prevention  | OpenZeppelin ReentrancyGuard, checks-effects-interactions pattern           |
| **Smart Contracts** | Static analysis        | Slither, Mythril                                                            |
| **Smart Contracts** | Fuzz testing           | Foundry invariant testing and stateful fuzzing                              |
| **Development**     | Secret scanning        | GitLeaks pre-commit hooks                                                   |
| **Development**     | Commit verification    | GPG signing required                                                        |
| **Relayer API**     | Webhook validation     | Telco HMAC signature verification                                           |
| **Relayer API**     | IDOR prevention        | Strict ownership enforcement on settlement queries                          |
| **Relayer API**     | Idempotency            | Deterministic hashing of Telco transaction IDs to prevent double-processing |
| **Container**       | Minimal attack surface | gcr.io/distroless/static-debian12 base image, no shell                      |
| **Container**       | Non-root execution     | UID 65532                                                                   |
| **Kubernetes**      | Secret management      | Vault agent sidecar injection for RPC keys and DB credentials               |
| **Kubernetes**      | Network isolation      | Default-deny network policies, explicit egress to RPC providers             |
| **CI/CD**           | Static analysis        | Semgrep, go vet, staticcheck                                                |
| **CI/CD**           | Dependency scanning    | govulncheck, Trivy (CRITICAL/HIGH fail build)                               |
| **CI/CD**           | Artifact integrity     | Cosign keyless signing for Docker images                                    |

### Project Directory Structure

```bash
openway/
├── .github/
│   └── workflows/
│       ├── ci.yml                  # Build, test, lint
│       ├── security-scan.yml       # SAST, SCA, Secret Scanning
│       └── release.yml             # Docker build, Cosign, K8s deploy
├── .vscode/
│   ├── settings.json               # Editor formatting and validation
│   └── extensions.json             # Recommended DevSecOps extensions
├── docs/
│   ├── architecture.md             # System design and data flow
│   ├── threat-model.md             # STRIDE analysis and mitigations
│   └── api-spec.yaml               # OpenAPI specification for Go backend
│
├── backend/                        # Go Backend Services
│   ├── cmd/
│   │   ├── api/main.go             # B2B Admin API entrypoint
│   │   ├── relayer/main.go         # Telco Webhook Ingestion entrypoint
│   │   └── worker/main.go          # Async Queue Processor entrypoint
│   ├── internal/
│   │   ├── api/                    # HTTP handlers, routing, response formatting
│   │   ├── auth/                   # JWT generation, API key validation
│   │   ├── telco/                  # Daraja/MTN MoMo HMAC & webhook logic
│   │   ├── web3/                   # RPC fallback, EIP-712 signing, ABI bindings
│   │   ├── queue/                  # Redis/Asynq task definitions & workers
│   │   ├── repository/             # PostgreSQL queries, idempotency checks
│   │   ├── service/                # Core business logic & settlement orchestration
│   │   └── telemetry/              # OpenTelemetry, Prometheus, structured logging
│   ├── pkg/
│   │   ├── crypto/                 # EIP-712 helpers, hashing utilities
│   │   └── errors/                 # Standardized domain error handling
│   ├── migrations/                 # PostgreSQL schema migrations (Goose/Migrate)
│   ├── go.mod
│   ├── go.sum
│   └── .golangci.yml               # Go linter configuration
│
├── contracts/                      # Solidity Smart Contracts (Foundry)
│   ├── src/
│   │   ├── EscrowRouter.sol        # Core settlement & EIP-712 verification
│   │   ├── interfaces/IEscrow.sol  # Contract interfaces
│   │   └── lib/EIP712.sol          # Signature verification library
│   ├── test/
│   │   ├── EscrowRouter.t.sol      # Unit and integration tests
│   │   └── invariants/             # Stateful fuzzing and invariant testing
│   ├── script/
│   │   ├── Deploy.s.sol            # Deployment automation
│   │   └── Interact.s.sol          # Local testing interactions
│   ├── foundry.toml                # Foundry configuration
│   └── remappings.txt              # Dependency mapping
│
├── frontend/                       # Next.js Frontend (App Router)
│   ├── app/
│   │   ├── layout.tsx              # Root layout, global providers
│   │   ├── page.tsx                # Landing/Role-selection redirect
│   │   ├── middleware.ts           # Role-based route protection (Admin vs Merchant)
│   │   ├── (admin)/                # Desktop B2B Admin Dashboard Route Group
│   │   │   ├── dashboard/page.tsx  # Air Traffic Control overview
│   │   │   ├── escrows/page.tsx    # Live escrow monitoring
│   │   │   ├── nodes/page.tsx      # RPC & Relayer health matrix
│   │   │   └── settings/page.tsx   # API keys, webhook configs
│   │   ├── (merchant)/             # Mobile PWA Merchant Portal Route Group
│   │   │   ├── portal/page.tsx     # Thumb-zone UI, balance overview
│   │   │   ├── scan/page.tsx       # QR Code scanner for settlements
│   │   │   └── history/page.tsx    # Recent payouts and statuses
│   │   └── (auth)/
│   │       └── login/page.tsx      # Authentication gateway
│   ├── components/
│   │   ├── ui/                     # Base design system (Shadcn/Radix)
│   │   ├── admin/                  # LiveTransactionStream, NodeHealthMatrix
│   │   ├── merchant/               # ScanQRButton, BalanceCard, StatusBadge
│   │   └── shared/                 # Header, Sidebar, MobileNav
│   ├── lib/
│   │   ├── web3/                   # Viem clients, contract ABIs, hooks
│   │   ├── api/                    # Axios/Fetch clients for Go backend
│   │   └── auth/                   # Session management utilities
│   ├── hooks/
│   │   └── useSettlementStream.ts  # WebSocket hooks for live updates
│   ├── public/
│   │   └── manifest.json           # PWA configuration for mobile merchants
│   ├── styles/
│   │   └── globals.css             # Design tokens (Navy, Emerald, Amber)
│   ├── next.config.js
│   ├── tailwind.config.ts
│   └── tsconfig.json
│
├── infra/                          # Infrastructure & Deployment
│   ├── docker/
│   │   ├── Dockerfile.api          # Hardened Distroless Go API
│   │   ├── Dockerfile.relayer      # Hardened Distroless Go Relayer
│   │   └── Dockerfile.frontend     # Hardened Node Next.js
│   ├── k8s/
│   │   ├── base/                   # Kustomize base manifests
│   │   └── overlays/prod/          # Prod patches, Vault injections
│   └── terraform/                  # Cloud infrastructure (Optional)
│
├── security-tests/                 # DAST & Auditing Configs
│   ├── zap/                        # OWASP ZAP API scan configs
│   └── slither/                    # Slither static analysis configs
│
├── .gitleaks.toml                  # Secret scanning rules
├── .trivyignore                    # Container scanning exceptions
├── .pre-commit-config.yaml         # Pre-commit hooks (Gitleaks, fmt)
├── docker-compose.yml              # Local dev (Postgres, Redis, Anvil)
├── Makefile                        # Dev, test, build, deploy commands
├── .gitignore
└── README.md
```

### DevSecOps Integration

This repository demonstrates the Modern DevSecOps Engineering Stack applied to Web3 Infrastructure:

- **Development Foundations:** Commit signing, pre-commit hooks, and secure smart contract review.
- **CI/CD Pipeline Security:** Least-privilege GitHub Actions, artifact signing, and environment isolation.
- **Static Analysis (SAST):** Semgrep for Go, Slither for Solidity, and custom rule writing.
- **Dynamic Analysis (DAST):** OWASP ZAP against the Relayer API and Foundry invariant testing.
- **Container & Kubernetes Security:** Distroless images, Vault injection, and strict network policies.
- **Observability & Incident Response:** Prometheus metrics, RPC fallback alerting, and runbooks.

---

### Contributing

This is an infrastructure project. Issues and Pull Requests are welcome, especially regarding:

- **Smart Contract Optimizations:** Gas optimizations and security patches.
- **Telco Integrations:** Additional Telco API sandbox integrations.
- **Documentation:** Clarifications and structural enhancements.
- **Kubernetes Deployments:** Helm chart implementations.

> ⚠️ **Note:** Please sign your commits (`git commit -S`) and ensure all local pre-commit hooks pass before pushing code.

---

### License

Distributed under the **MIT License**. See `LICENSE` for details.

---

### Connect

Built with **intent**. Secured by **design**.
