# Operações

Este documento fornece informações sobre monitoramento, dashboards e aspectos operacionais da API Shipping Calculator.

## Métricas

A aplicação expõe as seguintes métricas OpenTelemetry:

Todas as métricas seguem o padrão: `shipping.calculate[.suffix]`

### Contadores

- **`shipping.calculate`**: Contador de cálculos solicitados (Número total de requisições de cálculo de frete)
  - Caso de uso: Monitorar volume de requisições e padrões de tráfego
  - Limiar de alerta: Considere alertar se a taxa de requisições cair significativamente (possível degradação do serviço)

- **`shipping.calculate.error`**: Contador de erros (Número total de erros no cálculo de frete)
  - Caso de uso: Rastrear taxa de erro e identificar problemas
  - Limiar de alerta: Alertar se a taxa de erro exceder 5% do total de requisições

### Histogramas

- **`shipping.calculate.time`**: Tempo de resposta (Tempo gasto para calcular o frete em milissegundos)
  - Caso de uso: Monitorar desempenho e latência da API
  - Limiar de alerta: Alertar se a latência p95 exceder 500ms ou p99 exceder 1000ms

- **`shipping.calculate.cost.distribution`**: Distribuição dos custos calculados (Distribuição dos custos de frete calculados)
  - Caso de uso: Analisar padrões de custo e detectar anomalias
  - Limiar de alerta: Considere alertar se a distribuição de custos mostrar padrões inesperados

## Dashboards

Dashboards recomendados para criar:

1. **Dashboard de Volume de Requisições**
   - Total de requisições por minuto/hora
   - Tendências de taxa de requisições
   - Horários de pico de uso

2. **Dashboard de Taxa de Erro**
   - Contagem de erros ao longo do tempo
   - Percentual de taxa de erro
   - Breakdown de erros por tipo (erros de validação, erros de cálculo)

3. **Dashboard de Performance**
   - Latência P50, P95, P99
   - Distribuição de duração de requisições
   - Tendências de tempo de resposta

4. **Dashboard de Distribuição de Custos**
   - Histograma de custo de frete
   - Custo médio de frete
   - Distribuição de custos por tipo de serviço (padrão vs expresso)

## Monitores

### Monitor de Alta Taxa de Erro

- **Métrica**: `shipping.calculate.error` / `shipping.calculate`
- **Limiar**: Taxa de erro > 5% por 5 minutos
- **Propósito**: Detectar quando o serviço está experimentando altas taxas de erro, o que pode indicar problemas de validação, problemas de cálculo ou degradação do serviço

### Monitor de Alta Latência

- **Métrica**: `shipping.calculate.time` (p95)
- **Limiar**: Latência P95 > 500ms por 5 minutos
- **Propósito**: Identificar degradação de performance que pode impactar a experiência do usuário

### Monitor de Disponibilidade do Serviço

- **Métrica**: `shipping.calculate` (contagem de requisições)
- **Limiar**: Contagem de requisições cai para 0 por 10 minutos
- **Propósito**: Detectar indisponibilidade completa do serviço ou problemas de roteamento

## Logging

A aplicação usa logging estruturado com os seguintes níveis de log:

- **INFO**: Operações bem-sucedidas, detalhes de cálculo, processamento de requisições
- **ERROR**: Falhas de validação, erros de cálculo, erros do serviço

Todos os logs incluem informações de contexto como:
- CEPs de origem e destino da requisição
- Peso e volume do pacote
- Detalhes do cálculo (custo base, sobretaxas)
- Mensagens de erro e stack traces

## Health Checks

A aplicação pode expor endpoints de health check. Monitore esses endpoints para garantir a disponibilidade do serviço.

## Troubleshooting

### Problemas Comuns

1. **Alta Taxa de Erro**
   - Verifique erros de validação: CEPs inválidos, pesos negativos ou dimensões inválidas
   - Revise logs de erro para falhas de validação específicas
   - Verifique o formato dos dados de entrada dos clientes

2. **Alta Latência**
   - Verifique recursos do sistema (CPU, memória)
   - Revise a performance da lógica de cálculo
   - Verifique dependências externas que possam estar lentas

3. **Zero Requisições**
   - Verifique se o serviço está em execução e acessível
   - Verifique configuração de roteamento e load balancer
   - Revise conectividade de rede
