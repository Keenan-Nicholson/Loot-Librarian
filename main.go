package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	secret := os.Getenv("DISCORD_BOT_TOKEN")

	discord, err := discordgo.New("Bot " + secret)

	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.Contains(m.Content, "!loot") {

			item := m.Content[6:]
			url := "https://www.dnd5eapi.co/api/equipment/" + item
			method := "GET"

			client := &http.Client{}
			req, err := http.NewRequest(method, url, nil)

			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Add("Accept", "application/json")

			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			var equipment Equipment

			err = json.Unmarshal(body, &equipment)
			if err != nil {
				fmt.Println(err)
				return
			}

			weapon_category, weapon_range, damage_dice, damage_type, weapon_cost := equipment.Weapon.WeaponCategory, equipment.Weapon.WeaponRange, equipment.Weapon.Damage.DamageDice, equipment.Weapon.Damage.DamageType.Name, equipment.Weapon.Cost.Quantity
			armor_category, armor_class, str_minimum, stealth_disadvantage, weight, armor_cost := equipment.Armor.ArmorCategory, equipment.Armor.ArmorClass.Base, equipment.Armor.StrMinimum, equipment.Armor.StealthDisadvantage, equipment.Armor.Weight, equipment.Armor.Cost.Quantity

			fmt.Println(url)
			fmt.Println(equipment)
			fmt.Println(weapon_category, weapon_range, damage_dice, damage_type, weapon_cost)
			fmt.Println(armor_category, armor_class, str_minimum, stealth_disadvantage, weight, armor_cost)

			if weapon_category == "" {
				s.ChannelMessageSend(m.ChannelID,
					"Armor Category: "+armor_category+"\n"+
						"Armor Class: "+strconv.Itoa(armor_class)+"\n"+
						"Strength Minimum: "+strconv.Itoa(str_minimum)+"\n"+
						"Stealth Disadvantage: "+strconv.FormatBool(stealth_disadvantage)+"\n"+
						"Weight: "+strconv.Itoa(weight)+"\n"+
						"Cost: "+strconv.Itoa(armor_cost)+equipment.Armor.Cost.Unit)

			} else {
				s.ChannelMessageSend(m.ChannelID,
					"Weapon Category: "+weapon_category+"\n"+
						"Weapon Range: "+weapon_range+"\n"+
						"Damage Dice: "+damage_dice+"\n"+
						"Damage Type: "+damage_type+"\n"+
						"Cost: "+strconv.Itoa(weapon_cost)+equipment.Weapon.Cost.Unit)

			}

		}
	})

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
