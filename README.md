# GLPI-ClickUp Sync

Microserviço para sincronização automática de tickets do GLPI para o ClickUp.

## Configuração

Configure as seguintes variáveis de ambiente:

```bash
export DB_CONNECTION_STRING="usuario:senha@tcp(host:porta)/database"
export CLICKUP_API_KEY="api_key_aqui"
export CLICKUP_LIST_ID="id_da_lista_clickup"
```

## Execução

```bash
go build -o glpi-clickup .
./glpi-clickup
```

## Funcionalidades

- Sincronização automática de tickets novos do GLPI
- Conversão de conteúdo HTML para texto limpo
- Mapeamento de prioridades e status
- Processamento de dados de formulário com quebras de linha

## Estrutura

- `main.go` - Ponto de entrada da aplicação
- `config.go` - Configurações e variáveis de ambiente
- `models.go` - Estruturas de dados GLPI e ClickUp
- `clickup.go` - Cliente da API do ClickUp
- `sync.go` - Serviço de sincronização
- `scheduler.go` - Agendador de tarefas
