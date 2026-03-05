package com.jhopanstore.vpn

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.runtime.LaunchedEffect
import androidx.lifecycle.viewmodel.compose.viewModel
import com.jhopanstore.vpn.ui.MainScreen
import com.jhopanstore.vpn.ui.MainViewModel
import com.jhopanstore.vpn.ui.theme.JhopanStoreVPNTheme

class MainActivity : ComponentActivity() {

    private val vpnPermissionLauncher = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (result.resultCode == RESULT_OK) {
            // VPN permission granted — start connection
            startVpn()
        }
    }

    private var pendingViewModel: MainViewModel? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        setContent {
            JhopanStoreVPNTheme {
                val vm: MainViewModel = viewModel()
                pendingViewModel = vm

                // Load persisted settings once
                LaunchedEffect(Unit) {
                    vm.loadSettings(this@MainActivity)
                }

                MainScreen(
                    viewModel = vm,
                    onConnect = { requestVpnPermission(vm) },
                    onDisconnect = { disconnectVpn(vm) },
                    onImportClipboard = { importFromClipboard(vm) }
                )
            }
        }
    }

    override fun onPause() {
        super.onPause()
        pendingViewModel?.saveSettings(this)
    }

    private fun requestVpnPermission(vm: MainViewModel) {
        val intent = android.net.VpnService.prepare(this)
        if (intent != null) {
            pendingViewModel = vm
            vpnPermissionLauncher.launch(intent)
        } else {
            // Already have permission
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
}
