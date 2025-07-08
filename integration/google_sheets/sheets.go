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

const (
	incomeWRRange = "Finance!B2"
	expenseWRange = "Finance!D2"
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

		prepareSheets()
	})
	return err
}

func prepareSheets() error {
	ctx := context.Background()
	values := [][]interface{}{
		{"0"},
	}
	_, _ = clientInstance.Spreadsheets.Values.Update(spreadsheetId, incomeWRRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()

	_, _ = clientInstance.Spreadsheets.Values.Update(spreadsheetId, expenseWRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()

	return nil
}

func AdjustBalance(operationType string, amountString string) error {
	ctx := context.Background()

	var readRange string

	switch operationType {
	case "income":
		readRange = incomeWRRange
	case "expense":
		readRange = expenseWRange
	}

	response, err := clientInstance.Spreadsheets.Values.Get(spreadsheetId, readRange).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("ошибка чтения из Google Sheets: %v", err)
	}

	currentSum, err := strconv.ParseFloat(fmt.Sprint(response.Values[0][0]), 64)
	if err != nil {
		return fmt.Errorf("неправильный формат текущей суммы в таблице: %v", err)
	}

	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return fmt.Errorf("неправильный формат суммы: %v", err)
	}

	var result float64 = 0

	switch operationType {
	case "income":
		result = currentSum + amount
	case "expense":
		result = currentSum - amount
	}

	resultString := strconv.FormatFloat(result, 'f', -1, 64)

	values := [][]interface{}{
		{resultString},
	}

	_, err = clientInstance.Spreadsheets.Values.Update(spreadsheetId, readRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()

	if err != nil {
		log.Printf("Error output to Google Sheets: %v", err)
		return fmt.Errorf("error output to Google-Sheets")
	}

	return nil
}
