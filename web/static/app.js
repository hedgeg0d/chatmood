// ChatMood Web App JavaScript
class ChatMoodApp {
    constructor() {
        this.canvas = document.getElementById('stickerCanvas');
        this.ctx = this.canvas.getContext('2d');
        this.selectedMood = '';
        this.selectedEmoji = '😊';
        this.customText = '';
        this.textColor = '#000000';
        this.selectedEffect = 'none';
        this.backgroundColors = {
            happy: '#FFE066',
            sad: '#87CEEB',
            angry: '#FF6B6B',
            excited: '#FF8E53',
            calm: '#9ECAE1',
            love: '#FFB6C1',
            cool: '#98FB98',
            tired: '#DDA0DD'
        };

        this.init();
    }

    init() {
        this.initTelegramWebApp();
        this.setupEventListeners();
        this.populateMoods();
        this.populateEmojis();
        this.populateColors();
        this.populateEffects();
        this.updatePreview();
    }

    initTelegramWebApp() {
        if (window.Telegram && window.Telegram.WebApp) {
            const tg = window.Telegram.WebApp;
            tg.ready();
            tg.expand();

            // Set theme colors
            document.documentElement.style.setProperty('--tg-theme-bg-color', tg.themeParams.bg_color || '#ffffff');
            document.documentElement.style.setProperty('--tg-theme-text-color', tg.themeParams.text_color || '#000000');
            document.documentElement.style.setProperty('--tg-theme-button-color', tg.themeParams.button_color || '#667eea');
            document.documentElement.style.setProperty('--tg-theme-button-text-color', tg.themeParams.button_text_color || '#ffffff');
            document.documentElement.style.setProperty('--tg-theme-secondary-bg-color', tg.themeParams.secondary_bg_color || '#f8f9fa');
            document.documentElement.style.setProperty('--tg-theme-section-separator-color', tg.themeParams.section_separator_color || '#e1e8ed');

            // Enable back button
            tg.BackButton.show();
            tg.BackButton.onClick(() => {
                tg.close();
            });

            // Set main button
            tg.MainButton.setText('Share Sticker');
            tg.MainButton.color = tg.themeParams.button_color || '#667eea';
            tg.MainButton.textColor = tg.themeParams.button_text_color || '#ffffff';
            tg.MainButton.onClick(() => {
                this.shareSticker();
            });
        }
    }

    setupEventListeners() {
        // Text input
        document.getElementById('customText').addEventListener('input', (e) => {
            this.customText = e.target.value;
            this.updatePreview();
        });

        // Download button
        document.getElementById('downloadBtn').addEventListener('click', () => {
            this.downloadSticker();
        });

        // Share button
        document.getElementById('shareBtn').addEventListener('click', () => {
            this.shareSticker();
        });
    }

    populateMoods() {
        const moods = [
            { id: 'happy', emoji: '😊', label: 'Happy' },
            { id: 'sad', emoji: '😢', label: 'Sad' },
            { id: 'angry', emoji: '😠', label: 'Angry' },
            { id: 'excited', emoji: '🤩', label: 'Excited' },
            { id: 'calm', emoji: '😌', label: 'Calm' },
            { id: 'love', emoji: '😍', label: 'Love' },
            { id: 'cool', emoji: '😎', label: 'Cool' },
            { id: 'tired', emoji: '😴', label: 'Tired' }
        ];

        const moodGrid = document.getElementById('moodGrid');
        moods.forEach(mood => {
            const button = document.createElement('div');
            button.className = 'mood-btn';
            button.innerHTML = `
                <div class="mood-emoji">${mood.emoji}</div>
                <div class="mood-label">${mood.label}</div>
            `;

            button.addEventListener('click', () => {
                document.querySelectorAll('.mood-btn').forEach(btn => btn.classList.remove('active'));
                button.classList.add('active');
                this.selectedMood = mood.id;
                this.updatePreview();
            });

            moodGrid.appendChild(button);
        });

        // Select first mood by default
        moodGrid.firstElementChild.classList.add('active');
        this.selectedMood = moods[0].id;
    }

