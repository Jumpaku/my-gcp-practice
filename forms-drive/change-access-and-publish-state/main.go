package main

import (
	"context"
	"flag"
	"log"
	"os"

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
	var formId, publishState, access, domain string
	flag.StringVar(&formId, "form", "", "Form ID (Required)")
	flag.StringVar(&publishState, "publish", "", "Set publish state to 'not_published', 'accepting', 'not_accepting' (Required)")
	flag.StringVar(&access, "access", "", "Set responder to 'limited', 'domain', or 'anyone' (Required)")
	flag.StringVar(&domain, "domain", "", "Responder domain (Required if -access=domain is set)")

	flag.Parse()
	if formId == "" ||
		(publishState != "not_published" && publishState != "accepting" && publishState != "not_accepting") ||
		(access != "limited" && access != "domain" && access != "anyone") ||
		(access == "domain" && domain == "") {
		flag.Usage()
		log.Panicln("Error: All flags (-form, -publish, -responder) are required.")
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Panicf("Error: GOOGLE_APPLICATION_CREDENTIALS environment variable is not set.")
	}

	ctx := context.Background()

	client := must(google.DefaultClient(ctx,
		forms.FormsBodyScope,
		drive.DriveScope,
	))
	driveService := must(drive.NewService(ctx, option.WithHTTPClient(client)))
	{
		var permissions []*drive.Permission
		must(0, driveService.Permissions.
			List(formId).
			SupportsAllDrives(true).
			Pages(ctx, func(permission *drive.PermissionList) error {
				permissions = append(permissions, permission.Permissions...)
				return nil
			}))
		for _, permission := range permissions {
			if permission.Type == "anyone" && permission.Role == "reader" {
				must(0, driveService.Permissions.Delete(formId, permission.Id).SupportsAllDrives(true).Do())
			}
			if permission.Type == "domain" && permission.Role == "reader" {
				must(0, driveService.Permissions.Delete(formId, permission.Id).SupportsAllDrives(true).Do())
			}
		}
	}
	{
		switch access {
		case "limited":
		case "domain":
			must(driveService.Permissions.Create(formId, &drive.Permission{
				Type:   "domain",
				Role:   "reader",
				Domain: domain,
			}).
				SupportsAllDrives(true).
				Do())
		case "anyone":
			must(driveService.Permissions.Create(formId, &drive.Permission{
				Type: "anyone",
				Role: "reader",
			}).
				SupportsAllDrives(true).
				Do())
		}
	}
	{
		formsService := must(forms.NewService(ctx, option.WithHTTPClient(client)))
		must(formsService.Forms.SetPublishSettings(formId, &forms.SetPublishSettingsRequest{
			PublishSettings: &forms.PublishSettings{
				PublishState: &forms.PublishState{
					IsAcceptingResponses: publishState == "accepting",
					IsPublished:          publishState != "not_published",
				},
			},
			UpdateMask: "publish_state",
		}).Do())
	}
}
