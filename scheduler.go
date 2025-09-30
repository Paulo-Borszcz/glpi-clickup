package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Agendador struct {
	servicoSinc *ServicoSincronizacao
	intervalo   time.Duration
	canalParar  chan struct{}
}

func NovoAgendador(servicoSinc *ServicoSincronizacao, intervalo time.Duration) *Agendador {
	return &Agendador{
		servicoSinc: servicoSinc,
		intervalo:   intervalo,
		canalParar:  make(chan struct{}),
	}
}

func (a *Agendador) Iniciar(ctx context.Context) {
	log.Printf("Iniciando agendador com intervalo de %s", a.intervalo)
	
	ticker := time.NewTicker(a.intervalo)
	defer ticker.Stop()

	if err := a.servicoSinc.SincronizarTickets(); err != nil {
		log.Printf("Erro na sincronização inicial: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Agendador interrompido por contexto")
			return
		case <-a.canalParar:
			log.Println("Agendador interrompido")
			return
		case <-ticker.C:
			if err := a.servicoSinc.SincronizarTickets(); err != nil {
				log.Printf("Erro na sincronização: %v", err)
			}
		}
	}
}

func (a *Agendador) Parar() {
	close(a.canalParar)
}

func (a *Agendador) ObterStatus() string {
	contadorProcessados, ultimaSync := a.servicoSinc.ObterEstatisticas()
	return formatarStatus(contadorProcessados, ultimaSync, a.intervalo)
}

func formatarStatus(contadorProcessados int, ultimaSync time.Time, intervalo time.Duration) string {
	proximaSync := ultimaSync.Add(intervalo)
	tempoAteProxima := time.Until(proximaSync)
	
	if tempoAteProxima < 0 {
		return fmt.Sprintf("Sincronização atrasada - executando em breve... (%d tickets processados)", contadorProcessados)
	}
	
	return fmt.Sprintf("Próxima sincronização em %s (%d tickets processados)", formatarTempoRestante(tempoAteProxima), contadorProcessados)
}

func formatarTempoRestante(duracao time.Duration) string {
	segundos := int(duracao.Seconds())
	if segundos < 60 {
		return formatarSegundos(segundos)
	}
	
	minutos := segundos / 60
	segundosRestantes := segundos % 60
	return formatarMinutosESegundos(minutos, segundosRestantes)
}

func formatarSegundos(segundos int) string {
	if segundos == 1 {
		return "1 segundo"
	}
	return formatarPlural(segundos, "segundo")
}

func formatarMinutosESegundos(minutos, segundos int) string {
	minutosStr := formatarMinutos(minutos)
	if segundos == 0 {
		return minutosStr
	}
	segundosStr := formatarSegundos(segundos)
	return minutosStr + " e " + segundosStr
}

func formatarMinutos(minutos int) string {
	if minutos == 1 {
		return "1 minuto"
	}
	return formatarPlural(minutos, "minuto")
}

func formatarPlural(contador int, unidade string) string {
	return formatarNumero(contador) + " " + unidade + "s"
}

func formatarNumero(n int) string {
	return formatarInt(n)
}

func formatarInt(n int) string {
	return fmt.Sprintf("%d", n)
}