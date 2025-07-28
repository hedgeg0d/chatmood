# ChatMood API Documentation

## Overview

ChatMood provides a REST API for generating and sharing custom mood stickers in Telegram. The API is designed to be simple, fast, and easy to integrate with the Telegram Web Apps SDK.

**Base URL:** `https://your-domain.com/api`  
**Version:** 1.0  
**Content-Type:** `application/json`

## Authentication

Currently, the API doesn't require authentication for basic operations. However, some endpoints may validate Telegram Web App data in the future.

## Rate Limiting

- **API Endpoints:** 10 requests per second per IP
- **Webhook Endpoint:** 100 requests per second per IP

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Request limit per window
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Time when the rate limit resets

## Error Handling

All API responses follow a consistent error format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": "Additional error details (optional)"
  }
}
```

### HTTP Status Codes

- `200` - Success
- `400` - Bad Request (invalid parameters)
- `401` - Unauthorized (future use)
- `429` - Too Many Requests (rate limiting)
- `500` - Internal Server Error

## Endpoints

### 1. Generate Sticker

Create a custom mood sticker with specified parameters.

**Endpoint:** `POST /api/generate-sticker`

#### Request Body

```json
{
  "emoji": "😊",
  "text": "Happy Day!",
  "mood": "happy",
  "background": "#FFE066",
  "textColor": "#000000",
  "effect": "shadow"
}
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `emoji`   | string | No | Unicode emoji character |
| `text`    | string | No | Custom text (max 20 characters) |
| `mood`    | string | Yes | Mood category (see mood types) |
| `background` | string | No | Background color (hex format) |
| `textColor` | string | No | Text color (hex format) |
| `effect` | string | No | Text effect (see effect types) |

#### Mood Types

- `happy` - Happy/joyful mood
- `sad` - Sad/melancholic mood
- `angry` - Angry/frustrated mood
- `excited` - Excited/enthusiastic mood
- `calm` - Calm/peaceful mood
- `love` - Love/romantic mood
- `cool` - Cool/confident mood
- `tired` - Tired/exhausted mood

#### Effect Types

- `none` - No effect (default)
- `shadow` - Drop shadow effect
- `glow` - Glow effect
- `outline` - Text outline
- `gradient` - Gradient fill
- `rainbow` - Rainbow colors

#### Response

**Success (200):**
```json
{
  "success": true,
  "message": "Sticker generated successfully",
  "imageId": "sticker_abc123"
}
```

**Error (400):**
```json
{
  "error": {
    "code": "INVALID_PARAMETERS",
    "message": "Invalid mood type specified"
  }
}
```

#### Example

```bash
curl -X POST https://your-domain.com/api/generate-sticker \
  -H "Content-Type: application/json" \
  -d '{
    "emoji": "😊",
    "text": "Great Day!",
    "mood": "happy",
    "background": "#FFE066",
    "textColor": "#000000",
    "effect": "shadow"
  }'
```

### 2. Share Sticker

Share a generated sticker to Telegram chats.

**Endpoint:** `POST /api/share-sticker`

#### Request Body

```json
{
  "imageData": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJ...",
  "chatId": 123456789
}
```

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `imageData` | string | Yes | Base64 encoded PNG image data |
| `chatId` | integer | No | Telegram chat ID (if sharing to specific chat) |

#### Response

**Success (200):**
```json
{
  "success": true,
  "message": "Sticker shared successfully"
}
```

**Error (400):**
```json
{
  "error": {
    "code": "INVALID_IMAGE_DATA",
    "message": "Invalid base64 image data"
  }
}
```

#### Example

```bash
curl -X POST https://your-domain.com/api/share-sticker \
  -H "Content-Type: application/json" \
  -d '{
    "imageData": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJ...",
    "chatId": 123456789
  }'
```

## Webhook Endpoints

### Telegram Webhook

Receives updates from Telegram Bot API.

**Endpoint:** `POST /webhook`

This endpoint is automatically configured by the bot and handles incoming Telegram updates. It's not intended for direct API usage.

## Health Check

### Health Status

Check the health status of the application.

**Endpoint:** `GET /health`

#### Response

**Success (200):**
```
OK
```

## SDK Integration

### Telegram Web Apps

For integration with Telegram Web Apps, use the following JavaScript pattern:

```javascript
// Initialize Telegram Web App
const tg = window.Telegram.WebApp;
tg.ready();
tg.expand();

// Generate sticker
async function generateSticker(stickerData) {
  const response = await fetch('/api/generate-sticker', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(stickerData)
  });
  
  return response.json();
}

// Share sticker
async function shareSticker(imageData) {
  const response = await fetch('/api/share-sticker', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      imageData: imageData,
      chatId: tg.initDataUnsafe?.chat?.id
    })
  });
  
  if (response.ok) {
    tg.close();
  }
  
  return response.json();
}
```

## Best Practices

### Performance

1. **Image Optimization**: The API generates 512x512 PNG images optimized for stickers
2. **Caching**: Responses are cached where appropriate
3. **Compression**: All responses use gzip compression

### Security

1. **Input Validation**: All inputs are validated and sanitized
2. **Rate Limiting**: Prevent abuse with rate limiting
3. **CORS**: Configured for Telegram Web Apps origin

### Usage Guidelines

1. **Text Length**: Keep text under 20 characters for readability
2. **Color Contrast**: Ensure good contrast between text and background
3. **Emoji Support**: Use standard Unicode emojis for best compatibility
4. **Error Handling**: Always handle API errors gracefully

## Examples

### Complete Sticker Generation Flow

```javascript
// 1. Create sticker configuration
const stickerConfig = {
  emoji: "🎉",
  text: "Party Time!",
  mood: "excited",
  background: "#FF8E53",
  textColor: "#FFFFFF",
  effect: "glow"
};

// 2. Generate sticker
try {
  const result = await generateSticker(stickerConfig);
  if (result.success) {
    console.log('Sticker generated:', result.imageId);
    
    // 3. Get canvas data and share
    const canvas = document.getElementById('stickerCanvas');
    const imageData = canvas.toDataURL('image/png').split(',')[1];
    
    const shareResult = await shareSticker(imageData);
    if (shareResult.success) {
      console.log('Sticker shared successfully!');
    }
  }
} catch (error) {
  console.error('Error:', error);
}
```

### Error Handling Example

```javascript
async function handleStickerGeneration(config) {
  try {
    const response = await fetch('/api/generate-sticker', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(config)
    });
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error?.message || 'Generation failed');
    }
    
    return data;
  } catch (error) {
    // Handle network errors
    if (error.name === 'TypeError') {
      throw new Error('Network error - please check your connection');
    }
    
    // Re-throw API errors
    throw error;
  }
}
```

## Changelog

### Version 1.0 (2025-01-28)

- Initial API release
- Basic sticker generation
- Telegram Web App integration
- Share functionality

## Support

For issues, questions, or feature requests:

- **GitHub Issues**: [https://github.com/hedgeg0d/chatmood/issues](https://github.com/hedgeg0d/chatmood/issues)
- **Email**: Contact via GitHub profile
- **Documentation**: This file and README.md

## License

This API documentation is part of the ChatMood project, licensed under the MIT License.