package com.jhopanstore.vpn.service

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Intent
import android.net.VpnService
import android.os.Build
import android.os.ParcelFileDescriptor
import android.os.PowerManager
import android.system.Os
import android.system.OsConstants
import android.util.Log
import com.jhopanstore.vpn.MainActivity
import com.jhopanstore.vpn.R
import com.jhopanstore.vpn.core.Tun2socksManager
import com.jhopanstore.vpn.core.VlessConfig
import com.jhopanstore.vpn.core.XrayManager
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import libXray.LibXray
import java.io.IOException

/**
 * Android VpnService that creates a TUN interface and routes traffic
 * through Xray (via libXray in-process) + tun2socks bridge.
 *
 * Traffic flow:
 *   Apps → TUN fd → tun2socks → SOCKS5 → Xray (libXray) → internet
 *
 * Key improvement with libXray:
 *   - registerDialerController(protectFd) ensures Xray's outgoing sockets
 *     are protected from VPN routing (prevents routing loops)
 *   - No xray binary process to manage
 */
class JhopanVpnService : VpnService() {

    enum class VpnState { DISCONNECTED, CONNECTING, CONNECTED, FAILED }

    companion object {
        private const val TAG = "JhopanVpnService"
        private const val CHANNEL_ID = "jhopan_vpn_channel"
        private const val NOTIFICATION_ID = 1
        private const val ACTION_STOP = "ACTION_STOP"
        const val EXTRA_VLESS_URI = "vless_uri"
        const val EXTRA_DNS1 = "dns1"
        const val EXTRA_DNS2 = "dns2"
        const val EXTRA_AUTO_RECONNECT = "auto_reconnect"
        private const val MAX_RECONNECT_ATTEMPTS = 5

        @Volatile
        var isRunning = false
            private set

        // StateFlow untuk UI — ViewModel collect ini, tidak perlu polling
        private val _state = MutableStateFlow(VpnState.DISCONNECTED)
        val state: StateFlow<VpnState> = _state

        // Flag untuk cancel background thread saat stop dipanggil
        @Volatile
        var isStopping = false
            private set

        fun start(context: Context, vlessUri: String, dns1: String, dns2: String, autoReconnect: Boolean) {
            val intent = Intent(context, JhopanVpnService::class.java).apply {
                putExtra(EXTRA_VLESS_URI, vlessUri)
                putExtra(EXTRA_DNS1, dns1)
                putExtra(EXTRA_DNS2, dns2)
                putExtra(EXTRA_AUTO_RECONNECT, autoReconnect)
            }
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                context.startForegroundService(intent)
            } else {
                context.startService(intent)
            }
        }

