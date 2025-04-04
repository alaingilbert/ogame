package main

import (
	"fmt"
	"github.com/alaingilbert/ogame/pkg/wrapper"
	"log"
	"os"

	"github.com/abiosoft/ishell/v2"
)

const ascii = "" +
	"                         ▄▀▀▀▀█▄,\n" +
	"                     ,µ∩        ▀▓▀ªæ▄,\n" +
	" ▄▄▄▄ ,▄▄▄▄     ~^               ▐▓█  ▐▓\n" +
	"▐▓  ▓▓▓▓▄▄▄ ▀▓▓▓j▓▀▀▓▀▀▓ ▓▀▀▓L   ▐▓,-\"`\n" +
	"▐▓▄▄▓▓▓▓▄▄▓▐▓▄▄▓▐▓  ▓  ▓∩▓▓▄▄▄  ,▓\n" +
	"                             .~\"¬"

func main() {
	fmt.Println(ascii)
	log.SetFlags(log.Ltime | log.Lshortfile)
	universe := os.Getenv("UNIVERSE")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	LANGUAGE := os.Getenv("LANGUAGE")

	shell := ishell.New()
	shell.Println("logging to " + universe + " universe")

	bot, err := wrapper.New(nil, universe, username, password, LANGUAGE)
	if err != nil {
		log.Fatal(err)
	}

	shell.AddCmd(&ishell.Cmd{
		Name: "Planet",
		Help: "Planet infos",
		Func: func(c *ishell.Context) {
			planets, _ := bot.GetPlanets()
			p := planets[0]
			c.Printf("%s [%d:%d:%d]\n",
				p.Name, p.Coordinate.Galaxy, p.Coordinate.System, p.Coordinate.Position)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "Planets",
		Help: "List planets",
		Func: func(c *ishell.Context) {
			planets, _ := bot.GetPlanets()
			for _, p := range planets {
				c.Printf("%s (%d) [%d:%d:%d]\n",
					p.Name, p.ID, p.Coordinate.Galaxy, p.Coordinate.System, p.Coordinate.Position)
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "Player",
		Help: "Display player infos",
		Func: func(c *ishell.Context) {
			player := bot.GetCachedPlayer()
			c.Printf("%s %d (%d/%d) %d \n",
				player.PlayerName, player.Points, player.Rank, player.Total, player.HonourPoints)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "IsUnderAttack",
		Help: "IsUnderAttack",
		Func: func(c *ishell.Context) {
			c.Println(bot.IsUnderAttack())
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "GetResources",
		Help: "GetResources",
		Func: func(c *ishell.Context) {
			planets, _ := bot.GetPlanets()
			resources, _ := planets[0].GetResources()
			c.Printf("Metal: %d, Crystal: %d, Deuterium: %d, Energy: %d, Dark Matter: %d\n",
				resources.Metal, resources.Crystal, resources.Deuterium, resources.Energy, resources.Darkmatter)
		},
	})

	// run shell
	shell.Run()
}
