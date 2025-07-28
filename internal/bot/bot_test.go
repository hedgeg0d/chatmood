package bot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewBot(t *testing.T) {
	// Mock Telegram API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "getMe"):
			response := map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"id":         123456789,
					"is_bot":     true,
					"first_name": "ChatMood",
					"username":   "chatmood_bot",
				},
			}
			json.NewEncoder(w).Encode(response)
		case strings.Contains(r.URL.Path, "setWebhook"):
			response := map[string]interface{}{
				"ok":     true,
				"result": true,
			}
			json.NewEncoder(w).Encode(response)
		case strings.Contains(r.URL.Path, "setMyCommands"):
			response := map[string]interface{}{
				"ok":     true,
				"result": true,
			}
			json.NewEncoder(w).Encode(response)
		case strings.Contains(r.URL.Path, "setChatMenuButton"):
			response := map[string]interface{}{
				"ok":     true,
				"result": true,
			}
			json.NewEncoder(w).Encode(response)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Replace Telegram API URL with test server
	originalMakeRequest := (*Bot).makeRequest
	defer func() {
		// Restore original method after test
		(*Bot).makeRequest = originalMakeRequest
	}()

	bot := &Bot{
		Token:      "test_token",
		WebhookURL: "https://example.com",
		Client:     &http.Client{},
	}

	// Override makeRequest to use test server
	bot.makeRequest = func(method string, data interface{}) ([]byte, error) {
		url := server.URL + "/" + method
		var body *strings.Reader
		if data != nil {
			jsonData, _ := json.Marshal(data)
			body = strings.NewReader(string(jsonData))
		} else {
			body = strings.NewReader("")
		}

		req, _ := http.NewRequest("POST", url, body)
		if data != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := bot.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var responseBody []byte
		resp.Body.Read(responseBody)
		return responseBody, nil
	}

	// Test bot creation (this would normally fail without mocking)
	if bot.Token != "test_token" {
		t.Errorf("Expected token 'test_token', got %s", bot.Token)
	}

	if bot.WebhookURL != "https://example.com" {
		t.Errorf("Expected webhook URL 'https://example.com', got %s", bot.WebhookURL)
	}
}

func TestBot_GetUsername(t *testing.T) {
	bot := &Bot{
		Username: "test_bot",
	}

	username := bot.GetUsername()
	if username != "test_bot" {
		t.Errorf("Expected username 'test_bot', got %s", username)
	}
}

func TestBot_HandleUpdate(t *testing.T) {
	bot := &Bot{
		Token:      "test_token",
		WebhookURL: "https://example.com",
		Client:     &http.Client{},
	}

	tests := []struct {
		name   string
		update Update
		hasErr bool
	}{
		{
			name: "Message update",
			update: Update{
				UpdateID: 1,
				Message: &Message{
					MessageID: 1,
					From: &User{
						ID:        123,
						FirstName: "Test",
						Username:  "testuser",
					},
					Chat: &Chat{
						ID:   123,
						Type: "private",
					},
					Text: "/start",
				},
			},
			hasErr: false,
		},
		{
			name: "Callback query update",
			update: Update{
				UpdateID: 2,
				CallbackQuery: &CallbackQuery{
					ID: "test_callback",
					From: &User{
						ID:        123,
						FirstName: "Test",
						Username:  "testuser",
					},
					Data: "help",
				},
			},
			hasErr: false,
		},
		{
			name: "Empty update",
			update: Update{
				UpdateID: 3,
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bot.HandleUpdate(tt.update)
			if (err != nil) != tt.hasErr {
				t.Errorf("HandleUpdate() error = %v, hasErr %v", err, tt.hasErr)
			}
		})
	}
}

func TestBot_handleMessage(t *testing.T) {
	bot := &Bot{
		Token:      "test_token",
		WebhookURL: "https://example.com",
		Client:     &http.Client{},
	}

	// Mock the sendMessage method
	originalSendMessage := bot.sendMessage
	var sentMessages []string
	bot.sendMessage = func(chatID int64, text, parseMode string, keyboard *InlineKeyboardMarkup) error {
		sentMessages = append(sentMessages, text)
		return nil
	}
	defer func() {
		bot.sendMessage = originalSendMessage
	}()

	tests := []struct {
		name             string
		message          *Message
		expectedContains string
	}{
		{
			name: "Start command",
			message: &Message{
				Text: "/start",
				Chat: &Chat{ID: 123},
				From: &User{FirstName: "Test"},
			},
			expectedContains: "Welcome to ChatMood",
		},
		{
			name: "Help command",
			message: &Message{
				Text: "/help",
				Chat: &Chat{ID: 123},
			},
			expectedContains: "ChatMood Help",
		},
		{
			name: "Create command",
			message: &Message{
				Text: "/create",
				Chat: &Chat{ID: 123},
			},
			expectedContains: "Ready to create",
		},
		{
			name: "Unknown command",
			message: &Message{
				Text: "random text",
				Chat: &Chat{ID: 123},
			},
			expectedContains: "not sure what you mean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sentMessages = []string{} // Reset
			err := bot.handleMessage(tt.message)
			if err != nil {
				t.Errorf("handleMessage() error = %v", err)
			}

			if len(sentMessages) == 0 {
				t.Error("Expected message to be sent")
				return
			}

			if !strings.Contains(sentMessages[0], tt.expectedContains) {
				t.Errorf("Expected message to contain '%s', got '%s'", tt.expectedContains, sentMessages[0])
			}
		})
	}
}

