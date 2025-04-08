package solvers

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
)

var challenge int64 = -1

func DiscordSolver(token string, owner_id string) CaptchaCallback {
	return func(question, icons []byte) (int64, error) {
		bot, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Fatal(err)
		}

		// Init a DM, and interaction management
		channel, _ := bot.UserChannelCreate(owner_id)
		bot.AddHandler(HandleInteraction)

		go func() {
			bot.Open()
		}()

		// Decode image, compute min/max size
		question_image, _ := png.Decode(bytes.NewReader(question))
		icons_image, _ := png.Decode(bytes.NewReader(icons))
		question_bounds := question_image.Bounds()
		icons_bounds := icons_image.Bounds()
		top_left_position := image.Point{X: 0, Y: 0}
		result_width := int(math.Max(float64(question_bounds.Max.X), float64(icons_bounds.Max.X)))
		result_height := question_bounds.Max.Y + icons_bounds.Max.Y
		bottom_right_position := image.Point{X: result_width, Y: result_height}
		img := image.NewRGBA(image.Rectangle{Min: top_left_position, Max: bottom_right_position})

		// Generate a single image, as Discord does not support multiple image files in a same embed
		// Question first, start at 0;0
		for y := 0; y < question_bounds.Max.Y; y++ {
			for x := 0; x < question_bounds.Max.X; x++ {
				c := question_image.At(x, y)
				r, g, b, _ := c.RGBA()
				img.Set(x, y, color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 255})
			}
		}

		// Icons second, center the image
		icons_start_x := (question_bounds.Max.X - icons_bounds.Max.X) / 2
		for y := 0; y < icons_bounds.Max.Y; y++ {
			for x := 0; x < icons_bounds.Max.X; x++ {
				c := icons_image.At(x, y)
				r, g, b, _ := c.RGBA()
				img.Set(x+icons_start_x, (y + question_bounds.Max.Y), color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 255})
			}
		}

		buffer := new(bytes.Buffer)
		_ = jpeg.Encode(buffer, img, nil)

		msg := sendImageWithSelectMenu(bot, channel, buffer.Bytes())

		fmt.Print("Wait for reaction to solve challenge ... ")
		for challenge == -1 {
			time.Sleep(1000)
		}

		// Reset in case we have another challenge
		solution := challenge
		challenge = -1

		fmt.Printf("Selected answer : %d.\n", solution)

		// Self-cleaning history / image
		go func(msg *discordgo.Message) {
			time.Sleep(30 * time.Second)
			_ = bot.ChannelMessageDelete(msg.ChannelID, msg.ID)
			defer bot.Close()
		}(msg)

		return solution, nil
	}
}

// challenge is global due to laziness
func HandleInteraction(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		switch i.MessageComponentData().CustomID {
		case "btn_option1":
			challenge = 0
		case "btn_option2":
			challenge = 1
		case "btn_option3":
			challenge = 2
		case "btn_option4":
			challenge = 3
		default:
			challenge = -1
		}

		discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("You selected image number %d.", (challenge + 1)),
				Flags:   1 << 6, // ephemeral
			},
		})
	}
}

func sendImageWithSelectMenu(bot *discordgo.Session, channel *discordgo.Channel, img []byte) *discordgo.Message {
	buttons := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "1",
				Style:    discordgo.PrimaryButton,
				CustomID: "btn_option1",
			},
			discordgo.Button{
				Label:    "2",
				Style:    discordgo.SecondaryButton,
				CustomID: "btn_option2",
			},
			discordgo.Button{
				Label:    "3",
				Style:    discordgo.SuccessButton,
				CustomID: "btn_option3",
			},
			discordgo.Button{
				Label:    "4",
				Style:    discordgo.DangerButton,
				CustomID: "btn_option4",
			},
		},
	}

	message := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Select the image to answer the captcha !",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://image.png",
				},
			},
		},
		Components: []discordgo.MessageComponent{
			buttons,
		},
		Files: []*discordgo.File{
			{
				Name:   "image.png",
				Reader: bytes.NewReader(img),
			},
		},
	}

	msg, err := bot.ChannelMessageSendComplex(channel.ID, message)
	if err != nil {
		panic(err)
	}

	return msg
}
