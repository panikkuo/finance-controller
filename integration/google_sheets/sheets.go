package gsheets

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	clientInstance *sheets.Service
	spreadsheetId  string
	once           sync.Once
)

func InitClient(credentialsJSON []byte, _spreadsheetId string) error {
	var err error
	ctx := context.Background()
	once.Do(func() {
		spreadsheetId = _spreadsheetId
		config, e := google.JWTConfigFromJSON(credentialsJSON, sheets.SpreadsheetsScope)
		if e != nil {
			err = e
			return
		}
		httpClient := config.Client(ctx)
		clientInstance, err = sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	})
	return err
}

func AdjustBalance(operationType string, totalSum string) error {
	ctx := context.Background()

	writeRange := "Finance!A2:B2"

	sum, err := strconv.ParseFloat(totalSum, 64)
	if err != nil {
		return fmt.Errorf("неправильный формат суммы: %v", err)
	}

	values := [][]interface{}{
		{operationType, sum},
	}
	_, err = clientInstance.Spreadsheets.Values.Update(spreadsheetId, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()

	if err != nil {
		log.Printf("Error output to Google Sheets: %v", err)
		return fmt.Errorf("error output to Google-Sheets")
	}

	return nil
}
