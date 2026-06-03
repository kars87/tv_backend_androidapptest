package main

import (
	"bufio"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Antall simulerte apper/klienter du vil kaste mot serveren din
	TargetClients = 10000
	// Hvor raskt de skal koble seg på (i millisekunder mellom hver tilkobling)
	RampUpDelay = 1 * time.Millisecond
	// URL-en til Go-serveren din
	ServerURL = "http://localhost:8080/api/live-events"
)

func main() {
	var activeClients int64
	var failedClients int64
	var wg sync.WaitGroup

	fmt.Printf("Starter stresstest... Mål: %d simulerte klienter\n", TargetClients)
	startTime := time.Now()

	for i := 0; i < TargetClients; i++ {
		wg.Add(1)

		// Start en Goroutine per klient for å simulere parallellitet
		go func(id int) {
			defer wg.Done()

			// Opprett en HTTP-forespørsel
			req, err := http.NewRequest("GET", ServerURL, nil)
			if err != nil {
				atomic.AddInt64(&failedClients, 1)
				return
			}
			req.Header.Set("Accept", "text/event-stream")

			// Bruk en tilpasset HTTP-klient som ikkestenger tilkoblingen med en gang
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				atomic.AddInt64(&failedClients, 1)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				atomic.AddInt64(&failedClients, 1)
				return
			}

			// Klienten har koblet seg på vellykket!
			atomic.AddInt64(&activeClients, 1)

			// Les strømmen kontinuerlig (simulerer at appen lytter)
			reader := bufio.NewReader(resp.Body)
			for {
				_, err := reader.ReadString('\n')
				if err != nil {
					// Tilkoblingen ble brutt av en eller annen grunn
					atomic.AddInt64(&activeClients, -1)
					atomic.AddInt64(&failedClients, 1)
					return
				}
			}
		}(i)

		// Litt forsinkelse mellom hver pålogging for å ikke kvele nettverksstakken på PC-en din lokalt med en gang
		time.Sleep(RampUpDelay)

		// Print status hvert 500. klient
		if i%500 == 0 && i > 0 {
			fmt.Printf("[%s] Startet %d tråder... Aktive nå: %d\n", time.Since(startTime).Round(time.Millisecond), i, atomic.LoadInt64(&activeClients))
		}
	}

	// Hold stresstesteren i gang i 30 sekunder for å overvåke stabiliteten
	fmt.Println("\n--- Alle klienter har forsøkt å koble seg til! Overvåker i 30 sekunder ---")
	for i := 0; i < 6; i++ {
		time.Sleep(5 * time.Second)
		fmt.Printf("Status etter %ds -> Aktive app-klienter: %d | Feilede/Avbrutte: %d\n", (i+1)*5, atomic.LoadInt64(&activeClients), atomic.LoadInt64(&failedClients))
	}

	fmt.Println("Stresstest ferdig.")
}
