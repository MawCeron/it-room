package main

import (
	"log"

	"github.com/MawCeron/it-room/internal/db"
	"github.com/MawCeron/it-room/internal/ui"
)

func main() {
	// DB initialization
	d, err := db.New("itroom.db")
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer d.Close()

	app := ui.NewApp(d)
	if err := app.Run(); err != nil {
		log.Fatalf("ui error: %v", err)
	}
}
