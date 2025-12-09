package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Jumpaku/go-drivefs"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

func must[T any](v T, err error) T {
	if err != nil {
		log.Panic(err)
	}
	return v

}

func main() {
	// --- 1. Handle command-line arguments ---
	var sourceFormID, destinationFolderID, filename string
	flag.StringVar(&sourceFormID, "source", "", "ID of the source form to be copied (Required)")
	flag.StringVar(&destinationFolderID, "destination", "", "ID of the target folder (Required)")
	flag.StringVar(&filename, "filename", "", "New name for the copied form (Required)")

	flag.Parse()
	if sourceFormID == "" || destinationFolderID == "" || filename == "" {
		flag.Usage()
		log.Panicf("Error: All flags (-source, -destination, -filename) are required.")
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Panicf("Error: GOOGLE_APPLICATION_CREDENTIALS environment variable is not set.")
	}

	ctx := context.Background()

	// --- 2. Authentication (ADC) ---
	client := must(google.DefaultClient(ctx,
		forms.FormsBodyScope, // Forms (read/write)
		drive.DriveScope,     // Drive (full access, required for moving files)
	))

	var clonedFormID string
	{
		fs := drivefs.New(must(drive.NewService(ctx, option.WithHTTPClient(client))))
		copiedFile := must(fs.Copy(drivefs.FileID(destinationFolderID), drivefs.FileID(sourceFormID), filename))
		clonedFormID = string(copiedFile.ID)
	}
	{
		formsService := must(forms.NewService(ctx, option.WithHTTPClient(client)))
		form := must(formsService.Forms.Get(clonedFormID).Do())

		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		_ = e.Encode(form)
	}
	return

}
