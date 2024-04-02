package thp

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheets struct {
	SpreadsSvc *sheets.SpreadsheetsService
}

func MustGoogleAuthenticator(creds []byte) *GoogleSheets {
	ctx := context.Background()

	client, err := sheets.NewService(ctx, option.WithCredentialsJSON(creds))
	if err != nil {
		log.Fatal(err)
	}

	return &GoogleSheets{
		SpreadsSvc: sheets.NewSpreadsheetsService(client),
	}
}
