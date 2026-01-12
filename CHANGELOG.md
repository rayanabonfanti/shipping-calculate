# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), e este projeto adere
ao [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planejado

- Implementação de testes BDD (Behavior-Driven Development) com testes integrados
- Cenários de teste descritos em linguagem natural
- Validação de comportamentos end-to-end

## [0.1.0] - 2025-01-XX

### Adicionado

- Endpoint de cálculo de frete (`POST /calculate`)
- Validação de CEP brasileiro (formato de 8 dígitos)
- Validação de peso do pacote (deve ser positivo)
- Validação de dimensões do pacote (valores positivos, volume máximo de 15.000 cm³)
- Opção de frete padrão (entrega em 2 dias)
- Opção de frete expresso (entrega em 1 dia, sobretaxa de 50%)
- Cálculo de sobretaxa baseada em peso (10% do custo base por 0,5 kg)
- Cálculo de sobretaxa baseada em volume (5% do custo base por 1000 cm³)
- Integração de métricas OpenTelemetry:
  - `shipping.calculate`: Contador total de requisições de cálculo
  - `shipping.calculate.time`: Histograma de latência de cálculo
  - `shipping.calculate.cost.distribution`: Histograma de distribuição de custos de frete
  - `shipping.calculate.error`: Contador de erros
- Logging estruturado com propagação de contexto
- Tratamento de erros e validação abrangente
- Resposta inclui todas as opções de frete disponíveis com custos e prazos de entrega
- Cálculo de custo base baseado na distância entre CEPs
- Suporte a múltiplas opções de serviço na resposta
