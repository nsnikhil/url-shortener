package main

import (
	"log"
	"urlshortner/pkg/app"
	"urlshortner/pkg/store"
)

const (
	serveCommand    = "serve"
	migrateCommand  = "migrate"
	rollbackCommand = "rollback"
)

func commands() map[string]func() {
	return map[string]func(){
		serveCommand:    app.Start,
		migrateCommand:  store.RunMigrations,
		rollbackCommand: store.RollBackMigrations,
	}
}

func execute(cmd string) {
	run, ok := commands()[cmd]
	if !ok {
		log.Fatal("invalid command")
	}

	run()
}
