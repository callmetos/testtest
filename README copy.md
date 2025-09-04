# NavMate-Backend (template)

This is a minimal starter using **Gin + GORM + zap**, with `/health` and `/api/health`,
`docker-compose` (Postgres + MinIO), and GitHub Actions CI (`go test`).

## Quick start
```bash
cp .env.example .env
docker compose up -d postgres minio
go run ./cmd/main.go
curl -i http://localhost:8080/health
go test ./...
```
