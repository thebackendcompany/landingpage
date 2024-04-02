package emailleads

import (
	"context"
	"log"
	"thebackendcompany/pkg/thp"
	"time"

	"google.golang.org/api/sheets/v4"
)

type EmailLeadsSheets struct {
	svc      *thp.GoogleSheets
	sheetID  string
	colRange string
}

func NewEmailLeadsSheets(cred []byte, sheetID string) *EmailLeadsSheets {
	return &EmailLeadsSheets{
		svc:      thp.MustGoogleAuthenticator(cred),
		sheetID:  sheetID,
		colRange: "Sheet1!A:C",
	}
}

func (gs *EmailLeadsSheets) All(ctx context.Context) ([][]any, error) {
	handler := gs.svc.SpreadsSvc.Values.Get(gs.sheetID, gs.colRange)
	resp, err := handler.Do()
	if err != nil {
		log.Printf("gcloud read err %s %+v\n", gs.sheetID, err)
		return nil, err
	}

	return resp.Values, nil
}

func (gs *EmailLeadsSheets) Write(ctx context.Context, lead *SheetsEmailLeads) (int64, error) {
	values := sheets.ValueRange{
		Range:          gs.colRange,
		MajorDimension: "ROWS",
		Values: [][]interface{}{
			{lead.Email, lead.Body, lead.ContactedAt.Format(time.RFC3339)}, // %Y-%m-%d
		}}

	handler := gs.svc.SpreadsSvc.Values.Append(gs.sheetID, gs.colRange, &values).ValueInputOption("RAW")
	resp, err := handler.Do()

	if err != nil {
		log.Printf("gcloud write err %s %+v\n", gs.sheetID, err)
		return -1, err
	}

	return resp.Updates.UpdatedRows, nil
}
