# 💧 Sistema de Auditoria e Rastreamento de Carros-Pipa (MVP)

## 🎯 O Problema
Atualmente, a validação da entrega de água depende de hardwares físicos (leitura de dois cartões no caminhão). O modelo não permite auditoria em tempo real, dificulta a prevenção de fraudes de rota e restringe a transparência pública.

## 🚀 A Solução Proposta
Uma arquitetura de microsserviços orientada a eventos, utilizando **Geofencing (Cercas Virtuais)** e **Prova de Presença** para validar entregas matematicamente. O cidadão ganha um mapa em tempo real (Read-Only) e o Governo (Exército/Defesa Civil) ganha um painel de auditoria antifraude.

## ⚙️ Regras de Negócio Core

### 1. Validação Automática de Entrega (Algoritmo de 10 Minutos)
A entrega só muda para o status `ENTREGUE` se o banco de dados espacial registrar duas condições simultâneas:
- **Espacial:** Rastro do GPS do pipeiro a uma distância máxima de 50 metros da coordenada exata da cisterna cadastrada (usando PostGIS).
- **Temporal:** Presença ininterrupta dentro desse raio de 50m por no mínimo 10 minutos (tempo de esvaziamento da carga).

### 2. Transparência Pública vs. Ação Privada (Prevenção a Fraudes)
- **Modo Leitura (Público):** Qualquer cidadão com o app pode ver a localização dos caminhões em tempo real e o horário da última entrega na cisterna da sua região.
- **Modo Denúncia Segura:** Para evitar "trolls" e cliques falsos de "Não recebi a água", a denúncia exige vínculo. Se o sistema marcar a água como entregue, mas a cisterna estiver vazia, o usuário só pode abrir uma **Denúncia de Alta Prioridade (Red Flag)** se autenticar o relato com o CPF ou o Número do Cartão do Titular da Cisterna. Denúncias de terceiros entram como "Averiguação Comum".

## 🏗️ Arquitetura Técnica
- **Tracking Service (Go):** Recebe o fluxo contínuo de GPS e publica no Kafka em milissegundos.
- **Core Service (Go):** Consome o Kafka, calcula distância via PostGIS e valida os 10 minutos.
- **Map Service (Go):** API pública de alta performance lendo posições em tempo real do Redis.
- **Auth Service (Go):** Segurança e validação de JWT (Motorista, Titular da Cisterna, Admin).
- **Reporting Service (Go):** Recebe e classifica as denúncias para o painel de auditoria.
- **Infraestrutura:** PostgreSQL + PostGIS (Histórico oficial) e Apache Kafka (Mensageria) rodando via Docker.

## 🧩 Padrões de Projeto Identificados

| Padrão | Onde aparece | Por quê |
|---|---|---|
| Factory Method / Simple Factory | auth_service.go, handlers.go, gps.go, handlers.go, core_client.go | Há vários `New...` que criam objetos já prontos: `NewAuthService`, `NewHandler`, `NewGPSHandler`, `NewMapHandler`, `NewCoreClient`, `NewKafkaConsumer`, etc. |
| Adapter / Port-Adapter | core_client.go | O serviço de gestão usa uma interface `CoreClient` e uma implementação `coreClientImpl` para se adaptar ao acesso HTTP ao core-service. Também há esse estilo nas interfaces de repositório e nas implementações concretas. |
| Decorator | main.go, main.go, main.go, main.go | O middleware do Chi (Logger, Recoverer, RealIP) envolve o handler e acrescenta comportamento sem mudar o fluxo principal. Isso é muito próximo do Decorator. |
| Proxy | main.go, main.go | A função `requireRole` intercepta a requisição, valida o token e só então encaminha para o handler real. É um proxy/guard de acesso. |
| Chain of Responsibility | main.go, main.go, main.go | A cadeia de middlewares processa a requisição em sequência; cada camada pode parar a execução ou passar adiante. |
| Observer / Publisher-Subscriber | kafka_publisher.go, kafka_consumer.go, service.go | O auth-service publica eventos; o core-service e o logs-service consomem esses eventos de forma assíncrona. É um padrão de publicação/assinatura bem próximo do Observer. |
