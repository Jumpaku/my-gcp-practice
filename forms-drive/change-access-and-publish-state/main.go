package main

import (
	"context"
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
	fs := drivefs.New(must(drive.NewService(ctx, option.WithHTTPClient(client))))

	{
		perms := must(fs.PermList(drivefs.FileID(formId)))
		for _, perm := range perms {
			switch grantee := perm.Grantee().(type) {
			case drivefs.GranteeDomain, drivefs.GranteeAnyone:
				must(fs.PermDel(drivefs.FileID(formId), grantee))
				must(fs.PermDel(drivefs.FileID(formId), grantee))
			}
		}
		switch access {
		case "limited":
		case "domain":
			must(fs.PermSet(drivefs.FileID(formId), drivefs.DomainPermission(domain, drivefs.RoleReader, false)))
		case "anyone":
			must(fs.PermSet(drivefs.FileID(formId), drivefs.AnyonePermission(drivefs.RoleReader, false)))
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
