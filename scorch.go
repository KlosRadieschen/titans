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
	var displayNames = []string{"PC", "Xbox", "Playstations", "PS4", "PS4"}
	var data [][]string

	for i := range lowerCaseNames {
		file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), ",")
			if len(parts) == 3 && strings.ToLower(parts[2]) == lowerCaseNames[i] {
				member, err := s.GuildMember("1195135473006420048", parts[0])
				if err != nil {
					panic(parts[0] + ": " + err.Error())
				}
				rank, _ := s.State.Role("1195135473006420048", member.Roles[0])
				var battalionName string
				if len(member.Roles) < 2 {
					battalionName = ""
				} else {
					battalion, _ := s.State.Role("1195135473006420048", member.Roles[1])
					battalionName = battalion.Name
				}
				var row = []string{displayNames[i], member.Mention(), parts[1], rank.Name, battalionName}
				data = append(data, row)
				data = append(data, []string{""})
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

			var row = []string{parts[2], user.Mention(), parts[1], rank.Name, battalionName}
			data = append(data, row)
		}
	}

	var builder strings.Builder
	writer := &builder
	table := tablewriter.NewWriter(writer)
	//table.SetHeader([]string{"Platform", "User", "in-game name", "Rank", "Battalion"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()

	embed := &discordgo.MessageEmbed{
		Title:       "Registered members:",
		Description: builder.String(),
		Color:       0xF73718,
	}
	_, err := s.ChannelMessageEditEmbed("1196072273686315008", "1196079691577163798", embed)
	if err != nil {
		panic(err.Error())
	}
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	s.ChannelMessageSend("1195135473643958316", m.Mention()+", welcome to the AHA! Please consider using /register")
	s.GuildMemberRoleAdd(GuildID, m.User.ID, "1195136604373782658")
}
