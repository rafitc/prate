package main

import (
	"fmt"
	"os"
	msg "prate/message"
	"prate/terminal"
	"time"

	"golang.org/x/exp/rand"
)

func main() {
	fmt.Println("Starting Prate Chat")
	// Assume the user setup is done.
	// Do the server setup, ie open a port
	// Setup a gocoroutine to search for do nmcli for given port
	// Keep the list, and updated.
	// Get message from user and pass into the channel, which takes the msg and send things to the given port as tcp.
	// update everything on cli terminal ui
	newmsg := msg.Message{User: "rafi", Body: "hello"}
	fmt.Printf("%s", newmsg.ReadMsg())
	generateMessage()

	// Get the terminal UI here, so we can test atlreast the integration
	m := terminal.InitTerminal()
	rand.Seed(uint64(time.Now().Unix()))

	// Fire a goroutine
	go func() {
		for {
			pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond
			time.Sleep(pause)

			// Send the Bubble Tea program a message from outside the
			// tea.Program. This will block until it is ready to receive
			// messages.
			m.Send(terminal.NewMsg{User: "rafi", Body: fmt.Sprintf("current time %s", time.Now().String())})
		}
	}()
	// Start the Session
	_, err := m.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

// For PoC, contnuesly steam the message into a channel, and get that in terminal UI

func generateMessage() {
	fmt.Printf("I'll push some random messages \n")
}
