package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClienteClickUp struct {
	ChaveAPI string
	IDLista  string
	Cliente  *http.Client
}

func NovoClienteClickUp(chaveAPI, idLista string) *ClienteClickUp {
	return &ClienteClickUp{
		ChaveAPI: chaveAPI,
		IDLista:  idLista,
		Cliente:  &http.Client{},
	}
}

func (c *ClienteClickUp) CriarTarefa(ticket *TicketComObservador) error {
	var dataInicio *int64
	if ticket.DataCriacao.Valido {
		timestamp := ticket.DataCriacao.Tempo.Unix() * 1000
		dataInicio = &timestamp
	}

	tarefa := TarefaClickUp{
		Nome:              fmt.Sprintf("GLPI #%d - %s", ticket.ID, ticket.Nome),
		DescricaoMarkdown: c.formatarDescricao(ticket),
		Prioridade:        c.mapearPrioridade(ticket.Prioridade),
		DataInicio:        dataInicio,
	}

	payload, err := json.Marshal(tarefa)
	if err != nil {
		return fmt.Errorf("erro ao serializar tarefa: %w", err)
	}

	url := fmt.Sprintf("https://api.clickup.com/api/v2/list/%s/task", c.IDLista)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Authorization", c.ChaveAPI)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Cliente.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro na API do ClickUp (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *ClienteClickUp) formatarDescricao(ticket *TicketComObservador) string {
	descricao := fmt.Sprintf("**Ticket GLPI #%d**\n\n", ticket.ID)
	descricao += fmt.Sprintf("**Observador:** %s\n", ticket.Observador)
	descricao += fmt.Sprintf("**Status:** %s\n", c.obterNomeStatus(ticket.Status))
	descricao += fmt.Sprintf("**Prioridade:** %s\n", c.obterNomePrioridade(ticket.Prioridade))
	if ticket.DataCriacao.Valido {
		descricao += fmt.Sprintf("**Criado em:** %s\n\n", ticket.DataCriacao.Tempo.Format("02/01/2006 15:04"))
	} else {
		descricao += "**Criado em:** Não informado\n\n"
	}
	
	if ticket.ObterConteudoLimpo() != "" {
		descricao += "**Conteúdo:**\n"
		descricao += ticket.ObterConteudoLimpo()
		descricao += "\n\n"
	}
	
	descricao += fmt.Sprintf("**Link:** %s", ticket.Link)
	
	return descricao
}

func (c *ClienteClickUp) mapearPrioridade(prioridade int) int {
	switch prioridade {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return 3
	case 5:
		return 4
	default:
		return 3
	}
}

func (c *ClienteClickUp) obterNomeStatus(status int) string {
	switch status {
	case 1:
		return "Novo"
	case 2:
		return "Em andamento (atribuído)"
	case 3:
		return "Em andamento (planejado)"
	case 4:
		return "Pendente"
	case 5:
		return "Solucionado"
	case 6:
		return "Fechado"
	default:
		return "Desconhecido"
	}
}

func (c *ClienteClickUp) obterNomePrioridade(prioridade int) string {
	switch prioridade {
	case 1:
		return "Muito baixa"
	case 2:
		return "Baixa"
	case 3:
		return "Normal"
	case 4:
		return "Alta"
	case 5:
		return "Muito alta"
	default:
		return "Normal"
	}
}