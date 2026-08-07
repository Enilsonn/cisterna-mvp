# auth-service

Serviço de autenticação inicial com suporte a:
- login por email/senha;
- emissão de access token e refresh token;
- refresh de tokens;
- logout revogando refresh token;
- validação de access token.

## Como executar

```bash
cd auth-service
go run ./cmd/main.go
```

## Credenciais de exemplo

- admin@cisterna.com / admin123
- pipeiro@cisterna.com / pipeiro123
- cidadao@cisterna.com / cidadao123
