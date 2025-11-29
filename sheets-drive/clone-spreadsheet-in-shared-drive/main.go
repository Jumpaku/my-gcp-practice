package main

import (
	"context"
	"flag"
	"fmt"
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
	var sourceSpreadsheetID, destinationFolderID, filename string
	flag.StringVar(&sourceSpreadsheetID, "source", "", "ID of the source form to be copied (Required)")
	flag.StringVar(&destinationFolderID, "destination", "", "ID of the target folder (Required)")
	flag.StringVar(&filename, "filename", "", "New name for the copied form (Required)")

	flag.Parse()
	if sourceSpreadsheetID == "" || destinationFolderID == "" || filename == "" {
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

	{
		driveService := must(drive.NewService(ctx, option.WithHTTPClient(client)))
		copiedFile := must(driveService.Files.
			Copy(sourceSpreadsheetID, &drive.File{
				Name:    filename,
				Parents: []string{destinationFolderID},
			}).
			SupportsAllDrives(true).
			Do())
		fmt.Println(copiedFile.Id)
	}
}
