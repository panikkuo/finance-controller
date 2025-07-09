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

var (
	carcaseRange               = "!A1:D4"
	incomeWRRange              = "!B2"
	expenseWRange              = "!D2"
	incomeCategoryCol          = "!B"
	incomeCategoryColText      = "!A"
	incomeCategoryRow          = 4
	expenseCategoryCol         = "!D"
	expenseCategoryColText     = "!C"
	expenseCategoryRow         = 4
	expenseCategoryToRowCoords = make(map[string]string)
	incomeCategoryToRowCoords  = make(map[string]string)
)

func InitClient(credentialsJSON []byte, _spreadsheetId string, sheetName string) error {
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

		carcaseRange = sheetName + carcaseRange
		incomeWRRange = sheetName + incomeWRRange
		expenseWRange = sheetName + expenseWRange
		incomeCategoryCol = sheetName + incomeCategoryCol
		incomeCategoryColText = sheetName + incomeCategoryColText
		expenseCategoryCol = sheetName + expenseCategoryCol
		expenseCategoryColText = sheetName + expenseCategoryColText

		prepareSheets()
	})
	return err
}

func prepareSheets() error {
	ctx := context.Background()

	carcase := [][]interface{}{
		{nil, "Доходы", nil, "Расходы"},
		{"Total:", "0", "Total:", "0"},
		{nil, "Категория", nil, "Категория"},
		{nil, "0", nil, "0"},
	}

	_, _ = clientInstance.Spreadsheets.Values.Update(spreadsheetId, carcaseRange, &sheets.ValueRange{
		Values: carcase,
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()

	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          0,
					StartRowIndex:    1,
					EndRowIndex:      2,
					StartColumnIndex: 0,
					EndColumnIndex:   1,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: &sheets.Color{
							Red:   0.0,
							Green: 0.8,
							Blue:  0.3,
						},
					},
				},
				Fields: "userEnteredFormat.backgroundColor",
			},
		},
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          0,
					StartRowIndex:    1,
					EndRowIndex:      2,
					StartColumnIndex: 2,
					EndColumnIndex:   3,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: &sheets.Color{
							Red:   0.8,
							Green: 0.0,
							Blue:  0.3,
						},
					},
				},
				Fields: "userEnteredFormat.backgroundColor",
			},
		},
	}

	_, err := clientInstance.Spreadsheets.BatchUpdate(spreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		log.Fatalf("Не удалось обновить стили: %v", err)
	}

	return nil
}

func floatToStyleStringFormat(value float64) string { //0,1 -> 0.1
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func addToCell(amount float64, addCell string) error {
	ctx := context.Background()

	response, err := clientInstance.Spreadsheets.Values.Get(spreadsheetId, addCell).Context(ctx).Do()

	if err != nil {
		return fmt.Errorf("can't read from Google Sheets: %v", err)
	}

	currentAmount, err := strconv.ParseFloat(fmt.Sprint(response.Values[0][0]), 64)
	if err != nil {
		return fmt.Errorf("error format to cell %v: %v", addCell, err)
	}

	result := currentAmount + amount

	resultString := floatToStyleStringFormat(result)

	values := [][]interface{}{
		{resultString},
	}

	_, err = clientInstance.Spreadsheets.Values.Update(spreadsheetId, addCell, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()

	if err != nil {
		log.Printf("Error output to Google Sheets: %v", err)
		return fmt.Errorf("error output to Google-Sheets")
	}

	return nil
}

func getCategoryCell(category string, categoryMap *map[string]string, lastCol *string, lastRow *int, textCol *string) (string, string) {
	row, ok := (*categoryMap)[category]
	if !ok {
		lastRowString := strconv.Itoa(*lastRow)
		(*lastRow)++
		(*categoryMap)[category] = lastRowString
		return *lastCol + lastRowString, *textCol + lastRowString
	}
	return *lastCol + row, *textCol + row
}

func outTextToCell(text string, cell string) error {
	ctx := context.Background()

	values := [][]interface{}{
		{text},
	}

	_, err := clientInstance.Spreadsheets.Values.Update(spreadsheetId, cell, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()

	if err != nil {
		log.Printf("Error output to Google Sheets: %v", err)
		return fmt.Errorf("error output to Google-Sheets")
	}

	return nil
}

func AdjustBalance(operationType string, amountString string, category string) error {
	if category == "" {
		category = "Другое"
	}

	var addTotalCell string
	var addCategoryCell string
	var textCategoryCell string

	switch operationType {
	case "income":
		addTotalCell = incomeWRRange
		addCategoryCell, textCategoryCell = getCategoryCell(category, &incomeCategoryToRowCoords, &incomeCategoryCol, &incomeCategoryRow, &incomeCategoryColText)
	case "expense":
		addTotalCell = expenseWRange
		addCategoryCell, textCategoryCell = getCategoryCell(category, &expenseCategoryToRowCoords, &expenseCategoryCol, &expenseCategoryRow, &expenseCategoryColText)
	}

	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return fmt.Errorf("неправильный формат суммы: %v", err)
	}

	err = addToCell(amount, addTotalCell)

	if err != nil {
		return err
	}

	err = addToCell(amount, addCategoryCell)

	if err != nil {
		return err
	}

	err = outTextToCell(category, textCategoryCell)

	if err != nil {
		return err
	}

	return nil
}
