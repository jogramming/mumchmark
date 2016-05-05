// Simple mumble stress testing tool

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	"github.com/layeh/gumble/gumbleutil"
	_ "github.com/layeh/gumble/opus"
	"log"
	"net"
)

var (
	clients    []*gumble.Client
	numClient  = flag.Int("clients", 10, "Number of clients to spawn")
	addr       = flag.String("address", "localhost:64738", "Address of the mumble server")
	password   = flag.String("pw", "", "Password to connect with")
	curChannel = uint32(0)
)

func main() {
	flag.Parse()

	for i := 0; i < *numClient; i++ {
		err := spawnClient(fmt.Sprintf("Mumchmark_%d", i), *password, *addr)
		if err != nil {
			fmt.Printf("Failed to connect to server %q: %s\n", *addr, err.Error())
			return
		}
	}

	loop()

	for _, v := range clients {
		v.Disconnect()
	}
}

func loop() {
	for {
		fmt.Println("=======================")
		fmt.Println("What do you want to do?")
		fmt.Println("[q] Quit?")
		fmt.Println("[a] Play some audio (plays audio.mp3)")
		fmt.Println("[t] Send a text message")

		option := ""
		fmt.Scanln(&option)
		switch option {
		case "t":
			sendText()
		case "a":
			playAudio()
		case "q":
			fmt.Println("Quitting")
			return
		default:
			continue
		}
	}
}

func inputAmount() (int, bool) {
	fmt.Printf("Input a number: ")
	num := 0
	fmt.Scanln(&num)
	if num <= 0 {
		return len(clients), true
	} else if num > len(clients) {
		fmt.Println("num is >clients")
		return 0, false
	}
	return num, true
}

func sendText() {
	fmt.Println("What should the message be? (leave empty for \"testing\"")
	str := "testing"
	fmt.Scanln(&str)
	fmt.Printf("Message set to %q, How may clients should send this? leave empty for %d\n", str, len(clients))

	num, ok := inputAmount()
	if !ok {
		return
	}

	for i := 0; i < num; i++ {
		client := clients[i]
		channel := client.Channels[curChannel]
		channel.Send(str, false)
	}
}

func playAudio() {
	fmt.Printf("How many clients should play audio.mp3? leave empty for %d\n", len(clients))
	num, ok := inputAmount()
	if !ok {
		return
	}

	for i := 0; i < num; i++ {
		client := clients[i]
		ff := gumbleffmpeg.New(client, gumbleffmpeg.SourceFile("audio.mp3"))
		err := ff.Play()
		if err != nil {
			log.Println("Error playing audio: ", err)
		}
	}
}

func spawnClient(user, pw, server string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	config := &gumble.Config{
		Username:       user,
		Password:       pw,
		AudioInterval:  gumble.AudioDefaultInterval,
		AudioDataBytes: gumble.AudioDefaultDataBytes,
	}

	config.Attach(gumbleutil.Listener{
		TextMessage:   textMessageHandler,
		Connect:       connectHandler,
		Disconnect:    dcHandler,
		ChannelChange: channelChangeHandler,
	})
	client, err := gumble.DialWithDialer(new(net.Dialer), server, config, tlsConfig)
	if err != nil {
		return err
	}

	clients = append(clients, client)
	log.Printf("mumchmark user %q Connected to server %q!\n", user, server)
	return nil
}

func textMessageHandler(msg *gumble.TextMessageEvent) {
	//fmt.Printf("Received text message: %s\n", msg.Message)
}

func connectHandler(c *gumble.ConnectEvent) {
	msg := "(none)"
	if c.WelcomeMessage != nil {
		msg = *c.WelcomeMessage
	}
	fmt.Printf("Connected to server, welcome message: %q\n", msg)
}

func dcHandler(c *gumble.DisconnectEvent) {
	fmt.Println("Disconnected..", c.String)
}

func channelChangeHandler(c *gumble.ChannelChangeEvent) {
	curChannel = c.Channel.ID
}
