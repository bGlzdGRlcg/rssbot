package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bGlzdGRlcg/rssbot/internal/bot"
	"github.com/bGlzdGRlcg/rssbot/internal/core"
	"github.com/bGlzdGRlcg/rssbot/internal/log"
	"github.com/bGlzdGRlcg/rssbot/internal/scheduler"
)

func main() {
	appCore := core.NewCoreFormConfig()
	if err := appCore.Init(); err != nil {
		log.Fatal(err)
	}
	go handleSignal()
	b := bot.NewBot(appCore)

	task := scheduler.NewRssTask(appCore)
	task.Register(b)
	task.Start()
	b.Run()
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	os.Exit(0)
}
