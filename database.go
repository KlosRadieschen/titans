package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

func addHandlers() {
	commandHandlers["listbattalions"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_number, name, pk_userID FROM Battalion INNER JOIN Pilot ON pk_userID = fk_pilot_leads ORDER BY pk_number")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var number int
			var battalionName string
			var id string
			if err := rows.Scan(&number, &battalionName, &id); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			user, _ := s.User(id)
			resultString += fmt.Sprintf("%v Battalion: \"%v\", lead by "+user.Mention()+"\n", number, battalionName)
		}
		if resultString == "" {
			resultString = "No results"
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: resultString,
			},
		})
	}

	commandHandlers["show-fleet"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT fk_flagship_leads, carriers, battleships, heavyCruisers, lightCruisers, destroyers, frigates, dropships, transportShips FROM Fleet WHERE pkfk_battalion_owns=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].IntValue())
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var flagship string
			var carriers int
			var battleships int
			var heavyCruisers int
			var lightCruisers int
			var destroyers int
			var frigates int
			var dropships int
			var transportShips int
			if err := rows.Scan(&flagship, &carriers, &battleships, &heavyCruisers, &lightCruisers, &destroyers, &frigates, &dropships, &transportShips); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("Entire fleet of battalion %v:\nFlagship: %v, Carriers: %v, Battleships: %v, Heavy Cruiser: %v, Light Cruisers: %v, Destroyers: %v, Frigates: %v, Dropships: %v, Transport Ships: %v", i.ApplicationCommandData().Options[0].IntValue(), flagship, carriers, battleships, heavyCruisers, lightCruisers, destroyers, frigates, dropships, transportShips)
		}
		if resultString == "" {
			resultString = "No results"
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: resultString,
			},
		})
	}
}
