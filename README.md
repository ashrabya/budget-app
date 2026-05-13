# 💰 Cashflow — Budget Tracker

A full-stack budgeting app built with **Go** (stdlib only, no frameworks) and a vibrant modern HTML/JS UI, served from a single binary.

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

- **Dashboard** — Visual overview with income/expense donut charts and balance stats
- **Transactions** — Add, view, and delete income & expense entries with categories
- **Budgets** — Set monthly spending limits per category with progress bars
- **Persistence** — Data saved as a JSON file (no external DB required)
- **Single binary** — UI is embedded into the Go binary via `//go:embed`
- **Month filtering** — View data for any month

## 🗂️ Project Structure

```
budget-app/
├── cmd/server/           # Application entrypoint
│   └── main.go
├── internal/
│   ├── handlers/         # HTTP API handlers
│   ├── models/           # Domain types (Transaction, Budget, Summary)
│   └── storage/          # In-memory store with JSON file persistence
├── web/
│   └── index.html        # Embedded single-page UI
├── .github/workflows/    # GitHub Actions CI
├── Dockerfile
├── Makefile
└── go.mod
```

## 🚀 Quick Start

### Prerequisites
- Go 1.22+

### Run locally

```bash
git clone https://github.com/ashrabya/budget-app.git
cd budget-app
make run
# → http://localhost:8080
```

### Build binary

```bash
make build
./bin/budget-app
```

### Docker

```bash
make docker-build
make docker-run
# → http://localhost:8080
```

## ⚙️ Configuration

| Env Variable | Default              | Description                    |
|-------------|----------------------|--------------------------------|
| `PORT`      | `8080`               | HTTP port to listen on         |
| `DATA_FILE` | `budget-data.json`   | Path to JSON persistence file  |

```bash
PORT=9000 DATA_FILE=/var/data/budget.json ./bin/budget-app
```

## 🌐 API Reference

### Transactions

| Method | Path                        | Description          |
|--------|-----------------------------|----------------------|
| GET    | `/api/transactions`         | List all             |
| POST   | `/api/transactions`         | Create new           |
| GET    | `/api/transactions/:id`     | Get by ID            |
| PUT    | `/api/transactions/:id`     | Update               |
| DELETE | `/api/transactions/:id`     | Delete               |

**POST body example:**
```json
{
  "type": "expense",
  "amount": 42.50,
  "category": "Food & Dining",
  "description": "Lunch at cafe",
  "date": "2026-05-12T12:00:00Z"
}
```

### Budgets

| Method | Path              | Description |
|--------|-------------------|-------------|
| GET    | `/api/budgets`    | List all    |
| POST   | `/api/budgets`    | Create new  |
| DELETE | `/api/budgets/:id`| Delete      |

### Summary

| Method | Path                           | Description           |
|--------|--------------------------------|-----------------------|
| GET    | `/api/summary?month=2026-05`   | Monthly summary stats |

## 🚢 Deploy to Fly.io

```bash
fly launch
fly deploy
```

## 🔧 Development

```bash
make run      # Run with hot-reload via air (install separately)
make test     # Run tests
make build    # Compile binary
make clean    # Remove build artifacts
```

## 📄 License

MIT
