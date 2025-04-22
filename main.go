package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"gopkg.in/telegram-bot-api.v4"
)

// getTimeString возвращает текущее время в формате "HH:MM"
func getTimeString() string {
	return time.Now().Format("15:04")
}

// renderClockImage рендерит изображение часов с текущим временем поверх фона
func renderClockImage(bgPath, fontPath, outputPath string) error {
	bg, err := gg.LoadImage(bgPath)
	if err != nil {
		return fmt.Errorf("failed to load background: %v", err)
	}
	dc := gg.NewContextForImage(bg)
	if err := dc.LoadFontFace(fontPath, 120); err != nil {
		return fmt.Errorf("failed to load font: %v", err)
	}
	timeStr := getTimeString()
	dc.SetRGB(0.0, 1.0, 1.0)
	dc.DrawStringAnchored(timeStr, float64(dc.Width()/2), float64(dc.Height()/2), 0.5, 0.5)
	return dc.SavePNG(outputPath)
}

// updateBotAvatar обновляет аватарку бота
func updateBotAvatar(bot *tgbotapi.BotAPI, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	photoConfig := tgbotapi.NewSetChatPhoto(tgbotapi.ChatConfig{
		ChatID: bot.Self.ID,
	})
	photoConfig.Photo = tgbotapi.FileReader{
		Name:   filepath.Base(filePath),
		Reader: file,
		Size:   -1,
	}
	_, err = bot.Request(photoConfig)
	return err
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized as: %s", bot.Self.UserName)
	bgPath := "assets/background.png"
	fontPath := "assets/orbitron.ttf"
	outputPath := "avatar.png"
	for {
		err := renderClockImage(bgPath, fontPath, outputPath)
		if err != nil {
			log.Printf("Render error: %v", err)
		} else {
			err = updateBotAvatar(bot, outputPath)
			if err != nil {
				log.Printf("Update avatar error: %v", err)
			} else {
				log.Println("Avatar updated")
			}
		}
		time.Sleep(time.Minute)
	}
}
