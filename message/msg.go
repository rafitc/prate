/*
This is the package which act as a data layer.
Anything related to message will be here. It can be normal chat message, TCP connection, authentication related msgs etc.
*/

package message

import (
	"fmt"
	"log"
	"net"
)

type Speak struct {
	Ip   string
	Port uint16
	Message
}

type Message struct {
	User string
	Body string
}

func (m *Message) ReadMsg() string {
	return fmt.Sprintf("One new msg from %s : %s\n", m.User, m.Body)
}

// Function does simple TCP write in give IP and port

func (m *Speak) SendMessage() {
	// connect to given ip and port
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", m.Ip, m.Port))
	if err != nil {
		// Error in this specific ip,
		// since its a prate, not a critical info,
		// if its not able to deliver, leave it
		log.Printf("Error while sending to %s :- %s", m.Ip, err.Error())
	}
	defer conn.Close()

	_, err = conn.Write([]byte(m.Message.Body))
	if err != nil {
		// again some error. just leave it
		log.Printf("Error while writing msg to %s:%s", m.Ip, err.Error())
		return
	}
	// so all done, just leave it
}
