package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type ServicoSincronizacao struct {
	bd      *sqlx.DB
	clickup *ClienteClickUp
	estado  *EstadoSincronizacao
}

func NovoServicoSincronizacao(bd *sqlx.DB, clickup *ClienteClickUp) *ServicoSincronizacao {
	return &ServicoSincronizacao{
		bd:      bd,
		clickup: clickup,
		estado: &EstadoSincronizacao{
			UltimaSincronizacao: time.Now().Add(-24 * time.Hour),
			IDsProcessados:      make(map[int]bool),
		},
	}
}

func (s *ServicoSincronizacao) SincronizarTickets() error {
	log.Println("Iniciando sincronização de tickets...")

	tickets, err := s.obterTicketsNovos()
	if err != nil {
		return fmt.Errorf("erro ao buscar tickets: %w", err)
	}

	if len(tickets) == 0 {
		log.Println("Nenhum ticket novo encontrado")
		return nil
	}

	log.Printf("Encontrados %d tickets novos para sincronizar", len(tickets))

	for _, ticket := range tickets {
		if s.estado.IDsProcessados[ticket.ID] {
			log.Printf("Ticket %d já foi processado, pulando...", ticket.ID)
			continue
		}

		err := s.clickup.CriarTarefa(&ticket)
		if err != nil {
			log.Printf("Erro ao criar tarefa no ClickUp para ticket %d: %v", ticket.ID, err)
			continue
		}

		s.estado.IDsProcessados[ticket.ID] = true
		log.Printf("Ticket %d sincronizado com sucesso", ticket.ID)
	}

	s.estado.UltimaSincronizacao = time.Now()
	log.Printf("Sincronização concluída. Próxima sincronização em: %s", s.estado.UltimaSincronizacao.Add(30*time.Second).Format("15:04:05"))

	return nil
}

func (s *ServicoSincronizacao) obterTicketsNovos() ([]TicketComObservador, error) {
	query := `
		SELECT
			t.id,
			t.entities_id,
			t.name,
			t.date,
			t.closedate,
			t.solvedate,
			t.takeintoaccountdate,
			t.date_mod,
			t.users_id_lastupdater,
			t.status,
			t.users_id_recipient,
			t.requesttypes_id,
			t.content,
			t.urgency,
			t.impact,
			t.priority,
			t.itilcategories_id,
			t.type,
			t.global_validation,
			t.slas_id_ttr,
			t.slas_id_tto,
			t.slalevels_id_ttr,
			t.time_to_resolve,
			t.time_to_own,
			t.begin_waiting_date,
			t.sla_waiting_duration,
			t.ola_waiting_duration,
			t.olas_id_tto,
			t.olas_id_ttr,
			t.olalevels_id_ttr,
			t.ola_ttr_begin_date,
			t.internal_time_to_resolve,
			t.internal_time_to_own,
			t.waiting_duration,
			t.close_delay_stat,
			t.solve_delay_stat,
			t.takeintoaccount_delay_stat,
			t.actiontime,
			t.is_deleted,
			t.locations_id,
			t.validation_percent,
			t.date_creation,
			CASE
				WHEN g.name IS NOT NULL THEN g.name
				WHEN u.realname IS NOT NULL THEN CONCAT('Usuário: ', u.realname)
				ELSE 'Sem observador'
			END as observador,
			CONCAT('https://nexus.lojasmm.com.br/front/ticket.form.php?id=', t.id) as link
		FROM glpi_tickets t
		LEFT JOIN glpi_groups_tickets gt ON t.id = gt.tickets_id AND gt.type = 3
		LEFT JOIN glpi_groups g ON gt.groups_id = g.id
		LEFT JOIN glpi_tickets_users tu ON t.id = tu.tickets_id AND tu.type = 3
		LEFT JOIN glpi_users u ON tu.users_id = u.id
		WHERE gt.groups_id = 2
			AND t.status NOT IN (5,6)
			AND t.date_creation > ?
		ORDER BY t.date_creation DESC`

	var tickets []TicketComObservador
	ultimaSincStr := s.estado.UltimaSincronizacao.Format("2006-01-02 15:04:05")
	err := s.bd.Select(&tickets, query, ultimaSincStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}

	return tickets, nil
}

func (s *ServicoSincronizacao) ObterEstatisticas() (int, time.Time) {
	return len(s.estado.IDsProcessados), s.estado.UltimaSincronizacao
}