        fun stop(context: Context) {
            val intent = Intent(context, JhopanVpnService::class.java).apply {
                action = ACTION_STOP
            }
            context.startService(intent)
        }
    }

    private var tunFd: ParcelFileDescriptor? = null
    private var reconnectWakeLock: PowerManager.WakeLock? = null
    private var serviceWakeLock: PowerManager.WakeLock? = null

    // State for auto-reconnect
    private var lastVlessUri: String? = null
    private var lastDns1: String = "8.8.8.8"
    private var lastDns2: String = "8.8.4.4"
    private var autoReconnect = false
    private var reconnectAttempts = 0

    override fun onCreate() {
        super.onCreate()
        createNotificationChannel()
        
        // Acquire WakeLock untuk mencegah service mati saat layar mati
        val pm = getSystemService(POWER_SERVICE) as PowerManager
        serviceWakeLock = pm.newWakeLock(
            PowerManager.PARTIAL_WAKE_LOCK,
            "jhopanvpn:service"
        ).apply {
            setReferenceCounted(false)
            acquire()
        }
        Log.i(TAG, "Service WakeLock acquired")

        // Register VPN socket protection callback with libXray.
        // When Xray creates outgoing connections, protectFd() is called
        // so those sockets bypass the VPN TUN → prevents routing loops.
        LibXray.registerDialerController(object : libXray.DialerController {
            override fun protectFd(fd: Long): Boolean {
                val protected_ = protect(fd.toInt())
                if (!protected_) {
                    Log.w(TAG, "Failed to protect fd=$fd")
                }
                return protected_
            }
        })
        Log.i(TAG, "Registered DialerController for socket protection")
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        if (intent?.action == ACTION_STOP) {
            isStopping = true
            autoReconnect = false
            disconnect()
            stopSelf()
            return START_NOT_STICKY
        }

        val vlessUri = intent?.getStringExtra(EXTRA_VLESS_URI)
        if (vlessUri.isNullOrEmpty()) {
            stopSelf()
            return START_NOT_STICKY
        }

        val dns1 = intent.getStringExtra(EXTRA_DNS1) ?: "8.8.8.8"
        val dns2 = intent.getStringExtra(EXTRA_DNS2) ?: "8.8.4.4"
        autoReconnect = intent.getBooleanExtra(EXTRA_AUTO_RECONNECT, false)

        // Reset flag setiap kali ada koneksi baru
        isStopping = false
        _state.value = VpnState.CONNECTING

        startForeground(NOTIFICATION_ID, buildNotification("Connecting..."))

        Thread {
            connect(vlessUri, dns1, dns2)
        }.apply {
            isDaemon = true
            name = "vpn-connect"
        }.start()

        return START_STICKY
    }

    private fun connect(vlessUri: String, dns1: String, dns2: String) {
        try {
            // Save for reconnection
            lastVlessUri = vlessUri
            lastDns1 = dns1
            lastDns2 = dns2

            if (isStopping) return

            val parseResult = com.jhopanstore.vpn.core.VlessParser.parse(vlessUri)
            val cfg = parseResult.getOrElse {
                Log.e(TAG, "Failed to parse VLESS URI", it)
                _state.value = VpnState.FAILED
                updateNotification("Parse error")
                stopSelf()
                return
            }

            // Pre-resolve proxy server domain to IP BEFORE VPN TUN is established
            val resolvedIp = XrayManager.resolveDomain(cfg.address)
            Log.i(TAG, "Proxy server: ${cfg.address} -> ${resolvedIp ?: "unresolved (using domain)"}")

            // Set up death callback for auto-reconnect (iteratif, bukan rekursif)
            XrayManager.onProcessDied = {
                if (autoReconnect && lastVlessUri != null) {
                    Thread {
                        // Cegah CPU sleep saat proses reconnect di Doze mode
                        val pm = getSystemService(POWER_SERVICE) as PowerManager
                        val wl = pm.newWakeLock(PowerManager.PARTIAL_WAKE_LOCK, "jhopanvpn:reconnect")
                        wl.acquire(60_000L) // max 60 detik
                        reconnectWakeLock = wl
                        try {
                            var success = false
                            while (autoReconnect && reconnectAttempts < MAX_RECONNECT_ATTEMPTS && !success) {
                                reconnectAttempts++
                                val delay = 3000L * reconnectAttempts
                                Log.w(TAG, "Xray died, reconnecting in ${delay}ms (attempt $reconnectAttempts)")
                                updateNotification("Reconnecting ($reconnectAttempts)...")
                                Thread.sleep(delay)

                                // Bersihkan state lama
                                Tun2socksManager.stop()
                                try { tunFd?.close() } catch (_: Exception) {}
                                tunFd = null

                                // Restart Xray
                                val uri = lastVlessUri ?: break
                                val parsedCfg = com.jhopanstore.vpn.core.VlessParser.parse(uri).getOrNull() ?: break
                                val resolvedIp = XrayManager.resolveDomain(parsedCfg.address)
                                val xrayStarted = XrayManager.start(this, parsedCfg, lastDns1, lastDns2, resolvedIp)
                                if (!xrayStarted) continue

                                // Probe SOCKS5 port — proceed as soon as ready, max 5s
                                var portUp = false
                                for (probe in 0 until 20) {
                                    try {
                                        java.net.Socket("127.0.0.1", XrayManager.SOCKS_PORT).close()
                                        portUp = true
                                    } catch (_: Exception) {}
                                    if (portUp) break
                                    Thread.sleep(250)
                                }
                                if (!portUp) { XrayManager.stop(); continue }

                                // Re-establish TUN
                                val builder = Builder()
                                    .setSession("JhopanStoreVPN")
                                    .addAddress("10.0.0.2", 24)
                                    .addRoute("0.0.0.0", 0)
                                    .addDnsServer(lastDns1.ifBlank { "8.8.8.8" })
                                    .addDnsServer(lastDns2.ifBlank { "8.8.4.4" })
                                    .setMtu(1400)
                                builder.addDisallowedApplication(packageName)
                                tunFd = builder.establish()
                                if (tunFd == null) { XrayManager.stop(); continue }

                                val tun2socksOk = Tun2socksManager.start(this, tunFd!!.fd)
                                if (!tun2socksOk) { XrayManager.stop(); tunFd?.close(); tunFd = null; continue }

                                success = true
                            }
                            if (success) {
                                reconnectAttempts = 0
                                isRunning = true
                                _state.value = VpnState.CONNECTED
                                updateNotification("Connected")
                            } else {
                                Log.e(TAG, "Auto-reconnect exhausted")
                                disconnect()
                                stopSelf()
                            }
                        } finally {
                            wl.release()
                            reconnectWakeLock = null
                        }
                    }.apply {
                        isDaemon = true
                        name = "vpn-reconnect"
                    }.start()
                } else {
                    Log.e(TAG, "Auto-reconnect disabled or no URI")
                    disconnect()
                    stopSelf()
                }
            }

            // Start Xray core via libXray (in-process) — iterative retry, no stack accumulation
            if (isStopping) return

            var started = false
            while (!started && !isStopping) {
                started = XrayManager.start(this, cfg, dns1, dns2, resolvedIp)
                if (!started) {
                    if (autoReconnect && reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
                        reconnectAttempts++
                        val delay = 3000L * reconnectAttempts
                        Log.w(TAG, "Retrying Xray start in ${delay}ms (attempt $reconnectAttempts)")
                        updateNotification("Retry ($reconnectAttempts)...")
                        Thread.sleep(delay)
                    } else {
                        Log.e(TAG, "Failed to start Xray — giving up")
                        _state.value = VpnState.FAILED
                        updateNotification("Xray start failed")
                        stopSelf()
                        return
                    }
                }
            }
            if (isStopping) { XrayManager.stop(); return }

            // Probe SOCKS5 port until ready — max 10s (typically ready in 200-500ms)
            updateNotification("Waiting for Xray...")
            var portReady = false
            for (probe in 0 until 40) {
                if (isStopping) { XrayManager.stop(); return }
                try {
                    java.net.Socket("127.0.0.1", XrayManager.SOCKS_PORT).close()
                    portReady = true
                } catch (_: Exception) {}
                if (portReady) break
                Thread.sleep(250)
            }
            if (!portReady) {
                Log.e(TAG, "Xray SOCKS5 port not ready after 10s")
                XrayManager.stop()
                _state.value = VpnState.FAILED
                updateNotification("Xray timeout")
                stopSelf()
                return
            }

            // Establish TUN interface
            val builder = Builder()
                .setSession("JhopanStoreVPN")
                .addAddress("10.0.0.2", 24)
                .addRoute("0.0.0.0", 0)
                .addDnsServer(dns1.ifBlank { "8.8.8.8" })
                .addDnsServer(dns2.ifBlank { "8.8.4.4" })
                .setMtu(1400)

            // Exclude our own app as defense-in-depth (protectFd handles this too)
            builder.addDisallowedApplication(packageName)

            if (isStopping) {
                XrayManager.stop()
                return
            }

            tunFd = builder.establish()
            if (tunFd == null) {
                Log.e(TAG, "Failed to establish TUN interface")
                XrayManager.stop()
                _state.value = VpnState.FAILED
                updateNotification("TUN failed")
                stopSelf()
                return
            }

            val tunFdNum = tunFd!!.fd
            Log.d(TAG, "TUN fd number: $tunFdNum")

            // Clear O_CLOEXEC flag so tun2socks child process can inherit the fd
            try {
                val fileDescriptor = tunFd!!.fileDescriptor
                val flags = Os.fcntlInt(fileDescriptor, OsConstants.F_GETFD, 0)
                Os.fcntlInt(fileDescriptor, OsConstants.F_SETFD, flags and OsConstants.FD_CLOEXEC.inv())
                Log.d(TAG, "Cleared O_CLOEXEC on TUN fd")
            } catch (e: Exception) {
                Log.w(TAG, "Failed to clear O_CLOEXEC: ${e.message}")
            }

            // Start tun2socks to bridge TUN ↔ Xray SOCKS5 proxy
            val tun2socksStarted = Tun2socksManager.start(this, tunFdNum)
            if (!tun2socksStarted) {
                Log.e(TAG, "Failed to start tun2socks")
                XrayManager.stop()
                tunFd?.close()
                tunFd = null
                _state.value = VpnState.FAILED
                updateNotification("tun2socks failed")
                stopSelf()
                return
            }

            // Cek sekali lagi: kalau isStopping, jangan set CONNECTED
            if (isStopping) {
                disconnect()
                return
            }

            isRunning = true
            reconnectAttempts = 0
            _state.value = VpnState.CONNECTED
            Log.i(TAG, "VPN connected successfully (libXray + tun2socks)")
            updateNotification("Connected")

        } catch (e: Exception) {
            Log.e(TAG, "Connection error", e)
            if (!isStopping) _state.value = VpnState.FAILED
            disconnect()
            stopSelf()
        }
    }

    private fun disconnect() {
        isRunning = false
        XrayManager.onProcessDied = null

        // Release WakeLock jika masih aktif dari proses reconnect
        reconnectWakeLock?.let { if (it.isHeld) it.release() }
        reconnectWakeLock = null

        // Stop in reverse order: tun2socks first, then TUN, then Xray
        Tun2socksManager.stop()
        XrayManager.stop()

        try {
            tunFd?.close()
        } catch (e: IOException) {
            Log.e(TAG, "Error closing TUN", e)
        }
        tunFd = null

        _state.value = VpnState.DISCONNECTED
        stopForeground(STOP_FOREGROUND_REMOVE)
        Log.i(TAG, "VPN disconnected")
    }

    override fun onDestroy() {
        disconnect()
        super.onDestroy()
        
        // Release WakeLock saat service destroy
        serviceWakeLock?.let {
            if (it.isHeld) {
                it.release()
                Log.i(TAG, "Service WakeLock released")
            }
        }
        serviceWakeLock = null
    }

    override fun onRevoke() {
        autoReconnect = false
        disconnect()
        stopSelf()
        super.onRevoke()
    }

    // --- Notification helpers ---

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                "JhopanStoreVPN",
                // LOW: persistent notification without sound/vibration → less battery
                NotificationManager.IMPORTANCE_LOW
            ).apply {
                description = "VPN connection status"
                setShowBadge(false)
                enableLights(false)
                enableVibration(false)
            }
            val nm = getSystemService(NotificationManager::class.java)
            nm.createNotificationChannel(channel)
        }
    }

    private fun buildNotification(text: String): Notification {
        val pendingIntent = PendingIntent.getActivity(
            this, 0,
            Intent(this, MainActivity::class.java),
            PendingIntent.FLAG_IMMUTABLE or PendingIntent.FLAG_UPDATE_CURRENT
        )

        val stopIntent = PendingIntent.getService(
            this, 1,
            Intent(this, JhopanVpnService::class.java).apply { action = ACTION_STOP },
            PendingIntent.FLAG_IMMUTABLE or PendingIntent.FLAG_UPDATE_CURRENT
        )

        val builder = if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            Notification.Builder(this, CHANNEL_ID)
        } else {
            @Suppress("DEPRECATION")
            Notification.Builder(this)
                .setPriority(Notification.PRIORITY_HIGH)
        }

        return builder
            .setContentTitle("JhopanStoreVPN")
            .setContentText(text)
            .setSmallIcon(R.drawable.ic_vpn_key)
            .setContentIntent(pendingIntent)
            .addAction(
                Notification.Action.Builder(
                    null, "Disconnect", stopIntent
                ).build()
            )
            .setOngoing(true)
            .setCategory(Notification.CATEGORY_SERVICE)
            .setAutoCancel(false)
            .build()
    }

    private fun updateNotification(text: String) {
        val nm = getSystemService(NotificationManager::class.java)
        nm.notify(NOTIFICATION_ID, buildNotification(text))
    }
}
