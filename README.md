# ARCHIVAL VAULT // MONOSPACED RAW LEDGER & REGISTRY

```text
┌────────────────────────────────────────────────────────────────────────┐
│ ARCHIVAL VAULT :: RAW MONOSPACED WIREFRAME FULLSTACK ARCHITECTURE       │
│ Golang RESTful API + JWT Auth + Vue 3 TypeScript + TanStack Query      │
└────────────────────────────────────────────────────────────────────────┘
```

Minimalist, high-assurance technical record-keeping system and artifact vault. Designed with a **Raw Archival & Monospaced Wireframe** aesthetic (zero border-radius, carbon slate palette, classification stamp badges, live system telemetry).

Built to demonstrate all core competencies from **"FULLSTACK DEVELOPER GOLANG VUE"** (Fika Ridaul Maulayya):
- **RESTful API**: Standard HTTP routing in Go (`net/http.ServeMux` Go 1.22+ patterns), JSON marshalling, clean response wrappers.
- **JWT**: Secure JSON Web Token auth (`golang-jwt/jwt/v5`) with HS256 claims, bearer token parser, and protected middleware.
- **Vue 3**: Reactive Composition API with `<script setup lang="ts">`.
- **TypeScript**: Strict interfaces for users, records, payloads, telemetry, and API responses.
- **Auth**: End-to-end registration, login, profile fetch, session persistence, and auto-logout on token expiration.
- **Composables**: Custom composables (`useAuth`, `useRecords`, `useToast`).
- **TanStack Query**: `@tanstack/vue-query` (`useQuery`, `useMutation`, `useQueryClient`) for smart caching, 30s stale time, background refetch, and automatic query cache invalidation.

---

## Architecture & Integration

- **Backend (Golang 1.27)**: High-performance RESTful API listening on `http://127.0.0.1:8085`.
- **Database (FlyEnv MySQL 8.x)**: Port `3306`, database `archival_vault`. Auto-migrates `users` and `records` tables and auto-seeds initial data.
- **Frontend (Vue 3 + Vite + TypeScript)**: Built into `public/`, served directly through **Laravel Herd** at `http://archival-vault.test`.
- **Laravel Herd Gateway (`public/index.php`)**: Transparently proxies `/api/*` to the Go backend on port 8085 and serves the SPA for all client routes.

---

## Quick Start

### 1. Requirements
- MySQL (FlyEnv default on `127.0.0.1:3306`, root:root)
- Laravel Herd (parked at `C:\Users\Lenovo\Herd`)
- Go 1.22+ & Bun / Node.js

### 2. Run Backend
```bash
cd backend
go run main.go
# Or run precompiled binary:
./vault-server.exe
```

### 3. Open in Browser
- **Laravel Herd Site**: [http://archival-vault.test](http://archival-vault.test)
- **Direct Backend & Embedded UI**: [http://127.0.0.1:8085](http://127.0.0.1:8085)
- **Vite Dev Server (Frontend HMR)**:
  ```bash
  cd frontend
  bun run dev
  # -> http://localhost:5173
  ```

---

## Default Seed Credentials

| Identifier | Value |
|---|---|
| **Username** | `archivist` |
| **Password** | `password123` |
| **Role** | `ARCHIVIST` |
| **Database** | `archival_vault` (MySQL port 3306) |

*(A 1-click `[AUTO-FILL DEMO ARCHIVIST]` button is provided inside the Sign-in modal for instant access).*

---

## RESTful API Endpoints

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/api/health` | Public | Subsystem health status |
| `POST` | `/api/auth/register` | Public | Create operator identity and issue JWT |
| `POST` | `/api/auth/login` | Public | Verify credentials and return JWT |
| `GET` | `/api/auth/me` | Bearer JWT | Fetch current authenticated user claims |
| `GET` | `/api/records` | Public | List records with search, category, and classification filters |
| `GET` | `/api/records/{id}` | Public | Retrieve single archival record |
| `POST` | `/api/records` | Bearer JWT | Ingest and commit new record to ledger |
| `PUT` | `/api/records/{id}` | Bearer JWT | Patch existing archival record |
| `DELETE` | `/api/records/{id}` | Bearer JWT | Expunge record from registry |
| `GET` | `/api/stats` | Public | Realtime classification and telemetry counters |

---

## UI Aesthetic: Archival Raw & Monospaced Wireframe

- Strict monospace typeface (`JetBrains Mono`, `Courier New`, monospace).
- Carbon slate dark palette (`#090a0c`, `#111317`, `#262c35`, `#e6edf3`).
- 0px border radius across all buttons, inputs, modals, and container boxes.
- Classification stamp badges: `[CONFIDENTIAL]`, `[RESTRICTED]`, `[INTERNAL]`, `[PUBLIC]`.
- Live collapsible raw payload inspector with 1-click clipboard copy.
- Terminal-styled action prompts (`[+ COMMIT RECORD]`, `[VIEW RAW]`, `[CURL]`, `[SYNC]`).

---

## License

MIT © DevisRoomi
