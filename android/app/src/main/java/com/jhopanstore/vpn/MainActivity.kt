package com.jhopanstore.vpn

import android.Manifest
import android.content.ClipData
import android.content.Intent
import android.content.pm.PackageManager
import android.net.Uri
import android.os.Build
import android.os.Bundle
import android.os.PowerManager
import android.provider.Settings
import android.widget.Toast
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.core.content.ContextCompat
import androidx.lifecycle.viewmodel.compose.viewModel
import com.jhopanstore.vpn.ui.MainScreen
import com.jhopanstore.vpn.ui.MainViewModel
import com.jhopanstore.vpn.ui.theme.JhopanStoreVPNTheme

class MainActivity : ComponentActivity() {

    private val vpnPermissionLauncher = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (result.resultCode == RESULT_OK) {
            startVpn()
        }
    }

    private val notificationPermissionLauncher = registerForActivityResult(
        ActivityResultContracts.RequestPermission()
    ) { /* tidak perlu action apapun — user sudah lihat dialog */ }

    private var pendingViewModel: MainViewModel? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        setContent {
            JhopanStoreVPNTheme {
                val vm: MainViewModel = viewModel()
                pendingViewModel = vm

                var showBatteryOptDialog by remember { mutableStateOf(false) }

                LaunchedEffect(Unit) {
                    vm.loadSettings(this@MainActivity)
                    vm.checkHotspot()
                    vm.syncConnectionState()
                    val prefs = getSharedPreferences("app_prefs", MODE_PRIVATE)
                    val alreadyAsked = prefs.getBoolean("battery_opt_asked", false)
                    val pm = getSystemService(POWER_SERVICE) as PowerManager
                    val isIgnoring = pm.isIgnoringBatteryOptimizations(packageName)
                    // Hanya tampilkan dialog jika belum pernah ditanya DAN belum di-disable battery opt
                    // Sekali user pilih "Nanti" atau "Nonaktifkan", dialog tidak akan muncul lagi
                    if (!alreadyAsked && !isIgnoring) {
                        showBatteryOptDialog = true
                    }
                }

                if (showBatteryOptDialog) {
                    AlertDialog(
                        onDismissRequest = {
                            // Gunakan commit() agar langsung tersimpan (synchronous)
                            getSharedPreferences("app_prefs", MODE_PRIVATE).edit()
                                .putBoolean("battery_opt_asked", true)
                                .commit()
                            showBatteryOptDialog = false
                            requestNotificationPermission()
                        },
                        title = { Text("Optimalkan Performa VPN") },
                        text = { Text("Nonaktifkan pembatasan daya agar VPN tetap aktif saat layar mati dan tidak terputus tiba-tiba.") },
                        confirmButton = {
                            TextButton(onClick = {
                                // Gunakan commit() agar langsung tersimpan
                                getSharedPreferences("app_prefs", MODE_PRIVATE).edit()
                                    .putBoolean("battery_opt_asked", true)
                                    .commit()
                                showBatteryOptDialog = false
                                startActivity(
                                    Intent(Settings.ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS).apply {
                                        data = Uri.parse("package:$packageName")
                                    }
                                )
                                requestNotificationPermission()
                            }) { Text("Nonaktifkan") }
                        },
                        dismissButton = {
                            TextButton(onClick = {
                                // Gunakan commit() agar langsung tersimpan
                                getSharedPreferences("app_prefs", MODE_PRIVATE).edit()
                                    .putBoolean("battery_opt_asked", true)
                                    .commit()
                                showBatteryOptDialog = false
                                requestNotificationPermission()
                            }) { Text("Nanti") }
                        }
                    )
                }

                MainScreen(
                    viewModel = vm,
                    onConnect = { requestVpnPermission(vm) },
                    onDisconnect = { disconnectVpn(vm) },
                    onImportClipboard = { importFromClipboard(vm) },
                    onOpenHotspotSettings = { openHotspotSettings() },
                    onToggleProxy = { vm.toggleProxySharing(this@MainActivity) },
                    onCopyProxy = { copyProxyToClipboard(vm) },
                    onOpenBatterySettings = {
                        startActivity(
                            Intent(Settings.ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS).apply {
                                data = Uri.parse("package:$packageName")
                            }
                        )
                    }
                )
            }
        }
    }

    override fun onResume() {
        super.onResume()
        pendingViewModel?.checkHotspot()
        pendingViewModel?.syncConnectionState()
        val pm = getSystemService(POWER_SERVICE) as PowerManager
        val isIgnoring = pm.isIgnoringBatteryOptimizations(packageName)
        // Tampilkan warning banner hanya jika belum dibebaskan dari optimasi baterai
        // User bisa dismiss per-sesi dengan tombol ×
        pendingViewModel?.isBatteryOptimized = !isIgnoring
        requestNotificationPermission()
    }

    override fun onPause() {
        super.onPause()
        // Background save dengan apply() saat pause
        pendingViewModel?.saveSettings(this, immediate = false)
    }

    override fun onStop() {
        super.onStop()
        // Immediate save dengan commit() saat stop untuk mencegah data loss
        // jika app di-kill dari recent apps atau system
        pendingViewModel?.saveSettings(this, immediate = true)
    }

    override fun onDestroy() {
        super.onDestroy()
        // Final save dengan commit() saat destroy - last resort
        pendingViewModel?.saveSettings(this, immediate = true)
    }

    private fun requestNotificationPermission() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            if (ContextCompat.checkSelfPermission(this, Manifest.permission.POST_NOTIFICATIONS)
                != PackageManager.PERMISSION_GRANTED
            ) {
                notificationPermissionLauncher.launch(Manifest.permission.POST_NOTIFICATIONS)
            }
        }
    }

    private fun requestVpnPermission(vm: MainViewModel) {
        val intent = android.net.VpnService.prepare(this)
        if (intent != null) {
            pendingViewModel = vm
            vpnPermissionLauncher.launch(intent)
        } else {
            startVpn()
        }
    }

    private fun startVpn() {
        pendingViewModel?.connect(this)
    }

    private fun disconnectVpn(vm: MainViewModel) {
        vm.disconnect(this)
    }

    private fun importFromClipboard(vm: MainViewModel) {
        val clipboard = getSystemService(CLIPBOARD_SERVICE) as android.content.ClipboardManager
        val clip = clipboard.primaryClip
        if (clip != null && clip.itemCount > 0) {
            val text = clip.getItemAt(0).text?.toString() ?: ""
            vm.importVlessUri(text)
        }
    }

    private fun openHotspotSettings() {
        try {
            // Buka langsung halaman tethering/hotspot
            val intent = Intent(Settings.ACTION_WIRELESS_SETTINGS)
            startActivity(intent)
        } catch (e: Exception) {
            // Fallback ke settings utama
            startActivity(Intent(Settings.ACTION_SETTINGS))
        }
    }

    private fun copyProxyToClipboard(vm: MainViewModel) {
        // Use HTTP port (10809) for WiFi manual proxy, not SOCKS5 (10808)
        val text = "${vm.hotspotIp}:10809"
        val clipboard = getSystemService(CLIPBOARD_SERVICE) as android.content.ClipboardManager
        clipboard.setPrimaryClip(ClipData.newPlainText("Proxy", text))
        Toast.makeText(this, "Disalin: $text", Toast.LENGTH_SHORT).show()
    }
}
