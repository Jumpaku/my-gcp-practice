package main

import (
	"context"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

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
	var sheetID int64
	flag.StringVar(&spreadsheetID, "spreadsheet", "", "ID of the spreadsheet (Required)")
	flag.Int64Var(&sheetID, "sheet", 0, "ID of the sheet in the spreadsheet (Required)")
	flag.Parse()
	if spreadsheetID == "" || sheetID == 0 {
		flag.Usage()
		log.Panicf("Error: -spreadsheet and -sheet are required.")
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Panicf("Error: GOOGLE_APPLICATION_CREDENTIALS environment variable is not set.")
	}

	data := must(csv.NewReader(os.Stdin).ReadAll())

	ctx := context.Background()
	client := must(google.DefaultClient(ctx,
		sheets.SpreadsheetsScope,
		drive.DriveScope,
	))
	sheetsService := must(sheets.NewService(ctx, option.WithHTTPClient(client)))

	spreadsheet := must(sheetsService.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Do())
	{
		var sheetExists bool
		for _, sheet := range spreadsheet.Sheets {
			if sheet.Properties.SheetId == sheetID {
				sheetExists = true
			}
		}
		if !sheetExists {
			must(sheetsService.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
				Requests: []*sheets.Request{
					{
						AddSheet: &sheets.AddSheetRequest{
							Properties: &sheets.SheetProperties{
								Title:   strconv.FormatInt(sheetID, 10),
								SheetId: sheetID,
							},
						},
					},
				},
			}).Do())
		}

	}

	rows := []*sheets.RowData{}
	{
		for _, record := range data {
			values := []*sheets.CellData{}
			for _, value := range record {
				if strings.HasPrefix(value, "=") {
					values = append(values, &sheets.CellData{
						UserEnteredValue: &sheets.ExtendedValue{FormulaValue: &value},
					})
				} else {
					v := value[1:]
					values = append(values, &sheets.CellData{
						UserEnteredValue: &sheets.ExtendedValue{StringValue: &v},
					})
				}

			}
			rows = append(rows, &sheets.RowData{
				Values: values,
			})
		}
	}
	{
		must(sheetsService.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{
					UpdateCells: &sheets.UpdateCellsRequest{
						Fields: "*",
						Start:  &sheets.GridCoordinate{SheetId: sheetID},
						Rows:   rows,
					},
				},
			},
		}).Do())
	}
}
