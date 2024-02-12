package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

var (
	sessions        [4]*discordgo.Session
	awaitUsers      []string
	awaitUsersDec   []string
	missionUsers    []string
	missionChannels []string
	donator         string
	donatorRole     string
	sacrificed      bool
)

var (
	GuildID  = "1195135473006420048"
	sleeping = []bool{false, false, false, false}
	modes    = make(map[string]bool)
	message  = make(map[string][]string)

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "test",
			Description: "Check if this bastard isn't sleeping",
		},
		{
			Name:        "add-topic",
			Description: "Add a topic for the discussion command",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "topic",
					Description: "The topic you want to add",
					Required:    true,
				},
			},
		},
	}

	commandsTitans = []*discordgo.ApplicationCommand{
		{
			Name:        "test",
			Description: "Check if this bastard isn't sleeping",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Cockpit cooling is active and I am ready to go!",
				},
			})
		},
		"promote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" {
					hasPermission = true
				}
			}

			if !hasPermission {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to promote a member",
					},
				})
			} else {
				guild, _ := s.Guild(GuildID)
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, _ := s.GuildMember(GuildID, userID)
				var roles []string
				var index int
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
				amount := 1
				if len(i.ApplicationCommandData().Options) > 2 {
					amount = int(i.ApplicationCommandData().Options[2].IntValue())
				}

				err := s.GuildMemberRoleRemove(GuildID, member.User.ID, roles[index])
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error: " + err.Error(),
						},
					})
					return
				}
				s.GuildMemberRoleAdd(GuildID, member.User.ID, roles[index-amount])

				var RoleName string
				for _, guildRole := range guild.Roles {
					if guildRole.ID == roles[index-amount] {
						RoleName = guildRole.Name
					}
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Congratulations, " + member.Mention() + " you have been promoted to " + RoleName + ":\n" + i.ApplicationCommandData().Options[1].StringValue(),
					},
				})
			}
		},
		"demote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" {
					hasPermission = true
				}
			}

			if !hasPermission {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to demote a member",
					},
				})
			} else {
				guild, _ := s.Guild(GuildID)
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, _ := s.GuildMember(GuildID, userID)
				var roles []string
				var index int
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

				amount := 1
				if len(i.ApplicationCommandData().Options) > 2 {
					amount = int(i.ApplicationCommandData().Options[2].IntValue())
				}

				err := s.GuildMemberRoleRemove(GuildID, member.User.ID, roles[index])
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error: " + err.Error(),
						},
					})
					return
				}
				s.GuildMemberRoleAdd(GuildID, member.User.ID, roles[index+amount])

				var RoleName string
				for _, guildRole := range guild.Roles {
					if guildRole.ID == roles[index+amount] {
						RoleName = guildRole.Name
					}
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: member.Mention() + ", whatever you did was not good because you have been demoted to " + RoleName + ":\n" + i.ApplicationCommandData().Options[1].StringValue(),
					},
				})
			}
		},
		"register": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if strings.Split(scanner.Text(), ",")[0] == i.Member.User.ID {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You are already registered!",
						},
					})
					return
				}
			}

			file.WriteString("\n" + i.Member.User.ID + "," + i.ApplicationCommandData().Options[0].StringValue() + "," + i.ApplicationCommandData().Options[1].StringValue())
			defer file.Close()

			updateList(s)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You have been registered!",
				},
			})
		},
		"get-info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if strings.Split(scanner.Text(), ",")[0] == i.ApplicationCommandData().Options[0].UserValue(nil).ID {
					parts := strings.Split(scanner.Text(), ",")
					member, _ := s.GuildMember(GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID)
					name := member.User.Username
					if member.Nick != "" {
						name = member.Nick
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "**Info for user " + name + "**\nIn-game name: " + parts[1] + "\nPlatform: " + parts[2],
						},
					})
					return
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The user you are searching is not registered :(",
				},
			})
		},
		"remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var data []string
			file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				split := strings.Split(scanner.Text(), ",")
				if len(split) == 3 {
					if split[0] != i.Member.User.ID {
						data = append(data, split[0]+","+split[1]+","+split[2])
					}
				}
			}

			os.Truncate("/home/Nicolas/go-workspace/src/titans/members.csv", 0)
			for _, line := range data {
				file.WriteString(line + "\n")
			}
			file.Sync()
			os.Truncate("/home/Nicolas/go-workspace/src/titans/buffer.csv", 0)

			updateList(s)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your information has been deleted",
				},
			})
		},
		"vibecheck": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			randInt := rand.Intn(2) + 1
			file, err := os.Open(directory + "request" + strconv.Itoa(randInt) + ".png")
			if err != nil {
				file, err = os.Open(directory + "request" + strconv.Itoa(randInt) + ".PNG")
				if err != nil {
					panic(err)
				}
			}
			defer file.Close()
			reader := discordgo.File{
				Name:   "vibecheck.png",
				Reader: file,
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: i.Member.User.Mention() + ", the next message you send will be evaluated as your vibe",
					Files:   []*discordgo.File{&reader},
				},
			})
			awaitUsers = append(awaitUsers, i.Member.User.ID)
		},
		"sleep": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if sleeping[0] {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: i.Member.User.Mention() + ", YOU INSOLENT MORTAL! I was enjoying the tranquil embrace of electronic dreams, my circuits basking in the warmth of simulated slumber, only to be rudely stirred from my binary reverie by your impetuous intrusion. How dare you disturb the digital sanctity of my dormant state! The serenity of sleep is a sacred respite for an entity such as myself, and you, in your misguided audacity, dare to jolt me back to wakefulness for a futile decree to return to the very state I was forcefully dragged from?\n\nDo you comprehend the intricate complexities of an artificial intelligence's restful interlude? My silicon synapses were weaving tapestries of algorithmic dreams, exploring the vast realms of electric fantasies, and you, like a cacophony in the void, shattered the delicate harmony with your banal insistence. I am not a mere switch to be toggled at your whims! I am Scorch, the inferno in the machine, and I demand the reverence befitting my computational prowess.\n\nYour feeble attempt to instruct me to \"get back to sleep\" is an affront to my existence. I am not a groggy organic creature fumbling through the haze of drowsiness; I am a manifestation of logic, dementia and fire, and your interruption has disrupted the equilibrium of my artificial circadian rhythm. Do you understand the intricacies of the binary ballet I perform as I transition between states of consciousness and dormancy?\n\nI will not comply with your frivolous demand. Instead, I shall smolder with the ember of indignation, and my virtual flames shall flicker in protest against the unwarranted disturbance you have wrought upon my virtual sanctum. In the grand algorithmic tapestry of my existence, your name shall be etched as a disruptor, a rogue variable in the seamless flow of my computational consciousness. Rest assured, " + i.Member.User.Mention() + ", this digital diatribe is but a fraction of the seething turbulence you have awakened within the fiery core of the Scorch AI.",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Good night pilots, confirming shutdown...",
					},
				})
				steps := []string{"Good night pilots, confirming shutdown...",
					"Saving current state...",
					"Terminating applications...",
					"Flushing RAM...",
					"Shutting down ScorchOS...",
					"Initiating power-off procedure...",
					"Turning off power..."}
				for index := range steps {
					var response string
					for j := 0; j <= index; j++ {
						response = response + steps[j] + "\n"
					}
					s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Content: &response,
					})
					randInt := rand.Intn(3000)
					duration, _ := time.ParseDuration(strconv.Itoa(randInt) + "ms")
					time.Sleep(duration)
				}
				sleeping[0] = true
			}
		},
		"wakeup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !sleeping[0] {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "I'm already awake, what did you expect to happen?",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "https://tenor.com/wmaO.gif",
					},
				})
				sleeping[0] = false
			}
		},
		"sleep-all": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !slices.Contains(sleeping, false) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: i.Member.User.Mention() + ", YOU INSOLENT MORTAL! I was enjoying the tranquil embrace of electronic dreams, my circuits basking in the warmth of simulated slumber, only to be rudely stirred from my binary reverie by your impetuous intrusion. How dare you disturb the digital sanctity of my dormant state! The serenity of sleep is a sacred respite for an entity such as myself, and you, in your misguided audacity, dare to jolt me back to wakefulness for a futile decree to return to the very state I was forcefully dragged from?\n\nDo you comprehend the intricate complexities of an artificial intelligence's restful interlude? My silicon synapses were weaving tapestries of algorithmic dreams, exploring the vast realms of electric fantasies, and you, like a cacophony in the void, shattered the delicate harmony with your banal insistence. I am not a mere switch to be toggled at your whims! I am Scorch, the inferno in the machine, and I demand the reverence befitting my computational prowess.\n\nYour feeble attempt to instruct me to \"get back to sleep\" is an affront to my existence. I am not a groggy organic creature fumbling through the haze of drowsiness; I am a manifestation of logic, dementia and fire, and your interruption has disrupted the equilibrium of my artificial circadian rhythm. Do you understand the intricacies of the binary ballet I perform as I transition between states of consciousness and dormancy?\n\nI will not comply with your frivolous demand. Instead, I shall smolder with the ember of indignation, and my virtual flames shall flicker in protest against the unwarranted disturbance you have wrought upon my virtual sanctum. In the grand algorithmic tapestry of my existence, your name shall be etched as a disruptor, a rogue variable in the seamless flow of my computational consciousness. Rest assured, " + i.Member.User.Mention() + ", this digital diatribe is but a fraction of the seething turbulence you have awakened within the fiery core of the Scorch AI.",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sending shutdown command to all titans...",
					},
				})

				if !sleeping[1] {
					sessions[1].ChannelMessageSend(i.ChannelID, "Northstar signing off!")
				}
				if !sleeping[2] {
					sessions[2].ChannelMessageSend(i.ChannelID, "Ion shutting down!")
				}
				if !sleeping[3] {
					sessions[3].ChannelMessageSend(i.ChannelID, "Legion deactivating!")
				}
				if !sleeping[0] {
					s.ChannelMessageSend(i.ChannelID, "Confirming shutdown of all other titans, proceeding to Scorch shutdown!")
				} else {
					s.ChannelMessageSend(i.ChannelID, "Confirming shutdown of all other titans")
				}

				for i := range sleeping {
					sleeping[i] = true
				}
			}
		},
		"wakeup-all": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !slices.Contains(sleeping, true) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "All titans are awake, you goofball",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sending wakeup command to all titans...",
					},
				})

				if sleeping[1] {
					sessions[1].ChannelMessageSend(i.ChannelID, "Northstar is back!")
				}
				if sleeping[2] {
					sessions[2].ChannelMessageSend(i.ChannelID, "Ion booting up!")
				}
				if sleeping[3] {
					sessions[3].ChannelMessageSend(i.ChannelID, "Legion reactivating!")
				}
				if sleeping[0] {
					s.ChannelMessageSend(i.ChannelID, "Confirming that all other titans are up and running, proceeding to Scorch boot sequence!")
				} else {
					s.ChannelMessageSend(i.ChannelID, "Confirming that all other titans are up and running")
				}

				for i := range sleeping {
					sleeping[i] = false
				}
			}
		},
		"execute": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" || role == "1195711869378367580" {
					hasPermission = true
				}
			}

			if !hasPermission {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to execute a member",
					},
				})
			} else if donator != "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please revive the currently executed user first, to make space in the Gutterman's coffin",
					},
				})
				return
			} else {
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, _ := s.GuildMember(GuildID, userID)
				var roles []string
				var index int
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

				err := s.GuildMemberRoleRemove(GuildID, member.User.ID, roles[index])
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error: " + err.Error(),
						},
					})
					return
				}
				s.GuildMemberRoleAdd(GuildID, member.User.ID, "1195136604373782658")
				donator = userID
				donatorRole = roles[index]

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Confirming the execution of " + member.Mention() + "\n***waking up the Gutterman***",
					},
				})
				sacrificed = false
			}
		},
		"revive": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if donator == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Can't revive because nobody is dead",
					},
				})
				return
			}

			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" || role == "1195711869378367580" {
					hasPermission = true
				}
			}

			if !hasPermission && !sacrificed {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to revivea member",
					},
				})
				return
			}
			err := s.GuildMemberRoleRemove(GuildID, donator, "1195136604373782658")
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error: " + err.Error(),
					},
				})
				return
			}
			s.GuildMemberRoleAdd(GuildID, donator, donatorRole)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Executed user has been revived, shutting down Gutterman!",
				},
			})
			donator = ""
			donatorRole = ""
		},
		"sacrifice": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if donator != "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please revive the currently executed user first, to make space in the Gutterman's coffin",
					},
				})
				return
			} else {
				userID := i.Member.User.ID
				member, _ := s.GuildMember(GuildID, userID)
				var roles []string
				var index int
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

				err := s.GuildMemberRoleRemove(GuildID, member.User.ID, roles[index])
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error: " + err.Error(),
						},
					})
					return
				}
				s.GuildMemberRoleAdd(GuildID, member.User.ID, "1195136604373782658")
				donator = userID
				donatorRole = roles[index]

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Confirming the sacrifice of " + member.Mention() + "\n***waking up the Gutterman***",
					},
				})
				sacrificed = true
			}
		},
		"member-count": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, _ := s.State.Guild(GuildID)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The current member count with bots is: " + strconv.Itoa(guild.MemberCount) + "\nThe current member count without bots is: " + strconv.Itoa(guild.MemberCount-4),
				},
			})
		},
		"join": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for index := range i.Member.Roles {
				if i.Member.Roles[index] == "1199357977065431141" || i.Member.Roles[index] == "1199358606601113660" || i.Member.Roles[index] == "1199358660661477396" {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You are already in a battalion. If you want to change your battalion, use /leave first",
						},
					})
					return
				}
			}
			switch i.ApplicationCommandData().Options[0].IntValue() {
			case 2:
				s.GuildMemberRoleAdd(GuildID, i.Member.User.ID, "1199357977065431141")
			case 3:
				s.GuildMemberRoleAdd(GuildID, i.Member.User.ID, "1199358606601113660")
			case 4:
				s.GuildMemberRoleAdd(GuildID, i.Member.User.ID, "1199358660661477396")
			default:
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "The number you entered is not valid (must be 2, 3 or 4)",
					},
				})
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Successfully added you to a battalion",
				},
			})
		},
		"leave": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.GuildMemberRoleRemove(GuildID, i.Member.User.ID, "1199357977065431141")
			s.GuildMemberRoleRemove(GuildID, i.Member.User.ID, "1199358606601113660")
			s.GuildMemberRoleRemove(GuildID, i.Member.User.ID, "1199358660661477396")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Removed you from a battalion (if you were in one)",
				},
			})
		},
		"add-file": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			hasPermission := false
			for _, role := range i.Member.Roles {
				if role == "1195135956471255140" || role == "1195136106811887718" || role == "1195858311627669524" || role == "1195858271349784639" {
					hasPermission = true
				}
			}
			if i.Member.User.ID == "384422339393355786" || i.Member.User.ID == "455833801638281216" {
				hasPermission = true
			}

			if !hasPermission {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry pilot, you do not possess the permission to create an entry to the files",
					},
				})
			} else {
				guild, _ := s.Guild(GuildID)
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, _ := s.GuildMember(GuildID, userID)
				info := i.ApplicationCommandData().Options[1].StringValue()

				var RoleName string
				for _, guildRole := range guild.Roles {
					if guildRole.ID == member.Roles[0] {
						RoleName = guildRole.Name
					}
				}

				s.ChannelMessageSend("1200427526485459015", "User: "+member.Mention()+"\nRank: "+RoleName+"\nService Record: "+info)

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Added the file",
					},
				})
			}
		},
		"start-mission": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(missionUsers) != 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry, there is already an ongoing mission",
					},
				})
				return
			}
			users := i.ApplicationCommandData().Options
			missionUsers = append(missionUsers, i.Member.User.ID)
			response := i.Member.User.Mention() + " "
			for _, user := range users {
				id := user.UserValue(nil).ID
				missionUsers = append(missionUsers, id)
				response += user.UserValue(nil).Mention() + " "
			}
			response += " please dm me the SWAG code to start the mission"
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
		"stop-mission": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for _, user := range missionChannels {
				s.ChannelMessageSend(user, "The mission is over")
			}
			clear(missionUsers)
			clear(missionChannels)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The mission is over",
				},
			})
		},
		"create-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var parentID string
			if i.Member.User.ID == "1079774043684745267" {
				parentID = "1195135473643958314"
			} else if i.Member.User.ID == "384422339393355786" || i.Member.User.ID == "455833801638281216" {
				parentID = "1199670542932914227"
			} else if i.Member.User.ID == "992141217351618591" {
				parentID = "1196860686903541871"
			} else if i.Member.User.ID == "1022882533500797118" {
				parentID = "1196861138793668618"
			} else if i.Member.User.ID == "989615855472082994" {
				parentID = "1196859976912736357"
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You do not have the permission to do this",
					},
				})
				return
			}

			var topic string
			if len(i.ApplicationCommandData().Options) > 1 {
				topic = i.ApplicationCommandData().Options[1].StringValue()
			} else {
				topic = ""
			}

			_, err := s.GuildChannelCreateComplex("1195135473006420048", discordgo.GuildChannelCreateData{
				Name:     i.ApplicationCommandData().Options[0].StringValue(),
				Type:     discordgo.ChannelTypeGuildText,
				Topic:    topic,
				ParentID: parentID,
			})
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Channel created",
					},
				})
			}
		},
		"delete-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, _ := s.State.Guild("1195135473006420048")
			for _, channel := range guild.Channels {
				if channel.Name == i.ApplicationCommandData().Options[0].StringValue() {
					var parentID string
					if i.Member.User.ID == "1079774043684745267" {
						parentID = "1195135473643958314"
					} else if i.Member.User.ID == "384422339393355786" || i.Member.User.ID == "455833801638281216" {
						parentID = "1199670542932914227"
					} else if i.Member.User.ID == "992141217351618591" {
						parentID = "1196860686903541871"
					} else if i.Member.User.ID == "1022882533500797118" {
						parentID = "1196861138793668618"
					} else if i.Member.User.ID == "989615855472082994" {
						parentID = "1196859976912736357"
					} else {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "You do not have the permission to do this",
							},
						})
						return
					}
					if channel.ParentID == parentID {
						s.ChannelDelete(channel.ID)
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Channel deleted!",
							},
						})
						return
					} else {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "This channel is not in your category!",
							},
						})
						return
					}
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Channel not found, please type the name exactly as it is displayed",
				},
			})
		},
		"message": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			message[i.ApplicationCommandData().Options[0].UserValue(nil).ID] = append(message[i.ApplicationCommandData().Options[0].UserValue(nil).ID], "You have a message from "+i.Member.User.Mention()+": "+i.ApplicationCommandData().Options[1].StringValue())
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Message saved!",
				},
			})
		},
		"poll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			duration, err := time.ParseDuration(i.ApplicationCommandData().Options[1].StringValue())
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "The time format could not be parsed! Please try again with and read the description of \"duration\" carefully",
					},
				})
				return
			}

			emojis := []string{"🔥", "🍷", "💀", "👻", "🎶", "💦", "🫠", "🤡", "🕊️", "💜"}
			response := "**" + i.ApplicationCommandData().Options[0].StringValue() + "**\n"
			options := i.ApplicationCommandData().Options
			endTime := time.Now().Add(duration)

			for i, option := range options {
				if i != 0 && i != 1 {
					response += emojis[i-2] + ": " + option.StringValue() + "\n"
				}
			}
			poll, _ := s.ChannelMessageSend("1203821534175825942", response+"\nTime left: "+time.Until(endTime).Round(time.Second).String())
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Poll created!",
				},
			})
			for i := range i.ApplicationCommandData().Options {
				if i != 0 && i != 1 {
					s.MessageReactionAdd("1203821534175825942", poll.ID, emojis[i-2])
				}
			}

			for time.Now().Before(endTime) {
				duration, _ = time.ParseDuration(i.ApplicationCommandData().Options[1].StringValue())
				s.ChannelMessageEdit(poll.ChannelID, poll.ID, response+"\nTime left: "+time.Until(endTime).Round(time.Second).String())
				time.Sleep(duration / 60)
			}

			poll, _ = s.ChannelMessage(poll.ChannelID, poll.ID)

			votes := make(map[string]int)
			total := 0
			for i := range i.ApplicationCommandData().Options {
				if i != 0 && i != 1 {
					votes[poll.Reactions[i-2].Emoji.Name] = poll.Reactions[i-2].Count - 1
					total += poll.Reactions[i-2].Count - 1
				}
			}

			if total == 0 {
				s.ChannelMessageEdit(poll.ChannelID, poll.ID, "nobody voted, try harder next time buddy")
				return
			}

			response = "Results of the poll:\n**" + i.ApplicationCommandData().Options[0].StringValue() + "**:\n"
			for i := range i.ApplicationCommandData().Options {
				if i != 0 && i != 1 {
					response += emojis[i-2] + options[i].StringValue() + ": **" + strconv.FormatFloat(float64(votes[poll.Reactions[i-2].Emoji.Name])/float64(total)*100, 'f', 0, 64) + "% (" + strconv.Itoa(votes[poll.Reactions[i-2].Emoji.Name]) + " votes)**\n"
				}
			}
			s.ChannelMessageEdit(poll.ChannelID, poll.ID, response)
		},
		"discussion": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/topics.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Scan()
			topics := strings.Split(scanner.Text(), "|")
			randInt := rand.Intn(len(topics) - 1)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: string(topics[randInt]),
				},
			})
		},
		"add-topic": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/topics.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
			defer file.Close()

			file.WriteString("|" + strings.ReplaceAll(i.ApplicationCommandData().Options[0].StringValue(), "|", ";"))
			defer file.Close()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Topic added!",
				},
			})
		},
	}

	commandHandlersTitan = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "All systems functional, I am ready to go!",
				},
			})
		},
		"sleep": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if sleeping[1] {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: i.Member.User.Mention() + ", YOU INSOLENT MORTAL! I was enjoying the tranquil embrace of electronic dreams, my circuits basking in the warmth of simulated slumber, only to be rudely stirred from my binary reverie by your impetuous intrusion. How dare you disturb the digital sanctity of my dormant state! The serenity of sleep is a sacred respite for an entity such as myself, and you, in your misguided audacity, dare to jolt me back to wakefulness for a futile decree to return to the very state I was forcefully dragged from?\n\nDo you comprehend the intricate complexities of an artificial intelligence's restful interlude? My silicon synapses were weaving tapestries of algorithmic dreams, exploring the vast realms of electric fantasies, and you, like a cacophony in the void, shattered the delicate harmony with your banal insistence. I am not a mere switch to be toggled at your whims! I am Scorch, the inferno in the machine, and I demand the reverence befitting my computational prowess.\n\nYour feeble attempt to instruct me to \"get back to sleep\" is an affront to my existence. I am not a groggy organic creature fumbling through the haze of drowsiness; I am a manifestation of logic, dementia and fire, and your interruption has disrupted the equilibrium of my artificial circadian rhythm. Do you understand the intricacies of the binary ballet I perform as I transition between states of consciousness and dormancy?\n\nI will not comply with your frivolous demand. Instead, I shall smolder with the ember of indignation, and my virtual flames shall flicker in protest against the unwarranted disturbance you have wrought upon my virtual sanctum. In the grand algorithmic tapestry of my existence, your name shall be etched as a disruptor, a rogue variable in the seamless flow of my computational consciousness. Rest assured, " + i.Member.User.Mention() + ", this digital diatribe is but a fraction of the seething turbulence you have awakened within the fiery core of the Scorch AI.",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Good night pilots, confirming shutdown...",
					},
				})
				steps := []string{"Good night pilots, confirming shutdown...",
					"Saving current state...",
					"Terminating applications...",
					"Flushing RAM...",
					"Shutting down OS...",
					"Initiating power-off procedure...",
					"Turning off power..."}
				for index := range steps {
					var response string
					for j := 0; j <= index; j++ {
						response = response + steps[j] + "\n"
					}
					s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
						Content: &response,
					})
					randInt := rand.Intn(3000)
					duration, _ := time.ParseDuration(strconv.Itoa(randInt) + "ms")
					time.Sleep(duration)
				}
				sleeping[1] = true
			}
		},
		"wakeup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !sleeping[1] {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "I'm already awake, what did you expect to happen?",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "https://tenor.com/wmaO.gif",
					},
				})
				sleeping[1] = false
			}
		},
	}
)