    populateEmojis() {
        const emojiCategories = {
            faces: ['😀', '😃', '😄', '😁', '😆', '😅', '😂', '🤣', '😊', '😇', '🙂', '🙃', '😉', '😌', '😍', '🥰', '😘', '😗', '😙', '😚', '😋', '😛', '😝', '😜', '🤪', '🤨', '🧐', '🤓', '😎', '🤩', '🥳'],
            emotions: ['😏', '😒', '😞', '😔', '😟', '😕', '🙁', '☹️', '😣', '😖', '😫', '😩', '🥺', '😢', '😭', '😤', '😠', '😡', '🤬', '🤯', '😳', '🥵', '🥶', '😱', '😨', '😰', '😥', '😓', '🤗', '🤔'],
            gestures: ['🤭', '🤫', '🤥', '😶', '😐', '😑', '😬', '🙄', '😯', '😦', '😧', '😮', '😲', '🥱', '😴', '🤤', '😪', '😵', '🤐', '🥴', '🤢', '🤮', '🤧', '😷', '🤒', '🤕', '🤑', '🤠']
        };

        const allEmojis = [...emojiCategories.faces, ...emojiCategories.emotions, ...emojiCategories.gestures];
        const emojiGrid = document.getElementById('emojiGrid');

        allEmojis.forEach(emoji => {
            const button = document.createElement('div');
            button.className = 'emoji-btn';
            button.textContent = emoji;

            button.addEventListener('click', () => {
                document.querySelectorAll('.emoji-btn').forEach(btn => btn.classList.remove('selected'));
                button.classList.add('selected');
                this.selectedEmoji = emoji;
                this.updatePreview();

                // Show main button when emoji is selected
                if (window.Telegram && window.Telegram.WebApp) {
                    window.Telegram.WebApp.MainButton.show();
                }
            });

            emojiGrid.appendChild(button);
        });

        // Select first emoji by default
        emojiGrid.firstElementChild.classList.add('selected');
    }

    populateColors() {
        const colors = [
            '#000000', '#FFFFFF', '#FF0000', '#00FF00', '#0000FF', '#FFFF00',
            '#FF00FF', '#00FFFF', '#FFA500', '#800080', '#FFC0CB', '#A52A2A',
            '#808080', '#000080', '#008000', '#800000', '#FF6B6B', '#4ECDC4'
        ];

        const colorPicker = document.getElementById('colorPicker');
        colors.forEach(color => {
            const button = document.createElement('div');
            button.className = 'color-btn';
            button.style.backgroundColor = color;

            button.addEventListener('click', () => {
                document.querySelectorAll('.color-btn').forEach(btn => btn.classList.remove('selected'));
                button.classList.add('selected');
                this.textColor = color;
                this.updatePreview();
            });

            colorPicker.appendChild(button);
        });

        // Select first color by default
        colorPicker.firstElementChild.classList.add('selected');
        this.textColor = colors[0];
    }

    populateEffects() {
        const effects = [
            { id: 'none', label: 'None' },
            { id: 'shadow', label: 'Shadow' },
            { id: 'glow', label: 'Glow' },
            { id: 'outline', label: 'Outline' },
            { id: 'gradient', label: 'Gradient' },
            { id: 'rainbow', label: 'Rainbow' }
        ];

        const effectButtons = document.getElementById('effectButtons');
        effects.forEach(effect => {
            const button = document.createElement('div');
            button.className = 'effect-btn';
            button.textContent = effect.label;

            button.addEventListener('click', () => {
                document.querySelectorAll('.effect-btn').forEach(btn => btn.classList.remove('active'));
                button.classList.add('active');
                this.selectedEffect = effect.id;
                this.updatePreview();
            });

            effectButtons.appendChild(button);
        });

        // Select first effect by default
        effectButtons.firstElementChild.classList.add('active');
    }

    updatePreview() {
        const ctx = this.ctx;
        const canvas = this.canvas;

        // Clear canvas
        ctx.clearRect(0, 0, canvas.width, canvas.height);

        // Set background
        const bgColor = this.backgroundColors[this.selectedMood] || '#FFE066';
        const gradient = ctx.createRadialGradient(256, 256, 0, 256, 256, 300);
        gradient.addColorStop(0, bgColor);
        gradient.addColorStop(1, this.adjustBrightness(bgColor, -20));

        ctx.fillStyle = gradient;
        ctx.fillRect(0, 0, canvas.width, canvas.height);

        // Draw emoji
        ctx.font = '200px Arial';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';

        // Add emoji shadow
        ctx.fillStyle = 'rgba(0, 0, 0, 0.2)';
        ctx.fillText(this.selectedEmoji, 258, 258);

        // Draw main emoji
        ctx.fillStyle = '#000000';
        ctx.fillText(this.selectedEmoji, 256, 256);

        // Draw text if provided
        if (this.customText.trim()) {
            this.drawStyledText(ctx, this.customText.trim(), 256, 400);
        }

        // Add border
        ctx.strokeStyle = 'rgba(255, 255, 255, 0.3)';
        ctx.lineWidth = 4;
        ctx.strokeRect(2, 2, canvas.width - 4, canvas.height - 4);
    }

