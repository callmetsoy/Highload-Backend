# sis2-protocol

Сравнение gRPC (HTTP/2 + Protobuf) и REST (HTTP/1.1 + JSON) для внутреннего API сервиса коротких ссылок. 
Цель проекта — измерить разницу в latency (p50/p95/p99) и размере payload, а также продемонстрировать реализацию идемпотентности, таймаутов и ретраев.

## 🏗 Архитектура

Проект состоит из двух независимых серверов (REST и gRPC), которые используют общую бизнес-логику (`internal/domain`).

```mermaid
graph TD
    Client[Client / Benchmark] -->|HTTP POST| RestServer[REST Server :8080]
    Client -->|gRPC| GrpcServer[gRPC Server :9090]

    RestServer -->|JSON| Domain[Domain Logic]
    GrpcServer -->|Protobuf| Domain

    Domain -->|Idempotency Check| Store[(In-Memory Store)]

    style Client fill:#f9f,stroke:#333,stroke-width:2px
    style RestServer fill:#bbf,stroke:#333,stroke-width:1px
    style GrpcServer fill:#bbf,stroke:#333,stroke-width:1px
    style Domain fill:#bfb,stroke:#333,stroke-width:1px
    style Store fill:#fbb,stroke:#333,stroke-width:1px
