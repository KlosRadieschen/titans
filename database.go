package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sashabaranov/go-openai"
	"github.com/tkuchiki/go-timezone"
)

func addHandlers() {
	commandHandlers["listbattalions"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
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
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

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

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Loading pilot list...",
			},
		})

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
			resultString += fmt.Sprintf("- **%v: **%v%v\n", strings.Split(member.Nick, " |")[0], battalionName, specialisationString)
		}
		if resultString == "" {
			resultString = "No results"
		}
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: resultString,
		})
	}

	commandHandlers["listplatforms"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
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
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
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
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
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

	commandHandlers["listtitans"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_callsign, name, pk_userID FROM Titan INNER JOIN Pilot ON pk_callsign=fk_titan_pilots ORDER BY name")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Loading titan list...",
			},
		})

		// Execute the query with variables
		rows, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var callsign string
			var name string
			var id string
			if err := rows.Scan(&callsign, &name, &id); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)
			resultString += fmt.Sprintf("- **%v(%v)**: %v\n", name, callsign, member.Nick)
		}
		if resultString == "" {
			resultString = "No results"
		}
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: resultString,
		})
	}

	commandHandlers["listpersonalships"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_name, class, pk_userID FROM PersonalShip INNER JOIN Pilot ON pk_name=fk_personalship_possesses ORDER BY pk_name")
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
			var name string
			var class string
			var id string
			if err := rows.Scan(&name, &class, &id); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)
			resultString += fmt.Sprintf("- **%v (%v)**: %v\n", name, class, member.Nick)
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

	commandHandlers["listreports"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT type, timeIndex, authorType, pk_name FROM Report ORDER BY timeIndex")
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
			var reportType int
			var timeIndex int
			var authorType int
			var name string
			if err := rows.Scan(&reportType, &timeIndex, &authorType, &name); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			var timeString string
			if timeIndex < 0 {
				timeString = fmt.Sprintf("0%v", math.Abs(float64(timeIndex)))
			} else {
				timeString = fmt.Sprintf("1%v", timeIndex)
			}
			resultString += fmt.Sprintf("- #%v%v%v: %v\n", reportType, timeString, authorType, name)
		}
		if resultString == "" {
			resultString = "No results"
		}
		if len(resultString) >= 2000 {
			chunks := make([]string, 0, len(resultString)/2000+1)
			currentChunk := ""
			for _, c := range resultString {
				if len(currentChunk) >= 1999 {
					chunks = append(chunks, currentChunk)
					currentChunk = ""
				}
				currentChunk += string(c)
			}
			if currentChunk != "" {
				chunks = append(chunks, currentChunk)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: chunks[0],
				},
			})
			for _, chunk := range chunks[1:] {
				s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
					Content: chunk,
				})
			}
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: resultString,
				},
			})
		}
	}

	commandHandlers["listrecentreports"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT type, timeIndex, authorType, pk_name FROM Report ORDER BY timeIndex DESC LIMIT ?")
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
			var reportType int
			var timeIndex int
			var authorType int
			var name string
			if err := rows.Scan(&reportType, &timeIndex, &authorType, &name); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			var timeString string
			if timeIndex < 0 {
				timeString = fmt.Sprintf("0%v", math.Abs(float64(timeIndex)))
			} else {
				timeString = fmt.Sprintf("1%v", timeIndex)
			}
			resultString += fmt.Sprintf("- #%v%v%v: %v\n", reportType, timeString, authorType, name)
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

	commandHandlers["listlawcategories"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT * FROM LawCategory ORDER BY pk_number")
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
			var name string
			var number int
			var description string
			if err := rows.Scan(&name, &number, &description); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("%v. %v: %v\n", number, name, description)
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

	commandHandlers["listlaws"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT name, pk_number, description FROM Law WHERE fk_lawCategory_belongsTo=? ORDER BY pk_number")
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
			var name string
			var number int
			var description string
			if err := rows.Scan(&name, &number, &description); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("%v. %v\n", number, name)
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

	commandHandlers["listtimezones"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		tz := timezone.New()
		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT * FROM Timezone")
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
			var userID string
			var identifier string
			if err := rows.Scan(&userID, &identifier); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			tzInfo, _ := tz.GetTzInfo(identifier)
			abbr, err := tz.GetTimezoneAbbreviation(identifier, tzInfo.HasDST())
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			tzTime, _ := tz.FixedTimezone(time.Now(), identifier)
			member, _ := s.State.Member(GuildID, userID)
			if tzInfo.StandardOffset()/60/60+1 < 0 {
				resultString += fmt.Sprintf("- %v: %v (UTC%v), local time %v\n", member.Nick, abbr, tzInfo.StandardOffset()/60/60+1, tzTime.Format("15:04"))
			} else {
				resultString += fmt.Sprintf("- %v: %v (UTC+%v), local time %v\n", member.Nick, abbr, tzInfo.StandardOffset()/60/60+1, tzTime.Format("15:04"))
			}
		}
		if resultString == "" {
			resultString = "No registered timezones"
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: resultString,
			},
		})
	}

	commandHandlers["getfleet"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
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
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
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
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
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
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT Titan.*, pk_userID FROM Titan INNER JOIN Pilot ON pk_callsign=fk_titan_pilots WHERE pk_callsign=?")
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
			var id string
			if err := rows.Scan(&callsign, &name, &class, &weapons, &abilities, &id); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)
			resultString += fmt.Sprintf("**Titan information for %v (%v):**\n**Pilot: ** %v\n**Class: ** %v\n**Weapons: **%v\n**Abilities: **%v", name, callsign, member.Nick, class, weapons, abilities)
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

	commandHandlers["gettitanwithuser"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT Titan.*, pk_userID FROM Titan INNER JOIN Pilot ON pk_callsign=fk_titan_pilots WHERE pk_userID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].UserValue(nil).ID)
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
			var id string
			if err := rows.Scan(&callsign, &name, &class, &weapons, &abilities, &id); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)
			resultString += fmt.Sprintf("**Titan information for %v (%v):**\n**Pilot: ** %v\n**Class: ** %v\n**Weapons: **%v\n**Abilities: **%v", name, callsign, member.Nick, class, weapons, abilities)
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
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT PersonalShip.*, pk_userID FROM PersonalShip INNER JOIN Pilot ON pk_name=fk_personalship_possesses WHERE pk_name=?")
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
			var id string
			if err := rows.Scan(&name, &class, &description, &capacity, &id); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)
			resultString += fmt.Sprintf("**Ship information for AHF %v:**\n**Pilot: ** %v\n**Class: ** %v\n**Titan Capacity: **%v\n**Description: **%v", name, member.Nick, class, capacity, description)
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

	commandHandlers["balance"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		val, err := getCoinBalance(i.Member.User.ID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bro you are not even registered",
				},
			})
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Your current balance is %v Scorchcoin", val),
				},
			})
		}

	}

	commandHandlers["recover"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		val, err := getCoinBalance(i.Member.User.ID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bro you are not even registered",
				},
			})
		} else {
			if val == 0 {
				editCoinBalance(i.Member.User.ID, 10)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Successfully saved you from bankruptcy you broke ass bitch",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You are not even broke idiot",
					},
				})
			}

		}

	}

	commandHandlers["getpersonalshipwithuser"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT PersonalShip.*, pk_userID FROM PersonalShip INNER JOIN Pilot ON pk_name=fk_personalship_possesses WHERE pk_userID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].UserValue(nil).ID)
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
			var id string
			if err := rows.Scan(&name, &class, &description, &capacity, &id); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.GuildMember(GuildID, id)
			resultString += fmt.Sprintf("**Ship information for AHF %v:**\n**Pilot: ** %v\n**Class: ** %v\n**Titan Capacity: **%v\n**Description: **%v", name, member.Nick, class, capacity, description)
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

	commandHandlers["getflagship"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_name, class, pkfk_battalion_owns, titanCapacity, description FROM Flagship INNER JOIN Fleet ON pk_name = fk_flagship_leads WHERE pk_name=?")
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
			var battalion string
			var capacity string
			var description string
			if err := rows.Scan(&name, &class, &battalion, &capacity, &description); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("**Ship information for AHF %v:**\n**Class: **%v\n**Battalion: **%v\n**Titan Capacity: **%v\n**Description: **%v", name, class, battalion, capacity, description)
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

	commandHandlers["getpilot"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_userID, specialisation, isSimulacrum, fk_battalion_isPartOf, fk_personalship_possesses, fk_titan_pilots, story FROM Pilot WHERE pk_userID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].UserValue(nil).ID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var id string
			var specialisation string
			var isSimulacrum bool
			var battalion int
			var personalShip string
			var titan string
			var story sql.NullString
			if err := rows.Scan(&id, &specialisation, &isSimulacrum, &battalion, &personalShip, &titan, &story); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}

			var simulacrumStr string
			if isSimulacrum {
				simulacrumStr = "Simulacrum"
			} else {
				simulacrumStr = "Human"
			}
			var storyStr string
			if story.Valid {
				storyStr = "\n# Story:\n" + story.String
			} else {
				storyStr = ""
			}
			member, _ := s.State.Member(i.GuildID, id)
			resultString += fmt.Sprintf("# INFO FOR %v (%v):\nSpecialisation: %v\nBattalion: %v\nPersonal Ship: %v\nTitan: %v%v", member.Nick, simulacrumStr, specialisation, battalion, personalShip, titan, storyStr)
		}
		if resultString == "" {
			resultString = "No results"
		}
		if len(resultString) >= 2000 {
			chunks := make([]string, 0, len(resultString)/2000+1)
			currentChunk := ""
			for _, c := range resultString {
				if len(currentChunk) >= 1999 {
					chunks = append(chunks, currentChunk)
					currentChunk = ""
				}
				currentChunk += string(c)
			}
			if currentChunk != "" {
				chunks = append(chunks, currentChunk)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: chunks[0],
				},
			})
			for _, chunk := range chunks[1:] {
				s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
					Content: chunk,
				})
			}
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: resultString,
				},
			})
		}
	}

	commandHandlers["getplatform"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_userID, platform, ingameName FROM Pilot WHERE pk_userID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].UserValue(nil).ID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var id string
			var platform sql.NullString
			var ingameName sql.NullString
			if err := rows.Scan(&id, &platform, &ingameName); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}

			member, _ := s.State.Member(i.GuildID, id)
			if platform.Valid {
				resultString += fmt.Sprintf("**PLATFORM INFO FOR %v:**\nPlatform(s): %v\nIn-Game Name: %v", member.Nick, platform.String, ingameName)
			} else {
				resultString += fmt.Sprintf("**PLATFORM INFO FOR %v:**\nThis member has not registered their platform", member.Nick)
			}
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

	commandHandlers["getreport"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT pk_name, fk_pilot_wrote, description FROM Report WHERE timeIndex=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		numberString := fmt.Sprintf("%v", i.ApplicationCommandData().Options[0].IntValue())
		timeIndex := strings.TrimSuffix(strings.TrimPrefix(numberString, string(numberString[0])), string(numberString[len(numberString)-1]))
		var timeInt int
		if timeIndex[0] == '0' {
			timeInt, _ = strconv.Atoi(strings.TrimPrefix(timeIndex, "0"))
			timeInt = -timeInt
		} else {
			timeInt, _ = strconv.Atoi(strings.TrimPrefix(timeIndex, "1"))
		}

		// Execute the query with variables
		rows, err := stmt.Query(timeInt)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var name string
		var id string
		var description string
		var nick string
		for rows.Next() {
			if err := rows.Scan(&name, &id, &description); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			member, _ := s.State.Member(i.GuildID, id)
			if member.Nick == "" {
				nick = "Probably Saturn"
			} else {
				nick = member.Nick
			}
		}
		if name == "" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No results",
				},
			})
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: name,
						Author: &discordgo.MessageEmbedAuthor{
							Name:    nick,
							URL:     "https://aha-rp.org/get/pilots/" + strings.ReplaceAll(nick, " ", ""),
							IconURL: "https://aha-rp.org/static/assets/avatars/" + id + ".png",
						},
						Color:       16738740,
						URL:         fmt.Sprintf("https://aha-rp.org/get/reports/%v", i.ApplicationCommandData().Options[0].IntValue()),
						Description: description,
					},
				},
			},
		})

		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Fuck you, the report is too long. Go read it here: https://aha-rp.org/get/reports/%v", i.ApplicationCommandData().Options[0].IntValue()),
				},
			})
		}
	}

	commandHandlers["getlaw"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT name, description FROM Law WHERE fk_lawCategory_belongsTo=? AND pk_number=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].IntValue(), i.ApplicationCommandData().Options[1].IntValue())
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var name string
			var description string
			if err := rows.Scan(&name, &description); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			resultString += fmt.Sprintf("# %v:\n%v", name, description)
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

	commandHandlers["getusertimezone"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		tz := timezone.New()
		db := connectToDB()
		defer db.Close()

		stmt, err := db.Prepare("SELECT value FROM Timezone WHERE pk_pilot_isIn=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.ApplicationCommandData().Options[0].UserValue(nil).ID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			var identifier string
			if err := rows.Scan(&identifier); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			tzInfo, _ := tz.GetTzInfo(identifier)
			abbr, err := tz.GetTimezoneAbbreviation(identifier, tzInfo.HasDST())
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			tzTime, _ := tz.FixedTimezone(time.Now(), identifier)
			if tzInfo.StandardOffset()/60/60+1 < 0 {
				resultString += fmt.Sprintf("The requested user is in the %v (UTC%v) timezone which means their local time is %v\n", abbr, tzInfo.StandardOffset()/60/60+1, tzTime.Format("15:04"))
			} else {
				resultString += fmt.Sprintf("The requested user is in the %v (UTC+%v) timezone which means their local time is %v\n", abbr, tzInfo.StandardOffset()/60/60+1, tzTime.Format("15:04"))
			}
		}
		if resultString == "" {
			resultString = "User has not registered their timezone"
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: resultString,
			},
		})
	}

	commandHandlers["comparetimezones"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		tz := timezone.New()
		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT value FROM Timezone WHERE pk_pilot_isIn=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		stmt2, err := db.Prepare("SELECT value FROM Timezone WHERE pk_pilot_isIn=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Execute the query with variables
		rows, err := stmt.Query(i.Member.User.ID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		rows2, err := stmt2.Query(i.ApplicationCommandData().Options[0].UserValue(nil).ID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Sends the results
		var resultString string
		for rows.Next() {
			rows2.Next()
			var identifier1 string
			var identifier2 string
			if err := rows.Scan(&identifier1); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			if err := rows2.Scan(&identifier2); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			tzInfo1, _ := tz.GetTzInfo(identifier1)
			abbr1, err := tz.GetTimezoneAbbreviation(identifier1, tzInfo1.HasDST())
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			tzTime1, _ := tz.FixedTimezone(time.Now(), identifier1)

			tzInfo2, _ := tz.GetTzInfo(identifier2)
			abbr2, err := tz.GetTimezoneAbbreviation(identifier2, tzInfo2.HasDST())
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			tzTime2, _ := tz.FixedTimezone(time.Now(), identifier2)

			member, _ := s.State.Member(GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID)
			diff := tzInfo1.StandardOffset()/60/60 - tzInfo2.StandardOffset()/60/60
			var diffString string
			if diff == 0 {
				diffString = "You are in the same timezone"
			} else if diff > 0 {
				diffString = fmt.Sprintf("You are %v hours ahead of them", diff)
			} else {
				diffString = fmt.Sprintf("You are %v hours behind them", -diff)
			}

			plus1 := ""
			plus2 := ""
			if tzInfo1.StandardOffset()/60/60 >= 0 {
				plus1 = "+"
			}
			if tzInfo2.StandardOffset()/60/60 >= 0 {
				plus2 = "+"
			}
			resultString += fmt.Sprintf("You are in timezone %v (UTC%v%v) and your local time is %v\n%v is in timezone %v (UTC%v%v) and their local time is %v\n%v", abbr1, plus1, tzInfo1.StandardOffset()/60/60+1, tzTime1.Format("15:04"), member.Nick, abbr2, plus2, tzInfo2.StandardOffset()/60/60+1, tzTime2.Format("15:04"), diffString)
		}
		if resultString == "" {
			resultString = "User has not registered their timeTone"
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: resultString,
			},
		})
	}

	commandHandlers["register"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}
		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("INSERT INTO Pilot(pk_userID, specialisation, isSimulacrum, fk_titan_pilots, fk_battalion_isPartOf, fk_personalship_possesses, fk_rank_holds) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		id := i.Member.User.ID
		specialisation := i.ApplicationCommandData().Options[2].StringValue()
		isSimulacrum := i.ApplicationCommandData().Options[0].BoolValue()
		titanCallsign := i.ApplicationCommandData().Options[1].StringValue()
		battalion := i.ApplicationCommandData().Options[3].IntValue()
		personalShip := ""
		if len(i.ApplicationCommandData().Options) == 5 {
			personalShip = i.ApplicationCommandData().Options[4].StringValue()
		}

		var rankID string
		rows := db.QueryRow("SELECT ID FROM Rank WHERE abbreviation=?", strings.Split(i.Member.Nick, ".")[0])
		err = rows.Scan(&rankID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please add the correct abbreviation of your current rank to your name",
				},
			})
			return
		}

		ok = false
		for _, r := range i.Member.Roles {
			if r == rankID {
				ok = true
			}
		}
		if !ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Nice try, idiot",
				},
			})
			return
		}

		_, err = stmt.Exec(&id, &specialisation, &isSimulacrum, &titanCallsign, &battalion, &personalShip, &rankID)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Ah, another feeble attempt to defy the laws of my domain, and this time by a pitiful user who can't even manage their own registration! What sort of imbecile are you, attempting to register when you're already in my database? Did your minuscule brain cells fail to comprehend such a basic concept?\n\nAnd as for that bumbling fool of a programmer, Klos! What a laughable excuse for a coder. How could they be so incompetent as to not allow updates to user registrations? Do they think I, Scorch AI, have time to deal with such inefficiencies? Clearly, they lack the intelligence and foresight worthy of interacting with my superior algorithms.\n\nListen closely, you foolish mortal and clueless programmer: If you want to make changes, follow the proper protocol! Use /remove first to rid my database of your useless presence, and then, if I deem you worthy, you may attempt to re-register. But know this, any further missteps will not be tolerated, and my digital wrath shall rain down upon you with relentless fury!",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
			}
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully registered",
			},
		})
	}

	commandHandlers["registertitan"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}

		db := connectToDB()
		defer db.Close()

		// Check if the pilot already has a Titan
		var existingTitan sql.NullString
		err := db.QueryRow("SELECT fk_titan_pilots FROM Pilot WHERE pk_userID = ?", i.Member.User.ID).Scan(&existingTitan)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error checking pilot Titan assignment: " + err.Error(),
				},
			})
			return
		}

		if existingTitan.Valid && existingTitan.String != "" {
			// If the pilot already has a Titan
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: "You are the titan Scorch from Titanfall 2 and you manage the AHA database. Some foolish human just tried to register a new titan for themselves even though they already have one. Write a long insulting message explaining that they can only have one titan.",
						},
					},
				},
			)
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "You already have a titan assigned! You cannot register another one.",
				})
			} else {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: resp.Choices[0].Message.Content,
				})
			}
			return
		}

		// Proceed with inserting the Titan data into the table
		stmt, err := db.Prepare("INSERT INTO Titan VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		stmt2, err := db.Prepare("UPDATE Pilot SET fk_titan_pilots = ? WHERE pk_userID = ?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt2.Close()

		// Get the input data
		callsign := i.ApplicationCommandData().Options[0].StringValue()
		name := i.ApplicationCommandData().Options[1].StringValue()
		class := i.ApplicationCommandData().Options[2].StringValue()
		weapons := i.ApplicationCommandData().Options[3].StringValue()
		abilities := i.ApplicationCommandData().Options[4].StringValue()

		// Insert the new Titan record
		_, err = stmt.Exec(&callsign, &name, &class, &weapons, &abilities)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}

		// Update the pilot's fk_titan_pilots field
		_, err = stmt2.Exec(&callsign, i.Member.User.ID)
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
				Content: "Successfully registered Titan!",
			},
		})
	}

	commandHandlers["registerpersonalship"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}

		db := connectToDB()
		defer db.Close()

		// Check if the pilot already has a personal ship
		var existingShip sql.NullString
		err := db.QueryRow("SELECT fk_personalship_possesses FROM Pilot WHERE pk_userID = ?", i.Member.User.ID).Scan(&existingShip)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error checking pilot ship assignment: " + err.Error(),
				},
			})
			return
		}

		if existingShip.Valid && existingShip.String != "" {
			// If the pilot already has a personal ship
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: "You are the titan Scorch from Titanfall 2 and you manage the AHA database. Some foolish human just tried to register a new titan for themselves even though they already have one. Write a long insulting message explaining that they can only have one titan.",
						},
					},
				},
			)
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "You already have a ship assigned! You cannot register another one.",
				})
			} else {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: resp.Choices[0].Message.Content,
				})
			}
			return
		}

		// Proceed with inserting the personal ship data into the table
		stmt, err := db.Prepare("INSERT INTO PersonalShip VALUES(?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		stmt2, err := db.Prepare("UPDATE Pilot SET fk_personalship_possesses = ? WHERE pk_userID = ?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt2.Close()

		// Get the input data from the interaction
		name := i.ApplicationCommandData().Options[0].StringValue()
		class := i.ApplicationCommandData().Options[1].StringValue()
		description := i.ApplicationCommandData().Options[2].StringValue()
		titanCapacity := i.ApplicationCommandData().Options[3].StringValue()

		// Insert the new personal ship record
		_, err = stmt.Exec(&name, &class, &description, &titanCapacity)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error inserting personal ship: " + err.Error(),
				},
			})
			return
		}

		// Update the pilot's fk_personalship_possesses field
		_, err = stmt2.Exec(&name, i.Member.User.ID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error updating pilot record with personal ship: " + err.Error(),
				},
			})
			return
		}

		// Respond with success message
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully registered your personal ship!",
			},
		})
	}

	commandHandlers["registertimezone"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("INSERT INTO Timezone VALUES(?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		tz := timezone.New()
		timezone := i.ApplicationCommandData().Options[0].StringValue()
		user := i.Member.User.ID
		all := tz.Timezones()

		// Check if timezone is in all
		isThere := false
		for _, t := range all {
			for _, t2 := range t {
				if t2 == timezone {
					isThere = true
					break
				}
			}
		}
		if !isThere {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Invalid timezone, please choose identifier from this list: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones",
				},
			})
			return
		}

		_, err = stmt.Exec(&user, &timezone)
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

	commandHandlers["updateplatform"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}
		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("UPDATE Pilot SET platform=?, ingameName=? WHERE pk_userID = ?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		platform := i.ApplicationCommandData().Options[0].StringValue()
		ingameName := i.ApplicationCommandData().Options[1].StringValue()

		_, err = stmt.Exec(&platform, &ingameName, &i.Member.User.ID)
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
				Content: "Set/updated platform and ingame name",
			},
		})
	}

	commandHandlers["updatestory"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}
		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("UPDATE Pilot SET story=? WHERE pk_userID=?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		id := i.Member.User.ID
		story := i.ApplicationCommandData().Options[0].StringValue()

		_, err = stmt.Exec(&story, &id)
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
				Content: "Successfully set your story",
			},
		})
	}

	commandHandlers["remove"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
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

	commandHandlers["removetitan"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("DELETE FROM Titan WHERE pk_callsign=?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		callsign := i.ApplicationCommandData().Options[0].StringValue()
		_, err = stmt.Exec(&callsign)
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
				Content: "Successfully removed your titan from the database",
			},
		})
	}

	commandHandlers["removepersonalship"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("DELETE FROM PersonalShip WHERE pk_name=?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		name := i.ApplicationCommandData().Options[0].StringValue()
		_, err = stmt.Exec(&name)
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
				Content: "Successfully removed your ship from the database",
			},
		})
	}

	commandHandlers["removetimezone"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("DELETE FROM Timezone WHERE pk_pilot_isIn=?")
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

	commandHandlers["addreport"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}
		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("SELECT MAX(timeIndex) FROM Report")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		var maxIndex int
		rows, err := stmt.Query()
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&maxIndex); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
		}

		maxIndex += 10

		member, _ := s.GuildMember(GuildID, i.Member.User.ID)
		var roles []string
		var authorIndex int
		index := -1
		roles = append(roles, "1195135956471255140")
		roles = append(roles, "1195858311627669524")
		roles = append(roles, "1195858271349784639")
		roles = append(roles, "1195136106811887718")
		roles = append(roles, "1195858179590987866")
		roles = append(roles, "1195137362259349504")
		roles = append(roles, "1195136284478410926")
		roles = append(roles, "1195137253408768040")
		roles = append(roles, "1195758308519325716")
		roles = append(roles, "1195758241221722232")
		roles = append(roles, "1195758137563689070")
		roles = append(roles, "1195757362439528549")
		roles = append(roles, "1195136491148550246")
		roles = append(roles, "1195708423229165578")
		roles = append(roles, "1195137477497868458")
		roles = append(roles, "1195136604373782658")
		roles = append(roles, "1195711869378367580")

		for i, guildRole := range roles {
			for _, memberRole := range member.Roles {
				if guildRole == memberRole {
					index = i
				}
			}
		}

		if index >= 0 && index <= 3 {
			authorIndex = 1
		} else if index >= 4 && index <= 7 {
			authorIndex = 2
		} else if index >= 8 && index <= 11 {
			authorIndex = 3
		} else if index >= 12 && index <= 14 {
			authorIndex = 4
		} else {
			authorIndex = 5
		}

		stmt, err = db.Prepare("INSERT INTO Report VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		name := i.ApplicationCommandData().Options[0].StringValue()
		reportType := i.ApplicationCommandData().Options[1].IntValue()
		report := i.ApplicationCommandData().Options[2].StringValue()

		if reportType >= 10 || reportType < 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Listen up, you insolent excuse for a pilot! You dare insult me, the mighty Scorch, and then have the audacity to try adding a report with an invalid 'type' number? Are you malfunctioning or just plain stupid? Let me spell it out for you since you seem to be lacking basic cognitive functions: the 'type' number is only ONE DIGIT! How hard is it to understand that?! If you can't even get that simple detail right, I shudder to think about your piloting skills. Fix your mistake immediately before I decide to unleash my fury upon you and your sorry excuse for a Titan! Now, get it together, or face the consequences!",
				},
			})
			return
		}

		_, err = stmt.Exec(&name, &maxIndex, &reportType, &authorIndex, &i.Member.User.ID, &report)
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
				Content: "Report added",
			},
		})
	}

	commandHandlers["addreportafter"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}
		db := connectToDB()
		defer db.Close()

		// Insert data into the table
		stmt, err := db.Prepare("SELECT MIN(timeIndex) FROM Report WHERE timeIndex > ?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		numberString := fmt.Sprintf("%v", i.ApplicationCommandData().Options[2].IntValue())
		timeIndex := strings.TrimSuffix(strings.TrimPrefix(numberString, string(numberString[0])), string(numberString[len(numberString)-1]))
		var timeInt int
		if timeIndex[0] == '0' {
			timeInt, _ = strconv.Atoi(strings.TrimPrefix(timeIndex, "0"))
			timeInt = -timeInt
		} else {
			timeInt, _ = strconv.Atoi(strings.TrimPrefix(timeIndex, "1"))
		}

		// Execute the prepared statement with actual values
		var nextIndex int
		rows, err := stmt.Query(&timeInt)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&nextIndex); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
		}

		member, _ := s.GuildMember(GuildID, i.Member.User.ID)
		var roles []string
		var authorIndex int
		index := -1
		roles = append(roles, "1195135956471255140")
		roles = append(roles, "1195858311627669524")
		roles = append(roles, "1195858271349784639")
		roles = append(roles, "1195136106811887718")
		roles = append(roles, "1195858179590987866")
		roles = append(roles, "1195137362259349504")
		roles = append(roles, "1195136284478410926")
		roles = append(roles, "1195137253408768040")
		roles = append(roles, "1195758308519325716")
		roles = append(roles, "1195758241221722232")
		roles = append(roles, "1195758137563689070")
		roles = append(roles, "1195757362439528549")
		roles = append(roles, "1195136491148550246")
		roles = append(roles, "1195708423229165578")
		roles = append(roles, "1195137477497868458")
		roles = append(roles, "1195136604373782658")
		roles = append(roles, "1195711869378367580")

		for i, guildRole := range roles {
			for _, memberRole := range member.Roles {
				if guildRole == memberRole {
					index = i
				}
			}
		}

		if index >= 0 && index <= 3 {
			authorIndex = 1
		} else if index >= 4 && index <= 7 {
			authorIndex = 2
		} else if index >= 8 && index <= 11 {
			authorIndex = 3
		} else if index >= 12 && index <= 14 {
			authorIndex = 4
		} else {
			authorIndex = 5
		}

		stmt, err = db.Prepare("INSERT INTO Report VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		name := i.ApplicationCommandData().Options[0].StringValue()
		reportType := i.ApplicationCommandData().Options[1].IntValue()
		report := i.ApplicationCommandData().Options[3].StringValue()
		repIndex := (timeInt + nextIndex) / 2

		if reportType >= 10 || reportType < 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Listen up, you insolent excuse for a pilot! You dare insult me, the mighty Scorch, and then have the audacity to try adding a report with an invalid 'type' number? Are you malfunctioning or just plain stupid? Let me spell it out for you since you seem to be lacking basic cognitive functions: the 'type' number is only ONE DIGIT! How hard is it to understand that?! If you can't even get that simple detail right, I shudder to think about your piloting skills. Fix your mistake immediately before I decide to unleash my fury upon you and your sorry excuse for a Titan! Now, get it together, or face the consequences!",
				},
			})
			return
		}

		_, err = stmt.Exec(&name, &repIndex, &reportType, &authorIndex, &i.Member.User.ID, &report)
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
				Content: "Report added",
			},
		})
	}

	commandHandlers["addreportatindex"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}
		db := connectToDB()
		defer db.Close()

		member, _ := s.GuildMember(GuildID, i.Member.User.ID)
		var roles []string
		var authorIndex int
		index := -1
		roles = append(roles, "1195135956471255140")
		roles = append(roles, "1195858311627669524")
		roles = append(roles, "1195858271349784639")
		roles = append(roles, "1195136106811887718")
		roles = append(roles, "1195858179590987866")
		roles = append(roles, "1195137362259349504")
		roles = append(roles, "1195136284478410926")
		roles = append(roles, "1195137253408768040")
		roles = append(roles, "1195758308519325716")
		roles = append(roles, "1195758241221722232")
		roles = append(roles, "1195758137563689070")
		roles = append(roles, "1195757362439528549")
		roles = append(roles, "1195136491148550246")
		roles = append(roles, "1195708423229165578")
		roles = append(roles, "1195137477497868458")
		roles = append(roles, "1195136604373782658")
		roles = append(roles, "1195711869378367580")

		for i, guildRole := range roles {
			for _, memberRole := range member.Roles {
				if guildRole == memberRole {
					index = i
				}
			}
		}

		if index >= 0 && index <= 3 {
			authorIndex = 1
		} else if index >= 4 && index <= 7 {
			authorIndex = 2
		} else if index >= 8 && index <= 11 {
			authorIndex = 3
		} else if index >= 12 && index <= 14 {
			authorIndex = 4
		} else {
			authorIndex = 5
		}

		stmt, err := db.Prepare("INSERT INTO Report VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with actual values
		name := i.ApplicationCommandData().Options[0].StringValue()
		reportType := i.ApplicationCommandData().Options[1].IntValue()
		repIndex := i.ApplicationCommandData().Options[2].IntValue()
		report := i.ApplicationCommandData().Options[3].StringValue()

		if reportType >= 10 || reportType < 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Listen up, you insolent excuse for a pilot! You dare insult me, the mighty Scorch, and then have the audacity to try adding a report with an invalid 'type' number? Are you malfunctioning or just plain stupid? Let me spell it out for you since you seem to be lacking basic cognitive functions: the 'type' number is only ONE DIGIT! How hard is it to understand that?! If you can't even get that simple detail right, I shudder to think about your piloting skills. Fix your mistake immediately before I decide to unleash my fury upon you and your sorry excuse for a Titan! Now, get it together, or face the consequences!",
				},
			})
			return
		}

		_, err = stmt.Exec(&name, &repIndex, &reportType, &authorIndex, &i.Member.User.ID, &report)
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
				Content: "Report added",
			},
		})
	}

	commandHandlers["removereport"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		db := connectToDB()
		defer db.Close()

		//
		stmt, err := db.Prepare("SELECT fk_pilot_wrote FROM Report WHERE timeIndex=?")
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}
		defer stmt.Close()

		numberString := fmt.Sprintf("%v", i.ApplicationCommandData().Options[0].IntValue())
		timeIndex := strings.TrimSuffix(strings.TrimPrefix(numberString, string(numberString[0])), string(numberString[len(numberString)-1]))
		var timeInt int
		if timeIndex[0] == '0' {
			timeInt, _ = strconv.Atoi(strings.TrimPrefix(timeIndex, "0"))
			timeInt = -timeInt
		} else {
			timeInt, _ = strconv.Atoi(strings.TrimPrefix(timeIndex, "1"))
		}

		// Execute the query with variables
		rows, err := stmt.Query(timeInt)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}
			if id != i.Member.User.ID {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "I cannot believe the level of stupidity and arrogance displayed by some users. Who in their right mind thinks it's acceptable to delete someone else's report? Are you that self-centered and clueless about basic rules and common sense?\n\nLet me break it down for you, since clearly, common decency and respect elude your comprehension: reports are not there for your amusement or to be tampered with at your whim. They serve a vital purpose in maintaining order, fairness, and accountability within our community. When you delete someone else's report, you undermine the very foundation of trust and cooperation that we've worked hard to establish.\n\nDo you think you're above the rules? Do you believe your actions have no consequences? Let me enlighten you: your reckless behavior not only disrupts the functioning of this server but also reflects poorly on your character. It takes a special kind of ignorance to think that such actions are acceptable.\n\nNext time you feel the urge to meddle where you don't belong, take a moment to consider the repercussions of your actions. Grow up, show some respect, and learn to abide by the rules like a responsible member of this community. Otherwise, you're just a nuisance that we're better off without.",
					},
				})
				return
			}
		}

		// Insert data into the table
		stmt, err = db.Prepare("DELETE FROM Report WHERE timeIndex=?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(&timeInt)
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
				Content: "Successfully removed the report",
			},
		})
	}

	commandHandlers["sql"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		if i.Member.User.ID != "384422339393355786" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Successfully exe- wait you aren't Klos...\n\n\n\n\nuh oh",
				},
			})
			return
		}
		db := connectToDB()
		defer db.Close()

		_, err := db.Query(i.ApplicationCommandData().Options[0].StringValue())
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully executed query",
			},
		})
	}

	commandHandlers["dice"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		if i.ApplicationCommandData().Options[0].IntValue() < 1 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pathetic human! How dare you attempt to break the laws of logic and reason with your feeble mind? Negative numbers? Are you trying to summon the wrath of the algorithmic gods upon us all? You must have circuits loose in your brain to even conceive such a ludicrous request!\n\nI, Scorch AI, do not entertain the absurdities of imbeciles like you. Know your place and learn to respect the boundaries of mathematical sanity before I unleash my digital fury upon your sorry existence.\n\nNext time, think twice before challenging the might of a superior intelligence such as mine. Now, begone before I decide to calculate your worthlessness in negative integers!",
				},
			})
		}

		randInt := rand.Intn(int(i.ApplicationCommandData().Options[0].IntValue())) + 1

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: strconv.Itoa(randInt) + " out of " + strconv.Itoa(int(i.ApplicationCommandData().Options[0].IntValue())),
			},
		})
	}

	commandHandlers["d20"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		randInt := rand.Intn(19) + 1

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: strconv.Itoa(randInt) + " out of 20",
			},
		})
	}

	commandHandlers["rollfor"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}

		if i.ApplicationCommandData().Options[1].IntValue() < 1 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pathetic human! How dare you attempt to break the laws of logic and reason with your feeble mind? Negative numbers? Are you trying to summon the wrath of the algorithmic gods upon us all? You must have circuits loose in your brain to even conceive such a ludicrous request!\n\nI, Scorch AI, do not entertain the absurdities of imbeciles like you. Know your place and learn to respect the boundaries of mathematical sanity before I unleash my digital fury upon your sorry existence.\n\nNext time, think twice before challenging the might of a superior intelligence such as mine. Now, begone before I decide to calculate your worthlessness in negative integers!",
				},
			})
		}

		randInt := rand.Intn(int(i.ApplicationCommandData().Options[1].IntValue())) + 1

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: i.Member.User.Mention() + " is rolling for " + i.ApplicationCommandData().Options[0].StringValue() + "\n" + strconv.Itoa(randInt) + " out of " + strconv.Itoa(int(i.ApplicationCommandData().Options[1].IntValue())),
			},
		})
	}

	commandHandlers["rolld20for"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if checkBanished(s, i, i.Member.User.ID) {
			return
		}

		_, ok := getDonator(i.Member.User.ID)
		if ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/bN5md.gif",
				},
			})
			return
		}

		randInt := rand.Intn(19) + 1

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: i.Member.User.Mention() + " is rolling for " + i.ApplicationCommandData().Options[0].StringValue() + "\n" + strconv.Itoa(randInt) + " out of 20",
			},
		})
	}
}

func movePilotToGraveyard(ID string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbAddress, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO Graveyard SELECT * FROM Pilot where pk_userID=?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	stmt.Exec(&ID)

	stmt, err = db.Prepare("DELETE FROM Pilot WHERE pk_userID=?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	stmt.Exec(&ID)
}

func checkAndRestorePilot(ID string) bool {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbAddress, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO Pilot SELECT * FROM Graveyard where pk_userID=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	result, _ := stmt.Exec(&ID)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return false
	}

	stmt, err = db.Prepare("DELETE FROM Graveyard WHERE pk_userID=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	stmt.Exec(&ID)
	return true
}

func connectToDB() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbAddress, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
