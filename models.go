package main

import (
	"database/sql/driver"
	"fmt"
	"html"
	"strings"
	"time"
)

type TempoNulo struct {
	Tempo  time.Time
	Valido bool
}

func (tn *TempoNulo) Scan(valor interface{}) error {
	if valor == nil {
		tn.Tempo, tn.Valido = time.Time{}, false
		return nil
	}
	
	switch v := valor.(type) {
	case string:
		if v == "" || v == "0000-00-00 00:00:00" {
			tn.Tempo, tn.Valido = time.Time{}, false
			return nil
		}
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			tn.Tempo, tn.Valido = time.Time{}, false
			return err
		}
		tn.Tempo, tn.Valido = t, true
		return nil
	case []byte:
		s := string(v)
		if s == "" || s == "0000-00-00 00:00:00" {
			tn.Tempo, tn.Valido = time.Time{}, false
			return nil
		}
		t, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			tn.Tempo, tn.Valido = time.Time{}, false
			return err
		}
		tn.Tempo, tn.Valido = t, true
		return nil
	case time.Time:
		tn.Tempo, tn.Valido = v, true
		return nil
	default:
		return fmt.Errorf("cannot scan %T into TempoNulo", valor)
	}
}

func (tn TempoNulo) Value() (driver.Value, error) {
	if !tn.Valido {
		return nil, nil
	}
	return tn.Tempo, nil
}

type TicketGLPI struct {
	ID                         int       `json:"id" db:"id"`
	IDEntidades                int       `json:"entities_id" db:"entities_id"`
	Nome                       string    `json:"name" db:"name"`
	Data                       TempoNulo `json:"date" db:"date"`
	DataFechamento             TempoNulo `json:"closedate" db:"closedate"`
	DataResolucao              TempoNulo `json:"solvedate" db:"solvedate"`
	DataAceitacao              TempoNulo `json:"takeintoaccountdate" db:"takeintoaccountdate"`
	DataModificacao            TempoNulo `json:"date_mod" db:"date_mod"`
	IDUsuarioUltimaAtualizacao int       `json:"users_id_lastupdater" db:"users_id_lastupdater"`
	Status                     int       `json:"status" db:"status"`
	IDUsuarioDestinatario      int       `json:"users_id_recipient" db:"users_id_recipient"`
	IDTipoRequisicao           int       `json:"requesttypes_id" db:"requesttypes_id"`
	Conteudo                   string    `json:"content" db:"content"`
	Urgencia                   int       `json:"urgency" db:"urgency"`
	Impacto                    int       `json:"impact" db:"impact"`
	Prioridade                 int       `json:"priority" db:"priority"`
	IDCategoriasITIL           int       `json:"itilcategories_id" db:"itilcategories_id"`
	Tipo                       int       `json:"type" db:"type"`
	ValidacaoGlobal            int       `json:"global_validation" db:"global_validation"`
	IDSLATempResolucao         int       `json:"slas_id_ttr" db:"slas_id_ttr"`
	IDSLATempPropriedade       int       `json:"slas_id_tto" db:"slas_id_tto"`
	IDNivelSLATempResolucao    int       `json:"slalevels_id_ttr" db:"slalevels_id_ttr"`
	TempoParaResolucao         TempoNulo `json:"time_to_resolve" db:"time_to_resolve"`
	TempoParaPropriedade       TempoNulo `json:"time_to_own" db:"time_to_own"`
	DataInicioEspera           TempoNulo `json:"begin_waiting_date" db:"begin_waiting_date"`
	DuracaoEsperaSLA           int       `json:"sla_waiting_duration" db:"sla_waiting_duration"`
	DuracaoEsperaOLA           int       `json:"ola_waiting_duration" db:"ola_waiting_duration"`
	IDOLATempPropriedade       int       `json:"olas_id_tto" db:"olas_id_tto"`
	IDOLATempResolucao         int       `json:"olas_id_ttr" db:"olas_id_ttr"`
	IDNivelOLATempResolucao    int       `json:"olalevels_id_ttr" db:"olalevels_id_ttr"`
	DataInicioOLATempResolucao TempoNulo `json:"ola_ttr_begin_date" db:"ola_ttr_begin_date"`
	TempoInternoParaResolucao  TempoNulo `json:"internal_time_to_resolve" db:"internal_time_to_resolve"`
	TempoInternoParaPropriedade TempoNulo `json:"internal_time_to_own" db:"internal_time_to_own"`
	DuracaoEspera              int       `json:"waiting_duration" db:"waiting_duration"`
	EstatisticaAtrasoFechamento int      `json:"close_delay_stat" db:"close_delay_stat"`
	EstatisticaAtrasoResolucao int       `json:"solve_delay_stat" db:"solve_delay_stat"`
	EstatisticaAtrasoAceitacao int       `json:"takeintoaccount_delay_stat" db:"takeintoaccount_delay_stat"`
	TempoAcao                  int       `json:"actiontime" db:"actiontime"`
	Excluido                   int       `json:"is_deleted" db:"is_deleted"`
	IDLocalizacoes             int       `json:"locations_id" db:"locations_id"`
	PercentualValidacao        int       `json:"validation_percent" db:"validation_percent"`
	DataCriacao                TempoNulo `json:"date_creation" db:"date_creation"`
}

