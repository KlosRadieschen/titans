package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func reactReceived(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	msg, _ := s.ChannelMessage(r.ChannelID, r.MessageID)

	if (msg.Author.ID == "1062801024731054080" || msg.Author.ID == "1196526025211904110" || msg.Author.ID == "1196935886198276227" || msg.Author.ID == "1197159189265530920") && r.Emoji.Name == "❌" {
		s.ChannelMessageDelete(r.ChannelID, r.MessageID)
	}
}

func updateList(s *discordgo.Session) {
	var response string
	var lowerCaseNames = []string{"pc", "xbox", "playstation", "ps4", "ps5"}
	var displayNames = []string{"PC", "Xbox", "Playstations", "PS4", "PS4"}

	for i := range lowerCaseNames {
		appendEntryFromFile(response, lowerCaseNames[i], displayNames[i])
	}

	file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) == 3 && strings.ToLower(parts[2]) != "playstation" && strings.ToLower(parts[2]) != "pc" && strings.ToLower(parts[2]) != "xbox" && strings.ToLower(parts[2]) != "ps4" && strings.ToLower(parts[2]) != "ps5" {
			response += fmt.Sprintf("%s: %s\n", "**"+parts[2]+"**", parts[1])
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Registered members:",
		Description: response,
		Color:       0xF73718,
	}
	_, err := s.ChannelMessageEditEmbed("1196072273686315008", "1196079691577163798", embed)
	if err != nil {
		panic(err.Error())
	}
}

func appendEntryFromFile(response string, lowerCaseName string, displayName string) {
	file, _ := os.OpenFile("/home/Nicolas/go-workspace/src/titans/members.csv", os.O_APPEND|os.O_RDWR|os.O_SYNC, os.ModeAppend)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) == 3 && strings.ToLower(parts[2]) != "playstation" && strings.ToLower(parts[2]) != "pc" && strings.ToLower(parts[2]) != "xbox" && strings.ToLower(parts[2]) != "ps4" && strings.ToLower(parts[2]) != "ps5" {
			response += fmt.Sprintf("%s: %s\n", "**"+parts[2]+"**", parts[1])
		}
	}
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	s.ChannelMessageSend("1195135473643958316", m.Mention()+", welcome to the AHA! Please consider using /register")
	s.GuildMemberRoleAdd(GuildID, m.User.ID, "1195136604373782658")
}
