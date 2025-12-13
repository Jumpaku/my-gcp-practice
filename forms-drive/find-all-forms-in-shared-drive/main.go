package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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
	var folderID, filename, formID string
	flag.StringVar(&folderID, "folder", "", "Target shared drive folder ID")
	flag.StringVar(&filename, "filename", "", "Filename of the form created by service account (for shared drive storage)")
	flag.StringVar(&formID, "form", "", "Target Form ID")
	flag.Parse()
	if folderID == "" && formID == "" {
		flag.Usage()
		log.Panicf("Error: one of flags (-folder or -form) is required.")
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Panicf("Error: GOOGLE_APPLICATION_CREDENTIALS environment variable is not set.")
	}

	ctx := context.Background()
	client := must(google.DefaultClient(ctx,
		drive.DriveReadonlyScope,
		forms.FormsBodyReadonlyScope,
		forms.FormsResponsesReadonlyScope,
	))
	fs := drivefs.New(must(drive.NewService(ctx, option.WithHTTPClient(client))))

	formFiles := []drivefs.FileInfo{}
	{
		if formID == "" {
			query := []string{"mimeType='application/vnd.google-apps.form'", "trashed=false"}
			if folderID != "" {
				query = append(query, fmt.Sprintf("'%s' in parents", folderID))
			}
			if filename != "" {
				query = append(query, fmt.Sprintf("name contains '%s'", filename))
			}
			formFiles = must(fs.Query(strings.Join(query, " and ")))
		} else {
			formFiles = []drivefs.FileInfo{must(fs.Info(drivefs.FileID(formID)))}
		}
	}
	{
		formsService := must(forms.NewService(ctx, option.WithHTTPClient(client)))
		type FormInfo struct {
			File  drivefs.FileInfo
			Forms *forms.Form
		}
		var formInfo []FormInfo
		for _, formFile := range formFiles {
			form := must(formsService.Forms.Get(string(formFile.ID)).Do())
			formInfo = append(formInfo, FormInfo{
				File:  formFile,
				Forms: form,
			})
		}

		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		_ = e.Encode(formInfo)
	}
}
