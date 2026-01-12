# Observabilidade

Este documento descreve a estratégia de observabilidade da API Shipping Calculator, incluindo métricas, logs, traces e visualização.

## Visão Geral

A observabilidade é fundamental para entender o comportamento da aplicação em produção. Este projeto implementa observabilidade através dos três pilares:

1. **Métricas**: Medições numéricas do comportamento do sistema ao longo do tempo
2. **Logs**: Registros de eventos e atividades da aplicação
3. **Traces**: Rastreamento de requisições através de múltiplos serviços (planejado)

A aplicação utiliza OpenTelemetry como padrão para instrumentação e coleta de dados de observabilidade.

## Pilares da Observabilidade

### Métricas

As métricas fornecem insights sobre o desempenho e comportamento da aplicação. A aplicação expõe as seguintes métricas através do OpenTelemetry:

#### Convenção de Nomenclatura

Todas as métricas seguem o padrão: `shipping.calculate[.suffix]`

#### Métricas Disponíveis

- **`shipping.calculate`** (Int64Counter)
  - Contador de cálculos solicitados
  - Incrementado no início de cada requisição

- **`shipping.calculate.error`** (Int64Counter)
  - Contador de erros no cálculo de frete
  - Incrementado quando ocorrem erros de validação ou cálculo

- **`shipping.calculate.time`** (Int64Histogram)
  - Tempo de resposta em milissegundos
  - Registrado após conclusão bem-sucedida do cálculo

- **`shipping.calculate.cost.distribution`** (Float64Histogram)
  - Distribuição dos custos de frete calculados
  - Registrado após conclusão bem-sucedida do cálculo

Para mais detalhes sobre as métricas, consulte [metrics.md](./metrics.md).

### Logs

A aplicação utiliza logging estruturado com o zap logger, garantindo logs consistentes e pesquisáveis.

#### Estrutura de Logs

Todos os logs incluem:
- **Timestamp**: Data e hora do evento
- **Level**: Nível do log (INFO, WARN, ERROR)
- **Message**: Mensagem descritiva do evento
- **Context**: Campos estruturados adicionais

#### Campos de Contexto

- `correlation_id`: ID de correlação para rastrear requisições
- `trace_id`: ID de trace (quando disponível)
- `origem`: CEP de origem
- `destino`: CEP de destino
- `peso`: Peso do pacote
- `volume`: Volume do pacote
- `custo_base`: Custo base do frete
- `acréscimo_peso`: Sobretaxa por peso
- `acréscimo_volume`: Sobretaxa por volume
- `custo_envio`: Custo final do frete
- `tempo_estimado`: Tempo estimado de entrega

#### Níveis de Log

- **INFO**: Operações bem-sucedidas, detalhes de cálculo, processamento de requisições
- **WARN**: Solicitações com parâmetros inválidos, validações que falharam
- **ERROR**: Falhas de validação, erros de cálculo, erros do serviço

### Traces

Rastreamento distribuído está planejado para implementação futura, permitindo rastrear requisições através de múltiplos serviços e componentes.

## Configuração

### Variáveis de Ambiente

Configure estas variáveis de ambiente ao executar a aplicação:

```bash
# Opcional: Nome da aplicação
export APPLICATION_NAME="shipping-calculator"

# Opcional: Endpoint OpenTelemetry
export OTEL_EXPORTER_OTLP_ENDPOINT="http://otel-collector:4318"

# Opcional: Nome do serviço
export OTEL_SERVICE_NAME="shipping-calculator"
```

### OpenTelemetry Collector

A aplicação exporta métricas via OpenTelemetry Protocol (OTLP). Configure o OpenTelemetry Collector para receber e processar essas métricas:

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
  
  # Opcional: Exportar para outros backends
  # otlp:
  #   endpoint: "jaeger:4317"

service:
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [prometheus]
```

Configure o Prometheus para coletar do collector:

```yaml
scrape_configs:
  - job_name: 'otel-collector'
    static_configs:
      - targets: ['otel-collector:8889']
```

## Visualização com Grafana

O Grafana é a ferramenta recomendada para visualização de métricas e logs.

### Pré-requisitos

- Instância do Grafana (versão 8.0+)
- Prometheus configurado para coletar métricas
- Loki configurado para agregação de logs (opcional)

### Configuração de Data Sources

#### 1. Prometheus

1. No Grafana, vá em **Configuration** > **Data Sources** > **Add data source**
2. Selecione **Prometheus**
3. Configure a URL do Prometheus (ex: `http://prometheus:9090`)
4. Clique em **Save & Test**

#### 2. Loki (Opcional)

1. No Grafana, vá em **Configuration** > **Data Sources** > **Add data source**
2. Selecione **Loki**
3. Configure a URL do Loki (ex: `http://loki:3100`)
4. Clique em **Save & Test**

