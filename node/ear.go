package node

import (
	"fmt"
	"net"
)

// This the place where Node open its Ear and listen for message/connection from others nodes

type Ear struct {
	Port uint16
}

func (ear *Ear) OpenEar() error {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ear.Port))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer listener.Close()
	return nil
}
