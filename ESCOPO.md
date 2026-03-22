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
- **Infraestrutura:** PostgreSQL + PostGIS (Histórico oficial), Apache Kafka (Mensageria), Redis (Cache de tempo real) rodando via Docker.