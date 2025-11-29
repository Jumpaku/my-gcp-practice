package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
	// --- 1. Handle command-line arguments ---
	var formId string
	flag.StringVar(&formId, "form", "", "Target shared drive folder ID")

	flag.Parse()
	if formId == "" {
		flag.Usage()
		log.Panicf("Error: All flags (-form) are required.")
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
	{
		formsService := must(forms.NewService(ctx, option.WithHTTPClient(client)))
		form := must(formsService.Forms.Get(formId).Do())

		var responses []*forms.FormResponse
		must(0, formsService.Forms.Responses.
			List(formId).
			Pages(ctx, func(resp *forms.ListFormResponsesResponse) error {
				responses = append(responses, resp.Responses...)
				return nil
			}))

		type AnswerText []string
		type Answer struct {
			CreateTime     string
			LastUpdateTime string
			AnswerTexts    map[string]AnswerText
		}
		type Result struct {
			Questions map[string]string
			Answers   []Answer
		}
		result := Result{Questions: map[string]string{}}
		for _, item := range form.Items {
			if item.QuestionItem != nil && item.QuestionItem.Question != nil {
				result.Questions[item.QuestionItem.Question.QuestionId] = item.Title
			}
		}
		for _, response := range responses {
			answer := Answer{
				CreateTime:     response.CreateTime,
				LastUpdateTime: response.LastSubmittedTime,
				AnswerTexts:    map[string]AnswerText{},
			}
			for questionId := range result.Questions {
				textAnswers := response.Answers[questionId].TextAnswers
				if textAnswers == nil {
					continue
				}

				answerTexts := AnswerText{}
				for _, textAnswer := range textAnswers.Answers {
					answerTexts = append(answerTexts, textAnswer.Value)
				}

				answer.AnswerTexts[questionId] = answerTexts
			}
			result.Answers = append(result.Answers, answer)
		}

		buf := bytes.NewBuffer(nil)
		e := json.NewEncoder(buf)
		e.SetIndent("", "  ")
		_ = e.Encode(result)
		fmt.Println(buf.String())
	}
}
