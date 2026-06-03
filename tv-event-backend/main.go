package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Event definerer strukturen på live-oppdateringene som sendes til appen
type Event struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`  // f.eks. "GOAL", "AD_START", "AD_END", "YELLOW_CARD"
	Title     string `json:"title"` // f.eks. "MÅÅÅL! Brann scoring"
	Timestamp string `json:"timestamp"`
}

// Broker håndterer alle tilkoblede app-klienter
type Broker struct {
	clients    map[chan Event]bool
	newClients chan chan Event
	defClients chan chan Event
	mutex      sync.Mutex
}

func main() {
	broker := &Broker{
		clients:    make(map[chan Event]bool),
		newClients: make(chan chan Event),
		defClients: make(chan chan Event),
	}

	// Start brokeren i en egen tråd (Goroutine)
	go broker.listen()

	// Simuler en live-feed som genererer hendelser i bakgrunnen
	go simulateLiveEvents(broker)

	// API Endepunkt som Kotlin-appen skal koble seg til
	http.HandleFunc("/api/live-events", broker.serveHTTP)

	log.Println("Server startet på http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// serveHTTP oppretter en vedvarende tilkobling til klienten (SSE)
func (b *Broker) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// Sett korrekte headere for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // For testing

	// Hver klient får sin egen kanal
	clientChan := make(chan Event)
	b.newClients <- clientChan

	// Fjern klienten når tilkoblingen brytes
	defer func() {
		b.defClients <- clientChan
	}()

	// Lytt til "r.Context().Done()" for å fange opp at appen lukkes/kobler ut
	for {
		select {
		case event := <-clientChan:
			jsonData, _ := json.Marshal(event)
			// SSE formatet krever "data: <melding>\n\n"
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			w.(http.Flusher).Flush() // Tving dataene over nettverket umiddelbart
		case <-r.Context().Done():
			return
		}
	}
}

func (b *Broker) listen() {
	for {
		select {
		case client := <-b.newClients:
			b.mutex.Lock()
			b.clients[client] = true
			b.mutex.Unlock()
			log.Printf("Ny app koblet til. Totalt antall klienter: %d", len(b.clients))

		case client := <-b.defClients:
			b.mutex.Lock()
			delete(b.clients, client)
			close(client)
			b.mutex.Unlock()
			log.Printf("App koblet fra. Totalt antall klienter: %d", len(b.clients))
		}
	}
}

// simulateLiveEvents genererer falske TV-events hvert 5. sekund for å simulere en live-feed
func simulateLiveEvents(b *Broker) {
	mockEvents := []Event{
		{Type: "MATCH_START", Title: "Kampen har startet! SK Brann - Bodø/Glimt"},
		{Type: "GOAL", Title: "MÅÅÅL! SK Brann tar ledelsen 1-0!"},
		{Type: "YELLOW_CARD", Title: "Gult kort til Bodø/Glimt"},
		{Type: "AD_START", Title: "Kommersiell pause – starter reklamesending"},
		{Type: "AD_END", Title: "Reklame ferdig – Velkommen tilbake til kampen"},
	}

	i := 0
	for {
		time.Sleep(5 * time.Second)

		b.mutex.Lock()
		if len(b.clients) > 0 {
			event := mockEvents[i%len(mockEvents)]
			event.ID = time.Now().UnixNano()
			event.Timestamp = time.Now().Format("15:04:05")

			log.Printf("Pusher event til appene: %s", event.Title)

			// Send event til alle tilkoblede klienter
			for clientChan := range b.clients {
				clientChan <- event
			}
			i++
		}
		b.mutex.Unlock()
	}
}
