# Projeto Rate Limiter

## Visão Geral
Este projeto implementa um limitador de taxa usando Go e Redis. Ele inclui middleware para limitar o número de solicitações por segundo a uma API.

## Pré-requisitos
- Go (versão 1.16 ou superior)
- Docker (para executar o Redis)
- Redis (se não estiver usando Docker)

## Configuração

### 1. Clone o Repositório
```
git clone <URL_DO_REPOSITORIO>
cd <NOME_DO_REPOSITORIO>
```

### 2. Instale as Dependências
```
go mod tidy
```

### 3. Configure o Redis
- **Usando Docker:**
  ```
  docker run --name redis -p 6379:6379 -d redis
  ```

- **Sem Docker:**
  Certifique-se de que o Redis está instalado e em execução em `localhost:6379`.

### 4. Defina as Variáveis de Ambiente
Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:
```
REQUESTS_PER_SECOND=5
TIME_WINDOW=1s
BLOCK_DURATION=10s
REDIS_ADDR=localhost:6379
```

## Executando o Projeto

### 1. Inicie o Servidor
```
go run main.go
```

### 2. Teste o Rate Limiter
Use o arquivo `ratelimit.http` para enviar solicitações de teste. Você pode usar a ferramenta HTTP Client no GoLand para isso.

## Executando Testes
Para executar os testes automatizados, use o seguinte comando:
```
go test ./...
```

## Estrutura do Projeto
- `main.go`: Ponto de entrada da aplicação.
- `app/limiter/storage.go`: Interface de armazenamento.
- `app/limiter/redis_storage.go`: Implementação do armazenamento Redis.
- `app/middleware/ratelimit_test.go`: Testes automatizados para o limitador de taxa.
- `ratelimit.http`: Solicitações HTTP para testes manuais.

## Comandos Úteis

- **Parar o Container Redis:**
  ```
  docker stop redis
  ```

- **Remover o Container Redis:**
  ```
  docker rm redis
  ```

### Postman

- ** GET sem token **
  ```
  curl --location 'http://localhost:8080/full'
  ```

- ** POST novo token **
  ```
  curl --location 'http://localhost:8080/generate-token' \
  --header 'Content-Type: application/json' \
  --data '{
  "requests_per_second": 3
  }'
  ```

- ** GET com token **
  Substituir {{token}} pelo token gerado no passo anterior

```
curl --location 'http://localhost:8080/full' \
--header 'API_KEY: bearer {{token}}'
```

## Licença
Este projeto é licenciado sob a Licença MIT.