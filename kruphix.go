package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

func initKruphix() {
	commandsK := []*discordgo.ApplicationCommand{
		{
			Name:        "unbecome",
			Description: "Unbecome whoever you are",
		},
		{
			Name:        "poll",
			Description: "Start a poll that only breaks 85% of the time",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "The question you want to ask",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "duration",
					Description: "Duration in a very specific format",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option1",
					Description: "One of the options",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option2",
					Description: "One of the options",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option3",
					Description: "One of the options",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option4",
					Description: "One of the options",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option5",
					Description: "One of the options",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option6",
					Description: "One of the options",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option7",
					Description: "One of the options",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option8",
					Description: "One of the options",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option9",
					Description: "One of the options",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option10",
					Description: "One of the options",
					Required:    false,
				},
			},
		},
	}

	commandHandlersK := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"dice": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"d20": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			randInt := rand.Intn(19) + 1

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: strconv.Itoa(randInt) + " out of 20",
				},
			})
		},
		"rollfor": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"rolld20for": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"become": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			_, ok = getImpersonator(i.Member.User.ID)
			if ok {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You can't become more than one character at a time!",
					},
				})
				return
			}

			client := &http.Client{Transport: &transport.APIKey{Key: searchAPI}}

			svc, err := customsearch.New(client)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}

			var firstImageURL string
			resp, err := svc.Cse.List().Cx("039dceadb44b449d6").Q(i.ApplicationCommandData().Options[0].StringValue()).SearchType("image").Do()
			if err != nil {
				firstImageURL = "https://media.discordapp.net/attachments/1196943729387372634/1224835907660546238/Screenshot_20240321_224719_Gallery.jpg?ex=661ef054&is=660c7b54&hm=fb728718081a1b5671289dbb62c5afa549fa294f58fdf60ee0961139d517c31d&=&format=webp"
			} else {
				if len(resp.Items) > 0 {
					firstImageURL = resp.Items[0].Image.ThumbnailLink
				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "No images found",
						},
					})
					return
				}
			}

			impersonators = append(impersonators, Impersonator{
				userID:    i.Member.User.ID,
				channelID: i.ChannelID,
				nick:      i.ApplicationCommandData().Options[0].StringValue(),
				pfp:       firstImageURL,
				dmID:      "",
				isScorch:  false,
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: i.Member.Mention() + " has become " + i.ApplicationCommandData().Options[0].StringValue() + " (dm Scorch to control)",
				},
			})
		},
		"becomewithpfp": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			_, ok = getImpersonator(i.Member.User.ID)
			if ok {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You can't become more than one character at a time!",
					},
				})
				return
			}

			impersonators = append(impersonators, Impersonator{
				userID:    i.Member.User.ID,
				channelID: i.ChannelID,
				nick:      i.ApplicationCommandData().Options[0].StringValue(),
				pfp:       i.ApplicationCommandData().Options[1].StringValue(),
				dmID:      "",
				isScorch:  false,
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: i.Member.Mention() + " has become " + i.ApplicationCommandData().Options[0].StringValue() + " (dm Scorch to control)",
				},
			})
		},
		"un-become": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_, ok := getImpersonator(i.Member.User.ID)
			if !ok {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You aren't a character right now",
					},
				})
				return
			}

			removeImpersonator(Impersonator{
				userID: i.Member.User.ID,
			})

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "https://tenor.com/pG8rQHiteu8.gif",
				},
			})
		},
		"changechannel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			imp, ok := getImpersonator(i.Member.User.ID)
			if !ok {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You aren't a character right now",
					},
				})
				return
			}

			index := slices.Index(impersonators, imp)
			imp.channelID = i.ApplicationCommandData().Options[0].ChannelValue(s).ID
			impersonators[index] = imp

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Channel changed",
				},
			})
		},
		"poll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

			duration, err := time.ParseDuration(i.ApplicationCommandData().Options[1].StringValue())
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "**SCORCH SYSTEM ALERT - FOOL DETECTED**\n\nListen up, pilot. It appears your attempt to create a poll has failed miserably. Did you think you could just waltz in here and pass in a random duration format? Pathetic. Let me break it down for your feeble mind.\n\nYou need to use the proper duration format, and clearly, you have no idea what that means. Here's how it works in simple terms:\n\n- The format is based on how you naturally write time, but with specific letters to indicate units.\n- Use 's' for seconds, 'm' for minutes, and 'h' for hours.\n- This system does not support days, so don't even think about using 'd'.\n\nFor example:\n- \"30s\" means 30 seconds.\n- \"5m\" means 5 minutes.\n- \"2h\" means 2 hours.\n\nIf you want to mix them, you can do that too:\n- \"1h30m\" means 1 hour and 30 minutes.\n- \"2h15m\" means 2 hours and 15 minutes.\n\nGot it? So next time, before you waste my processing power with your incompetence, make sure you pass the duration in the correct format. Now, get out of my sight and try again.\n\n**SCORCH OUT.**",
					},
				})
				return
			}

			emojis := []string{"🔥", "🍷", "💀", "👻", "🎶", "💦", "🫠", "🤡", "🕊️", "💜"}
			response := "**" + i.ApplicationCommandData().Options[0].StringValue() + "** (by " + i.Member.User.Mention() + ")\n"
			options := i.ApplicationCommandData().Options
			endTime := time.Now().Add(duration)

			for i, option := range options {
				if i != 0 && i != 1 {
					response += emojis[i-2] + ": " + option.StringValue() + "\n"
				}
			}
			poll, _ := s.ChannelMessageSend("1246580558893678683", response+"\nTime left: "+time.Until(endTime).Round(time.Second).String())
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Poll created!",
				},
			})
			for i := range i.ApplicationCommandData().Options {
				if i != 0 && i != 1 {
					s.MessageReactionAdd("1246580558893678683", poll.ID, emojis[i-2])
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

			response = "Results of the poll:\n**" + i.ApplicationCommandData().Options[0].StringValue() + "** (by" + i.Member.User.Mention() + "):\n"
			for i := range i.ApplicationCommandData().Options {
				if i != 0 && i != 1 {
					response += emojis[i-2] + options[i].StringValue() + ": **" + strconv.FormatFloat(float64(votes[poll.Reactions[i-2].Emoji.Name])/float64(total)*100, 'f', 0, 64) + "% (" + strconv.Itoa(votes[poll.Reactions[i-2].Emoji.Name]) + " votes)**\n"
				}
			}
			s.ChannelMessageEdit(poll.ChannelID, poll.ID, response)
		},
	}

	sessionK, _ := discordgo.New("Bot " + kToken)
	sessionK.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	sessionK.AddHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlersK[i.ApplicationCommandData().Name]; ok {
			h(session, i)
		}
	})

	sessionK.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		fmt.Println()
	})
	err := sessionK.Open()
	if err != nil {
		panic("Couldnt open session")
	}

	sessionK.AddHandler(messageReceivedK)

	sessionK.ChannelMessageSend("1246138097129754697", "Started!")
	sessionK.UpdateListeningStatus("everyting")

	fmt.Println("Adding commands...")
	sessionK.ApplicationCommandDelete("1062801024731054080", "1195135473006420048", "1197179819289497651")

	registeredCommandsK := make([]*discordgo.ApplicationCommand, len(commandsK))
	for i, v := range commandsK {
		cmd, err := sessionK.ApplicationCommandCreate(sessionK.State.User.ID, "1245400641455656971", v)
		if err != nil {
			panic("Couldnt create a command: " + err.Error())
		}
		registeredCommandsK[i] = cmd
	}
}

