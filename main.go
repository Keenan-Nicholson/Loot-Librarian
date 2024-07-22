package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func runBot() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	secret := os.Getenv("DISCORD_BOT_TOKEN")

	discord, err := discordgo.New("Bot " + secret)
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer discord.Close()

	fmt.Println("Bot is running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Bot is stopping!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!loot") {
		item := m.Content[6:]
		result, err := getLootInfo(item)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching loot information.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, result)
	}
}

func getLootInfo(item string) (string, error) {
	url := "https://www.dnd5eapi.co/api/equipment/" + item
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var equipment Equipment
	err = json.Unmarshal(body, &equipment)
	if err != nil {
		return "", err
	}

	weapon_category, weapon_range, damage_dice, damage_type, weapon_cost := equipment.WeaponCategory, equipment.WeaponRange, equipment.Damage.DamageDice, equipment.Damage.DamageType.Name, equipment.Cost.Quantity
	armor_category, armor_class, str_minimum, stealth_disadvantage, weight, armor_cost := equipment.ArmorCategory, equipment.ArmorClass.Base, equipment.StrMinimum, equipment.StealthDisadvantage, equipment.ArmorWeight, equipment.Cost.Quantity

	if weapon_category == "" {
		return fmt.Sprintf("Armor Category: %s\nArmor Class: %d\nStrength Minimum: %d\nStealth Disadvantage: %t\nWeight: %d\nCost: %d%s",
			armor_category, armor_class, str_minimum, stealth_disadvantage, weight, armor_cost, equipment.Cost.Unit), nil
	}
	return fmt.Sprintf("Weapon Category: %s\nWeapon Range: %s\nDamage Dice: %s\nDamage Type: %s\nCost: %d%s",
		weapon_category, weapon_range, damage_dice, damage_type, weapon_cost, equipment.Cost.Unit), nil
}

func testGetLootInfo(item string) {
	result, err := getLootInfo(item)
	if err != nil {
		fmt.Println("Error fetching loot information:", err)
	} else {
		fmt.Println(result)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		if len(os.Args) != 3 {
			fmt.Println("Usage: go run main.go test <item>")
			return
		}
		item := os.Args[2]
		testGetLootInfo(item)
		return
	}

	runBot()
}
