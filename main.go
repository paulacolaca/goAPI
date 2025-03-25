package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/enescakir/emoji"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const monitoramentos = 1
const delay = 60

type MonitoramentoAPIs struct {
	url            string
	status         int
	tempoResposta  float64
	dataRequisicao string
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	exibeIntroducao()

	for {
		exibeMenu()

		comando := leComando()
		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			exibirLogs()
		case 3:
			fmt.Println("Saindo do programa, até a próxima!", emoji.GrinningFaceWithSmilingEyes)
			os.Exit(0)
		default:
			fmt.Println("Comando inválido, por favor escolha uma opção entre 1 e 3.")
		}
	}
}

func exibeIntroducao() {
	versao := 1.0
	fmt.Println("Olá, usuário(a)!", emoji.WavingHand)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("1 – Iniciar monitoramento de APIs")
	fmt.Println("2 – Exibir Logs")
	fmt.Println("3 – Sair do programa")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("")

	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println(emoji.MagnifyingGlassTiltedLeft, "Monitoramento iniciado...")
	listaAPIS, _ := buscarAPIsBD()

	for i := 0; i < monitoramentos; i++ {
		for j, api := range listaAPIS {

			fmt.Println("Testando API", j, ":", api)
			checarAPI(api)
		}
		time.Sleep(delay * time.Minute)
		fmt.Print("")
	}
}

func buscarAPIsBD() ([]string, error) {

	connString := os.Getenv("DB_Token")

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println("Erro ao conectar:", err)
	}

	defer conn.Close(context.Background())

	sql := "SELECT lk_url FROM listaapis"

	rows, err := conn.Query(context.Background(), sql)

	if err != nil {
		fmt.Println("Erro ao executar o SELECT:", err)
	}
	defer rows.Close()

	var urls []string

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func checarAPI(api string) MonitoramentoAPIs {
	inicio := time.Now()
	resp, err := http.Get(api)
	tempoResposta := time.Since(inicio).Seconds() * 1000

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		monitoramento := MonitoramentoAPIs{
			url:            api,
			tempoResposta:  0,
			dataRequisicao: time.Now().Format("02/01/2006 15:04:05"),
			status:         0,
		}
		return monitoramento
	}
	defer resp.Body.Close()

	monitoramento := MonitoramentoAPIs{
		url:            api,
		tempoResposta:  tempoResposta,
		dataRequisicao: time.Now().Format("02/01/2006 15:04:05"),
	}

	if resp.StatusCode == 200 {
		fmt.Println(emoji.CheckMark, " API:", api, "está respondendo!")
		monitoramento.status = 1
		registraLog(monitoramento)
	} else {
		fmt.Println(emoji.CrossMark, " API:", api, "não está respondendo! Status Code:", resp.StatusCode)
		monitoramento.status = 0
		registraLog(monitoramento)
	}

	return monitoramento
}

func registraLog(api MonitoramentoAPIs) {

	connString := os.Getenv("DB_Token")

	conn, err := pgx.Connect(context.Background(), connString)

	if err != nil {
		fmt.Println("Erro ao conectar:", err)
	}

	defer conn.Close(context.Background())

	sql := "INSERT INTO logapi (lk_url, vl_status,vl_temporesposta, dt_requisicao) VALUES ($1, $2, $3, $4)"
	_, err = conn.Exec(context.Background(), sql, api.url, api.status, api.tempoResposta, api.dataRequisicao)

	if err != nil {
		fmt.Println("Erro ao executar o INSERT:", err)
	}
}

func exibirLogs() {

	connString := os.Getenv("DB_Token")

	conn, err := pgx.Connect(context.Background(), connString)

	if err != nil {
		fmt.Println("Erro ao conectar:", err)
	}

	defer conn.Close(context.Background())

	sql := "SELECT lk_url, vl_status, vl_temporesposta, dt_requisicao FROM logapi WHERE dt_requisicao::date = now()::date"

	rows, err := conn.Query(context.Background(), sql)

	if err != nil {
		fmt.Println("Erro ao executar o SELECT:", err)
	}

	defer rows.Close()

	var logResultados []string

	for rows.Next() {
		var url, tempoResposta string
		var status int
		var dataRequisicao time.Time

		if err := rows.Scan(&url, &status, &tempoResposta, &dataRequisicao); err != nil {
			fmt.Println("Ocorreu um erro:", err)
		}

		statusIcone := emoji.GreenCircle

		if status == 0 {
			statusIcone = emoji.RedCircle
		}

		logRow := fmt.Sprintf("%v URL: %s | Status: %d | Tempo de resposta: %sms | Data da requisição: %s ", statusIcone, url, status, strings.Split(tempoResposta, ".")[0], dataRequisicao.Format("2006-01-02 15:04:05"))

		logResultados = append(logResultados, logRow+"\n")
	}

	if len(logResultados) > 0 {
		fmt.Println(logResultados)
	} else {
		fmt.Println("Não há registros de monitoramentos realizados hoje", time.Now().Format("02/01/2006"))
	}

	fmt.Println("")
}
