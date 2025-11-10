package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

// ★★★ Edit here... ★★★
const (
	// 1. The exact title of the form you want to search for
	FORM_TITLE = "(Enter the exact title of the form here)"

	// 2. The ID of the folder where the form is stored
	TARGET_FOLDER_ID = "(Enter the folder ID where the form is here)"
)

// ★★★ End of section to edit ★★★

func main() {
	ctx := context.Background()

	// --- 1. Authentication (ADC) ---
	client, err := google.DefaultClient(ctx,
		drive.DriveReadonlyScope,
		forms.FormsBodyReadonlyScope,
		forms.FormsResponsesReadonlyScope,
	)
	if err != nil {
		log.Fatalf("Failed to create authentication client: %v\n (Is the GOOGLE_APPLICATION_CREDENTIALS environment variable set?)", err)
	}

	// --- 2. Initialize each API service ---
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to initialize Drive service: %v", err)
	}
	formsService, err := forms.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to initialize Forms service: %v", err)
	}
	log.Println("Authentication complete. Starting form search...")

	// --- 3. Step 1: Search for the form using the Drive API ---
	formId, err := findFormInSharedDrive(driveService, FORM_TITLE, TARGET_FOLDER_ID)
	if err != nil {
		log.Fatalf("Failed to search for form: %v", err)
	}
	log.Printf("Form found (ID: %s)", formId)

	// --- 4. Step 2: Retrieve question details using the Forms API ---
	questionMap, err := getFormQuestions(formsService, formId)
	if err != nil {
		log.Fatalf("Failed to retrieve form questions: %v", err)
	}

	// --- 5. Step 3: Retrieve all responses from all pages using the Forms API ---
	var allResponses []*forms.FormResponse
	var pageToken string

	log.Println("Retrieving all response pages from the form...")

	for {
		resp, err := formsService.Forms.Responses.List(formId).PageToken(pageToken).Do()
		if err != nil {
			log.Fatalf("Failed to retrieve form responses (PageToken: %s): %v", pageToken, err)
		}

		if len(resp.Responses) > 0 {
			allResponses = append(allResponses, resp.Responses...)
		}

		// Get the next page token
		pageToken = resp.NextPageToken
		if pageToken == "" {
			// Exit loop if there are no more pages
			break
		}
		log.Printf("%d responses retrieved, loading next page...", len(allResponses))
	}

	log.Printf("All %d responses retrieved. Starting aggregation...", len(allResponses))

	// --- 6. Step 4: Aggregate in Go program ---
	// Structure: map[QuestionID] -> map[AnswerValue] -> Count
	aggregation := make(map[string]map[string]int)

	for _, response := range allResponses {
		for questionId, answer := range response.Answers {
			if answer.TextAnswers == nil {
				continue
			}
			if _, ok := aggregation[questionId]; !ok {
				aggregation[questionId] = make(map[string]int)
			}
			for _, textAnswer := range answer.TextAnswers.Answers {
				aggregation[questionId][textAnswer.Value]++
			}
		}
	}

	// --- 7. Step 5: Display aggregation results ---
	log.Println("=============================")
	log.Println("       Aggregation Results")
	log.Println("=============================")
	for qId, counts := range aggregation {
		questionTitle, ok := questionMap[qId]
		if !ok {
			questionTitle = qId
		}
		log.Printf("\n▼ Question: %s (%s)\n", questionTitle, qId)
		for answerValue, count := range counts {
			log.Printf("  - '%s': %d responses\n", answerValue, count)
		}
	}
}

// (findFormInSharedDrive 関数 - 変更なし)
func findFormInSharedDrive(driveService *drive.Service, title string, folderId string) (string, error) {
	// ... (前回のコードと同じ) ...
	query := fmt.Sprintf(
		"mimeType='application/vnd.google-apps.form' and name='%s' and '%s' in parents and trashed=false",
		title, folderId,
	)
	fileList, err := driveService.Files.List().
		Q(query).
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Fields("files(id)").
		Do()
	if err != nil {
		return "", fmt.Errorf("Failed to search with Drive API: %v", err)
	}
	if len(fileList.Files) == 0 {
		return "", errors.New("No form found matching the specified criteria.")
	}
	if len(fileList.Files) > 1 {
		log.Printf("Warning: Multiple forms matched the criteria. Using the first one (ID: %s).", fileList.Files[0].Id)
	}
	return fileList.Files[0].Id, nil
}

// (getFormQuestions 関数 - 変更なし)
func getFormQuestions(formsService *forms.Service, formId string) (map[string]string, error) {
	// ... (前回のコードと同じ) ...
	form, err := formsService.Forms.Get(formId).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve form with Forms API: %v", err)
	}
	qMap := make(map[string]string)
	for _, item := range form.Items {
		if item.QuestionItem != nil && item.QuestionItem.Question != nil {
			qMap[item.QuestionItem.Question.QuestionId] = item.Title
		}
	}
	return qMap, nil
}
