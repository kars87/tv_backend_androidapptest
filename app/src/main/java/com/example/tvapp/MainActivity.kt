package com.example.tvapp

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewmodel.compose.viewModel
import kotlinx.coroutines.MainScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch



// Enkel ViewModel som samler opp events i en liste
class LiveEventViewModel : ViewModel() {
    private val repository = EventRepository()
    private val _events = MutableStateFlow<List<LiveEvent>>(emptyList())
    val events: StateFlow<List<LiveEvent>> = _events.asStateFlow()

    init {
        // Start lytting i en Coroutine med en gang ViewModelen opprettes
        MainScope().launch {
            repository.listenToLiveEvents().collect { newEvent ->
                // Legg det nye eventet øverst i listen
                _events.value = listOf(newEvent) + _events.value
            }
        }
    }
}

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MaterialTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = Color(0xFF121212) // Mørk TV 2-aktig bakgrunn
                ) {
                    LiveDashboardScreen()
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LiveDashboardScreen(viewModel: LiveEventViewModel = viewModel()) {
    val eventList by viewModel.events.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("TV 2 Play - Live Feed", color = Color.White, fontWeight = FontWeight.Bold) },
                colors = TopAppBarDefaults.topAppBarColors(containerColor = Color(0xFFFF5722)) // TV 2 Rød/Orange
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .background(Color(0xFF1C1C1E))
        ) {
            if (eventList.isEmpty()) {
                Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    CircularProgressIndicator(color = Color(0xFFFF5722))
                }
            } else {
                LazyColumn(
                    modifier = Modifier.fillMaxSize(),
                    contentPadding = PaddingValues(16.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    items(eventList) { event ->
                        EventCard(event = event)
                    }
                }
            }
        }
    }
}

@Composable
fun EventCard(event: LiveEvent) {
    val cardColor = when (event.type) {
        "MATCH_START" -> Color(0xFF1976D2) // Blå for start
        "GOAL" -> Color(0xFF2E7D32)        // Grønn for mål
        "AD_START" -> Color(0xFFC62828)    // Rød for reklame
        "YELLOW_CARD" -> Color(0xFFFBC02D)  // Gul for kort
        else -> Color(0xFF2C2C2E)          // Grå standard
    }

    Card(
        shape = RoundedCornerShape(8.dp),
        modifier = Modifier.fillMaxWidth(),
        colors = CardDefaults.cardColors(containerColor = cardColor)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = event.type,
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Bold,
                    color = if (event.type == "YELLOW_CARD") Color.Black else Color.White
                )
                Text(
                    text = event.timestamp,
                    fontSize = 12.sp,
                    color = if (event.type == "YELLOW_CARD") Color.DarkGray else Color.LightGray
                )
            }
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = event.title,
                fontSize = 16.sp,
                fontWeight = FontWeight.Medium,
                color = if (event.type == "YELLOW_CARD") Color.Black else Color.White
            )
        }
    }
}