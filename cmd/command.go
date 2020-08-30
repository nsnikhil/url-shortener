package main

import "log"

const (
	serveCommand    = "serve"
	migrateCommand  = "migrate"
	rollbackCommand = "rollback"
)

func commands() map[string]func() {
	return map[string]func(){
		serveCommand:    serve,
		migrateCommand:  runMigrations,
		rollbackCommand: rollBackMigrations,
	}
}

func execute(cmd string) {
	run, ok := commands()[cmd]
	if !ok {
		log.Fatal("invalid command")
	}

	run()
}
