package main

import (
	"context"
	"fmt"
	"github.com/RomanLevBy/BurgersAPI/internal/app"
	"log"
)

func main() {
	a, err := app.NewApp(context.Background())
	if err != nil {
		log.Fatalf("Failed to init app: %s", err.Error())
	}

	err = a.Run(context.Background())
	if err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}

	fmt.Println("Hello burger API")
}
