package com.example.barkdroid

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import cn.jpush.android.api.JPushInterface
import com.example.barkdroid.network.ApiService
import com.example.barkdroid.push.PushEventBus

class MainActivity : ComponentActivity() {

    private lateinit var apiService: ApiService

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Init JPush (AppKey from BuildConfig, injected at compile time)
        JPushInterface.init(this)
        JPushInterface.setDebugMode(BuildConfig.DEBUG)

        // Setup network
        apiService = ApiService(
            baseUrl = BuildConfig.SERVER_BASE_URL,
            token = BuildConfig.SERVER_TOKEN,
        )

        // Listen for registration ID
        PushEventBus.registrationId.observe(this) { regId ->
            apiService.register(regId) { deviceKey ->
                runOnUiThread {
                    _deviceKey.value = deviceKey
                    _status.value = "registered"
                }
            }
        }

        setContent {
            BarkdroidApp(
                status = _status.value,
                deviceKey = _deviceKey.value,
                regId = _regId.value,
                serverUrl = BuildConfig.SERVER_BASE_URL,
            )
        }
    }

    // Observable state
    private val _status = mutableStateOf("initializing")
    private val _deviceKey = mutableStateOf<String?>(null)
    private val _regId = mutableStateOf<String?>(null)

    init {
        PushEventBus.registrationId.observeForever { _regId.value = it }
    }
}

@Composable
fun BarkdroidApp(
    status: String,
    deviceKey: String?,
    regId: String?,
    serverUrl: String,
) {
    MaterialTheme {
        Scaffold { padding ->
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .padding(24.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center,
            ) {
                Text(
                    text = "barkdroid",
                    fontSize = 28.sp,
                    fontWeight = FontWeight.Bold,
                )
                Spacer(modifier = Modifier.height(32.dp))

                // Status card
                Card(modifier = Modifier.fillMaxWidth()) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text("Status: $status", fontWeight = FontWeight.Medium)
                        Spacer(modifier = Modifier.height(8.dp))
                        Text("Server: $serverUrl", fontSize = 12.sp, color = MaterialTheme.colorScheme.onSurfaceVariant)
                        regId?.let {
                            Spacer(modifier = Modifier.height(4.dp))
                            Text("RegID: ${it.take(16)}...", fontSize = 12.sp, color = MaterialTheme.colorScheme.onSurfaceVariant)
                        }
                    }
                }

                // Device key (the important part — user copies this)
                if (deviceKey != null) {
                    Spacer(modifier = Modifier.height(24.dp))
                    Card(
                        modifier = Modifier.fillMaxWidth(),
                        colors = CardDefaults.cardColors(
                            containerColor = MaterialTheme.colorScheme.primaryContainer,
                        ),
                    ) {
                        Column(
                            modifier = Modifier.padding(16.dp),
                            horizontalAlignment = Alignment.CenterHorizontally,
                        ) {
                            Text("Your Device Key", fontSize = 12.sp)
                            Spacer(modifier = Modifier.height(8.dp))
                            Text(
                                text = deviceKey,
                                fontSize = 24.sp,
                                fontWeight = FontWeight.Bold,
                                color = MaterialTheme.colorScheme.onPrimaryContainer,
                            )
                            Spacer(modifier = Modifier.height(8.dp))
                            Text(
                                text = "curl $serverUrl/$deviceKey/Title/Body",
                                fontSize = 11.sp,
                                color = MaterialTheme.colorScheme.onPrimaryContainer.copy(alpha = 0.7f),
                            )
                        }
                    }
                }
            }
        }
    }
}
