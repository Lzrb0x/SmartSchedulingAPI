# SmartSchedulingAPI MVP Plan

## Overview
Este documento descreve a arquitetura inicial, modelo de dados e roadmap para o MVP da SmartSchedulingAPI, garantindo base sólida para evolução futura em SaaS multi-tenant.

## Arquitetura
- **Stack:** Go 1.25+, Gin para HTTP, sqlx + pgx para Postgres, envconfig para configuração.
- **Camadas:** `cmd/api` (bootstrap), `internal/config`, `internal/database`, `internal/server` (routes/middlewares), `internal/domain` (entidades/regra), `internal/auth`, `internal/tenant`.
- **Tenancy:** coluna `tenant_id` em todas as tabelas com tabela `tenants` dedicada. Middleware JWT injeta tenant no contexto. Futuro: isolamento por schema/RLS.
- **Autenticação:** JWT com roles (`owner`, `employee`, `client`), middlewares para autorização.

## Modelo de Dados (resumo)
- `tenants(id, name, status)`
- `users(tenant_id, role, email uniq por tenant)`
- `barbershops(tenant_id, características)`
- `services(tenant_id, barbershop_id, nome uniq por barbearia, preço em centavos, duração)`
- `professionals(usuario funcionário/dono vinculado)`
- `agendas` (blocos de trabalho/folga por profissional)
- `bookings` (agendamentos com checagem de conflito via índice único + validação transacional)

## Endpoints Prioritários
1. `/auth/login`, `/auth/register`
2. `/barbearias`, `/servicos`, `/profissionais`
3. `/agendas` para horários/folgas
4. `/agendamentos` CRUD + consultas por cliente/profissional
5. `/health` & `/api/tenants/current` para verificação

## Roadmap de Implementação
1. Configuração de projeto (Go module, dependências, config/env, Docker futuro).
2. Migrações com goose: `0001_init.sql` já define schema base.
3. Infra de auth (hash de senha, JWT, middleware tenant).
4. CRUDs principais + camada de serviços (validar regras de negócio).
5. Regras anti-conflito no `BookingService` (queries com `FOR UPDATE`).
6. Testes unitários/integrados.
7. Documentação (OpenAPI) + plano de escalabilidade (onboarding tenant, billing, observabilidade).

## Próximos Passos
- Implementar lógica real de usuários/tenants no handler de auth.
- Criar services/repositories utilizando sqlx e estruturar camadas.
- Expandir documentação com OpenAPI inicial e diagramas.
