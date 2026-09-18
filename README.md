# GoSocial

API REST para fórum desenvolvida em Go, com arquitetura baseada em hexagonal, autenticação via PASETO, cache Redis e containerizado com Docker.

---

## Stack

| Camada | Tecnologia |
|--------|------------|
| Linguagem | Go 1.27 |
| Roteamento | chi |
| Banco de Dados | PostgreSQL 17 |
| Cache | Redis 7 |
| Autenticação | PASETO v2 (local) |
| Hash de Senha | Argon2id |
| Logging | zap |
| Documentação | Swagger |
| Containerização | Docker + Docker Compose |
| CI/CD | GitHub Actions |

---

## Arquitetura

O projeto segue os princípios da arquitetura hexagonal (Ports & Adapters), com separação clara entre domínio, aplicação e adaptadores.

```
GoSocial/
├── cmd/
│   ├── api/                    # Ponto de entrada da API
│   └── migrate/                # Migrações do banco
├── internal/
│   ├── auth/                   # Autenticação (PASETO, Argon2id)
│   ├── cache/                  # Camada de cache (Redis)
│   ├── db/                     # Conexão com banco
│   ├── env/                    # Configuração via variáveis de ambiente
│   ├── mailer/                 # Envio de emails
│   └── store/                  # Repositórios (PostgreSQL, mocks)
├── docs/                       # Documentação Swagger
├── scripts/                    # Scripts auxiliares
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```

### Fluxo de Dependência

```
Handler HTTP → Service → Repository (interface) → PostgreSQL
                              ↓
                          Cache (Redis)
```

O domínio define as interfaces. Os adaptadores implementam. A injeção é feita manualmente no `main.go`.

---

## Funcionalidades

- CRUD de usuários
- CRUD de posts
- Sistema de seguidores (follow/unfollow)
- Feed personalizado
- Comentários
- Ativação de conta por token
- Autenticação via PASETO
- Cache de usuários com Redis
- Autorização baseada em papéis (RBAC)
- Rate limiting
- Graceful shutdown
- Métricas com expvar
- Migrações versionadas

---

## Segurança

| Recurso | Implementação |
|---------|---------------|
| Hash de senha | Argon2id (64 MiB, 3 iterações, 2 threads, salt de 16 bytes) |
| Token de autenticação | PASETO v2.local (XChaCha20-Poly1305) |
| Token de convite | SHA-256 sobre UUID v4 |
| Comparação de senha | Constant-time (subtle.ConstantTimeCompare) |
| Autorização | RBAC com níveis (user, moderator, admin) |
| Headers de autenticação | Bearer token |

---

## Pré-requisitos

- Go 1.27 ou superior
- Docker e Docker Compose
- Make (opcional)

---

## Como executar

### Com Docker Compose

```bash
docker compose up --build
```

A API estará disponível em `http://localhost:8080`.

### Localmente

1. Suba a aplicação:

```bash
docker compose up -d 
```

2. Configure as variáveis de ambiente:

```bash
cp .env.example .env
# edite o .env com suas configurações
```

3. Execute as migrações:

```bash
make migrateup
```

4. Inicie a API:

```bash
make run
```

---

## Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `ADDR` | Endereço da API | `:8080` |
| `DB_ADDR` | String de conexão do PostgreSQL | — |
| `DB_MAX_OPEN_CONN` | Máximo de conexões abertas | `30` |
| `DB_MAX_IDLE_CONN` | Máximo de conexões ociosas | `30` |
| `DB_MAX_IDLE_TIME` | Tempo máximo de conexão ociosa | `15m` |
| `REDIS_ADDR` | Endereço do Redis | — |
| `REDIS_PASSWORD` | Senha do Redis | — |
| `REDIS_DB` | Índice do banco Redis | `0` |
| `REDIS_ENABLED` | Habilita o cache | `false` |
| `AUTH_TOKEN_SECRET` | Chave simétrica do PASETO (base64) | — |
| `AUTH_TOKEN_EXP` | Expiração do token | `15m` |
| `AUTH_BASIC_USER` | Usuário do Basic Auth (health check) | — |
| `AUTH_BASIC_PASS` | Senha do Basic Auth | — |
| `SENDGRID_API_KEY` | Chave da API SendGrid | — |
| `MAIL_FROM_EMAIL` | Email remetente | — |
| `ENV` | Ambiente (`development`, `production`) | `development` |
| `API_URL` | URL pública da API | — |
| `FRONTEND_URL` | URL do frontend | — |

