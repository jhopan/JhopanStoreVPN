import java.net.URL
import java.util.zip.ZipInputStream

// libXray.aar is pre-built in app/libs/ — contains Xray-core as a Go library
// with native .so for arm64-v8a, armeabi-v7a, x86, x86_64

plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
}

android {
    namespace = "com.jhopanstore.vpn"
    compileSdk = 34

    defaultConfig {
        applicationId = "com.jhopanstore.vpn"
        minSdk = 24
        targetSdk = 34
        versionCode = 1
        versionName = "1.0.0"

        vectorDrawables {
            useSupportLibrary = true
        }

        ndk {
            // All supported ABIs — splits{} below produces per-ABI APKs + universal
            abiFilters += listOf("arm64-v8a", "armeabi-v7a", "x86_64", "x86")
        }
    }

    externalNativeBuild {
        cmake {
            path = file("src/main/cpp/CMakeLists.txt")
        }
    }

    buildTypes {
        debug {
            isMinifyEnabled = false
            applicationIdSuffix = ".debug"
        }
        release {
            isMinifyEnabled = true
            isShrinkResources = true
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
            // Use debug signing for now (replace with proper keystore for production)
            signingConfig = signingConfigs.getByName("debug")
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    kotlinOptions {
        jvmTarget = "17"
    }

    buildFeatures {
        compose = true
    }

    composeOptions {
        kotlinCompilerExtensionVersion = "1.5.5"
    }

    packaging {
        resources {
            excludes += "/META-INF/{AL2.0,LGPL2.1}"
        }
        jniLibs {
            useLegacyPackaging = true
        }
    }

    // ─── Per-ABI APK Splits ───────────────────────────────────────────────────
    // Produces separate lightweight APKs per architecture + one universal APK.
    //
    // Output APKs (./gradlew assembleRelease):
    //   app-arm64-v8a-release.apk    ~most Android phones 2016+  (recommended)
    //   app-armeabi-v7a-release.apk  ~older 32-bit ARM phones
    //   app-x86_64-release.apk       ~emulators, some Chromebooks
    //   app-x86-release.apk          ~legacy x86 emulators
    //   app-universal-release.apk    ~all ABIs bundled (largest, sideload-friendly)
    //
    // For Play Store multi-APK, assign unique versionCodes per ABI:
    //   arm64-v8a=4x, x86_64=3x, armeabi-v7a=2x, x86=1x  (x = base versionCode)
    splits {
        abi {
            isEnable = true
            reset()
            include("arm64-v8a", "armeabi-v7a", "x86_64", "x86")
            isUniversalApk = true
        }
    }
}

dependencies {
    // libXray — Xray-core as in-process Go library (replaces xray binary)
    implementation(files("libs/libXray.aar"))

    // Core Android
    implementation("androidx.core:core-ktx:1.12.0")
    implementation("androidx.lifecycle:lifecycle-runtime-ktx:2.6.2")
    implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.6.2")
    implementation("androidx.activity:activity-compose:1.8.1")
    // Jetpack Compose
    implementation(platform("androidx.compose:compose-bom:2023.10.01"))
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-graphics")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")
    implementation("androidx.compose.material:material-icons-extended")

    debugImplementation("androidx.compose.ui:ui-tooling")
    debugImplementation("androidx.compose.ui:ui-test-manifest")
}

// ═══════════════════════════════════════════════════════════════
//  Xray-core is now provided by libXray.aar (in-process Go lib).
//  No binary download needed.
// ═══════════════════════════════════════════════════════════════

// ═══════════════════════════════════════════════════════════════
//  Download tun2socks binary (bridges TUN ↔ SOCKS5 proxy).
//  Run manually:  ./gradlew downloadTun2socks
// ═══════════════════════════════════════════════════════════════
val tun2socksVersion = "v2.6.0"

tasks.register("downloadTun2socks") {
    description = "Download tun2socks binary from GitHub releases into jniLibs"
    group = "xray"

    val jniLibsDir = file("src/main/jniLibs")

    doLast {
        mapOf(
            "arm64-v8a"   to "linux-arm64",
            "armeabi-v7a" to "linux-armv7"
        ).forEach { (abi, archName) ->
            val dir = File(jniLibsDir, abi).also { it.mkdirs() }
            val target = File(dir, "libtun2socks.so")

            if (target.exists()) {
                println("✓  tun2socks $abi already present (${target.length() / 1024} KB)")
                return@forEach
            }

            val url = "https://github.com/xjasonlyu/tun2socks/releases/download/$tun2socksVersion/tun2socks-$archName.zip"
            println("⬇  Downloading tun2socks $abi from $url …")

            try {
                val tmp = File.createTempFile("tun2socks-$abi", ".zip")

                val connection = URL(url).openConnection()
                connection.connectTimeout = 30_000
                connection.readTimeout = 120_000
                connection.connect()
                connection.getInputStream().use { src ->
                    tmp.outputStream().use { dst -> src.copyTo(dst) }
                }

                var found = false
                ZipInputStream(tmp.inputStream()).use { zis ->
                    var entry = zis.nextEntry
                    while (entry != null) {
                        if (!entry.isDirectory && entry.name.startsWith("tun2socks")) {
                            target.outputStream().use { out -> zis.copyTo(out) }
                            found = true
                            break
                        }
                        entry = zis.nextEntry
                    }
                }
                tmp.delete()

                if (found) {
                    target.setExecutable(true)
                    println("✓  tun2socks $abi saved (${target.length() / 1024} KB)")
                } else {
                    println("✗  tun2socks binary not found in zip for $abi")
                }
            } catch (e: Exception) {
                println("✗  Failed to download tun2socks $abi: ${e.message}")
                println("   The app will try downloading at runtime instead.")
            }
        }
    }
}

// Automatically download tun2socks before build
tasks.matching { it.name == "preBuild" }.configureEach {
    dependsOn("downloadTun2socks")
}
