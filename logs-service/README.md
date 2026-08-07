# logs-service

Serviço de auditoria orientado a eventos para registrar eventos de posição recebidos via Kafka.

## Como executar

1. Suba a infraestrutura com Docker Compose:
   ```bash
   docker compose -f docker-compose.infra.yaml up -d
   ```
2. Execute o serviço:
   ```bash
   cd logs-service
   go run ./cmd/main.go
   ```

## Comportamento atual

- consome mensagens do tópico `truck_coordinates`;
- registra cada evento em um arquivo local `logs.txt`;
- imprime no console os eventos recebidos para observação simples.
