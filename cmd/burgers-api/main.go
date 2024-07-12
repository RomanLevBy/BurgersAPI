package main

import (
	"context"
	"github.com/RomanLevBy/BurgersAPI/internal/app"
	"log"
)

func main() {
	a, err := app.NewApp(context.Background())
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run(context.Background())
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}
