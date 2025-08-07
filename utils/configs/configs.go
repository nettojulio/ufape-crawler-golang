package configs

import (
	"fmt"
	"time"
)

var (
	MaxDepth       uint
	RemoveFragment bool
	LowerCaseURLs  bool
	RequestTimeout time.Duration
	FliterLinks    bool
	InitialURL     string
)

func DisplayConfigs() {
	split := "============================================================"
	fmt.Println(split)
	fmt.Printf("Configuração da Aplicação:\n")
	fmt.Printf("\t - Profundidade Máxima: %d\n", MaxDepth)
	fmt.Printf("\t - Remover Fragmentos: %v\n", RemoveFragment)
	fmt.Printf("\t - URLs em Minúsculas: %v\n", LowerCaseURLs)
	fmt.Printf("\t - Timeout da Requisição: %s\n", RequestTimeout)
	fmt.Printf("\t - Filtrar links já visitados: %v\n", FliterLinks)
	fmt.Printf("\t - URL Inicial: %s\n", InitialURL)
	fmt.Println(split + "\n")
}
