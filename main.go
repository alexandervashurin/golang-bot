package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

type TarotCard struct {
	Name         string `json:"name"`
	ImagePath    string `json:"image_path"`
	Interpretation string `json:"interpretation"`
}

var tarotCards []TarotCard

func loadTarotCards(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &tarotCards)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	err = loadTarotCards("tarot_interpretations.json")
	if err != nil {
		panic(err)
	}

	b, err := bot.New(os.Getenv("EXAMPLE_TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", handleStart)
	b.RegisterHandler(bot.HandlerTypeCallbackQuery, handleCallback)

	fmt.Println("Bot is running...")
	b.Start(context.Background())
}

func handleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome to the Tarot Bot! Click the button to draw a card.",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "🃏 Draw a Card",
						CallbackData: "draw_card",
					},
				},
			},
		},
	})
}

func handleCallback(ctx context.Context, b *bot.Bot, callback *models.CallbackQuery) {
	if callback.Data == "draw_card" {
		rand.Seed(time.Now().UnixNano())
		card := tarotCards[rand.Intn(len(tarotCards))]

		if callback.Message != nil && callback.Message.Chat != nil {
			b.SendPhoto(ctx, &bot.SendPhotoParams{
				ChatID: callback.Message.Chat.ID,
				Photo:  models.FileURL(card.ImagePath),
				Caption: fmt.Sprintf("Card: %s\n\nInterpretation: %s", card.Name, card.Interpretation),
			})
		}
	}
}
