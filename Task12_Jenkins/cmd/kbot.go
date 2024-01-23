/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	// TeleToken bot
	TeleToken = os.Getenv("TELE_TOKEN")
)

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"go"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started", appVersion)
		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		kbot.Handle("/start", func(m telebot.Context) error {
			menu := &telebot.ReplyMarkup{
				ReplyKeyboard: [][]telebot.ReplyButton{
					{{Text: "Hello"}, {Text: "Help"}},
					{{Text: "Kyiv"}, {Text: "Boston"}, {Text: "London"}},
					{{Text: "Vienna"}, {Text: "Tbilisi"}, {Text: "Vancouver"}},
				},
			}
			return m.Send("Welcome to Kbot!", menu)
		})

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			switch m.Text() {
			case "Hello":
				return m.Send(fmt.Sprintf("Hi! I'm Kbot %s! And I know what time it is!", appVersion))
			case "Help":
				return m.Send("This is the help message. Here you can find out the current time in the locations of your partners and team members: Kyiv, Boston, London, Vienna, Tbilisi or Vancouver")
			case "Kyiv":
				return m.Send("Current time in Kyiv: " + getTime("Kyiv"))
			case "Boston":
				return m.Send("Current time in Boston: " + getTime("Boston"))
			case "London":
				return m.Send("Current time in London: " + getTime("London"))
			case "Vienna":
				return m.Send("Current time in Vienna: " + getTime("Vienna"))
			case "Tbilisi":
				return m.Send("Current time in Tbilisi: " + getTime("Tbilisi"))
			case "Vancouver":
				return m.Send("Current time in Vancouver: " + getTime("Vancouver"))
			default:
				return m.Send("Unknown command. Please try again.")
			}
		})

		kbot.Start()
	},
}

func getTime(location string) string {
	var locName string
	switch location {
	case "Kyiv":
		locName = "Europe/Kiev"
	case "Boston":
		locName = "America/New_York"
	case "London":
		locName = "Europe/London"
	case "Vienna":
		locName = "Europe/Vienna"
	case "Tbilisi":
		locName = "Asia/Tbilisi"
	case "Vancouver":
		locName = "America/Vancouver"
	default:
		return "Invalid location"
	}

	loc, err := time.LoadLocation(locName)
	if err != nil {
		return "Invalid location"
	}
	return time.Now().In(loc).Format("15:04:05")
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
