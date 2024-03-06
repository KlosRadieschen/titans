package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

var (
	queue        sync.Mutex
	con          *discordgo.VoiceConnection
	queueCounter = 0
)

func reactReceived(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	msg, _ := s.ChannelMessage(r.ChannelID, r.MessageID)

	if (msg.Author.ID == "1062801024731054080" || msg.Author.ID == "1196526025211904110" || msg.Author.ID == "1196935886198276227" || msg.Author.ID == "1197159189265530920") && r.Emoji.Name == "❌" {
		s.ChannelMessageDelete(r.ChannelID, r.MessageID)
	}
}

func updateList(s *discordgo.Session) {
	var lowerCaseNames = []string{"pc", "xbox", "playstation", "ps4", "ps5"}
	var displayNames = []string{"PC", "Xbox", "Playstation", "PS4", "PS5"}
	var data [][]string

	for i := range lowerCaseNames {
		file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), ",")
			if len(parts) == 3 && strings.ToLower(parts[2]) == lowerCaseNames[i] {
				member, _ := s.GuildMember("1195135473006420048", parts[0])
				rank := ""
				battalion := ""
				for _, roleID := range member.Roles {
					role, _ := s.State.Role("1195135473006420048", roleID)
					if strings.Contains(role.Name, "Battalion") || strings.Contains(role.Name, "Operative") {
						battalion = role.Name
					} else if role.Name != "PC" && role.Name != "Xbox" && role.Name != "PlayStation" {
						rank = role.Name
					}
				}
				var row = []string{displayNames[i], member.User.Username, parts[1], rank, battalion}
				data = append(data, row)
			}
		}
	}

	file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) == 3 && strings.ToLower(parts[2]) != "playstation" && strings.ToLower(parts[2]) != "pc" && strings.ToLower(parts[2]) != "xbox" && strings.ToLower(parts[2]) != "ps4" && strings.ToLower(parts[2]) != "ps5" {
			user, _ := s.GuildMember("1195135473006420048", parts[0])
			rank, _ := s.State.Role("1195135473006420048", user.Roles[0])
			var battalionName string
			if len(user.Roles) < 2 {
				battalionName = ""
			} else {
				battalion, _ := s.State.Role("1195135473006420048", user.Roles[1])
				battalionName = battalion.Name
			}

			var row = []string{parts[2], user.User.Username, parts[1], rank.Name, battalionName}
			data = append(data, row)
		}
	}

	var builder strings.Builder
	writer := &builder
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Platform", "User", "in-game name", "Rank", "Battalion"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()

	tableString := builder.String()
	drawText(SplitLines(tableString))

	file, err := os.Open(directory + "table.png")
	if err != nil {
		file, err = os.Open(directory + "table.png")
		if err != nil {
			panic(err)
		}
	}
	defer file.Close()

	s.State.MaxMessageCount = 100
	channel, _ := s.State.Channel("1196072273686315008")
	for _, msg := range channel.Messages {
		err := s.ChannelMessageDelete(channel.ID, msg.ID)
		if err != nil {
			panic(err.Error())
		}
	}

	reader := discordgo.File{
		Name:   "table.png",
		Reader: file,
	}
	messageSend := &discordgo.MessageSend{
		Files: []*discordgo.File{&reader},
	}
	s.ChannelMessageSendComplex(channel.ID, messageSend)
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	s.ChannelMessageSend("1195135473643958316", m.Mention()+", welcome to the AHA discord server. Here is some useful information to get you started.\nThe AHA, or Anti-Horny Alliance, is a faction dedicated to the extermination of the horny. Our enemies include the PHC (Pro-Horny Coalition), which we defeated in what is referred to as the first war, and the PLR (Pro-Lewd Coalition), which is a remnant of the PHC.\n**We are divided in four different battalions, lead by the highest ranking members:**\n- The 1st battalion is directly controlled by our General and leader, Samp. They control our main base on the planet Harmony and our flagship, the AHF Resolute.\n- The 2nd battalion is lead by Lieutenant General TU-8271 and controls the outpost on the planet Typhon and their main ship, the AHF Midas.\n- The 3rd battalion is lead by Vice Admiral Storm and controls the outpost on the planet Orthros and their main ship, the AHF Rift Walker.\n- The 4th battalion is lead by Commander Klos and controlls the planet Laythe and their main ship, the AHF Meruda.\n- The SWAG, or Special Warfare Assault Group is lead by Captain Voodoo-6. They don't control a planet and stay on Harmony but they do have a ship called the AHF Infiltrator.\n**Additional info about the server:**\n- Don't invite anyone unless you have approval from the General\n- ***NEVER*** post about the AHA on r/titanfall\n- Commissars are outside the normal ranking system and have the ability to execute (mute) members if they misbehave. They only report directly to the General\n- We have 4 titan AIs as bots, but Scorch is the main one. To see what he can do, type / and click on his icon (some commands might even be useful). You can also write to any of them in the titan AI channel by mentioning their name in your message.\n- You will hear the term \"simulacrum\" a lot. Simulacrums are robot bodies with a human mind inside of them.")
	s.GuildMemberRoleAdd(GuildID, m.User.ID, "1195136604373782658")
}

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func handlesoundEffect(s *discordgo.Session, m *discordgo.MessageCreate) {
	ref := m.Reference()

	s.ChannelMessageDelete(m.ChannelID, m.ID)
	files, err := os.ReadDir("/home/Nicolas/go-workspace/src/knilchbot/sfx")
	if err != nil {
		fmt.Println("Oops! Something went wrong:", err)
		return
	}
	commands := make(map[int]string)
	counter := 1
	for _, file := range files {
		commands[counter] = file.Name()
		counter++
	}
	vs, err := s.State.VoiceState("1195135473006420048", m.Author.ID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "You have to be in a voicechat", ref)
	} else {
		sfxNumber, err := strconv.Atoi(m.Content)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, "Thats not even a number", ref)
			return
		}
		if sfxNumber > len(commands) {
			s.ChannelMessageSendReply(m.ChannelID, "That number doesn't exist", ref)
			return
		}
		queueCounter++
		queue.Lock()
		if con == nil {
			con, _ = s.ChannelVoiceJoin("1195135473006420048", vs.ChannelID, false, false)
			for !con.Ready {
				time.Sleep(100 * time.Millisecond)
			}
		}
		dgvoice.PlayAudioFile(con, "/home/Nicolas/go-workspace/src/knilchbot/sfx/"+commands[sfxNumber], make(chan bool))
		queueCounter--
		queue.Unlock()
		if queueCounter == 0 {
			con.Disconnect()
			con = nil
		}
	}
}
