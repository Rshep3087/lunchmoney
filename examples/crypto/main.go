package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/icco/lunchmoney"
)

func main() {
	apiKey := os.Getenv("LUNCHMONEY_API_KEY")
	if apiKey == "" {
		log.Fatal("LUNCHMONEY_API_KEY environment variable is required")
	}

	client, err := lunchmoney.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Get all crypto assets
	cryptoAssets, err := client.GetCrypto(ctx)
	if err != nil {
		log.Fatalf("Failed to get crypto assets: %v", err)
	}

	fmt.Printf("Found %d crypto assets:\n", len(cryptoAssets))
	for _, crypto := range cryptoAssets {
		fmt.Printf("  - %s: %s %s (Source: %s, Status: %s)\n",
			crypto.Name,
			crypto.Balance,
			crypto.Currency,
			crypto.Source,
			crypto.Status,
		)
		if crypto.InstitutionName != nil {
			fmt.Printf("    Institution: %s\n", *crypto.InstitutionName)
		}
		if crypto.DisplayName != nil {
			fmt.Printf("    Display Name: %s\n", *crypto.DisplayName)
		}
		if crypto.ToBase != nil {
			fmt.Printf("    Value in base currency: %.2f\n", *crypto.ToBase)
		}
		fmt.Println()
	}
}
