# Test Harness

Esta pasta reúne scripts e uma coleção de requisições para simular o fluxo do sistema sem depender do frontend.

## Estrutura
- `admin_flow.py`: cria um usuário admin, faz login, cadastra cisternas, caminhões, pipeiros e entregas.
- `pipeiro_flow.py`: cria um usuário pipeiro, faz login, envia posições GPS a cada 2 segundos.
- `citizen_flow.py`: cria um usuário cidadão e faz login.
- `gps_points.json`: arquivo com posições de exemplo para simular GPS.
- `manual-test-collection.json`: coleção pronta para importar no Postman ou Insomnia.

## Requisitos
- Python 3.10+
- requests

Instale:
```bash
pip install requests
```

## Como usar a coleção manual
1. Inicie os serviços de auth, management, tracking e logs.
2. Importe o arquivo `manual-test-collection.json` no Postman ou Insomnia.
3. Defina os valores das variáveis de ambiente:
   - `base_url_auth`: `http://localhost:8082`
   - `base_url_management`: `http://localhost:8081`
   - `base_url_tracking`: `http://localhost:8080`
4. Execute primeiro as requisições de auth para obter `access_token` e `refresh_token`.
5. Depois, use o token gerado nas requisições protegidas do management e tracking.

## Fluxo recomendado
1. Auth → register/login
2. Management → criar cisterna, pipeiro, caminhão e entrega
3. Tracking → enviar posições GPS
4. Logs → acompanhar os eventos no console
