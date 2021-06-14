package main

import (
	"fmt"
	"sanguessuga"
)

func main() {
	fmt.Println(sanguessuga.FetchReports("https://www.gov.br/inep/pt-br/acesso-a-informacao/dados-abertos/microdados/enade", "a"))
}
