package com.jhopanstore.vpn.ui

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ContentPaste
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.jhopanstore.vpn.R

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MainScreen(
    viewModel: MainViewModel,
    onConnect: () -> Unit,
    onDisconnect: () -> Unit,
    onImportClipboard: () -> Unit
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("JhopanStoreVPN", fontWeight = FontWeight.Bold) },
                actions = {
                    IconButton(onClick = onImportClipboard) {
                        Icon(Icons.Default.ContentPaste, contentDescription = "Paste")
                    }
                    IconButton(onClick = { viewModel.showSettings = !viewModel.showSettings }) {
                        Icon(Icons.Default.Settings, contentDescription = "Settings")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.surface,
                    titleContentColor = MaterialTheme.colorScheme.onSurface,
                    actionIconContentColor = MaterialTheme.colorScheme.onSurface
                )
            )
        }
    ) { padding ->
        if (viewModel.showSettings) {
            SettingsScreen(
                viewModel = viewModel,
                onClose = { viewModel.showSettings = false },
                modifier = Modifier.padding(padding)
            )
        } else {
            MainContent(
                viewModel = viewModel,
                onConnect = onConnect,
                onDisconnect = onDisconnect,
                modifier = Modifier.padding(padding)
            )
        }
    }
}

@Composable
private fun MainContent(
    viewModel: MainViewModel,
    onConnect: () -> Unit,
    onDisconnect: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(horizontal = 24.dp, vertical = 16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Spacer(modifier = Modifier.height(12.dp))

        // Logo
        Image(
            painter = painterResource(id = R.drawable.ic_logo),
            contentDescription = "Logo",
            modifier = Modifier.size(120.dp)
        )

        Spacer(modifier = Modifier.height(28.dp))

        // Address field
        Text(
            text = "Address",
            fontWeight = FontWeight.Bold,
            fontSize = 16.sp,
            modifier = Modifier.fillMaxWidth()
        )
        Spacer(modifier = Modifier.height(4.dp))
        OutlinedTextField(
            value = viewModel.address,
            onValueChange = { viewModel.address = it },
            placeholder = { Text("example.com:443") },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            enabled = !viewModel.isConnected && !viewModel.isConnecting,
            shape = RoundedCornerShape(8.dp)
        )

        Spacer(modifier = Modifier.height(12.dp))

        // UUID field
        Text(
            text = "UUID",
            fontWeight = FontWeight.Bold,
            fontSize = 16.sp,
            modifier = Modifier.fillMaxWidth()
        )
        Spacer(modifier = Modifier.height(4.dp))
        OutlinedTextField(
            value = viewModel.uuid,
            onValueChange = { viewModel.uuid = it },
            placeholder = { Text("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx") },
            modifier = Modifier.fillMaxWidth(),
            singleLine = true,
            enabled = !viewModel.isConnected && !viewModel.isConnecting,
            shape = RoundedCornerShape(8.dp)
        )

        Spacer(modifier = Modifier.height(36.dp))

        // Connect / Disconnect buttons side-by-side
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            Button(
                onClick = onConnect,
                modifier = Modifier
                    .weight(1f)
                    .height(48.dp),
                shape = RoundedCornerShape(8.dp),
                enabled = !viewModel.isConnected && !viewModel.isConnecting,
                colors = ButtonDefaults.buttonColors(
                    containerColor = MaterialTheme.colorScheme.primary
                )
            ) {
                if (viewModel.isConnecting) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(20.dp),
                        color = MaterialTheme.colorScheme.onPrimary,
                        strokeWidth = 2.dp
                    )
                } else {
                    Text("CONNECT", fontWeight = FontWeight.Bold)
                }
            }

            OutlinedButton(
                onClick = onDisconnect,
                modifier = Modifier
                    .weight(1f)
                    .height(48.dp),
                shape = RoundedCornerShape(8.dp),
                enabled = viewModel.isConnected
            ) {
                Text("DISCONNECT", fontWeight = FontWeight.Bold)
            }
        }

        Spacer(modifier = Modifier.height(28.dp))

        // Status
        Text(
            text = viewModel.statusText,
            fontSize = 16.sp,
            fontWeight = FontWeight.Medium,
            textAlign = TextAlign.Center,
            color = when {
                viewModel.isConnected -> MaterialTheme.colorScheme.primary
                viewModel.isConnecting -> MaterialTheme.colorScheme.tertiary
                else -> MaterialTheme.colorScheme.onSurface
            }
        )

        Spacer(modifier = Modifier.height(8.dp))

        Text(
            text = "Ping: ${viewModel.pingResult}",
            fontSize = 14.sp,
            textAlign = TextAlign.Center,
            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
        )
    }
}
