package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("Iniciando microserviço GLPI-ClickUp...")

	configuracao := carregarConfiguracao()
	
	if configuracao.ClickUp.ChaveAPI == "" {
		log.Fatal("CLICKUP_API_KEY é obrigatório")
	}
	if configuracao.ClickUp.IDLista == "" {
		log.Fatal("CLICKUP_LIST_ID é obrigatório")
	}
	
	log.Printf("Configurado para List ID: %s", configuracao.ClickUp.IDLista)

	bd, err := conectarBD(configuracao.BancoDados.StringConexao)
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}
	defer bd.Close()

	clienteClickup := NovoClienteClickUp(configuracao.ClickUp.ChaveAPI, configuracao.ClickUp.IDLista)
	servicoSinc := NovoServicoSincronizacao(bd, clienteClickup)
	agendador := NovoAgendador(servicoSinc, 30*time.Second)

	ctx, cancelar := context.WithCancel(context.Background())
	defer cancelar()

	go func() {
		agendador.Iniciar(ctx)
	}()

	log.Println("Microserviço iniciado com sucesso!")
	log.Println("Para parar o serviço, pressione Ctrl+C")

	canalSinal := make(chan os.Signal, 1)
	signal.Notify(canalSinal, syscall.SIGINT, syscall.SIGTERM)

	<-canalSinal
	log.Println("Recebido sinal de parada, finalizando...")
	
	agendador.Parar()
	cancelar()
	
	log.Println("Microserviço finalizado")
}

func conectarBD(stringConexao string) (*sqlx.DB, error) {
	bd, err := sqlx.Connect("mysql", stringConexao)
	if err != nil {
		return nil, err
	}

	bd.SetMaxOpenConns(10)
	bd.SetMaxIdleConns(5)
	bd.SetConnMaxLifetime(time.Hour)

	if err := bd.Ping(); err != nil {
		return nil, err
	}

	log.Println("Conectado ao banco de dados com sucesso")
	return bd, nil
}
