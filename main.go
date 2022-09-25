package main

import (
	"flag"
	"log"
	"rkperes/discalendar/server"
)

var (
	Token string
	AppID string
)

func init() {
	flag.StringVar(&AppID, "a", "", "Application ID")
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	s := server.NewServer(AppID, Token)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
	s.Wait()
	if err := s.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