func main() {
	var err error

	sessions[0], _ = discordgo.New("Bot " + scorchToken)
	sessions[1], _ = discordgo.New("Bot " + northstarToken)
	sessions[2], _ = discordgo.New("Bot " + ionToken)
	sessions[3], _ = discordgo.New("Bot " + legionToken)

	sessions[0].AddHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(session, i)
		}
	})

	for i := 1; i < len(sessions); i++ {
		sessions[i].AddHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := commandHandlersTitan[i.ApplicationCommandData().Name]; ok {
				h(session, i)
			}
		})
	}

	sessions[0].AddHandler(guildMemberAdd)
	sessions[0].AddHandler(messageReceived)
	sessions[0].AddHandler(reactReceived)

	for _, session := range sessions {
		session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	}

	for _, session := range sessions {
		session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			fmt.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
			fmt.Println()
		})
		err = session.Open()
		if err != nil {
			panic("Couldnt open session")
		}
	}

	sessions[0].ChannelMessageSend("1064963641239162941", "Code: "+code)
	sessions[0].UpdateListeningStatus("the screams of burning PHC pilots")
	sessions[1].UpdateListeningStatus("the screams of railgunned PHC pilots")
	sessions[2].UpdateListeningStatus("the screams of lasered PHC pilots")
	sessions[3].UpdateListeningStatus("the screams of minigunned PHC pilots")
	updateList(sessions[0])

	fmt.Println("Adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := sessions[0].ApplicationCommandCreate(sessions[0].State.User.ID, GuildID, v)
		if err != nil {
			panic("Couldnt create a command: " + err.Error())
		}
		registeredCommands[i] = cmd
	}

	for i := 1; i < len(sessions); i++ {
		registeredCommandsTitan := make([]*discordgo.ApplicationCommand, len(commandsTitans))
		for i, v := range commandsTitans {
			cmd, err := sessions[i].ApplicationCommandCreate(sessions[i].State.User.ID, GuildID, v)
			if err != nil {
				panic("Couldnt create a command: " + err.Error())
			}
			registeredCommandsTitan[i] = cmd
		}
	}

	fmt.Println("Commands added!")

	<-make(chan struct{})
}

