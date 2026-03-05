package com.jhopanstore.vpn.ui

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ContentCopy
import androidx.compose.material.icons.filled.ContentPaste
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material.icons.filled.Warning
import androidx.compose.material.icons.filled.WifiTethering
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
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
    onImportClipboard: () -> Unit,
    onOpenHotspotSettings: () -> Unit,
    onToggleProxy: () -> Unit,
    onCopyProxy: () -> Unit,
    onOpenBatterySettings: () -> Unit
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("JhopanStoreVPN", fontWeight = FontWeight.Bold) },
                actions = {
                    IconButton(
                        onClick = onImportClipboard,
                        enabled = !viewModel.isConnected
                    ) {
                        Icon(Icons.Default.ContentPaste, contentDescription = "Paste")
                    }
                    IconButton(onClick = { viewModel.showHotspot = !viewModel.showHotspot; viewModel.showSettings = false }) {
                        Icon(
                            Icons.Default.WifiTethering,
                            contentDescription = "Bagikan VPN",
                            tint = if (viewModel.isProxySharingActive) Color(0xFF4CAF50)
                                   else LocalContentColor.current
                        )
                    }
                    IconButton(onClick = { viewModel.showSettings = !viewModel.showSettings; viewModel.showHotspot = false }) {
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
        when {
            viewModel.showHotspot -> HotspotScreen(
                viewModel = viewModel,
                onClose = { viewModel.showHotspot = false },
                onToggleProxy = onToggleProxy,
                onCheckHotspot = { viewModel.checkHotspot() },
                modifier = Modifier.padding(padding)
            )
            viewModel.showSettings -> SettingsScreen(
                viewModel = viewModel,
                onClose = { viewModel.showSettings = false },
                modifier = Modifier.padding(padding)
            )
            else -> MainContent(
                viewModel = viewModel,
                onConnect = onConnect,
                onDisconnect = onDisconnect,
                onCopyProxy = onCopyProxy,
                onOpenBatterySettings = onOpenBatterySettings,
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
    onCopyProxy: () -> Unit,
    onOpenBatterySettings: () -> Unit,
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

        // Warning: battery optimization masih aktif
        if (viewModel.isBatteryOptimized) {
            Surface(
                color = Color(0xFFFFF3CD),
                shape = RoundedCornerShape(8.dp),
                modifier = Modifier.fillMaxWidth()
            ) {
                Row(
                    modifier = Modifier.padding(horizontal = 12.dp, vertical = 10.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        Icons.Default.Warning,
                        contentDescription = null,
                        tint = Color(0xFFF57C00),
                        modifier = Modifier.size(18.dp)
                    )
                    Spacer(Modifier.width(8.dp))
                    Text(
                        text = "Pembatasan daya aktif. VPN bisa terputus saat layar mati.",
                        fontSize = 12.sp,
                        color = Color(0xFF5D4037),
                        modifier = Modifier.weight(1f)
                    )
                    Spacer(Modifier.width(4.dp))
                    TextButton(
                        onClick = onOpenBatterySettings,
                        contentPadding = PaddingValues(horizontal = 6.dp, vertical = 0.dp)
                    ) {
                        Text("Perbaiki", fontSize = 12.sp, color = Color(0xFFF57C00))
                    }
                }
            }
            Spacer(modifier = Modifier.height(12.dp))
        }

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

        // Proxy sharing info row — tampil hanya saat VPN konek + proxy aktif
        if (viewModel.isConnected && viewModel.isProxySharingActive && viewModel.hotspotIp.isNotEmpty()) {
            Spacer(modifier = Modifier.height(16.dp))
            Surface(
                color = Color(0xFF4CAF50).copy(alpha = 0.12f),
                shape = RoundedCornerShape(8.dp),
                modifier = Modifier.fillMaxWidth()
            ) {
                Row(
                    modifier = Modifier.padding(horizontal = 12.dp, vertical = 10.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        Icons.Default.WifiTethering,
                        contentDescription = null,
                        tint = Color(0xFF4CAF50),
                        modifier = Modifier.size(18.dp)
                    )
                    Spacer(Modifier.width(8.dp))
                    Text(
                        text = "Proxy: ${viewModel.hotspotIp}:10808",
                        fontSize = 13.sp,
                        fontWeight = FontWeight.Medium,
                        modifier = Modifier.weight(1f)
                    )
                    IconButton(
                        onClick = onCopyProxy,
                        modifier = Modifier.size(32.dp)
                    ) {
                        Icon(
                            Icons.Default.ContentCopy,
                            contentDescription = "Salin",
                            modifier = Modifier.size(16.dp)
                        )
                    }
                }
            }
        }
    }
}