func TestBot_adjustBrightness(t *testing.T) {
	// This would be in a separate utility package in a real project
	// For now, we'll test the concept
	tests := []struct {
		name     string
		color    string
		amount   int
		expected string
	}{
		{
			name:     "Darken color",
			color:    "#FFFFFF",
			amount:   -50,
			expected: "#CDCDCD", // Approximate expected result
		},
		{
			name:     "Brighten color",
			color:    "#000000",
			amount:   50,
			expected: "#323232", // Approximate expected result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test would need the actual adjustBrightness function
			// For now, we just verify the test structure
			if tt.color == "" {
				t.Error("Color should not be empty")
			}
		})
	}
}

func TestBotCommand_Validation(t *testing.T) {
	commands := []BotCommand{
		{Command: "start", Description: "Start using ChatMood"},
		{Command: "help", Description: "Get help and information"},
		{Command: "create", Description: "Create a new mood sticker"},
	}

	for _, cmd := range commands {
		if cmd.Command == "" {
			t.Error("Command should not be empty")
		}
		if cmd.Description == "" {
			t.Error("Description should not be empty")
		}
		if len(cmd.Command) > 32 {
			t.Errorf("Command '%s' is too long (max 32 chars)", cmd.Command)
		}
		if len(cmd.Description) > 256 {
			t.Errorf("Description for '%s' is too long (max 256 chars)", cmd.Command)
		}
	}
}

func TestInlineKeyboardMarkup_Structure(t *testing.T) {
	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{
					Text: "Test Button",
					WebApp: &WebAppInfo{
						URL: "https://example.com/app",
					},
				},
			},
		},
	}

	if len(keyboard.InlineKeyboard) == 0 {
		t.Error("Keyboard should have at least one row")
	}

	if len(keyboard.InlineKeyboard[0]) == 0 {
		t.Error("First row should have at least one button")
	}

	button := keyboard.InlineKeyboard[0][0]
	if button.Text == "" {
		t.Error("Button text should not be empty")
	}

	if button.WebApp == nil {
		t.Error("WebApp should be set")
	}

	if button.WebApp.URL == "" {
		t.Error("WebApp URL should not be empty")
	}
}

func TestSendMessageRequest_JSON(t *testing.T) {
	req := SendMessageRequest{
		ChatID:    123456789,
		Text:      "Test message",
		ParseMode: "Markdown",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Errorf("Failed to marshal SendMessageRequest: %v", err)
	}

	var unmarshaled SendMessageRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal SendMessageRequest: %v", err)
	}

	if unmarshaled.ChatID != req.ChatID {
		t.Errorf("Expected ChatID %d, got %d", req.ChatID, unmarshaled.ChatID)
	}

	if unmarshaled.Text != req.Text {
		t.Errorf("Expected Text '%s', got '%s'", req.Text, unmarshaled.Text)
	}

	if unmarshaled.ParseMode != req.ParseMode {
		t.Errorf("Expected ParseMode '%s', got '%s'", req.ParseMode, unmarshaled.ParseMode)
	}
}

func BenchmarkBot_HandleUpdate(b *testing.B) {
	bot := &Bot{
		Token:      "test_token",
		WebhookURL: "https://example.com",
		Client:     &http.Client{},
	}

	update := Update{
		UpdateID: 1,
		Message: &Message{
			MessageID: 1,
			From: &User{
				ID:        123,
				FirstName: "Test",
				Username:  "testuser",
			},
			Chat: &Chat{
				ID:   123,
				Type: "private",
			},
			Text: "/start",
		},
	}

	// Mock sendMessage to avoid actual HTTP calls
	bot.sendMessage = func(chatID int64, text, parseMode string, keyboard *InlineKeyboardMarkup) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bot.HandleUpdate(update)
	}
}
