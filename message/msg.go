/*
This is the package which act as a data layer.
Anything related to message will be here. It can be normal chat message, TCP connection, authentication related msgs etc.
*/

package message

import "fmt"

type Message struct {
	User string
	Body string
}

func (m *Message) ReadMsg() string {
	return fmt.Sprintf("One new msg from %s : %s\n", m.User, m.Body)
}