    drawStyledText(ctx, text, x, y) {
        ctx.font = 'bold 48px Roboto, Arial, sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';

        // Apply effects
        switch (this.selectedEffect) {
            case 'shadow':
                ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
                ctx.fillText(text, x + 3, y + 3);
                break;

            case 'glow':
                ctx.shadowColor = this.textColor;
                ctx.shadowBlur = 10;
                ctx.shadowOffsetX = 0;
                ctx.shadowOffsetY = 0;
                break;

            case 'outline':
                ctx.strokeStyle = this.adjustBrightness(this.textColor, -40);
                ctx.lineWidth = 4;
                ctx.strokeText(text, x, y);
                break;

            case 'gradient':
                const textGradient = ctx.createLinearGradient(0, y - 24, 0, y + 24);
                textGradient.addColorStop(0, this.textColor);
                textGradient.addColorStop(1, this.adjustBrightness(this.textColor, -30));
                ctx.fillStyle = textGradient;
                ctx.fillText(text, x, y);
                return;

            case 'rainbow':
                this.drawRainbowText(ctx, text, x, y);
                return;
        }

        // Draw main text
        ctx.fillStyle = this.textColor;
        ctx.fillText(text, x, y);

        // Reset shadow
        ctx.shadowBlur = 0;
    }

    drawRainbowText(ctx, text, x, y) {
        const colors = ['#FF0000', '#FF7F00', '#FFFF00', '#00FF00', '#0000FF', '#4B0082', '#9400D3'];
        const charWidth = ctx.measureText(text).width / text.length;
        const startX = x - (ctx.measureText(text).width / 2);

        for (let i = 0; i < text.length; i++) {
            ctx.fillStyle = colors[i % colors.length];
            ctx.fillText(text[i], startX + (i * charWidth) + (charWidth / 2), y);
        }
    }

    adjustBrightness(color, amount) {
        const hex = color.replace('#', '');
        const num = parseInt(hex, 16);
        const r = Math.max(0, Math.min(255, (num >> 16) + amount));
        const g = Math.max(0, Math.min(255, (num >> 8 & 0x00FF) + amount));
        const b = Math.max(0, Math.min(255, (num & 0x0000FF) + amount));
        return `#${(r << 16 | g << 8 | b).toString(16).padStart(6, '0')}`;
    }

    async downloadSticker() {
        const link = document.createElement('a');
        link.download = `chatmood-sticker-${Date.now()}.png`;
        link.href = this.canvas.toDataURL();
        link.click();

        // Show feedback
        this.showNotification('Sticker downloaded successfully! 📥');
    }

    async shareSticker() {
        try {
            const imageData = this.canvas.toDataURL('image/png').split(',')[1];

            const shareData = {
                imageData: imageData
            };

            // If in Telegram WebApp, try to get chat info
            if (window.Telegram && window.Telegram.WebApp) {
                const tg = window.Telegram.WebApp;
                if (tg.initDataUnsafe && tg.initDataUnsafe.start_param) {
                    shareData.chatId = parseInt(tg.initDataUnsafe.start_param);
                }
            }

            const response = await fetch('/api/share-sticker', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(shareData)
            });

            if (response.ok) {
                this.showNotification('Sticker shared successfully! 🎉');

                // Close the WebApp if in Telegram
                if (window.Telegram && window.Telegram.WebApp) {
                    setTimeout(() => {
                        window.Telegram.WebApp.close();
                    }, 1500);
                }
            } else {
                throw new Error('Failed to share sticker');
            }
        } catch (error) {
            console.error('Error sharing sticker:', error);
            this.showNotification('Failed to share sticker. Please try again.', 'error');
        }
    }

    showNotification(message, type = 'success') {
        // Create notification element
        const notification = document.createElement('div');
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            left: 50%;
            transform: translateX(-50%);
            background: ${type === 'error' ? '#ff4757' : '#2ed573'};
            color: white;
            padding: 1rem 2rem;
            border-radius: 25px;
            font-weight: 500;
            box-shadow: 0 4px 15px rgba(0,0,0,0.2);
            z-index: 1000;
            animation: slideInDown 0.3s ease;
        `;
        notification.textContent = message;

        // Add animation styles
        const style = document.createElement('style');
        style.textContent = `
            @keyframes slideInDown {
                from {
                    opacity: 0;
                    transform: translateX(-50%) translateY(-20px);
                }
                to {
                    opacity: 1;
                    transform: translateX(-50%) translateY(0);
                }
            }
        `;
        document.head.appendChild(style);

        document.body.appendChild(notification);

        // Remove notification after 3 seconds
        setTimeout(() => {
            notification.style.animation = 'slideInDown 0.3s ease reverse';
            setTimeout(() => {
                if (notification.parentNode) {
                    notification.parentNode.removeChild(notification);
                }
                if (style.parentNode) {
                    style.parentNode.removeChild(style);
                }
            }, 300);
        }, 3000);

        // Show haptic feedback if available
        if (window.Telegram && window.Telegram.WebApp) {
            window.Telegram.WebApp.HapticFeedback.notificationOccurred(
                type === 'error' ? 'error' : 'success'
            );
        }
    }
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new ChatMoodApp();
});
