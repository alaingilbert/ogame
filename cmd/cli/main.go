package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alaingilbert/ogame"
	"gopkg.in/abiosoft/ishell.v2"
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

	bot, err := ogame.New(universe, username, password, LANGUAGE)
	if err != nil {
		log.Fatal(err)
	}

	shell.AddCmd(&ishell.Cmd{
		Name: "Planet",
		Help: "Planet infos",
		Func: func(c *ishell.Context) {
			p := bot.GetPlanets()[0]
			c.Printf("%s [%d:%d:%d]\n",
				p.Name, p.Coordinate.Galaxy, p.Coordinate.System, p.Coordinate.Position)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "Planets",
		Help: "List planets",
		Func: func(c *ishell.Context) {
			for _, p := range bot.GetPlanets() {
				c.Printf("%s (%d) [%d:%d:%d]\n",
					p.Name, p.ID, p.Coordinate.Galaxy, p.Coordinate.System, p.Coordinate.Position)
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "Player",
		Help: "Display player infos",
		Func: func(c *ishell.Context) {
			c.Printf("%s %d (%d/%d) %d \n",
				bot.Player.PlayerName, bot.Player.Points, bot.Player.Rank, bot.Player.Total, bot.Player.HonourPoints)
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
			resources, _ := bot.GetPlanets()[0].GetResources()
			c.Printf("Metal: %d, Crystal: %d, Deuterium: %d, Energy: %d, Dark Matter: %d\n",
				resources.Metal, resources.Crystal, resources.Deuterium, resources.Energy, resources.Darkmatter)
		},
	})

	// run shell
	shell.Run()
}
