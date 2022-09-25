package server

import (
	"fmt"
	"os"
	"os/signal"
	"rkperes/discalendar/handle"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	dg *discordgo.Session
	sc chan os.Signal

	controller *handle.Controller

	appID string
	token string
}

func NewServer(appID string, token string) *Server {
	return &Server{
		appID: appID,
		token: token,
	}
}

func (s *Server) Run() error {
	if err := s.init(); err != nil {
		return err
	}

	s.sc = make(chan os.Signal, 1)
	signal.Notify(s.sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	if err := s.dg.Open(); err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}
	return nil
}

func (s *Server) init() error {
	dg, err := discordgo.New("Bot " + s.token)
	if err != nil {
		return fmt.Errorf("failed to create Discord session: %w", err)
	}
	s.dg = dg

	// register command handlers
	s.controller = handle.NewController(s.appID)
	if err := s.controller.RegisterCommands(dg); err != nil {
		return fmt.Errorf("failed to init controller: %w", err)
	}
	return nil
}

func (s *Server) Wait() {
	<-s.sc
}

func (s *Server) Shutdown() error {
	close(s.sc)
	return s.dg.Close()
}
