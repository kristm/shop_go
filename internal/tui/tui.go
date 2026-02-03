package tui

import (
	"fmt"
	"log"
	"os"
	"shop_go/internal/config"
	"shop_go/internal/models"
	"shop_go/internal/tui/modes"

	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	ProductUpdate Mode = iota
)

func Run(mode Mode, param any) {
	cfg, err := config.LoadConfig(".env")
	_ = models.ConnectDatabase(&cfg)
	if err != nil {
		log.Printf("ERROR LOADING CONFIG")
	}

	var initialModel any

	switch mode {
	case ProductUpdate:
		initialModel, err = models.GetProductBySku(param.(string))
		if err != nil {
			log.Printf("Product Not Found!")
			os.Exit(0)
		}

		categories := make(map[int]string)
		categoriesObj, err := models.GetCategories()
		if err != nil {
			log.Printf("Error getting categories")
			os.Exit(1)
		}
		for _, category := range categoriesObj {
			categories[category.Id] = category.Name
		}
	}

	model := initialModel.(models.Product)
	if _, err := tea.NewProgram(modes.ProductModel(&model)).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(0)
	}
}
