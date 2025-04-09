package solvers

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"time"

	"github.com/bwmarrin/discordgo"
)

// DiscordSolver manual discord solver for gameforge challenge solving
// Create an application here -> https://discord.com/developers/applications
// Get the "token" from the "Bot" tab
// Go to "OAuth2" tab, in "scopes" select "Bot" and "applications.commands", in "Bot Permissions" select "Send Messages"
// Copy the "Generated URL" and invite your bot to a server that you own.
// "ownerID" get your user ID from your profile in the Discord application
func DiscordSolver(token string, ownerID string) CaptchaCallback {
	return func(ctx context.Context, question, icons []byte) (int64, error) {
		bot, err := discordgo.New("Bot " + token)
		if err != nil {
			return -1, err
		}

		// Init a DM, and interaction management
		channel, err := bot.UserChannelCreate(ownerID)
		if err != nil {
			return -1, err
		}

		answerCh := make(chan int64)

		rmHandlerFn := bot.AddHandler(handleInteraction(ctx, answerCh))
		defer rmHandlerFn()

		if err := bot.Open(); err != nil {
			return -1, err
		}
		defer bot.Close()

		embedImg, err := buildEmbedImg(question, icons)
		if err != nil {
			return -1, err
		}

		msg, err := sendImageWithSelectMenu(bot, channel, embedImg)
		if err != nil {
			return -1, err
		}

		fmt.Print("Wait for reaction to solve challenge ... ")
		var answer int64
		select {
		case answer = <-answerCh:
		case <-ctx.Done():
			return -1, ctx.Err()
		}
		fmt.Printf("Selected answer : %d.\n", answer)

		// Self-cleaning history / image
		go func(msg *discordgo.Message) {
			select {
			case <-time.After(30 * time.Second):
			case <-ctx.Done():
				return
			}
			_ = bot.ChannelMessageDelete(msg.ChannelID, msg.ID)
		}(msg)

		return answer, nil
	}
}

// Generate a single image, as Discord does not support multiple image files in a same embed
func buildEmbedImg(question, icons []byte) (out []byte, err error) {
	const topPadding = 5
	questionImage, _ := png.Decode(bytes.NewReader(question))
	iconsImage, _ := png.Decode(bytes.NewReader(icons))
	questionBounds := questionImage.Bounds()
	iconsBounds := iconsImage.Bounds()
	resultWidth := max(questionBounds.Max.X, iconsBounds.Max.X)
	resultHeight := questionBounds.Max.Y + iconsBounds.Max.Y + topPadding
	bottomRightPosition := image.Point{X: resultWidth, Y: resultHeight}
	img := image.NewRGBA(image.Rectangle{Max: bottomRightPosition})
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: 0, G: 0, B: 0, A: 255}}, image.Point{}, draw.Src)
	questionBounds = questionBounds.Add(image.Point{Y: topPadding})
	draw.Draw(img, questionBounds, questionImage, image.Point{}, draw.Over)
	iconsStartX := (questionBounds.Max.X - iconsBounds.Max.X) / 2
	iconsBounds = iconsBounds.Add(image.Point{X: iconsStartX, Y: questionBounds.Max.Y})
	draw.Draw(img, iconsBounds, iconsImage, image.Point{}, draw.Src)
	buffer := new(bytes.Buffer)
	if err = png.Encode(buffer, img); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func handleInteraction(ctx context.Context, answerCh chan int64) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			var answer int64
			switch i.MessageComponentData().CustomID {
			case "btn_option1":
				answer = 0
			case "btn_option2":
				answer = 1
			case "btn_option3":
				answer = 2
			case "btn_option4":
				answer = 3
			default:
				answer = -1
			}
			_ = discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You selected image number %d.", answer+1),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			select {
			case answerCh <- answer:
			case <-ctx.Done():
				return
			}
		}
	}
}

func sendImageWithSelectMenu(bot *discordgo.Session, channel *discordgo.Channel, img []byte) (*discordgo.Message, error) {
	buttons := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "1", Style: discordgo.SecondaryButton, CustomID: "btn_option1"},
			discordgo.Button{Label: "2", Style: discordgo.SecondaryButton, CustomID: "btn_option2"},
			discordgo.Button{Label: "3", Style: discordgo.SecondaryButton, CustomID: "btn_option3"},
			discordgo.Button{Label: "4", Style: discordgo.SecondaryButton, CustomID: "btn_option4"},
		},
	}
	message := &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{{Title: "Select the image to answer the captcha !", Image: &discordgo.MessageEmbedImage{URL: "attachment://image.png"}}},
		Components: []discordgo.MessageComponent{buttons},
		Files:      []*discordgo.File{{Name: "image.png", Reader: bytes.NewReader(img)}},
	}
	return bot.ChannelMessageSendComplex(channel.ID, message)
}