---

## Geração da Chave PASETO

A chave simétrica deve ter 32 bytes. Gere uma com:

```bash
openssl rand -base64 32
```

Cole o resultado em `AUTH_TOKEN_SECRET`.

---

## Endpoints

### Health

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/v1/health` | Verifica saúde da API (Basic Auth) |
| GET | `/v1/metrics` | Verifica as métricas do backend |
| GET | `/v1/swagger/index.html` | Acessa a documentação no Swagger |

### Autenticação

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/v1/authentication/user` | Registra novo usuário |
| POST | `/v1/authentication/token` | Gera token de acesso |

### Usuários

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/v1/users/{userID}` | Busca usuário por ID |
| PUT | `/v1/users/{userID}` | Atualiza usuário |
| DELETE | `/v1/users/{userID}` | Remove usuário |
| PUT | `/v1/users/{userID}/follow` | Segue usuário |
| PUT | `/v1/users/{userID}/unfollow` | Deixa de seguir |
| PUT | `/v1/users/activate/{token}` | Ativa conta por token |
| GET | `/v1/users/feed` | Feed personalizado |

### Posts

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/v1/posts` | Cria post |
| GET | `/v1/posts/{postID}` | Busca post |
| PUT | `/v1/posts/{postID}` | Atualiza post |
| DELETE | `/v1/posts/{postID}` | Remove post |

Documentação completa disponível em `/swagger/index.html`.

---

## Testes

```bash
# Todos os testes
make test

# Com cobertura
make test-coverage

# Apenas um pacote
go test ./internal/auth/ -v
```

### Cobertura Atual

O projeto utiliza mocks para testar handlers e serviços sem dependência de banco ou cache.

---

## Benchmark

Teste de carga realizado com `autocannon` (1000 conexões, 20 segundos):

| Métrica | Com Redis | Sem Redis |
|---------|-----------|-----------|
| Requisições/segundo | 36.232 | 21.234 |
| Latência média | 27,46 ms | 47,70 ms |
| Latência p99 | 56 ms | 100 ms |
| Throughput | 11,7 MB/s | 7,79 MB/s |

O cache Redis proporciona ganho de aproximadamente 70% em throughput e 42% de redução na latência.

---

## Makefile

| Comando | Descrição |
|---------|-----------|
| `make migrate (nome-migration) ` | Gera um novo arquivo de migração
| `make test` | Executa os testes |
| `make migrateup` | Aplica migrações |
| `make migratedown` | Reverte migrações |
| `make seed` | Popula o banco com dados de teste |
| `make gen-docs` | Gera a documentação Swagger |

---

## CI/CD

O projeto utiliza GitHub Actions para:

- Verificação de dependências (`go mod verify`)
- Análise estática (`go vet`, `staticcheck`)
- Execução de testes (`go test -race`)
- Build da aplicação

---

## Decisões de Projeto

### Por que PASETO em vez de JWT

PASETO elimina as brechas de configuração do JWT (como `alg: none` e confusão de algoritmos). A segurança é garantida por padrão, sem necessidade de flags adicionais.

### Por que Argon2id em vez de bcrypt ou scrypt

Argon2id é o vencedor do Password Hashing Competition (2015), com resistência superior a ataques de GPU e ASIC devido ao uso intensivo de memória.

### Por que cache Redis

Com carga alta (1000 conexões simultâneas), o Redis reduz a carga no PostgreSQL e proporciona ganho significativo de throughput e latência.

---

## Próximos Passos

- Testes de integração com banco real (testcontainers)
- Implementação de Autenticação oauth2
- Criação de frontend em framework JS/TS
- Refatoração geral do sistema
- Suporte a HTTPS

---

## Licença

Este projeto está licenciado sob a MIT License. Consulte o arquivo LICENSE para mais detalhes.

---

## Contato

Alvaro Lucio

- GitHub: [github.com/alvarolucio2007](https://github.com/alvarolucio2007)
- Repositório: [github.com/alvarolucio2007/GoSocial](https://github.com/alvarolucio2007/GoSocial)