### Dashboards Recomendados

#### Dashboard de Visão Geral

1. **Taxa de Requisições** - Séries temporais mostrando requisições por segundo
2. **Taxa de Erro** - Séries temporais mostrando percentual de erro
3. **Tempo de Resposta** - Séries temporais com percentis P50, P95, P99
4. **Distribuição de Custos** - Histograma mostrando distribuição de custos de frete
5. **Total de Requisições** - Painel de estatística mostrando contagem total
6. **Total de Erros** - Painel de estatística mostrando contagem total de erros

#### Dashboard de Performance

1. **Percentis de Latência** - P50, P95, P99, P99.9
2. **Distribuição de Duração de Requisições** - Histograma
3. **Throughput** - Requisições por segundo/minuto
4. **Breakdown de Erros** - Gráfico de pizza ou barras por tipo de erro

#### Dashboard de Análise de Custos

1. **Custo Médio de Frete** - Séries temporais
2. **Distribuição de Custos** - Histograma
3. **Percentis de Custos** - P50, P95, P99
4. **Custo por Tipo de Serviço** - Padrão vs Expresso

### Consultas de Exemplo

#### Métricas (PromQL)

**Taxa de Requisições:**
```promql
rate(shipping_calculate_total[5m])
```

**Taxa de Erro:**
```promql
rate(shipping_calculate_error_total[5m]) / rate(shipping_calculate_total[5m]) * 100
```

**Latência P95:**
```promql
histogram_quantile(0.95, rate(shipping_calculate_time_bucket[5m]))
```

#### Logs (LogQL)

**Logs de Erro:**
```logql
{job="shipping-calculator"} |= "error" | json
```

**Logs de Requisições:**
```logql
{job="shipping-calculator"} |= "Solicitação de cálculo" | json
```

**Filtrar por Correlation ID:**
```logql
{job="shipping-calculator"} | json | correlation_id="<correlation_id>"
```

## Alertas

### Alerta de Alta Taxa de Erro

**Condição:** Taxa de erro > 5% por 5 minutos

**Query:**
```promql
rate(shipping_calculate_error_total[5m]) / rate(shipping_calculate_total[5m]) > 0.05
```

### Alerta de Alta Latência

**Condição:** Latência P95 > 500ms por 5 minutos

**Query:**
```promql
histogram_quantile(0.95, rate(shipping_calculate_time_bucket[5m])) > 500
```

### Alerta de Serviço Indisponível

**Condição:** Sem requisições por 10 minutos

**Query:**
```promql
rate(shipping_calculate_total[10m]) == 0
```

## Boas Práticas

### Métricas

- Use nomes descritivos e consistentes
- Inclua labels relevantes para filtragem e agregação
- Evite alta cardinalidade em labels
- Documente o significado de cada métrica

### Logs

- Use logging estruturado com campos consistentes
- Inclua correlation_id em todos os logs relacionados a uma requisição
- Evite logs excessivos em loops ou operações de alto volume
- Use níveis de log apropriados (INFO, WARN, ERROR)

### Performance

- Minimize o overhead de instrumentação
- Use amostragem para traces quando necessário
- Configure retenção apropriada para logs e métricas
- Monitore o próprio sistema de observabilidade

## Troubleshooting

### Métricas Não Aparecendo

1. Verifique a configuração do endpoint OpenTelemetry
3. Verifique se as métricas estão sendo exportadas (verifique logs do collector)
4. Certifique-se de que o Prometheus está coletando do endpoint correto

### Logs Não Aparecendo

1. Verifique se o Loki está em execução e acessível
2. Verifique a configuração do Promtail
3. Verifique se o formato do log corresponde às expectativas do Loki
4. Verifique os logs da aplicação para campos correlation_id e trace_id

### Dashboard Não Carregando

1. Verifique se o data source está configurado corretamente
2. Verifique se os nomes das métricas correspondem
3. Verifique se o intervalo de tempo é apropriado
4. Verifique se há erros de sintaxe nas consultas

## Testes e Observabilidade

A observabilidade também é importante durante o desenvolvimento e testes. O projeto está planejado para implementar testes BDD (Behavior-Driven Development) com testes integrados, que permitirão validar o comportamento da aplicação de forma mais descritiva e próxima à linguagem de negócio.

Os testes integrados podem ser instrumentados para coletar métricas e logs durante a execução, facilitando a identificação de problemas e a validação do comportamento esperado.

## Documentação Relacionada

- [Métricas](./metrics.md) - Documentação detalhada sobre as métricas expostas
- [Operações](./operations.md) - Guia de monitoramento e operações
- [README.md](../README.md) - Visão geral do projeto
