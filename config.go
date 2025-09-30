package main

import (
	"os"
)

type Configuracao struct {
	BancoDados ConfigBancoDados
	ClickUp    ConfigClickUp
}

type ConfigBancoDados struct {
	StringConexao string
}

type ConfigClickUp struct {
	ChaveAPI string
	IDLista  string
}

func carregarConfiguracao() Configuracao {
	return Configuracao{
		BancoDados: ConfigBancoDados{
			StringConexao: obterEnvOuPadrao("DB_CONNECTION_STRING", ""),
		},
		ClickUp: ConfigClickUp{
			ChaveAPI: obterEnvOuPadrao("CLICKUP_API_KEY", ""),
			IDLista:  obterEnvOuPadrao("CLICKUP_LIST_ID", "901319796950"),
		},
	}
}

func obterEnvOuPadrao(chave, valorPadrao string) string {
	if valor := os.Getenv(chave); valor != "" {
		return valor
	}
	return valorPadrao
}
