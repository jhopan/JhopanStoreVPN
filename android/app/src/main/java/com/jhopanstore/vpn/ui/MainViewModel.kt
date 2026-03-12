package com.jhopanstore.vpn.ui

import android.app.Application
import android.content.Context
import android.util.Log
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.jhopanstore.vpn.core.VlessConfig
import com.jhopanstore.vpn.core.VlessParser
import com.jhopanstore.vpn.core.XrayManager
import com.jhopanstore.vpn.service.JhopanVpnService
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.collectLatest
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.net.Inet4Address
import java.net.HttpURLConnection
import java.net.InetSocketAddress
import java.net.NetworkInterface
import java.net.Proxy
import java.net.URL

class MainViewModel(application: Application) : AndroidViewModel(application) {
    private val appContext = application.applicationContext

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
    
    // Flag untuk menandai sedang restart VPN (jangan reset proxy state)
    private var isRestarting = false

    init {
        // Collect StateFlow dari service — langsung update UI tanpa polling
        viewModelScope.launch {
            JhopanVpnService.state.collectLatest { state ->
                when (state) {
                    JhopanVpnService.VpnState.CONNECTED -> {
                        isConnected = true
                        isConnecting = false
                        statusText = "Connected"
                        isRestarting = false  // Reset flag setelah berhasil connect
                        startPingLoop()
                    }
                    JhopanVpnService.VpnState.CONNECTING -> {
                        isConnecting = true
                        isConnected = false
                        statusText = "Connecting..."
                    }
                    JhopanVpnService.VpnState.FAILED -> {
                        isConnected = false
                        isConnecting = false
                        statusText = "Connection failed"
                        pingResult = "-"
                        isRestarting = false  // Reset flag
                        // Hanya reset proxy state jika BUKAN sedang restart
                        if (!isRestarting) {
                            isProxySharingActive = false
                            XrayManager.hotspotSharing = false
                        }
                    }
                    JhopanVpnService.VpnState.DISCONNECTED -> {
                        // Hanya update jika memang sedang connected/connecting
                        // Hindari reset state saat app baru buka (initial DISCONNECTED)
                        if (isConnected || isConnecting) {
                            isConnected = false
                            isConnecting = false
                            statusText = "Disconnected"
                            pingResult = "-"
                            // Jangan reset proxy state jika sedang restart VPN
                            if (!isRestarting) {
                                isProxySharingActive = false
                                XrayManager.hotspotSharing = false
                            }
                        }
                    }
                }
            }
        }
    }

    // Hotspot sharing state
    var showHotspot by mutableStateOf(false)
    var isHotspotDetected by mutableStateOf(false)
    var hotspotIp by mutableStateOf("")
    var isProxySharingActive by mutableStateOf(false)

    // Battery optimization warning
    var isBatteryOptimized by mutableStateOf(false)

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

        // State CONNECTING akan diset oleh StateFlow collector saat service kirim sinyal
        isConnecting = true
        statusText = "Connecting..."

        val uri = buildVlessUri()
        JhopanVpnService.start(context, uri, dns1, dns2, autoReconnect)

