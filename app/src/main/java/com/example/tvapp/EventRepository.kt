package com.example.tvapp

import com.google.gson.Gson
import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.callbackFlow
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.Response
import okhttp3.sse.EventSource
import okhttp3.sse.EventSourceListener
import okhttp3.sse.EventSources
import java.util.concurrent.TimeUnit

class EventRepository {
    private val client = OkHttpClient.Builder()
        .connectTimeout(5, TimeUnit.SECONDS)
        .readTimeout(10, TimeUnit.MINUTES) // Viktig: Hold tilkoblingen åpen
        .build()

    private val gson = Gson()

    fun listenToLiveEvents(): Flow<LiveEvent> = callbackFlow {
        // 10.0.2.2 er den magiske IP-en Android-emulatoren bruker for å treffe "localhost" på PC-en din
        val request = Request.Builder()
            .url("http://10.0.2.2:8080/api/live-events")
            .header("Accept", "text/event-stream")
            .build()

        val listener = object : EventSourceListener() {
            override fun onEvent(eventSource: EventSource, id: String?, type: String?, data: String) {
                try {
                    // Pars JSON-strengen fra Go til et LiveEvent-objekt
                    val event = gson.fromJson(data, LiveEvent::class.java)
                    // Skyv eventet inn i Flow-strømmen
                    trySend(event)
                } catch (e: Exception) {
                    e.printStackTrace()
                }
            }

            override fun onFailure(eventSource: EventSource, t: Throwable?, response: Response?) {
                // Her kunne vi lagt inn logikk for å prøve å koble til på nytt (reconnect)
                close(t)
            }
        }

        val eventSource = EventSources.createFactory(client)
            .newEventSource(request, listener)

        // Steng tilkoblingen hvis brukeren lukker appen eller forlater skjermen
        awaitClose {
            eventSource.cancel()
        }
    }
}