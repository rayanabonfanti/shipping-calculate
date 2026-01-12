# Shipping Calculator API

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

API REST em Go para cálculo de custos de frete e prazos de entrega baseado em dimensões, peso e tipo de entrega.

## Visão Geral

Esta aplicação fornece uma API REST para calcular custos de frete para pacotes. Suporta:
- Opções de frete padrão e expresso
- Sobretaxas baseadas em peso
- Sobretaxas baseadas em volume
- Validação de CEP brasileiro
- Validação de dimensões de pacote
- Coleta de métricas OpenTelemetry

## Funcionalidades

- **Cálculo de Frete**: Calcula custos de frete baseado em peso, volume e tipo de entrega
- **Múltiplas Opções de Frete**: Entrega padrão (2 dias) e expressa (1 dia)
- **Validação de Entrada**: Valida CEPs brasileiros, dimensões de pacote e peso
- **Telemetria**: Métricas OpenTelemetry para monitoramento e observabilidade
- **Logging Estruturado**: Logging abrangente usando zap logger

## Requisitos

- Go 1.24 ou superior
- Docker (opcional, para deploy containerizado)

## Instalação

### Desenvolvimento Local

1. Clone o repositório:
```bash
git clone https://github.com/rayanabonfanti/shipping-calculate.git
cd shipping-calculator
```

2. Instale as dependências:
```bash
go mod download
```

3. Execute a aplicação:
```bash
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080` (ou na porta especificada na variável de ambiente `PORT`).

### Docker

Construa e execute com Docker:

```bash
docker build -t shipping-calculator .
docker run -p 8080:8080 shipping-calculator
```

## Endpoints da API

### POST /calculate

Calcula o custo de frete e o tempo de entrega para um pacote.

**Corpo da Requisição:**
```json
{
  "origin_zipcode": "01310-100",
  "destination_zipcode": "04547-130",
  "weight": 2.5,
  "dimensions": {
    "length": 30.0,
    "width": 20.0,
    "height": 15.0
  },
  "is_express": false
}
```

**Resposta (200 OK):**
```json
{
  "shipping_cost": 1100.0,
  "estimated_delivery_time": "2 dias",
  "available_services": ["standard", "express"],
  "shipping_options": [
    {
      "service": "standard",
      "cost": 1100.0,
      "time": "2 dias"
    },
    {
      "service": "express",
      "cost": 1650.0,
      "time": "1 dia"
    }
  ]
}
```

**Regras de Validação:**
- `origin_zipcode` e `destination_zipcode`: Devem estar no formato de CEP brasileiro válido (8 dígitos)
- `weight`: Deve ser maior que 0 (em kg)
- `dimensions`: Todas as dimensões devem ser positivas e o volume não deve exceder 15.000 cm³

**Fórmula de Preço:**
- Custo base: 10,00 BRL (1000 centavos)
- Sobretaxa de peso: 10% do custo base por 0,5 kg
- Sobretaxa de volume: 5% do custo base por 1000 cm³
- Sobretaxa expressa: 50% do subtotal (padrão + peso + volume)

## Configuração

A aplicação pode ser configurada usando variáveis de ambiente:

- `PORT`: Porta do servidor (padrão: 8080)
- `APPLICATION_NAME`: Nome da aplicação para métricas (padrão: shipping-calculator)
- `OTEL_EXPORTER_OTLP_ENDPOINT`: URL do endpoint OTLP do OpenTelemetry para exportar métricas
- `OTEL_SERVICE_NAME`: Nome do serviço para atributos de recurso do OpenTelemetry

## Testes

Execute os testes:
```bash
go test ./...
```

Execute os testes com cobertura:
```bash
go test -cover ./...
```

### Testes BDD e Integrados

O projeto implementa testes unitários usando a biblioteca `testify` e está planejado para implementar testes BDD (Behavior-Driven Development) com testes integrados. Os testes BDD permitirão validar o comportamento da aplicação de forma mais descritiva e próxima à linguagem de negócio, facilitando a comunicação entre desenvolvedores e stakeholders.

**Planejado:**
- Implementação de testes BDD usando frameworks como Gherkin/Cucumber ou bibliotecas Go específicas
- Testes integrados que validam o fluxo completo da aplicação
- Cenários de teste descritos em linguagem natural
- Validação de comportamentos end-to-end

## Estrutura do Projeto

```
.
├── cmd/
│   └── api/
│       └── main.go          # Ponto de entrada da aplicação
├── internal/
│   ├── handler/             # Handlers HTTP
│   ├── logger/              # Utilitários de logging
│   ├── model/               # Modelos de dados
│   ├── service/             # Lógica de negócio
│   └── validator/           # Validação de entrada
├── telemetry/               # Métricas e observabilidade
├── docs/                    # Documentação
├── Dockerfile               # Arquivo de build Docker
└── go.mod                   # Definição do módulo Go
```

## Tecnologias

- **Go**: Linguagem de programação
- **Chi**: Roteador HTTP
- **Zap**: Logging estruturado
- **OpenTelemetry**: Métricas e observabilidade
- **Testify**: Framework de testes

## Documentação Adicional

- [Métricas](./docs/metrics.md) - Documentação sobre as métricas expostas
- [Operações](./docs/operations.md) - Guia de monitoramento e operações
- [Observabilidade](./docs/observability.md) - Estratégia completa de observabilidade, métricas, logs e visualização

## Licença

Este projeto é open source e está disponível sob a Licença MIT.

## Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para enviar um Pull Request.
