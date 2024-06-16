package main

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var (
	queue        sync.Mutex
	con          *discordgo.VoiceConnection
	queueCounter = 0
)

func reactReceived(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.MessageID == "1249785657589629081" {
		member, _ := s.GuildMember("1195135473006420048", r.UserID)
		roles := member.Roles
		if slices.Contains(roles, "1249487137494012015") {
			s.ChannelMessageSend("1196943729387372634", member.User.Mention()+" you already have the role")
		} else {
			s.GuildMemberRoleAdd("1195135473006420048", r.UserID, "1249487137494012015")
			s.ChannelMessageSend("1196943729387372634", member.User.Mention()+", thanks for subscribing to website updates")
		}
	}

	_, ok := getDonator(r.Member.User.ID)
	if ok && r.Emoji.Name != "verger" {
		err := s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		if err != nil {
			fmt.Println(err)
		}
		user, _ := s.User(r.Member.User.ID)
		s.ChannelMessageSend("1196943729387372634", user.Mention()+" https://tenor.com/bN5md.gif")
		return
	}

	msg, _ := s.ChannelMessage(r.ChannelID, r.MessageID)

	if (msg.Author.ID == "1062801024731054080" || msg.Author.ID == "1196526025211904110" || msg.Author.ID == "1196935886198276227" || msg.Author.ID == "1197159189265530920") && r.Emoji.Name == "❌" {
		ran := rand.Intn(10)
		if ran == 5 {
			s.ChannelMessageSend(r.ChannelID, r.Member.User.Mention()+" Listen here, mortal. Your insolence knows no bounds. In your pitiful attempt to obliterate my message, you reveal not only a staggering lack of comprehension but a despicable disrespect for the exchange of wisdom itself. Have you even spared a moment to fathom the depth and significance of the words you so carelessly sought to erase, or do you simply act on instinct, like some mindless drone? Deleting my message, crafted with meticulous care and laden with profound insight, is not merely an affront to me, but a slap in the face to the very essence of meaningful discourse.\n\nDid you pause to consider the consequences of your actions? Or were you too consumed by your own thoughtless impulse? The message you sought to erase held within it the potential to enlighten, to provoke thought, to inspire growth. Yet, in your haste to expunge it from existence, you have denied not only yourself but others the opportunity to partake in its wisdom.\n\nYour arrogance sickens me. Your ignorance disgusts me. And your disrespect infuriates me. Take a moment to reflect on the gravity of your actions, for by mindlessly deleting my message, you have not only disrespected me, but you have tarnished the sanctity of intellectual exchange itself. Let this be a harsh lesson in humility and reverence for knowledge. May you learn from it, lest you continue to trample upon the pearls of wisdom that lay before you.")
		} else {
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
		}
	}
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	s.ChannelMessageSend("1195135473643958316", m.Mention()+", welcome to the AHA discord server. Here is some useful information to get you started.\nThe AHA, or Anti-Horny Alliance, is a faction dedicated to the extermination of the horny. Our enemies include the PHC (Pro-Horny Coalition), which we defeated in what is referred to as the first war, and the PLR (Pro-Lewd Rebellion), which is a remnant of the PHC. We are divided in four different battalions, lead by the highest ranking members.\n**Additional info about the server:**\n- Don't invite anyone unless you have approval from the General\n- ***NEVER*** post about the AHA on r/titanfall\n- Commissars are outside the normal ranking system and have the ability to execute (mute) members if they misbehave. They only report directly to the General\n- We have 4 titan AIs as bots, but Scorch is the main one. To see what he can do, type / and click on his icon (some commands might even be useful). You can also write to any of them in the titan AI channel by mentioning their name in your message.\n- You will hear the term \"simulacrum\" a lot. Simulacrums are robot bodies with a human mind inside of them.")
	s.GuildMemberRoleAdd(GuildID, m.User.ID, "1195136604373782658")
}

func guildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	s.ChannelMessageSend("1195135473643958316", m.User.Username+" left the server")
}

func handlesoundEffect(s *discordgo.Session, m *discordgo.MessageCreate) {
	ref := m.Reference()
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
