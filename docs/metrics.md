# Métricas

Este documento descreve todas as métricas expostas pela API Shipping Calculator.

## Visão Geral

A aplicação usa OpenTelemetry para expor métricas que rastreiam o desempenho, confiabilidade e padrões de uso do serviço de cálculo de frete.

## Convenção de Nomenclatura de Métricas

Todas as métricas seguem o padrão: `shipping.calculate[.suffix]`

## Métricas

### Contadores

#### `shipping.calculate`

- **Tipo**: Int64Counter
- **Descrição**: Contador de cálculos solicitados (Número total de requisições de cálculo de frete)
- **Casos de Uso**:
  - Monitorar volume de requisições e padrões de tráfego
  - Rastrear tendências de uso do serviço
  - Identificar horários de pico
- **Limiar de Alerta**: Considere alertar se a taxa de requisições cair significativamente (possível degradação do serviço)

#### `shipping.calculate.error`

- **Tipo**: Int64Counter
- **Descrição**: Contador de erros (Número total de erros no cálculo de frete)
- **Casos de Uso**:
  - Rastrear taxa de erro e identificar problemas
  - Monitorar confiabilidade do serviço
  - Detectar problemas de validação ou cálculo
- **Limiar de Alerta**: Alertar se a taxa de erro exceder 5% do total de requisições

### Histogramas

#### `shipping.calculate.time`

- **Tipo**: Int64Histogram
- **Descrição**: Tempo de resposta (Tempo gasto para calcular o frete em milissegundos)
- **Casos de Uso**:
  - Monitorar desempenho e latência da API
  - Rastrear tendências de tempo de resposta
  - Identificar degradação de performance
- **Limiar de Alerta**: Alertar se a latência p95 exceder 500ms ou p99 exceder 1000ms

#### `shipping.calculate.cost.distribution`

- **Tipo**: Float64Histogram
- **Descrição**: Distribuição dos custos calculados (Distribuição dos custos de frete calculados)
- **Casos de Uso**:
  - Analisar padrões de custo e detectar anomalias
  - Monitorar tendências de custo de frete
  - Identificar cálculos de custo incomuns
- **Limiar de Alerta**: Considere alertar se a distribuição de custos mostrar padrões inesperados

## Configuração

### Variáveis de Ambiente

- `APPLICATION_NAME`: Nome da aplicação para OpenTelemetry. Padrão: "shipping-calculator" se não definido.
- `OTEL_EXPORTER_OTLP_ENDPOINT`: URL do endpoint OTLP para exportar métricas (ex: `http://otel-collector:4318`)
- `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`: Endpoint opcional específico para métricas (sobrescreve OTEL_EXPORTER_OTLP_ENDPOINT para métricas)
- `OTEL_SERVICE_NAME`: Nome do serviço para atributos de recurso do OpenTelemetry

## Coleta de Métricas

As métricas são coletadas e exportadas automaticamente através do SDK do OpenTelemetry. As métricas são registradas nos seguintes locais:

- **Contador de Requisições**: Incrementado no início de cada requisição de cálculo
- **Contador de Erros**: Incrementado quando ocorrem erros de validação ou cálculo
- **Histograma de Tempo**: Registrado após conclusão bem-sucedida do cálculo
- **Histograma de Distribuição de Custos**: Registrado após conclusão bem-sucedida do cálculo

## Integração

Essas métricas podem ser integradas com qualquer plataforma de observabilidade compatível com OpenTelemetry (Prometheus, Grafana, Datadog, etc.) e podem ser visualizadas em dashboards e usadas para alertas. Para mais informações sobre operações e monitoramento, consulte [operations.md](./operations.md) e [observability.md](./observability.md).
