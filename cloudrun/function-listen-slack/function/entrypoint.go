package function

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type SlackRequest struct {
	Type      string          `json:"type"`
	Challenge string          `json:"challenge"`
	Event     json.RawMessage `json:"event"`
}

func HandleSlackEvent(w http.ResponseWriter, r *http.Request) {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if signingSecret == "" {
		log.Println("Error: SLACK_SIGNING_SECRET is not set")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !isValidSignature(w, r, signingSecret) {
		return
	}

	var req SlackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("json decode error: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Type == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(req.Challenge))
		return
	}

	log.Printf("Verified Event received: %v", string(req.Event))
	w.WriteHeader(http.StatusOK)
}

func isValidSignature(w http.ResponseWriter, r *http.Request, secret string) bool {
	slackSignature := r.Header.Get("X-Slack-Signature")
	slackTimestamp := r.Header.Get("X-Slack-Request-Timestamp")

	if slackSignature == "" || slackTimestamp == "" {
		http.Error(w, "Missing Slack headers", http.StatusUnauthorized)
		return false
	}

	ts, err := strconv.ParseInt(slackTimestamp, 10, 64)
	if err != nil {
		http.Error(w, "Invalid timestamp", http.StatusUnauthorized)
		return false
	}
	if time.Now().Unix()-ts > 60*5 {
		http.Error(w, "Request too old", http.StatusUnauthorized)
		return false
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return false
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	baseString := fmt.Sprintf("v0:%s:%s", slackTimestamp, string(bodyBytes))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(baseString))
	expectedSignature := "v0=" + hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedSignature), []byte(slackSignature)) {
		log.Printf("Signature mismatch! Expected: %s, Got: %s", expectedSignature, slackSignature)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return false
	}

	return true
}
