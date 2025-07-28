package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Bot struct {
	Token      string
	Username   string
	WebhookURL string
	Client     *http.Client
}

type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from,omitempty"`
	Chat      *Chat  `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text,omitempty"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"message,omitempty"`
	Data    string   `json:"data"`
}

type User struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type SendMessageRequest struct {
	ChatID      int64                 `json:"chat_id"`
	Text        string                `json:"text"`
	ParseMode   string                `json:"parse_mode,omitempty"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string      `json:"text"`
	URL          string      `json:"url,omitempty"`
	CallbackData string      `json:"callback_data,omitempty"`
	WebApp       *WebAppInfo `json:"web_app,omitempty"`
}

type WebAppInfo struct {
	URL string `json:"url"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type SetMyCommandsRequest struct {
	Commands []BotCommand `json:"commands"`
}

type SetChatMenuButtonRequest struct {
	ChatID     int64       `json:"chat_id,omitempty"`
	MenuButton *MenuButton `json:"menu_button"`
}

type MenuButton struct {
	Type   string      `json:"type"`
	Text   string      `json:"text,omitempty"`
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

func NewBot(token, webhookURL string) (*Bot, error) {
	bot := &Bot{
		Token:      token,
		WebhookURL: webhookURL,
		Client:     &http.Client{},
	}

	if err := bot.getMe(); err != nil {
		return nil, fmt.Errorf("failed to get bot info: %w", err)
	}

	if webhookURL != "" {
		if err := bot.setWebhook(); err != nil {
			return nil, fmt.Errorf("failed to set webhook: %w", err)
		}
	}

	if err := bot.setCommands(); err != nil {
		log.Printf("Warning: failed to set commands: %v", err)
	}

	if err := bot.setMenuButton(); err != nil {
		log.Printf("Warning: failed to set menu button: %v", err)
	}

	return bot, nil
}

func (b *Bot) GetUsername() string {
	return b.Username
}

func (b *Bot) getMe() error {
	resp, err := b.makeRequest("getMe", nil)
	if err != nil {
		return err
	}

	var result struct {
		OK     bool `json:"ok"`
		Result User `json:"result"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return err
	}

	if !result.OK {
		return fmt.Errorf("API request failed")
	}

	b.Username = result.Result.Username
	return nil
}

func (b *Bot) setWebhook() error {
	data := map[string]string{
		"url": b.WebhookURL + "/webhook",
	}

	_, err := b.makeRequest("setWebhook", data)
	return err
}

func (b *Bot) setCommands() error {
	commands := []BotCommand{
		{Command: "start", Description: "Start using ChatMood"},
		{Command: "help", Description: "Get help and information"},
		{Command: "create", Description: "Create a new mood sticker"},
	}

	req := SetMyCommandsRequest{Commands: commands}
	_, err := b.makeRequest("setMyCommands", req)
	return err
}

func (b *Bot) setMenuButton() error {
	menuButton := &MenuButton{
		Type: "web_app",
		Text: "🎭 Create Sticker",
		WebApp: &WebAppInfo{
			URL: b.WebhookURL + "/app",
		},
	}

	req := SetChatMenuButtonRequest{
		MenuButton: menuButton,
	}

	_, err := b.makeRequest("setChatMenuButton", req)
	return err
}

func (b *Bot) HandleUpdate(update Update) error {
	if update.Message != nil {
		return b.handleMessage(update.Message)
	}

	if update.CallbackQuery != nil {
		return b.handleCallbackQuery(update.CallbackQuery)
	}

	return nil
}

func (b *Bot) handleMessage(message *Message) error {
	switch {
	case strings.HasPrefix(message.Text, "/start"):
		return b.sendWelcomeMessage(message.Chat.ID, message.From)
	case strings.HasPrefix(message.Text, "/help"):
		return b.sendHelpMessage(message.Chat.ID)
	case strings.HasPrefix(message.Text, "/create"):
		return b.sendCreateMessage(message.Chat.ID)
	default:
		return b.sendDefaultMessage(message.Chat.ID)
	}
}

func (b *Bot) handleCallbackQuery(query *CallbackQuery) error {
	log.Printf("Received callback query: %s", query.Data)
	return nil
}

func (b *Bot) sendWelcomeMessage(chatID int64, user *User) error {
	name := user.FirstName
	if name == "" {
		name = "there"
	}

	text := fmt.Sprintf(`🎭 *Welcome to ChatMood, %s!*

Express your emotions with custom mood stickers!

✨ *What you can do:*
• Create personalized mood stickers
• Choose from various emotions and styles
• Add custom text and effects
• Share directly in your chats

🚀 *Get started:*
Tap the button below or use the menu to create your first sticker!`, name)

	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{
					Text: "🎨 Create Your First Sticker",
					WebApp: &WebAppInfo{
						URL: b.WebhookURL + "/app",
					},
				},
			},
			{
				{
					Text:         "ℹ️ How it works",
					CallbackData: "help",
				},
			},
		},
	}

	return b.sendMessage(chatID, text, "Markdown", keyboard)
}

func (b *Bot) sendHelpMessage(chatID int64) error {
	text := `🆘 *ChatMood Help*

*How to create stickers:*
1️⃣ Choose your mood category
2️⃣ Pick an emoji that fits your vibe
3️⃣ Add custom text (optional)
4️⃣ Select colors and effects
5️⃣ Preview your creation
6️⃣ Share it in your chats!

*Tips:*
• Keep text short for better readability
• Experiment with different color combinations
• Try various effects to make your sticker unique

*Commands:*
/start - Welcome message
/help - This help message
/create - Open sticker creator

Need more help? Contact @hedgeg0d`

	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{
					Text: "🎨 Create Sticker",
					WebApp: &WebAppInfo{
						URL: b.WebhookURL + "/app",
					},
				},
			},
		},
	}

	return b.sendMessage(chatID, text, "Markdown", keyboard)
}

func (b *Bot) sendCreateMessage(chatID int64) error {
	text := `🎨 *Ready to create?*

Open the ChatMood app to start creating your custom mood stickers!

Your creativity awaits! ✨`

	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{
					Text: "🚀 Open ChatMood App",
					WebApp: &WebAppInfo{
						URL: b.WebhookURL + "/app",
					},
				},
			},
		},
	}

	return b.sendMessage(chatID, text, "Markdown", keyboard)
}

func (b *Bot) sendDefaultMessage(chatID int64) error {
	text := `🤔 I'm not sure what you mean, but I'm here to help you create amazing mood stickers!

Try these commands:
/start - Get started
/help - Learn how to use ChatMood
/create - Open the sticker creator`

	keyboard := &InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{
					Text: "🎭 Create Sticker",
					WebApp: &WebAppInfo{
						URL: b.WebhookURL + "/app",
					},
				},
			},
		},
	}

	return b.sendMessage(chatID, text, "", keyboard)
}

func (b *Bot) sendMessage(chatID int64, text, parseMode string, keyboard *InlineKeyboardMarkup) error {
	req := SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   parseMode,
		ReplyMarkup: keyboard,
	}

	_, err := b.makeRequest("sendMessage", req)
	return err
}

func (b *Bot) SendSticker(chatID int64, stickerData []byte) error {
	text := "🎉 Your custom sticker has been created! Unfortunately, Telegram API doesn't support sending custom stickers directly, but you can save and share your creation!"

	return b.sendMessage(chatID, text, "", nil)
}

func (b *Bot) makeRequest(method string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.Token, method)

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
