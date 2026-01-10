package function

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestHandleSlackEvent(t *testing.T) {
	testSecret := "test-signing-secret"
	_ = os.Setenv("SLACK_SIGNING_SECRET", testSecret)

	jsonBody := []byte(`{"type": "url_verification", "challenge": "success-challenge-123"}`)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sigBaseString := fmt.Sprintf("v0:%s:%s", timestamp, string(jsonBody))

	mac := hmac.New(sha256.New, []byte(testSecret))
	mac.Write([]byte(sigBaseString))
	signature := "v0=" + hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Slack-Request-Timestamp", timestamp)
	req.Header.Set("X-Slack-Signature", signature)

	w := httptest.NewRecorder()

	HandleSlackEvent(w, req)

	// 7. 結果確認
	resp := w.Result()
	fmt.Printf("\n--- Execution Result ---\n")
	fmt.Printf("Status Code : %d\n", resp.StatusCode)
	fmt.Printf("Response    : %s\n", w.Body.String())

	if resp.StatusCode == http.StatusOK && w.Body.String() == "success-challenge-123" {
		fmt.Println("✅ Success: Signature verified and challenge returned.")
	} else {
		fmt.Println("❌ Failed: Something went wrong.")
	}
}
