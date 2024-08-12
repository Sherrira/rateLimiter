# Rate Limiter

Este projeto é uma implementação de um rate limiter de requisições HTTP utilizando Redis.

## Pré-requisitos

Para executar este projeto, você precisará ter instalado em sua máquina:

- Docker
- Docker Compose
- Make

## Variáveis de Ambiente

As seguintes variáveis de ambiente podem ser configuradas para ajustar o comportamento do rate limiter, configuradas via `docker-compose.yml` ou no arquivo `.env` (execução local):

- `IP_RATE_LIMIT`: Número máximo de requisições permitidas por IP (padrão: `10`)
- `TOKEN_LIMITS=`: Número máximo de requisições permitidas por token (exemplo: `"abc123:100,def456:20"`)
- `BLOCK_DURATION`: Duração do bloqueio após exceder o limite (padrão: `1m`)
- `REDIS_ADDR`: Endereço do servidor Redis (padrão: `localhost:6379`)
- `REDIS_PASSWORD`: Senha do Redis (padrão: vazio)
- `REDIS_DB`: Banco de dados do Redis a ser utilizado (padrão: `0`)

## Comandos do Makefile

O `Makefile` fornece comandos para facilitar a manipulação dos containers Docker. Aqui estão os comandos disponíveis:

- **Buildar todos os containers:**
```sh
$ make build
```

- **Subir todos os containers:**
```sh
$ make up
```

- **Buildar e subir todos os containers:**
```sh
$ make up-build
```

- **Parar todos os containers:**
```sh
$ make down
```

- **Subir um serviço específico:**
```sh
$ make up-<service_name>
```

- **Buildar e subir um serviço específico:**
```sh
$ make up-build-<service_name>
```

- **Parar um serviço específico:**
```sh
$ make down-<service_name>
```

Substitua `<service_name>` pelo nome do serviço que você deseja manipular, como `redis` ou `app`.

## Testando a Implementação

Você pode testar este projeto de duas maneiras:

### 1) Teste automatizado

O arquivo `cmd/rateLimiter_test.go` possui um teste automatizado que executa requisições e testa os retornos de acordo com as configurações default parametrizadas em `cmd/.env`. Basta executar o teste com o ambiente docker no ar. Para isso, execute na raiz do projeto o comando:

```sh
$ make up-build
```

Em seguida, execute o teste automatizado:
```sh
$ go test -v ./cmd/rateLimiter_test.go 
```
OBS: Caso você execute o teste novamente, ele trará um erro, justamente validando os valores esperados, confirmando que o middleware bloqueou as requisições do teste. 
Saída:
```sh
$ go test -v ./cmd/rateLimiter_test.go 
=== RUN   TestRateLimiter
    rateLimiter_test.go:31: 
                Error Trace:    ~/rateLimiter/cmd/rateLimiter_test.go:31
                Error:          Not equal: 
                                expected: 100
                                actual  : 0
                Test:           TestRateLimiter
    rateLimiter_test.go:32: 
                Error Trace:    ~/rateLimiter/cmd/rateLimiter_test.go:32
                Error:          Not equal: 
                                expected: 1
                                actual  : 101
                Test:           TestRateLimiter
    rateLimiter_test.go:50: 
                Error Trace:    ~/rateLimiter/cmd/rateLimiter_test.go:50
                Error:          Not equal: 
                                expected: 10
                                actual  : 0
                Test:           TestRateLimiter
    rateLimiter_test.go:51: 
                Error Trace:    ~/rateLimiter/cmd/rateLimiter_test.go:51
                Error:          Not equal: 
                                expected: 1
                                actual  : 11
                Test:           TestRateLimiter
--- FAIL: TestRateLimiter (0.09s)
FAIL
FAIL    command-line-arguments  0.093s
FAIL
```
Dessa forma, basta esperar 1 minuto (ou o tempo configurado) para que as requisições possam ser aceitas novamente.

### 2) Teste de carga

Para testar a implementação, você pode usar uma ferramenta de teste de carga. Recomendamos clonar o repositório do GitHub chamado [github.com/Sherrira/loadTester](https://github.com/Sherrira/loadTester). Siga a documentação do repositório para mais detalhes, mas para ser rápido, você pode construir a imagem Docker e executar o teste de carga com os seguintes comandos:

#### Passo 1: Clonar o Repositório
```sh
$ git clone https://github.com/Sherrira/loadTester.git
$ cd loadTester
```

#### Passo 2: Construir a Imagem Docker
```sh
$ docker build -t loadTester .
```

#### Passo 3: Executar o Teste de Carga
```sh
$ docker run --rm --network ratelimiter_network loadtester --url=http://app:8080 --requests=500 --concurrency=50 --method=GET --headers="API_KEY:abc123"
```

Este comando irá enviar 500 requisições concorrentes para o endpoint da aplicação, com uma concorrência de 50 requisições simultâneas, usando o método GET e incluindo o cabeçalho API_KEY:abc123.

