package main

import (
	"log"
	"time"

	"github.com/gofinance/ib"
)

func main() {
	opts := ib.EngineOptions{
		Gateway:          "localhost:7496",
		DumpConversation: true,
	}

	engine, err := ib.NewEngine(opts)
	if err != nil {
		log.Fatalf("ib.NewEngine: %s", err)
	}
	defer engine.Stop()

	manager, err := ib.NewPrimaryAccountManager(engine)
	if err != nil {
		log.Fatalf("ib.NewPrimaryAccountManager: %s", err)
	}
	defer manager.Close()

	minUpdates := 1
	updates, err := ib.SinkManager(manager, 30*time.Second, minUpdates)
	if err != nil {
		log.Fatalf("ib.SinkManager returned an error after %d updates: %v", updates, err)
	}

	if updates < minUpdates {
		log.Fatalf("ib.SinkManager returned %d updates (expected >= %d)", updates, minUpdates)
	}

	portfolio := manager.Portfolio()
	for k, v := range portfolio {
		log.Print("++++++++:", k.AccountCode, k.ContractID)
		log.Printf("\t: %d, %s", v.Position, v.Contract.LocalSymbol)
	}

	if b, ok := <-manager.Refresh(); ok {
		log.Fatalf("Expected the refresh channel to be closed, but got %t", b)
	}
}
