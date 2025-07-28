# ChatMood - Quick Start Guide

Get ChatMood up and running in 5 minutes! 🚀

## Prerequisites

- Go 1.24.5+ installed
- Docker and Docker Compose installed
- Telegram Bot Token (get one from [@BotFather](https://t.me/BotFather))

## 🚀 Quick Setup

### 1. Clone and Setup

```bash
git clone https://github.com/hedgeg0d/chatmood.git
cd chatmood
make setup
```

### 2. Configure Environment

Edit the `.env` file with your bot token:

```bash
# Required
TELEGRAM_BOT_TOKEN=your_bot_token_here
WEBHOOK_URL=https://your-domain.com

# Optional
PORT=8080
DEBUG=true
LOG_LEVEL=debug
```

### 3. Run Locally

Choose one of these methods:

#### Option A: Go Development Server
```bash
make dev
```

#### Option B: Docker Development
```bash
make docker-dev
```

### 4. Test Your Bot

1. Open your browser to `http://localhost:8080`
2. Message your bot on Telegram: `/start`
3. Click "Create Sticker" button to open the web app

## 🔧 Development Commands

```bash
# Start development server with hot reload
make dev

# Run tests
make test

# Check code quality
make lint

# Format code
make format

# Build for production
make build

# View project info
make info
```

## 📱 Using the Web App

1. **Choose Mood**: Select from 8 different mood categories
2. **Pick Emoji**: Browse and select from 90+ emojis
3. **Add Text**: Enter custom text (optional, max 20 chars)
4. **Customize**: Choose colors and effects
5. **Preview**: See real-time canvas preview
6. **Share**: Download or share directly to Telegram

## 🎨 Available Features

### Moods
- 😊 Happy - Bright yellow background
- 😢 Sad - Calm blue background  
- 😠 Angry - Bold red background
- 🤩 Excited - Vibrant orange background
- 😌 Calm - Peaceful blue background
- 😍 Love - Pink background
- 😎 Cool - Green background
- 😴 Tired - Purple background

### Effects
- **None**: Clean simple text
- **Shadow**: Drop shadow effect
- **Glow**: Glowing text
- **Outline**: Text with border
- **Gradient**: Gradient fill
- **Rainbow**: Multi-color text

## 🐳 Production Deployment

### Quick Deploy with Docker

```bash
# Setup production environment
cp .env.example .env
# Edit .env with production values

# Deploy
./deploy.sh
```

### Manual Production Setup

```bash
# Build
make build

# Run production containers
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Check status
make docker-logs
```

## 🔗 Important URLs

- **Web App**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/health`
- **API Docs**: See `API.md`
- **Webhook**: `http://localhost:8080/webhook`

## 🐛 Troubleshooting

### Bot Not Responding
```bash
# Check bot token
echo $TELEGRAM_BOT_TOKEN

# Check logs
make docker-logs

# Test webhook
curl -X POST http://localhost:8080/health
```

### Build Issues
```bash
# Clean and rebuild
make clean
make deps
make build
```

### Port Already in Use
```bash
# Change port in .env
PORT=8081

# Or kill existing process
lsof -ti:8080 | xargs kill -9
```

## 📚 Next Steps

- Read [README.md](README.md) for detailed documentation
- Check [API.md](API.md) for API reference
- See [CONTRIBUTING.md](CONTRIBUTING.md) to contribute
- Report issues on [GitHub](https://github.com/hedgeg0d/chatmood/issues)

## 💡 Pro Tips

1. **Development**: Use `make dev` for hot reload during development
2. **Testing**: Run `make test-coverage` to see test coverage
3. **Production**: Use `./deploy.sh` for automated production deployment
4. **Monitoring**: Access Grafana at `http://localhost:3000` in production
5. **SSL**: The deploy script can auto-setup Let's Encrypt certificates

## 🎯 Example API Usage

```javascript
// Generate a sticker
const response = await fetch('/api/generate-sticker', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    emoji: '🎉',
    text: 'Party!',
    mood: 'excited',
    effect: 'glow'
  })
});

const result = await response.json();
console.log('Generated:', result.imageId);
```

## 🆘 Need Help?

- 📖 Check documentation files
- 🐛 [Report issues](https://github.com/hedgeg0d/chatmood/issues)
- 💬 Start a [discussion](https://github.com/hedgeg0d/chatmood/discussions)
- ⭐ Star the repo if you like it!

---

Made with ❤️ by [hedgeg0d](https://github.com/hedgeg0d)

Happy sticker creating! 🎭✨