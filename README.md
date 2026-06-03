# TV - Live Event Streams (Proof of Concept)

Dette prosjektet er en teknisk demonstrasjon (PoC) av en moderne, reaktiv arkitektur for håndtering av sanntids-events i en strømmetjeneste (f.eks. måloppdateringer, live-statistikk og reklame-triggers under direktesendte eventer). 

Systemet består av **Go-backend** og en moderne **Android-app skrevet i Kotlin og Jetpack Compose**, bundet sammen av en vedvarende datastrøm via Server-Sent Events (SSE).

---

## 🛠️ Arkitektur og Teknologivalg

Prosjektet er bevisst designet for å speile utfordringene en storskala strømmetjeneste som TV tjenester står overfor under høytrafikk-eventer (som Premier League eller landskamper).

### 1. Go Backend (Server)
Valget falt på **Go (Golang)** på grunn av språkets ekstreme ytelse og effektive minnehåndtering ved massiv parallellitet (concurrency).
* **Server-Sent Events (SSE):** Brukes i stedet for WebSockets for envejs sanntidsstrømming av metadata til klientene, noe som reduserer protokoll-overhead og holder ressursbruken minimal.
* **Goroutines & Channels:** Hver tilkoblet app-klient tildeles en lettvektstråd (Goroutine). Koordinering og distribusjon av live-events skjer trådsikkert via Go Channels og en sentralisert `Broker`-struktur.

### 2. Android App (Frontend)
Appen er bygget etter moderne standarder for Android-utvikling, med fokus på reaktivitet og lav batteri-/ressursbruk under kontinuerlig nettverkslytting.
* **Jetpack Compose:** Brukes for et 100 % deklarativt UI. Grensesnittet reagerer og animeres umiddelbart basert på tilstandsendringer.
* **Kotlin Coroutines & Flows:** Nettverkslaget bruker `callbackFlow` fra OkHttp for å transformere den kontinuerlige SSE-strømmen til en reaktiv Kotlin-strøm. Dette sikrer at appen aldri blokkerer hovedtråden (UI-tråden).
* **Unidirectional Data Flow (UDF):** En dedikert `ViewModel` samler opp hendelsene og eksponerer en uforanderlig (immutable) tilstand til Compose-skjermen via `StateFlow`.

---

## 🚀 Slik kjører du prosjektet lokalt

Prosjektet er utviklet og testet på et Windows/Linux-miljø.

### 1. Start Go-backenden
Naviger til mappen med Go-koden og kjør serveren:
```bash
go run main.go

Serveren starter opp på http://localhost:8080 og vil begynne å simulere live TV-events (kampstart, mål, gule kort, reklamepauser) i bakgrunnen så fort en klient kobler seg på.
2. Start Android-appen

    Åpne prosjektet i Android Studio.

    Sørg for at du har en aktiv Android-emulator (f.eks. en Pixel-enhet).

    Trykk på Run (Shift + F10).

Merk: Appen er konfigurert til å bruke lokal-IP-en 10.0.2.2 for å kommunisere med Go-serveren som kjører på vertsmaskinen.
🏋️ Lasttesting og Skalerbarhet (300 000+ simuleringer)

For å verifisere at arkitekturen er klar for "TV 2 Play-skala", ble det utviklet et eget stresstest-script i Go (stresstest.go).
