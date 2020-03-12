package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const delay = 5
const monitoramentos = 3

func main() {
	exibeMenu()
	comando := leComando()

	switch comando {
	case 1:
		iniciarMonitoramento()
	case 2:
		color.Cyan("Exibindo logs...")
		imprimeLogs()
	case 0:
		color.Cyan("Saindo do programa.")
		os.Exit(0)
	default:
		color.Cyan("Não conheço este comando.")
		os.Exit(-1)
	}
}

func iniciarMonitoramento() {
	color.White("Monitorando...")

	sites := leSitesDoArquivo()

	for i := 0; i < monitoramentos; i++ {
		for i, site := range sites {
			color.White("Testando site %d - %s:", i, site)
			testaSite(site)
		}

		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
	fmt.Println("")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	color.Yellow("O comando escolhido foi %d:", comandoLido)
	fmt.Println("")

	return comandoLido
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		color.Red("Ocorreu um erro: %s", err)
	}

	if resp.StatusCode == 200 {
		color.Green("Site: %s foi carregado com sucesso!", site)
		registraLog(site, true)
	} else {
		color.Red("Site: %s está com problemas. Status Code: %d", site, resp.StatusCode)
		registraLog(site, false)
	}
}

func leSitesDoArquivo() []string {
	var sites []string

	arquivo, err := os.Open("sites.txt")

	if err != nil {
		color.Red("Ocorreu um erro: %s", err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		color.Red("Ocorreu um erro: %s", err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006") + " - " + site +
		"- online -" + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		color.Red("Ocorreu um erro: %s", err)
	}

	fmt.Println(arquivo)
}

func exibeMenu() {
	color.Magenta("1- Iniciar Monitoramento")
	color.Magenta("2- Exibir Logs")
	color.Magenta("0- Sair do Programa")
}
