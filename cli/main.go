package main

import (
	"fmt"
	"log"
	"os"
	"sanguessuga"
)

func main() {
	reports, _ := sanguessuga.ScrapeReports("https://www.gov.br/inep/pt-br/acesso-a-informacao/dados-abertos/microdados/censo-da-educacao-superior", "a")

	err := os.MkdirAll("tmp", 0755)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to cread tmp folder to download reports: %q", err))
	}

	fmt.Println(reports[0])

	fmt.Println(sanguessuga.DownloadReport(reports[0], "tmp"))
}
