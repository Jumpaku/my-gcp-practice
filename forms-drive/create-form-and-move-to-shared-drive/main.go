package main

import (
	"context"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

// ★★★ Edit here... ★★★
const (
	// The folder ID in the shared drive where the form will be moved
	// This is the part of the URL (folders/...) when you open the folder in your browser
	TARGET_FOLDER_ID = "YOUR_SHARED_DRIVE_FOLDER_ID_HERE"

	// Title of the form to be created
	NEW_FORM_TITLE = "Form created by service account (for shared drive storage)"
)

// ★★★ End of section to edit ★★★

func main() {
	ctx := context.Background()

	// --- 1. Authentication (ADC and scope settings) ---
	// Use Application Default Credentials (ADC)
	// Request permissions (scopes) for both Forms API and Drive API
	client, err := google.DefaultClient(ctx,
		forms.FormsBodyScope, // Forms (read/write)
		drive.DriveScope,     // Drive (full access, required for moving files)
	)
	if err != nil {
		log.Fatalf("Failed to create authentication client: %v\n (Is the GOOGLE_APPLICATION_CREDENTIALS environment variable set?)", err)
	}

	// --- 2. Initialize each API service ---
	formsService, err := forms.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to initialize Forms service: %v", err)
	}

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to initialize Drive service: %v", err)
	}

	log.Println("Authentication and service initialization completed.")

	// --- 3. Step 1: Create a form with the Forms API ---
	log.Printf("Creating '%s'...", NEW_FORM_TITLE)

	newForm := &forms.Form{
		Info: &forms.Info{
			Description:   "",
			DocumentTitle: "",
			Title:         NEW_FORM_TITLE,
		},
		Settings: &forms.FormSettings{
			EmailCollectionType: "",
			QuizSettings:        nil,
		},
	}
	createdForm, err := formsService.Forms.Create(newForm).Do()
	if err != nil {
		log.Fatalf("Failed to create form with Forms API: %v", err)
	}
	formId := createdForm.FormId
	log.Printf("Form created (ID: %s)", formId)

	// --- 4. Step 2: Move the form with the Drive API ---
	log.Printf("Moving form to folder (ID: %s)...", TARGET_FOLDER_ID)

	// 4a. Get the current parent (My Drive)
	file, err := driveService.Files.Get(formId).Fields("parents").Do()
	if err != nil {
		log.Fatalf("Failed to get current parent (My Drive) of form: %v", err)
	}
	if len(file.Parents) == 0 {
		log.Fatalf("No parent found for the form.")
	}
	previousParent := file.Parents[0]

	// 4b. Update the file (add and remove parent = move)
	_, err = driveService.Files.Update(formId, nil).
		AddParents(TARGET_FOLDER_ID).  // Add destination folder (in shared drive)
		RemoveParents(previousParent). // Remove source parent (My Drive)
		SupportsAllDrives(true).       // Required for shared drive operations
		Do()
	if err != nil {
		log.Fatalf("Failed to move form with Drive API: %v\n (Does the service account have 'Content Manager' permission on the shared drive?)", err)
	}

	log.Printf("--- Done ---")
	log.Printf("Form (ID: %s) was successfully moved to folder (ID: %s).", formId, TARGET_FOLDER_ID)
}
