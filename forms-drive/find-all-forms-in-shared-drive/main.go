package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

	// --- 2. Authentication (ADC) ---
	client := must(google.DefaultClient(ctx,
		drive.DriveReadonlyScope,
		forms.FormsBodyReadonlyScope,
		forms.FormsResponsesReadonlyScope,
	))

	formFiles := map[string]*drive.File{}
	{
		driveService := must(drive.NewService(ctx, option.WithHTTPClient(client)))
		if formID == "" {
			query := []string{"mimeType='application/vnd.google-apps.form'", "trashed=false"}
			if folderID != "" {
				query = append(query, fmt.Sprintf("'%s' in parents", folderID))
			}
			if filename != "" {
				query = append(query, fmt.Sprintf("name contains '%s'", filename))
			}
			must(0, driveService.Files.List().
				Q(strings.Join(query, " and ")).
				SupportsAllDrives(true).                  // ★ 共有ドライブの検索に必須
				IncludeItemsFromAllDrives(true).          // ★ 共有ドライブの検索に必須
				Fields("nextPageToken, files(id, name)"). // ID と 名前 を取得
				Pages(ctx, func(fileList *drive.FileList) error {
					for _, file := range fileList.Files {
						formFiles[file.Id] = file
					}
					return nil
				}))
			for _, file := range formFiles {
				must(0, driveService.Permissions.List(file.Id).
					SupportsAllDrives(true).
					IncludePermissionsForView("published").
					Pages(ctx, func(p *drive.PermissionList) error {
						formFiles[file.Id].Permissions = p.Permissions
						return nil
					}))
			}
		} else {
			file := must(driveService.Files.Get(formID).
				SupportsAllDrives(true).
				IncludePermissionsForView("published").
				Do())
			formFiles[formID] = file

			must(0, driveService.Permissions.List(file.Id).
				SupportsAllDrives(true).
				Pages(ctx, func(p *drive.PermissionList) error {
					formFiles[file.Id].Permissions = p.Permissions
					return nil
				}))
		}
	}
	{
		formsService := must(forms.NewService(ctx, option.WithHTTPClient(client)))
		type FormInfo struct {
			File  *drive.File
			Forms *forms.Form
		}
		var formInfo []FormInfo
		for _, formFile := range formFiles {
			form := must(formsService.Forms.Get(formFile.Id).Do())
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
