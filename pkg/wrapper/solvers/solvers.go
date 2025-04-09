package solvers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"github.com/alaingilbert/ogame/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// TelegramSolver ...
func TelegramSolver(tgBotToken string, tgChatID int64) gameforge.CaptchaSolver {
	return func(ctx context.Context, question, icons []byte) (int64, error) {
		tgBot, err := tgbotapi.NewBotAPI(tgBotToken)
		if err != nil {
			return -1, err
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0", "0"),
			tgbotapi.NewInlineKeyboardButtonData("1", "1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		))
		questionImgOrig, err := png.Decode(bytes.NewReader(question))
		if err != nil {
			return -1, err
		}
		bounds := questionImgOrig.Bounds()
		lowRight := bounds.Max
		img := image.NewRGBA(image.Rectangle{Max: lowRight})
		draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{A: 255}}, image.Point{}, draw.Src)
		draw.Draw(img, bounds, questionImgOrig, image.Point{}, draw.Over)
		buf := bytes.NewBuffer(nil)
		if err := png.Encode(buf, img); err != nil {
			return -1, err
		}
		questionImg := tgbotapi.FileBytes{Name: "question", Bytes: buf.Bytes()}
		iconsImg := tgbotapi.FileBytes{Name: "icons", Bytes: icons}
		_, _ = tgBot.Send(tgbotapi.NewPhotoUpload(tgChatID, questionImg))
		_, _ = tgBot.Send(tgbotapi.NewPhotoUpload(tgChatID, iconsImg))
		msg := tgbotapi.NewMessage(tgChatID, "Pick one")
		msg.ReplyMarkup = keyboard
		sentMsg, _ := tgBot.Send(msg)
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updatesCh, err := tgBot.GetUpdatesChan(u)
		if err != nil {
			return -1, err
		}
		for {
			select {
			case <-ctx.Done():
				return -1, ctx.Err()
			case update, ok := <-updatesCh:
				if !ok {
					return -1, errors.New("failed to get answer")
				}
				if update.CallbackQuery != nil && update.CallbackQuery.Message.MessageID == sentMsg.MessageID {
					_, _ = tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
					_, _ = tgBot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "got "+update.CallbackQuery.Data))
					v, err := utils.ParseI64(update.CallbackQuery.Data)
					if err != nil {
						return -1, err
					}
					return v, nil
				}
			}
		}
	}
}

// NinjaSolver direct integration of ogame.ninja captcha auto solver service
func NinjaSolver(apiKey string) gameforge.CaptchaSolver {
	return func(ctx context.Context, question, icons []byte) (int64, error) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("question", "question.png")
		if err != nil {
			return -1, err
		}
		if _, err = io.Copy(part, bytes.NewReader(question)); err != nil {
			return -1, err
		}
		part1, err := writer.CreateFormFile("icons", "icons.png")
		if err != nil {
			return -1, err
		}
		if _, err := io.Copy(part1, bytes.NewReader(icons)); err != nil {
			return -1, err
		}
		if err := writer.Close(); err != nil {
			return -1, err
		}

		req, err := http.NewRequest(http.MethodPost, "https://www.ogame.ninja/api/v1/captcha/solve", body)
		if err != nil {
			return -1, err
		}
		req.Header.Add("Content-Type", writer.FormDataContentType())
		req.Header.Set("NJA_API_KEY", apiKey)
		req.WithContext(ctx)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return -1, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			by, err := io.ReadAll(resp.Body)
			if err != nil {
				return -1, err
			}
			return -1, errors.New("failed to auto solve captcha: " + string(by))
		}
		by, err := io.ReadAll(resp.Body)
		if err != nil {
			return -1, err
		}
		var answerJson struct {
			Answer int64 `json:"answer"`
		}
		if err := json.Unmarshal(by, &answerJson); err != nil {
			return -1, err
		}
		return answerJson.Answer, nil
	}
}

// ManualSolver manually solve the captcha
func ManualSolver() gameforge.CaptchaSolver {
	return func(ctx context.Context, question, icons []byte) (int64, error) {
		saveImg := func(imgBytes []byte, fileName string) error {
			img, err := png.Decode(bytes.NewReader(imgBytes))
			if err != nil {
				return err
			}
			imgFile, err := os.Create(fileName)
			if err != nil {
				return err
			}
			defer imgFile.Close()
			return png.Encode(imgFile, img)
		}
		if err := saveImg(question, "question.png"); err != nil {
			return -1, err
		}
		if err := saveImg(icons, "icons.png"); err != nil {
			return -1, err
		}
		var answer int64
		fmt.Print("Answer: ")
		if _, err := fmt.Scan(&answer); err != nil {
			return -1, err
		}
		return answer, nil
	}
}
