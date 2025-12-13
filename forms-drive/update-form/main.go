package main

import (
	"context"
	"encoding/json"
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
	var formId, title, description string
	flag.StringVar(&formId, "form", "", "Form ID (Required)")
	flag.StringVar(&title, "title", "", "Form Title")
	flag.StringVar(&description, "description", "", "Form Description")
	flag.Parse()
	if formId == "" || title == "" || description == "" {
		flag.Usage()
		log.Panicf("Error: All flags (-form, -title, -description) are required.")
	}
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Panicf("Error: GOOGLE_APPLICATION_CREDENTIALS environment variable is not set.")
	}

	ctx := context.Background()
	client := must(google.DefaultClient(ctx,
		forms.FormsBodyScope,
		drive.DriveScope,
	))

	formsService := must(forms.NewService(ctx, option.WithHTTPClient(client)))
	form := must(formsService.Forms.Get(formId).Do())

	updates := []*forms.Request{}
	if title != "" {
		updates = append(updates, &forms.Request{
			UpdateFormInfo: &forms.UpdateFormInfoRequest{
				Info:       &forms.Info{Title: title},
				UpdateMask: "title",
			},
		})
	}
	if description != "" {
		updates = append(updates, &forms.Request{
			UpdateFormInfo: &forms.UpdateFormInfoRequest{
				Info:       &forms.Info{Description: description},
				UpdateMask: "description",
			},
		})
	}

	res := must(formsService.Forms.BatchUpdate(form.FormId, &forms.BatchUpdateFormRequest{Requests: updates}).Do())
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	_ = e.Encode(res.Form)
}