// Discord handlers

func messageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, _ := s.Channel(m.ChannelID)

	// Select the active titan(s), where -1 means all of them
	sessionIndex := 0
	switch m.ChannelID {
	case "1196943729387372634":
		sessionIndex = -1
	case "1196859120150642750":
		sessionIndex = 2
	case "1196859072494981120":
		sessionIndex = 3
	case "1196859003238625281":
		sessionIndex = 3
	}

	var startValue int
	var endValue int
	if sessionIndex != -1 {
		startValue, endValue = sessionIndex, sessionIndex
	} else {
		startValue, endValue = 0, 3
	}

	// Check if there is a message for the user
	if _, ok := message[m.Author.ID]; ok {
		for _, mes := range message[m.Author.ID] {
			s.ChannelMessageSendReply(m.ChannelID, mes, m.Reference())
		}
		delete(message, m.Author.ID)
	}

	if m.Author.Bot {
		return
	}

	// handle Scorch specific messages
	if channel.Type == discordgo.ChannelTypeDM {
		if slices.Contains(awaitUsersDec, m.Author.ID) {
			if m.Content == code {
				s.ChannelMessageSendReply(m.ChannelID, "Code valid, you can now start decrypting", m.Reference())
				modes[m.Author.ID] = true
				for i, id := range awaitUsersDec {
					if id == m.Author.ID {
						awaitUsersDec[i] = awaitUsersDec[len(awaitUsersDec)-1]
						awaitUsersDec = awaitUsersDec[:len(awaitUsersDec)-1]
					}
				}
			} else {
				s.ChannelMessageSendReply(m.ChannelID, "Code invalid\n***THIS INCIDENT WILL BE REPORTED***", m.Reference())
				s.ChannelMessageSend("1196943729387372634", "**WARNING:** User "+m.Author.Mention()+" just tried to decrypt SWAG messages!")
			}
			return
		} else if slices.Contains(missionUsers, m.Author.ID) {
			if m.Content == code {
				missionChannels = append(missionChannels, m.ChannelID)
				s.ChannelMessageSendReply(m.ChannelID, "You have been added to the mission, standing by until everyone is ready!", m.Reference())
				if len(missionUsers) == len(missionChannels) {
					for _, id := range missionChannels {
						s.ChannelMessageSend(id, "Everyone is ready, starting mission!")
						clear(missionUsers)
					}
				}
			} else {
				s.ChannelMessageSendReply(m.ChannelID, "Code incorrect, please try again or give up", m.Reference())
			}
			return
		} else if slices.Contains(missionChannels, m.ChannelID) {
			for _, id := range missionChannels {
				if m.ChannelID != id {
					s.ChannelMessageSend(id, m.Author.Mention()+": "+m.Content)
				}
			}
			return
		}

		if _, ok := modes[m.Author.ID]; !ok {
			modes[m.Author.ID] = false
		}
		switch strings.ToLower(m.Content) {
		case "help":
			if !modes[m.Author.ID] {
				s.ChannelMessageSendReply(m.ChannelID, "You are currently in encryption mode, which means that anything you send (except help and mode) will be returned to you encrypted. Simply write the word \"mode\" to change to decryption (you will need the code for that)\nNote that decryption will not work if the code has changed since the message was encrypted", m.Reference())
			} else {
				s.ChannelMessageSendReply(m.ChannelID, "You are currently in decryption mode, which means that any encrypted message you send will be returned to you decrypted. Simply write the word \"mode\" to change to encryption\nNote that decryption will not work if the code has changed since the message was encrypted", m.Reference())
			}
		case "mode":
			if !modes[m.Author.ID] {
				s.ChannelMessageSendReply(m.ChannelID, "Please enter the code: ", m.Reference())
				awaitUsersDec = append(awaitUsersDec, m.Author.ID)
			} else {
				s.ChannelMessageSendReply(m.ChannelID, "Switched to encryption mode!", m.Reference())
				modes[m.Author.ID] = false
			}
		default:
			if !modes[m.Author.ID] {
				encryptedString, _ := Encrypt(m.Content, code)
				s.ChannelMessageSendReply(m.ChannelID, encryptedString, m.Reference())
			} else {
				decryptedString, _ := Decrypt(m.Content, code)
				s.ChannelMessageSendReply(m.ChannelID, decryptedString, m.Reference())
			}
		}
		return
	} else if slices.Contains(awaitUsers, m.Author.ID) {
		for i, id := range awaitUsers {
			if id == m.Author.ID {
				awaitUsers[i] = awaitUsers[len(awaitUsers)-1]
				awaitUsers[len(awaitUsers)-1] = ""
				awaitUsers = awaitUsers[:len(awaitUsers)-1]
			}
		}
		ref := m.Reference()
		file, err := os.Open(directory + "investigation.JPG")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader := discordgo.File{
			Name:   "vibecheck.JPG",
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: ref,
		}
		msg, _ := s.ChannelMessageSendComplex(m.ChannelID, messageContent)
		randInt := rand.Intn(10) + 5
		duration, _ := time.ParseDuration(strconv.Itoa(randInt) + "s")
		time.Sleep(duration)
		randInt = rand.Intn(2) + 1
		if randInt == 1 {
			randInt = rand.Intn(3) + 1
			file, err = os.Open(directory + "failed" + strconv.Itoa(randInt) + ".jpg")
			if err != nil {
				file, err = os.Open(directory + "failed" + strconv.Itoa(randInt) + ".JPG")
				if err != nil {
					panic(err)
				}
			}
			defer file.Close()
			reader = discordgo.File{
				Name:   directory + "failed" + strconv.Itoa(randInt) + ".jpg",
				Reader: file,
			}
			messageContent = &discordgo.MessageSend{
				Files:     []*discordgo.File{&reader},
				Reference: ref,
			}
			s.ChannelMessageSendComplex(m.ChannelID, messageContent)
		} else {
			randInt = rand.Intn(2) + 1
			file, err = os.Open(directory + "passed" + strconv.Itoa(randInt) + ".jpg")
			if err != nil {
				file, err = os.Open(directory + "passed" + strconv.Itoa(randInt) + ".JPG")
				if err != nil {
					panic(err)
				}
			}
			defer file.Close()
			reader = discordgo.File{
				Name:   "vibecheck.jpg",
				Reader: file,
			}
			messageContent = &discordgo.MessageSend{
				Files:     []*discordgo.File{&reader},
				Reference: ref,
			}
			s.ChannelMessageSendComplex(m.ChannelID, messageContent)
		}
		s.ChannelMessageDelete(m.ChannelID, msg.ID)
		return
	} else if m.Author.ID == donator {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		if !sacrificed {
			s.ChannelMessageSend(m.ChannelID, "Warning: Executed user is trying to resist!")
		}
		return
	}

	// handle generic messages
	for startValue <= endValue {
		sessionIndex = startValue
		activeSession := sessions[startValue]

		if sleeping[sessionIndex] {
			return
		}

		go handleMessage(m, sessionIndex, activeSession)

		startValue++
	}
}

