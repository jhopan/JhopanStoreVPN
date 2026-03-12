package com.jhopanstore.vpn.ui

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
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

@Composable
fun SettingsScreen(
    viewModel: MainViewModel,
    onClose: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        // Title
        Text(
            text = "Settings",
            fontSize = 22.sp,
            fontWeight = FontWeight.Bold,
            modifier = Modifier.fillMaxWidth()
        )

        Spacer(modifier = Modifier.height(16.dp))

        // Logo
        Image(
            painter = painterResource(id = R.drawable.ic_logo),
            contentDescription = "Logo",
            modifier = Modifier.size(100.dp)
        )

        Spacer(modifier = Modifier.height(20.dp))

        // ═══════════════════════════════
        //  Connection
        // ═══════════════════════════════
        SectionHeader("Connection")

        SettingRow("Path:") {
            OutlinedTextField(
                value = viewModel.path,
                onValueChange = { viewModel.path = it },
                modifier = Modifier.fillMaxWidth(),
                singleLine = true,
                enabled = !viewModel.isConnected
            )
        }

        SettingRow("SNI:") {
            OutlinedTextField(
                value = viewModel.sni,
                onValueChange = { viewModel.sni = it },
                modifier = Modifier.fillMaxWidth(),
                singleLine = true,
                enabled = !viewModel.isConnected
            )
        }

        SettingRow("Host:") {
            OutlinedTextField(
                value = viewModel.host,
                onValueChange = { viewModel.host = it },
                modifier = Modifier.fillMaxWidth(),
                singleLine = true,
                enabled = !viewModel.isConnected
            )
        }

        Spacer(modifier = Modifier.height(12.dp))

        // ═══════════════════════════════
        //  DNS
        // ═══════════════════════════════
        SectionHeader("DNS")

        SettingRow("DNS 1:") {
            OutlinedTextField(
                value = viewModel.dns1,
                onValueChange = { viewModel.dns1 = it },
                modifier = Modifier.fillMaxWidth(),
                singleLine = true,
                enabled = !viewModel.isConnected,
                isError = viewModel.dns1.isNotBlank() && !isValidIpOrEmpty(viewModel.dns1),
                supportingText = if (viewModel.dns1.isNotBlank() && !isValidIpOrEmpty(viewModel.dns1)) {
                    { Text("Format tidak valid (contoh: 8.8.8.8)", color = MaterialTheme.colorScheme.error, fontSize = 11.sp) }
                } else null
            )
        }

        SettingRow("DNS 2:") {
            OutlinedTextField(
                value = viewModel.dns2,
                onValueChange = { viewModel.dns2 = it },
                modifier = Modifier.fillMaxWidth(),
                singleLine = true,
                enabled = !viewModel.isConnected,
                isError = viewModel.dns2.isNotBlank() && !isValidIpOrEmpty(viewModel.dns2),
                supportingText = if (viewModel.dns2.isNotBlank() && !isValidIpOrEmpty(viewModel.dns2)) {
                    { Text("Format tidak valid (contoh: 8.8.4.4)", color = MaterialTheme.colorScheme.error, fontSize = 11.sp) }
                } else null
            )
        }

        Spacer(modifier = Modifier.height(12.dp))

        // ═══════════════════════════════
        //  HTTP Ping
        // ═══════════════════════════════
        SectionHeader("HTTP Ping")

        SettingRow("Ping URL:") {
            OutlinedTextField(
                value = viewModel.pingUrl,
                onValueChange = { viewModel.pingUrl = it },
                modifier = Modifier.fillMaxWidth(),
                singleLine = true
            )
        }

        Spacer(modifier = Modifier.height(12.dp))

        // ═══════════════════════════════
        //  Behavior
        // ═══════════════════════════════
        SectionHeader("Behavior")

        CheckboxRow(
            label = "Auto Reconnect",
            checked = viewModel.autoReconnect,
            onCheckedChange = { viewModel.autoReconnect = it }
        )

        CheckboxRow(
            label = "Allow Insecure TLS (skip verify)",
            checked = viewModel.allowInsecure,
            onCheckedChange = { viewModel.allowInsecure = it },
            enabled = !viewModel.isConnected
        )

        Spacer(modifier = Modifier.height(28.dp))

        // Close button
        Button(
            onClick = onClose,
            shape = RoundedCornerShape(8.dp),
            modifier = Modifier.widthIn(min = 120.dp)
        ) {
            Text("Close")
        }

        Spacer(modifier = Modifier.height(16.dp))
    }
}

// ── Helpers ──────────────────────────────────────────────────────

/** Returns true for blank (use default) or valid dotted-decimal IPv4 address. */
private fun isValidIpOrEmpty(ip: String): Boolean {
    if (ip.isBlank()) return true
    val parts = ip.trim().split(".")
    if (parts.size != 4) return false
    return parts.all { part -> part.toIntOrNull()?.let { it in 0..255 } ?: false }
}

@Composable
private fun SectionHeader(title: String) {
    Text(
        text = "\u2014 $title \u2014",
        fontWeight = FontWeight.Bold,
        fontSize = 14.sp,
        textAlign = TextAlign.Center,
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 8.dp),
        color = MaterialTheme.colorScheme.onSurface
    )
}

@Composable
private fun SettingRow(
    label: String,
    content: @Composable () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            text = label,
            modifier = Modifier.width(76.dp),
            fontSize = 14.sp,
            fontWeight = FontWeight.Medium
        )
        Spacer(modifier = Modifier.width(8.dp))
        Box(modifier = Modifier.weight(1f)) {
            content()
        }
    }
}

@Composable
private fun CheckboxRow(
    label: String,
    checked: Boolean,
    onCheckedChange: (Boolean) -> Unit,
    enabled: Boolean = true
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 2.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Checkbox(
            checked = checked,
            onCheckedChange = onCheckedChange,
            enabled = enabled
        )
        Spacer(modifier = Modifier.width(4.dp))
        Text(text = label, fontSize = 14.sp)
    }
}