type TicketComObservador struct {
	TicketGLPI
	Observador string `db:"observador"`
	Link       string `db:"link"`
}

func (t *TicketGLPI) ObterConteudoLimpo() string {
	conteudo := html.UnescapeString(t.Conteudo)
	
	if indiceFooter := strings.Index(conteudo, "<footer"); indiceFooter != -1 {
		conteudo = conteudo[:indiceFooter]
	}
	
	conteudo = strings.ReplaceAll(conteudo, "<br>", "\n")
	conteudo = strings.ReplaceAll(conteudo, "<br/>", "\n")
	conteudo = strings.ReplaceAll(conteudo, "<br />", "\n")
	
	for strings.Contains(conteudo, "<") && strings.Contains(conteudo, ">") {
		inicio := strings.Index(conteudo, "<")
		fim := strings.Index(conteudo[inicio:], ">")
		if fim != -1 {
			conteudo = conteudo[:inicio] + conteudo[inicio+fim+1:]
		} else {
			break
		}
	}
	
	conteudo = strings.ReplaceAll(conteudo, "61) ", "6\n1) ")
	conteudo = strings.ReplaceAll(conteudo, "62) ", "6\n2) ")
	conteudo = strings.ReplaceAll(conteudo, "63) ", "6\n3) ")
	conteudo = strings.ReplaceAll(conteudo, "64) ", "6\n4) ")
	conteudo = strings.ReplaceAll(conteudo, "65) ", "6\n5) ")
	conteudo = strings.ReplaceAll(conteudo, "66) ", "6\n6) ")
	conteudo = strings.ReplaceAll(conteudo, "67) ", "6\n7) ")
	conteudo = strings.ReplaceAll(conteudo, "68) ", "6\n8) ")
	conteudo = strings.ReplaceAll(conteudo, "69) ", "6\n9) ")
	
	conteudo = strings.ReplaceAll(conteudo, "71) ", "7\n1) ")
	conteudo = strings.ReplaceAll(conteudo, "72) ", "7\n2) ")
	conteudo = strings.ReplaceAll(conteudo, "73) ", "7\n3) ")
	conteudo = strings.ReplaceAll(conteudo, "74) ", "7\n4) ")
	conteudo = strings.ReplaceAll(conteudo, "75) ", "7\n5) ")
	conteudo = strings.ReplaceAll(conteudo, "76) ", "7\n6) ")
	conteudo = strings.ReplaceAll(conteudo, "77) ", "7\n7) ")
	conteudo = strings.ReplaceAll(conteudo, "78) ", "7\n8) ")
	conteudo = strings.ReplaceAll(conteudo, "79) ", "7\n9) ")
	
	for i := 0; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			padrao := fmt.Sprintf("%d%d) ", i, j)
			substituicao := fmt.Sprintf("%d\n%d) ", i, j)
			conteudo = strings.ReplaceAll(conteudo, padrao, substituicao)
		}
	}
	
	for j := 1; j <= 9; j++ {
		padrao := fmt.Sprintf("s%d) ", j)
		substituicao := fmt.Sprintf("s\n%d) ", j)
		conteudo = strings.ReplaceAll(conteudo, padrao, substituicao)
		
		padrao = fmt.Sprintf("a%d) ", j)
		substituicao = fmt.Sprintf("a\n%d) ", j)
		conteudo = strings.ReplaceAll(conteudo, padrao, substituicao)
		
		padrao = fmt.Sprintf("m%d) ", j)
		substituicao = fmt.Sprintf("m\n%d) ", j)
		conteudo = strings.ReplaceAll(conteudo, padrao, substituicao)
		
		padrao = fmt.Sprintf("o%d) ", j)
		substituicao = fmt.Sprintf("o\n%d) ", j)
		conteudo = strings.ReplaceAll(conteudo, padrao, substituicao)
		
		padrao = fmt.Sprintf("e%d) ", j)
		substituicao = fmt.Sprintf("e\n%d) ", j)
		conteudo = strings.ReplaceAll(conteudo, padrao, substituicao)
		
		padrao = fmt.Sprintf("d%d) ", j)
		substituicao = fmt.Sprintf("d\n%d) ", j)
		conteudo = strings.ReplaceAll(conteudo, padrao, substituicao)
	}
	
	conteudo = strings.ReplaceAll(conteudo, "Seção1) ", "Seção\n1) ")
	
	conteudo = strings.ReplaceAll(conteudo, "\n\n\n", "\n\n")
	for strings.Contains(conteudo, "\n\n\n") {
		conteudo = strings.ReplaceAll(conteudo, "\n\n\n", "\n\n")
	}
	
	return strings.TrimSpace(conteudo)
}

type TarefaClickUp struct {
	Nome              string `json:"name"`
	DescricaoMarkdown string `json:"markdown_description"`
	Prioridade        int    `json:"priority"`
	DataInicio        *int64 `json:"start_date,omitempty"`
}

type EstadoSincronizacao struct {
	UltimaSincronizacao time.Time
	IDsProcessados      map[int]bool
}