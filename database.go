package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

func addHandlers() {
	commandHandlers["list"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmt.Println("list")
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/test.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Prepare a parameterized SQL query
		stmt, err := db.Prepare("SELECT pk_* FROM ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].StringValue)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var pk interface{}
			if err := rows.Scan(&pk); err != nil {
				log.Fatal(err)
			}
			resultString += fmt.Sprintf("%v ", pk)
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
