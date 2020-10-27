package main

import (
	"github.com/AlexBrin/monigo/monigo"
)

func main() {
	var err error
	err = monigo.LoadConfig("monigo.json")
	if err != nil {
		monigo.LogError("Load config err: %s", err)
		return
	}

	monigo.SetLogOptions(monigo.Config.Log)

	err = monigo.DbConnect()
	if err != nil {
		monigo.LogCritical("DB connecting err: %s", err)
		return
	}

	monigo.Start()
}