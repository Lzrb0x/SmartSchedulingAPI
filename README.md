# SmartSchedulingAPI

MVP de uma API Go (Gin + sqlx) pensada para evoluir em um SaaS multi-tenant destinado a barbearias. Inclui configuração básica, migração inicial e documentação do plano de implementação.

## Estrutura
- `cmd/api`: ponto de entrada HTTP.
- `internal/config`: carregamento via envconfig.
- `internal/database`: conexão Postgres com sqlx/pgx.
- `internal/server`: inicialização Gin e rotas básicas.
- `internal/auth`: handlers e middleware JWT (stubs).
- `internal/tenant`: helpers para contexto de tenant.
- `internal/domain`: entidades de domínio.
- `migrations`: goose migrations (schema base multi-tenant).
- `docs`: plano do MVP e roadmap.
- `tools/goose.conf`: configuração para executar migrations (usa `DATABASE_URL`).

## Como Rodar
1. Defina `APP_DB_URL` ou use default (`postgres://postgres:postgres@localhost:5432/smartscheduling?sslmode=disable`).
2. Rode migrations: `DATABASE_URL=$APP_DB_URL goose -dir migrations up`.
3. Inicie a API: `go run ./cmd/api`.

Endpoints disponíveis (MVP):
- `GET /health`
- `POST /api/auth/login`
- `POST /api/auth/register`
- `GET /api/tenants/current` (requer JWT; placeholder até auth completo).

Confira `docs/README.md` para detalhes do plano, regras de negócio e roadmap para evoluir em SaaS multi-tenant.
