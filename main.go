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
)

type DamageType struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Damage struct {
	DamageDice string     `json:"damage_dice"`
	DamageType DamageType `json:"damage_type"`
}

type Cost struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type Range struct {
	Normal int `json:"normal"`
}

type Content struct {
	WeaponCategory string `json:"weapon_category"`
	WeaponRange    string `json:"weapon_range"`
	Damage         Damage `json:"damage"`
	Cost           Cost   `json:"cost"`
	Range          Range  `json:"range"`
}

func main() {
	discord, err := discordgo.New("Bot " + "MTI2Mjg2OTgwMzkzMzIzNzQxMg.GDSw3L.tj_MI9iL8ozHisBZv79s4Jm3aAjhex3y1RRYfw")

	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.Contains(m.Content, "loot lib get") {

			item := m.Content[13:]
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

			var content Content

			err = json.Unmarshal(body, &content)
			if err != nil {
				fmt.Println(err)
				return
			}

			weapon_category, weapon_range, damage_dice, damage_type, cost := content.WeaponCategory, content.WeaponRange, content.Damage.DamageDice, content.Damage.DamageType.Name, content.Cost.Quantity

			s.ChannelMessageSend(m.ChannelID, "Weapon Category: "+weapon_category+"\n"+"Weapon Range: "+weapon_range+"\n"+"Damage Dice: "+damage_dice+"\n"+"Damage Type: "+damage_type+"\n"+"Cost: "+strconv.Itoa(cost)+content.Cost.Unit)
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
