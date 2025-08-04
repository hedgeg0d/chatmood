# ChatMood 🎭

A Telegram Mini App for creating custom mood stickers. Express your emotions with personalized stickers and share them in your chats!

## 🚀 Features

- **Custom Sticker Creation**: Generate personalized stickers with emojis, text, and effects
- **Mood Expression**: Choose from various moods and emotions
- **Canvas Rendering**: Real-time sticker preview using HTML5 Canvas
- **Telegram Integration**: Seamless sharing directly to Telegram chats
- **Material Design**: Clean, modern UI following Material Design principles
- **Free Resources**: Uses OpenMoji emojis and Google Fonts

## 🛠️ Tech Stack

- **Backend**: Go 1.24.5
- **Frontend**: React with TypeScript
- **Canvas**: HTML5 Canvas API for sticker rendering
- **UI**: Material-UI components with custom styling
- **Bot**: Telegram Bot API integration
- **Fonts**: Google Fonts (Roboto, Material Icons)
- **Emojis**: OpenMoji collection

## 📱 Screenshots
<img width="515" height="823" alt="image" src="https://github.com/user-attachments/assets/a1a089ca-be2f-4a89-b7cc-847b56115292" />
<img width="507" height="831" alt="image" src="https://github.com/user-attachments/assets/84f8049d-66a9-4938-b83d-bc87880d866e" />
<img width="516" height="831" alt="image" src="https://github.com/user-attachments/assets/0def46a2-0782-4e9d-bf78-64f439b7a4e6" />
<img width="513" height="832" alt="image" src="https://github.com/user-attachments/assets/1abd03cb-18aa-4506-89a9-7a8ac9d668c5" />

## 🔧 Setup & Installation

### Prerequisites

- Go 1.24.5 or higher
- Node.js 18+ and npm
- Telegram Bot Token

### Environment Variables

Create a `.env` file in the root directory:

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
PORT=8080
WEBHOOK_URL=https://your-domain.com
```

### Backend Setup

1. Clone the repository:
```bash
git clone https://github.com/hedgeg0d/chatmood.git
cd chatmood
```

2. Install Go dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run cmd/server/main.go
```

### Frontend Setup

1. Navigate to the web directory:
```bash
cd web
```

2. Install dependencies:
```bash
npm install
```

3. Start development server:
```bash
npm run dev
```

### Production Build

```bash
npm run build
go build -o chatmood cmd/server/main.go
./chatmood
```

## 🤖 Bot Setup

1. Create a new bot via [@BotFather](https://t.me/BotFather)
2. Get your bot token and add it to `.env`
3. Set up webhook URL for your deployed application
4. Configure bot commands and menu button

## 🎨 How It Works

1. **Choose Mood**: Select from predefined mood categories
2. **Customize**: Add text, choose colors, and apply effects
3. **Preview**: Real-time canvas preview of your sticker
4. **Share**: Send directly to Telegram chats or save locally



## 📁 Project Structure

```
chatmood/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── bot/            # Telegram bot logic
│   └── server/         # HTTP server and handlers
├── web/
│   ├── src/            # React application source
│   ├── static/         # Static assets
│   └── templates/      # HTML templates
├── assets/             # Project assets and resources
└── README.md
```

## 🌟 Features Roadmap

- [x] Basic sticker generation
- [x] Telegram Web App integration
- [x] Material Design UI
- [ ] Sticker history and favorites
- [ ] More emoji collections
- [ ] Custom background patterns
- [ ] Animation effects
- [ ] Multi-language support

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**hedgeg0d** - [GitHub](https://github.com/hedgeg0d)

## 🙏 Acknowledgments

- [OpenMoji](https://openmoji.org/) for free emoji collection
- [Material-UI](https://mui.com/) for React components
- [Telegram](https://core.telegram.org/bots/webapps) for Web Apps SDK
- [Google Fonts](https://fonts.google.com/) for typography

---

Made with ❤️ for the Telegram community
