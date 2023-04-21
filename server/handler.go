package main

import "fmt"

type Server struct {

}

func (srvr *Server) DeliverValue(msg string, _ *interface{}) error {
	fmt.Println("Registered an incoming messsage")
	return nil
}