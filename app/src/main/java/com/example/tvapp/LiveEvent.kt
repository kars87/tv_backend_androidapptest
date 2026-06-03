package com.example.tvapp

data class LiveEvent(
    val id: Long,
    val type: String,      // "MATCH_START", "GOAL", "AD_START", etc.
    val title: String,     // "MÅÅÅL! SK Brann..."
    val timestamp: String  // "15:04:05"
)