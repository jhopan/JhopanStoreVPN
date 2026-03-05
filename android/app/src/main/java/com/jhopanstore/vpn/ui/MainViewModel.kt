package com.jhopanstore.vpn.ui

import android.content.Context
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.jhopanstore.vpn.core.VlessConfig
import com.jhopanstore.vpn.core.VlessParser
import com.jhopanstore.vpn.service.JhopanVpnService
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.net.HttpURLConnection
import java.net.InetSocketAddress
import java.net.Proxy
import java.net.URL

class MainViewModel : ViewModel() {

    // Connection fields — address includes port, e.g. "example.com:443"
    var address by mutableStateOf("")
    var uuid by mutableStateOf("")

    // Settings
    var path by mutableStateOf("/vless")
    var sni by mutableStateOf("")
    var host by mutableStateOf("")
    var dns1 by mutableStateOf("8.8.8.8")
    var dns2 by mutableStateOf("8.8.4.4")
    var allowInsecure by mutableStateOf(true)
    var autoReconnect by mutableStateOf(true)
    var pingUrl by mutableStateOf("https://dns.google")

    // State
    var isConnected by mutableStateOf(false)
    var isConnecting by mutableStateOf(false)
    var statusText by mutableStateOf("Disconnected")
    var pingResult by mutableStateOf("-")
    var showSettings by mutableStateOf(false)

    /** Split "host:port" input; defaults port to 443 */
    private fun parseAddress(): Pair<String, Int> {
        val trimmed = address.trim()
        val lastColon = trimmed.lastIndexOf(':')
        return if (lastColon > 0) {
            val host = trimmed.substring(0, lastColon)
            val port = trimmed.substring(lastColon + 1).toIntOrNull() ?: 443
            Pair(host, port)
        } else {
            Pair(trimmed, 443)
        }
    }

    private fun buildVlessUri(): String {
        val (addr, port) = parseAddress()
        val actualSni = sni.ifEmpty { addr }
        val actualHost = host.ifEmpty { addr }
        return "vless://$uuid@$addr:$port?type=ws&security=tls&path=${
            java.net.URLEncoder.encode(path, "UTF-8")
        }&sni=$actualSni&host=$actualHost&allowInsecure=$allowInsecure#JhopanStoreVPN"
    }

    fun importVlessUri(uri: String) {
        val result = VlessParser.parse(uri)
        result.onSuccess { cfg ->
            address = "${cfg.address}:${cfg.port}"
            uuid = cfg.uuid
            path = cfg.path
            sni = cfg.sni
            host = cfg.host
            allowInsecure = cfg.allowInsecure
            statusText = "Imported"
        }.onFailure {
            statusText = "Invalid VLESS URI"
        }
    }

    fun connect(context: Context) {
        if (address.isEmpty() || uuid.isEmpty()) {
            statusText = "Enter address and UUID"
            return
        }

        isConnecting = true
        statusText = "Connecting..."

        val uri = buildVlessUri()
        JhopanVpnService.start(context, uri, dns1, dns2, autoReconnect)

        // Poll for connection status
        viewModelScope.launch {
            var attempts = 0
            while (attempts < 20) {
                delay(1000)
                if (JhopanVpnService.isRunning) {
                    isConnected = true
                    isConnecting = false
                    statusText = "Connected"
                    startPingLoop()
                    return@launch
                }
                attempts++
            }
            isConnecting = false
            isConnected = JhopanVpnService.isRunning
            statusText = if (isConnected) "Connected" else "Connection failed"
        }
    }

    fun disconnect(context: Context) {
        JhopanVpnService.stop(context)
        isConnected = false
        isConnecting = false
        statusText = "Disconnected"
        pingResult = "-"
    }

    // --- Persistence ---

    fun saveSettings(context: Context) {
        context.getSharedPreferences("vpn_settings", Context.MODE_PRIVATE).edit()
            .putString("address", address)
            .putString("uuid", uuid)
            .putString("path", path)
            .putString("sni", sni)
            .putString("host", host)
            .putString("dns1", dns1)
            .putString("dns2", dns2)
            .putString("pingUrl", pingUrl)
            .putBoolean("allowInsecure", allowInsecure)
            .putBoolean("autoReconnect", autoReconnect)
            .apply()
    }

    fun loadSettings(context: Context) {
        val p = context.getSharedPreferences("vpn_settings", Context.MODE_PRIVATE)
        address = p.getString("address", "") ?: ""
        uuid = p.getString("uuid", "") ?: ""
        path = p.getString("path", "/vless") ?: "/vless"
        sni = p.getString("sni", "") ?: ""
        host = p.getString("host", "") ?: ""
        dns1 = p.getString("dns1", "8.8.8.8") ?: "8.8.8.8"
        dns2 = p.getString("dns2", "8.8.4.4") ?: "8.8.4.4"
        pingUrl = p.getString("pingUrl", "https://dns.google") ?: "https://dns.google"
        allowInsecure = p.getBoolean("allowInsecure", true)
        autoReconnect = p.getBoolean("autoReconnect", true)
    }

    // --- Ping ---

    private fun startPingLoop() {
        viewModelScope.launch {
            while (isConnected) {
                doPing()
                delay(5000)
            }
        }
    }

    private suspend fun doPing() {
        withContext(Dispatchers.IO) {
            try {
                val proxy = Proxy(
                    Proxy.Type.SOCKS,
                    InetSocketAddress("127.0.0.1", com.jhopanstore.vpn.core.XrayManager.SOCKS_PORT)
                )
                val url = URL(pingUrl)
                val conn = url.openConnection(proxy) as HttpURLConnection
                conn.connectTimeout = 5000
                conn.readTimeout = 5000
                conn.requestMethod = "HEAD"

                val start = System.currentTimeMillis()
                conn.connect()
                val code = conn.responseCode
                val elapsed = System.currentTimeMillis() - start
                conn.disconnect()

                pingResult = if (code in 200..399) "${elapsed} ms" else "Error $code"
            } catch (_: Exception) {
                pingResult = "Timeout"
            }
        }
    }
}
