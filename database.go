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
		stmt, err := db.Prepare("SELECT pk_number, name, pk_userID, pk_name FROM Battalion INNER JOIN Pilot ON pk_userID = fk_pilot_leads INNER JOIN Fleet ON pk_number = pkfk_battalion_owns INNER JOIN Flagship ON fk_flagship_leads = pk_name ORDER BY pk_number")
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
			var flagship string
			if err := rows.Scan(&number, &battalionName, &id, &flagship); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)
			resultString += fmt.Sprintf("%v Battalion: \"%v\", lead by **"+member.Nick+"** on the **AHF %v**\n\n", number, battalionName, flagship)
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

	commandHandlers["listpilots"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_userID, specialisation, fk_battalion_isPartOf FROM Pilot ORDER BY fk_battalion_isPartOf")
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
			var id string
			var specialisation sql.NullString
			var battalion sql.NullInt64
			if err := rows.Scan(&id, &specialisation, &battalion); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)
			battalionName := "no"
			if battalion.Valid {
				if battalion.Int64 > 0 {
					battalionName = fmt.Sprintf("%v. battalion", battalion.Int64)
				} else {
					battalionName = "SWAG"
				}
			}
			specialisationString := ""
			if specialisation.Valid {
				specialisationString = ", " + specialisation.String
			}
			resultString += fmt.Sprintf("- **%v: **%v%v\n", member.Nick, battalionName, specialisationString)
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

	commandHandlers["listplatforms"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_userID, platform, ingameName FROM Pilot WHERE platform != '' ORDER BY platform")
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
			var id string
			var platform string
			var ingameName string
			if err := rows.Scan(&id, &platform, &ingameName); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)

			resultString += fmt.Sprintf("**%v:**\nPlatform: %v, Ingame name: %v\n\n", member.Nick, platform, ingameName)
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

	commandHandlers["listbases"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT Base.pk_name, size, fk_planet_isOn, fk_battalion_controls FROM Base INNER JOIN Planet ON fk_planet_isOn = Planet.pk_name ORDER BY fk_battalion_controls")
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
			var baseName string
			var size string
			var planetName string
			var battalion int
			if err := rows.Scan(&baseName, &size, &planetName, &battalion); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("**%v:**\n%v on %v, controlled by %v. battalion\n\n", baseName, size, planetName, battalion)
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

	commandHandlers["listplanets"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_name, fk_system_isInside, fk_battalion_controls FROM Planet ORDER BY fk_battalion_controls")
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
			var planetName string
			var system string
			var battalion int
			if err := rows.Scan(&planetName, &system, &battalion); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("**%v:**\nPlanet in the %v system, controlled by %v. battalion\n\n", planetName, system, battalion)
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

	commandHandlers["getfleet"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			resultString += fmt.Sprintf("**Entire fleet of battalion %v**:\nFlagship: %v\nCarriers: %v\nBattleships: %v\nHeavy Cruiser: %v\nLight Cruisers: %v\nDestroyers: %v\nFrigates: %v\nDropships: %v\nTransport Ships: %v", i.ApplicationCommandData().Options[0].IntValue(), flagship, carriers, battleships, heavyCruisers, lightCruisers, destroyers, frigates, dropships, transportShips)
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

	commandHandlers["getplanet"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_name, environment, fk_system_isInside, fk_battalion_controls FROM Planet WHERE pk_name=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].StringValue())
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var planetName string
			var environment string
			var inSystem string
			var battalion string
			if err := rows.Scan(&planetName, &environment, &inSystem, &battalion); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("**Planet information**:\n%v is a planet inside the %v system and is controlled by the %v. battalion\n**Description:**\n%v", planetName, inSystem, battalion, environment)
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

	commandHandlers["getplanet"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_name, environment, fk_system_isInside, fk_battalion_controls FROM Planet WHERE pk_name=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].StringValue())
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var planetName string
			var environment string
			var inSystem string
			var battalion string
			if err := rows.Scan(&planetName, &environment, &inSystem, &battalion); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("**Planet information**:\n%v is a planet inside the %v system and is controlled by the %v. battalion\n**Description:**\n%v", planetName, inSystem, battalion, environment)
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

	commandHandlers["gettitan"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT * FROM Titan WHERE pk_callsign=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].StringValue())
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var callsign string
			var name string
			var class string
			var weapons string
			var abilities string
			if err := rows.Scan(&callsign, &name, &class, &weapons, &abilities); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("**Titan information for %v (%v):**\n**Class: ** %v\n**Weapons: **%v\n**Abilities: **%v", name, callsign, class, weapons, abilities)
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

	commandHandlers["getpersonalship"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT * FROM PersonalShip WHERE pk_name=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].StringValue())
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var name string
			var class string
			var description string
			var capacity string
			if err := rows.Scan(&name, &class, &description, &capacity); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("**Ship information for AHF %v:**\n**Class: ** %v\n**Titan Capacity: **%v\n**Description: **%v", name, class, capacity, description)
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

	// commanderHandler to INSERT a new pilot
	commandHandlers["register"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("INSERT INTO Pilot(pk_userID, rank, specialisation, isSimulacrum, fk_titan_pilots, fk_battalion_isPartOf, fk_personalship_possesses) VALUES(?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		id := i.Member.User.ID
		rank, _ := s.State.Role(GuildID, i.Member.Roles[0])
		rankName := rank.Name
		specialisation := i.ApplicationCommandData().Options[2].StringValue()
		isSimulacrum := i.ApplicationCommandData().Options[0].BoolValue()
		titanCallsign := i.ApplicationCommandData().Options[1].StringValue()
		battalion := i.ApplicationCommandData().Options[3].StringValue()
		personalShip := ""
		if len(i.ApplicationCommandData().Options) == 5 {
			personalShip = i.ApplicationCommandData().Options[4].StringValue()
		}
		_, err = stmt.Exec(&id, &rankName, &specialisation, &isSimulacrum, &titanCallsign, &battalion, &personalShip)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully registered",
			},
		})
	}

	commandHandlers["remove"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		db, err := sql.Open("sqlite3", "/home/Nicolas/go-workspace/src/titans/AHA.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("DELETE FROM Pilot WHERE pk_userID=?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		id := i.Member.User.ID
		_, err = stmt.Exec(&id)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully removed you from the database",
			},
		})
	}
}