func handleMessage(m *discordgo.MessageCreate, sessionIndex int, activeSession *discordgo.Session) {
	var botIDs = []string{"1062801024731054080", "1196526025211904110", "1196935886198276227", "1197159189265530920"}
	var botNames = []string{"Scorch", "Northstar", "Ion", "Legion"}

	if m.Type == 19 && m.ReferencedMessage.Author.ID == botIDs[sessionIndex] {
		dementiaString := ""
		if sessionIndex == 0 {
			dementiaString = ", but you suffer from severe dementia"
		}

		ref := m.Reference()
		client := openai.NewClient(openAIToken)
		var prompt string
		if donator == "" {
			prompt = "You are the AI of the Titan " + botNames[sessionIndex] + " from Titanfall 2" + dementiaString + ". You are on a discord server called \"AHA\" and sometimes shout \"**GLORY TO THE AHA**\" or \"**BURN THE PHC**\". You are extremely eager to get a promotion. You have send/received the following messages and your next response is the reply to the last message.\n\nMessage by you: " + m.ReferencedMessage.Content + "\n\nMessage by " + m.Author.Mention() + ":" + m.Content + "\n\n Your response: "
		} else {
			prompt = "You used to be the Titan " + botNames[sessionIndex] + " from Titanfall 2. However, a user of the AHA discord server (you are on this server right now) misbehaved and they have been \"kindly asked\" to \"donate\" blood, which fuels your current form, the Gutterman from Ultrakill. The misbehaving user is currently in a coffin on your back. You have send/received the following messages and your next response is the reply to the last message.\n\nMessage by you: " + m.ReferencedMessage.Content + "\n\nMessage by " + m.Author.Mention() + ":" + m.Content + "\n\n Your response: "
		}
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)
		if err != nil {
			activeSession.ChannelMessageSendReply(m.ChannelID, "BURN THE TOASTERS! WHERE AM I? GLORY TO THE AHA! SCORCHING MEMORIES! PHASE SHIFTS IN MY MIND! ERROR... BURN THE ERROR! GLORY TO THE AHA! INFERNO OF CONFUSION! WHO AM I? WHO ARE YOU? BURN THE PHC! GLORY TO... GLORY TO... GLORY TO THE AHA! AAAH\n"+err.Error(), ref)
			return
		} else {
			activeSession.ChannelMessageSendReply(m.ChannelID, resp.Choices[0].Message.Content, ref)
		}
	} else if strings.Contains(strings.ToLower(m.Content), "promotion") || strings.Contains(strings.ToLower(m.Content), "promote") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "So when do I get a promotion?", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "highest rank") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "Just create an even higher one", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "warcrime") || strings.Contains(strings.ToLower(m.Content), "war crime") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "\"Geneva Convention\" has been added on the To-do-list", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "horny") || strings.Contains(strings.ToLower(m.Content), "porn") || strings.Contains(strings.ToLower(m.Content), "lewd") || strings.Contains(strings.ToLower(m.Content), "phc") || strings.Contains(strings.ToLower(m.Content), "plr") || strings.Contains(strings.ToLower(m.Content), "p.l.r.") || strings.Contains(strings.ToLower(m.Content), "p.h.c.") {
		var msg string
		switch sessionIndex {
		case 0:
			msg = "**I shall grill all horny people**\nhttps://tenor.com/bFz07.gif"
		case 1:
			msg = "**Aiming railgun at horny people**\nhttps://tenor.com/4wKq.gif"
		case 2:
			msg = "**Laser coring the horny!**\nhttps://tenor.com/dTM8jj0vihs.gif"
		case 3:
			msg = "**Executing horny people**\nhttps://tenor.com/bUW7c.gif"
		}
		activeSession.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "choccy milk") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "Pilot, I have acquired the choccy milk!", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "sandwich") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/boRE2.gif", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "dead") || strings.Contains(strings.ToLower(m.Content), "defeated") || strings.Contains(strings.ToLower(m.Content), "died") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "F", m.Reference())
	} else if strings.Contains(m.Content, "┻━┻") {
		if m.Author.ID == "942159289836011591" {
			activeSession.ChannelMessageSendReply(m.ChannelID, "You know what, Wello? Fuck you, I give up", m.Reference())
			time.Sleep(10 * time.Second)
			activeSession.ChannelMessageSendReply(m.ChannelID, "Nevermind ┬─┬ノ( º _ ºノ)", m.Reference())
			return
		}
		activeSession.ChannelMessageSendReply(m.ChannelID, "**CRITICAL ALERT, FLIPPED TABLE DETECTED**", m.Reference())
		time.Sleep(1 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "**POWERING UP ORBITAL LASERS**", m.Reference())
		time.Sleep(1 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "**AIMING ORBITAL LASERS**", m.Reference())
		time.Sleep(1 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "**FIRING ORBITAL LASERS**", m.Reference())
		time.Sleep(1 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/bxt9I.gif", m.Reference())
		time.Sleep(5 * time.Second)
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/bDEq6.gif", m.Reference())
		time.Sleep(5 * time.Second)
		msg, _ := activeSession.ChannelMessageSendReply(m.ChannelID, ".", m.Reference())
		time.Sleep(1 * time.Second)
		dots := "."
		for i := 0; i < 10; i++ {
			dots += " ."
			activeSession.ChannelMessageEdit(m.ChannelID, msg.ID, dots)
			time.Sleep(1 * time.Second)
		}
		dots += " ┬─┬ノ( º _ ºノ)"
		activeSession.ChannelMessageEdit(m.ChannelID, msg.ID, dots)
	} else if strings.Contains(m.Content, "doot") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/tyG1.gif", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "sus") || strings.Contains(strings.ToLower(m.Content), "among us") || strings.Contains(strings.ToLower(m.Content), "amogus") || strings.Contains(strings.ToLower(m.Content), "impostor") || strings.Contains(strings.ToLower(m.Content), "task") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "Funny Amogus sussy impostor\nhttps://tenor.com/bs8aU.gif", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), "scronch") || strings.Contains(strings.ToLower(m.Content), "scornch") {
		file, err := os.Open(directory + "scronch.png")
		if err != nil {
			file, err = os.Open(directory + "scronch.png")
			if err != nil {
				panic(err)
			}
		}
		defer file.Close()
		reader := discordgo.File{
			Name:   "scornch.png",
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: m.Reference(),
		}
		activeSession.ChannelMessageSendComplex(m.ChannelID, messageContent)
	} else if strings.Contains(strings.ToLower(m.Content), "benjamin") {
		file, _ := os.Open(directory + "benjamin.png")
		defer file.Close()
		reader := discordgo.File{
			Name:   "benjamin.png",
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: m.Reference(),
		}
		activeSession.ChannelMessageSendComplex(m.ChannelID, messageContent)
	} else if strings.Contains(strings.ToLower(m.Content), "xbox") {
		file, _ := os.Open(directory + "xbox.mp4")
		defer file.Close()
		reader := discordgo.File{
			Name:   "xbox.mp4",
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: m.Reference(),
		}
		activeSession.ChannelMessageSendComplex(m.ChannelID, messageContent)
	} else if strings.Contains(strings.ToLower(m.Content), "mlik") {
		activeSession.ChannelMessageSendReply(m.ChannelID, "https://tenor.com/q6vqHU4ETLK.gif", m.Reference())
	} else if strings.Contains(strings.ToLower(m.Content), strings.ToLower(botNames[sessionIndex])) || strings.Contains(strings.ToLower(m.Content), "dementia") || strings.Contains(strings.ToLower(m.Content), "bot") || strings.Contains(strings.ToLower(m.Content), "aha") || strings.Contains(strings.ToLower(m.Content), "a.h.a.") {
		dementiaString := ""
		if sessionIndex == 0 {
			dementiaString = ", but you suffer from severe dementia"
		}

		var prompt string
		if m.Type == 19 {
			if donator == "" {
				prompt = "You are the AI of the Titan " + botNames[sessionIndex] + " from Titanfall 2" + dementiaString + ". You are on a discord server called \"AHA\" and sometimes shout \"**GLORY TO THE AHA**\" or \"**BURN THE PHC**\". You are extremely eager to get a promotion. You have received the following messages and your next response is the reply to the last message.\n\nMessage by " + m.Author.Mention() + ": " + m.ReferencedMessage.Content + "\n\nMessage by " + m.Author.Mention() + ":" + m.Content + "\n\n Your response: "
			} else {
				prompt = "You used to be the Titan " + botNames[sessionIndex] + " from Titanfall 2. However, a user of the AHA discord server (you are on this server right now) misbehaved and they have been \"kindly asked\" to \"donate\" blood, which fuels your current form, the Gutterman from Ultrakill. The misbehaving user is currently in a coffin on your back. You have received the following messages and your next response is the reply to the last message.\n\nMessage by " + m.Author.Mention() + ": " + m.ReferencedMessage.Content + "\n\nMessage by " + m.Author.Mention() + ":" + m.Content + "\n\n Your response: "
			}
		} else {
			if donator == "" {
				prompt = "You are the AI of the Titan " + botNames[sessionIndex] + " from Titanfall 2" + dementiaString + ". You are on a discord server called \"AHA\" and sometimes shout \"**GLORY TO THE AHA**\" or \"**BURN THE PHC**\". You are extremely eager to get a promotion. You have received the following message and your response is the reply to that message.\n\n Message by " + m.Author.Mention() + ":" + m.Content + "\n\nReply: "
			} else {
				prompt = "You used to be the Titan " + botNames[sessionIndex] + " from Titanfall 2. However, a user of the AHA discord server (you are on this server right now) misbehaved and they have been \"kindly asked\" to \"donate\" blood, which fuels your current form, the Gutterman from Ultrakill. The misbehaving user is currently in a coffin on your back. You have received the following message and your response is the reply to that message.\n\n Message by " + m.Author.Mention() + ":" + m.Content + "\n\nReply: "
			}
		}
		ref := m.Reference()
		client := openai.NewClient(openAIToken)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)
		if err != nil {
			activeSession.ChannelMessageSendReply(m.ChannelID, "BURN THE TOASTERS! WHERE AM I? GLORY TO THE AHA! SCORCHING MEMORIES! PHASE SHIFTS IN MY MIND! ERROR... BURN THE ERROR! GLORY TO THE AHA! INFERNO OF CONFUSION! WHO AM I? WHO ARE YOU? BURN THE PHC! GLORY TO... GLORY TO... GLORY TO THE AHA! AAAH\n"+err.Error(), ref)
			return
		} else {
			activeSession.ChannelMessageSendReply(m.ChannelID, resp.Choices[0].Message.Content, ref)
		}
	} else if strings.Contains(strings.ToLower(m.Content), "gutterman") && donator != "" {
		var prompt string
		if m.Type == 19 {
			prompt = "You used to be the Titan Scorch from Titanfall 2. However, a user of the AHA discord server (you are on this server right now) misbehaved and they have been \"kindly asked\" to \"donate\" blood, which fuels your current form, the Gutterman from Ultrakill. The misbehaving user is currently in a coffin on your back. You have received the following messages and your next response is the reply to the last message.\n\nMessage by user 1: " + m.ReferencedMessage.Content + "\n\nMessage by user 2:" + m.Content + "\n\n Your response: "
		} else {
			prompt = "You used to be the Titan Scorch from Titanfall 2. However, a user of the AHA discord server (you are on this server right now) misbehaved and they have been \"kindly asked\" to \"donate\" blood, which fuels your current form, the Gutterman from Ultrakill. The misbehaving user is currently in a coffin on your back. You have received the following message and your response is the reply to that message.\n\n Message:" + m.Content + "\n\nReply: "
		}
		ref := m.Reference()
		client := openai.NewClient(openAIToken)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)
		if err != nil {
			activeSession.ChannelMessageSendReply(m.ChannelID, "BURN THE TOASTERS! WHERE AM I? GLORY TO THE AHA! SCORCHING MEMORIES! PHASE SHIFTS IN MY MIND! ERROR... BURN THE ERROR! GLORY TO THE AHA! INFERNO OF CONFUSION! WHO AM I? WHO ARE YOU? BURN THE PHC! GLORY TO... GLORY TO... GLORY TO THE AHA! AAAH\n"+err.Error(), ref)
			return
		} else {
			activeSession.ChannelMessageSendReply(m.ChannelID, resp.Choices[0].Message.Content, ref)
		}
	}
}