        // Safety timeout: jika service hang total (tidak emit state apapun),
        // paksa reset setelah 30 detik agar field tidak terkunci selamanya
        viewModelScope.launch {
            delay(30_000)
            if (isConnecting && !isConnected) {
                isConnecting = false
                statusText = "Connection timeout"
            }
        }
    }

    fun disconnect(context: Context) {
        JhopanVpnService.stop(context)
        isConnected = false
        isConnecting = false
        statusText = "Disconnected"
        pingResult = "-"
        // reset proxy sharing when VPN disconnects
        isProxySharingActive = false
        XrayManager.hotspotSharing = false
    }

    // --- Hotspot Sharing ---

    /** Cek apakah hotspot aktif berdasarkan NetworkInterface. Panggil dari onResume. */
    fun checkHotspot() {
        viewModelScope.launch(Dispatchers.IO) {
            val ip = detectHotspotIp()
            withContext(Dispatchers.Main) {
                hotspotIp = ip ?: ""
                isHotspotDetected = ip != null
                // Kalau hotspot dimatikan saat proxy sharing aktif, reset
                if (!isHotspotDetected && isProxySharingActive) {
                    isProxySharingActive = false
                    XrayManager.hotspotSharing = false
                }
            }
        }
    }

    private fun detectHotspotIp(): String? {
        return try {
            val interfaces = NetworkInterface.getNetworkInterfaces()?.toList() ?: return null
            for (iface in interfaces) {
                if (!iface.isUp || iface.isLoopback) continue
                val name = iface.name.lowercase()
                
                // Skip interface yang PASTI BUKAN hotspot:
                // - tun: VPN tunnel interface
                // - rmnet/r_rmnet/ccmni: Data seluler interface (Qualcomm/Mediatek)
                // - ppp: Point-to-point protocol
                // - dummy: Virtual dummy interface
                // - v4-/clat: IPv4/IPv6 translation layer
                if (name.startsWith("tun") || name.startsWith("rmnet") ||
                    name.startsWith("ppp") || name.startsWith("dummy") ||
                    name.startsWith("v4-") || name.startsWith("clat") ||
                    name.startsWith("ccmni") || name.startsWith("r_rmnet")) continue
                
                for (addr in iface.inetAddresses.toList()) {
                    if (addr !is Inet4Address || addr.isLoopbackAddress) continue
                    val ip = addr.hostAddress ?: continue
                    
                    // Terima semua IP private range yang umum dipakai hotspot:
                    // 192.168.x.x (paling umum untuk hotspot)
                    // 10.x.x.x (beberapa device pakai range ini)
                    // 172.16.x.x - 172.31.x.x (jarang tapi mungkin)
                    val isPrivateRange = ip.startsWith("192.168.") || 
                                        ip.startsWith("10.") ||
                                        ip.matches(Regex("^172\\.(1[6-9]|2[0-9]|3[0-1])\\..*"))
                    
                    if (isPrivateRange) {
                        // Return IP ASLI device, JANGAN ubah jadi .1!
                        // Device hotspot Android ITU SENDIRI yang jadi gateway & proxy server
                        // Device lain harus connect ke IP ini, bukan gateway .1
                        return ip
                    }
                }
            }
            null
        } catch (e: Exception) { 
            Log.e("MainViewModel", "Error detecting hotspot IP", e)
            null 
        }
    }

    /** Toggle proxy sharing on/off. Restart VPN dengan binding baru. */
    fun toggleProxySharing(context: Context) {
        if (!isConnected) return
        isProxySharingActive = !isProxySharingActive
        XrayManager.hotspotSharing = isProxySharingActive
        Log.d("MainViewModel", "Toggle proxy sharing: $isProxySharingActive")
        restartVpn(context)
    }

    private fun restartVpn(context: Context) {
        isRestarting = true  // Set flag sebelum restart
        Log.d("MainViewModel", "Restarting VPN (isRestarting=$isRestarting, proxyActive=$isProxySharingActive)")
        JhopanVpnService.stop(context)
        isConnecting = true
        statusText = "Reconnecting..."
        viewModelScope.launch {
            delay(1500)
            val uri = buildVlessUri()
            JhopanVpnService.start(context, uri, dns1, dns2, autoReconnect)
            var attempts = 0
            while (attempts < 15) {
                delay(1000)
                if (JhopanVpnService.isRunning) {
                    isConnected = true
                    isConnecting = false
                    statusText = "Connected"
                    isRestarting = false  // Reset flag setelah berhasil
                    Log.d("MainViewModel", "VPN restarted successfully (proxyActive=$isProxySharingActive)")
                    return@launch
                }
                attempts++
            }
            isConnecting = false
            isConnected = JhopanVpnService.isRunning
            statusText = if (isConnected) "Connected" else "Reconnect failed"
            isRestarting = false  // Reset flag meski gagal
        }
    }

    // --- Persistence ---

    fun saveSettings(context: Context, immediate: Boolean = false) {
        val editor = context.getSharedPreferences("vpn_settings", Context.MODE_PRIVATE).edit()
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
        
        // Use commit() untuk immediate save (synchronous) agar tidak hilang saat app di-kill
        // Use apply() untuk background save (asynchronous) saat tidak urgent
        if (immediate) {
            editor.commit()
            Log.d("MainViewModel", "Settings saved immediately (commit)")
        } else {
            editor.apply()
            Log.d("MainViewModel", "Settings saved in background (apply)")
        }
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

    /** Sync isConnected dengan status service yang sebenarnya. Panggil dari onResume.
     *  StateFlow sudah handle update real-time, ini hanya untuk kasus
     *  app dibuka ulang saat VPN service masih jalan di background.
     */
    fun syncConnectionState() {
        val running = JhopanVpnService.isRunning
        // Hanya sync jika ada ketidaksesuaian dan TIDAK sedang dalam proses connecting
        if (running && !isConnected && !isConnecting) {
            isConnected = true
            isConnecting = false
            statusText = "Connected"
            startPingLoop()
        } else if (!running && isConnected) {
            // Service mati tapi UI masih showing connected
            isConnected = false
            isConnecting = false
            statusText = "Disconnected"
            pingResult = "-"
            isProxySharingActive = false
            XrayManager.hotspotSharing = false
        }
        // Jika isConnecting == true dan running == false: biarkan StateFlow yang handle
    }

    // --- Ping ---
    // Fires a short burst of 3 pings right after connect, reports the best result,
    // then stops — no background coroutine or CPU/radio wake after that.

    private fun startPingLoop() {
        viewModelScope.launch(Dispatchers.IO) {
            var best: Long = Long.MAX_VALUE
            var gotResult = false
            repeat(3) { attempt ->
                if (!isConnected) return@launch
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

                    if (code in 200..399 && elapsed < best) {
                        best = elapsed
                        gotResult = true
                        pingResult = "$elapsed ms"
                    }
                } catch (_: Exception) {
                    // ignore individual ping failure; report timeout only if all three fail
                }
                if (attempt < 2) delay(800)
            }
            if (!gotResult) pingResult = "Timeout"
            // Burst complete — no further background work
        }
    }
}
