package solvers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alaingilbert/ogame/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// CaptchaCallback ...
type CaptchaCallback func(question, icons []byte) (int64, error)

// TelegramSolver ...
func TelegramSolver(tgBotToken string, tgChatID int64) CaptchaCallback {
	return func(question, icons []byte) (int64, error) {
		tgBot, _ := tgbotapi.NewBotAPI(tgBotToken)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0", "0"),
			tgbotapi.NewInlineKeyboardButtonData("1", "1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		))
		questionImgOrig, _ := png.Decode(bytes.NewReader(question))
		bounds := questionImgOrig.Bounds()
		upLeft := image.Point{X: 0, Y: 0}
		lowRight := bounds.Max
		img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
		for y := 0; y < lowRight.Y; y++ {
			for x := 0; x < lowRight.X; x++ {
				c := questionImgOrig.At(x, y)
				r, g, b, _ := c.RGBA()
				img.Set(x, y, color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 255})
			}
		}
		buf := bytes.NewBuffer(nil)
		_ = png.Encode(buf, img)
		questionImg := tgbotapi.FileBytes{Name: "question", Bytes: buf.Bytes()}
		iconsImg := tgbotapi.FileBytes{Name: "icons", Bytes: icons}
		_, _ = tgBot.Send(tgbotapi.NewPhotoUpload(tgChatID, questionImg))
		_, _ = tgBot.Send(tgbotapi.NewPhotoUpload(tgChatID, iconsImg))
		msg := tgbotapi.NewMessage(tgChatID, "Pick one")
		msg.ReplyMarkup = keyboard
		_, _ = tgBot.Send(msg)
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, _ := tgBot.GetUpdatesChan(u)
		for update := range updates {
			if update.CallbackQuery != nil {
				_, _ = tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				_, _ = tgBot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "got "+update.CallbackQuery.Data))
				v, err := utils.ParseI64(update.CallbackQuery.Data)
				if err != nil {
					return 0, err
				}
				return v, nil
			}
		}
		return 0, errors.New("failed to get answer")
	}
}

// NinjaSolver direct integration of ogame.ninja captcha auto solver service
func NinjaSolver(apiKey string) CaptchaCallback {
	return func(question, icons []byte) (int64, error) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("question", "question.png")
		_, _ = io.Copy(part, bytes.NewReader(question))
		part1, _ := writer.CreateFormFile("icons", "icons.png")
		_, _ = io.Copy(part1, bytes.NewReader(icons))
		_ = writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "https://www.ogame.ninja/api/v1/captcha/solve", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		req.Header.Set("NJA_API_KEY", apiKey)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			by, err := io.ReadAll(resp.Body)
			if err != nil {
				return 0, errors.New("failed to auto solve captcha: " + err.Error())
			}
			return 0, errors.New("failed to auto solve captcha: " + string(by))
		}
		by, _ := io.ReadAll(resp.Body)
		var answerJson struct {
			Answer int64 `json:"answer"`
		}
		if err := json.Unmarshal(by, &answerJson); err != nil {
			return 0, errors.New("failed to auto solve captcha: " + err.Error())
		}
		return answerJson.Answer, nil
	}
}

// ManualSolver manually solve the captcha
func ManualSolver() CaptchaCallback {
	return func(question, icons []byte) (int64, error) {
		questionImg, err := png.Decode(bytes.NewReader(question))
		if err != nil {
			return -1, err
		}
		questionFile, err := os.Create("question.png")
		if err != nil {
			return -1, err
		}
		defer questionFile.Close()
		if err := png.Encode(questionFile, questionImg); err != nil {
			return -1, err
		}
		iconsImg, err := png.Decode(bytes.NewReader(icons))
		if err != nil {
			return -1, err
		}
		iconsFile, err := os.Create("icons.png")
		if err != nil {
			return -1, err
		}
		defer iconsFile.Close()
		if err := png.Encode(iconsFile, iconsImg); err != nil {
			return -1, err
		}
		var answer int64
		fmt.Print("Answer: ")
		_, _ = fmt.Scan(&answer)
		return answer, nil
	}
}
