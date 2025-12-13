package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func must[T any](v T, err error) T {
	if err != nil {
		log.Panic(err)
	}
	return v

}

func main() {
	var spreadsheetID string
	flag.StringVar(&spreadsheetID, "spreadsheet", "", "ID of the spreadsheet (Required)")
	flag.Parse()
	if spreadsheetID == "" {
		flag.Usage()
		log.Panicf("Error: All flags (-source, -destination, -filename) are required.")
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Panicf("Error: GOOGLE_APPLICATION_CREDENTIALS environment variable is not set.")
	}

	ctx := context.Background()
	client := must(google.DefaultClient(ctx,
		sheets.SpreadsheetsScope,
		drive.DriveScope,
	))

	type Sheet struct {
		SheetId       int64
		Title         string
		FormattedData [][]string
	}
	type Spreadsheet struct {
		SpreadsheetId string
		Sheets        []Sheet
	}
	sheetsService := must(sheets.NewService(ctx, option.WithHTTPClient(client)))
	spreadsheet := must(sheetsService.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Do())

	s := Spreadsheet{SpreadsheetId: spreadsheet.SpreadsheetId}
	for _, sheet := range spreadsheet.Sheets {
		data := [][]string{}
		if len(sheet.Data) == 0 {
			continue
		}
		for _, r := range sheet.Data[0].RowData {
			row := []string{}
			for _, v := range r.Values {
				row = append(row, v.FormattedValue)
			}
			data = append(data, row)
		}
		s.Sheets = append(s.Sheets, Sheet{
			SheetId:       sheet.Properties.SheetId,
			Title:         sheet.Properties.Title,
			FormattedData: data,
		})
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	_ = e.Encode(s)
}
