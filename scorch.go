package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
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
	s.ChannelMessageSend("1195135473643958316", m.Mention()+", welcome to the AHA! Please consider using /register")
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