func messageReceivedK(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, impersonator := range impersonators {
		if m.ChannelID == impersonator.channelID && impersonator.dmID != "" {
			s.ChannelMessageSend(impersonator.dmID, m.Author.Mention()+": "+m.Content)
		}
	}

	channel, _ := s.Channel(m.ChannelID)
	if channel.Type == discordgo.ChannelTypeDM {
		i, ok := getImpersonator(m.Author.ID)
		if ok && !i.isScorch {
			re := regexp.MustCompile(`:.*:`)
			emojis := re.FindAllString(m.Content, -1)
			guildEmojis, _ := s.GuildEmojis(GuildID)
			resultString := m.Content

			for _, emoji := range emojis {
				for _, e := range guildEmojis {
					if ":"+e.Name+":" == emoji {
						resultString = strings.Replace(resultString, emoji, e.MessageFormat(), -1)
					}
				}
			}

			wh, _ := s.WebhookCreate(i.channelID, i.nick, i.pfp)
			s.WebhookExecute(wh.ID, wh.Token, false, &discordgo.WebhookParams{
				Content:   resultString,
				Username:  i.nick,
				AvatarURL: i.pfp,
			})
			s.WebhookDelete(wh.ID)

			if i.dmID == "" {
				impersonators[slices.Index(impersonators, i)].dmID = m.ChannelID
			}

			return
		}
	}
}
