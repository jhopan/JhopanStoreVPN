-keepclassmembers class * implements android.net.VpnService { *; }
-keep class com.jhopanstore.vpn.** { *; }

# libXray (Go mobile bridge) — jangan obfuscate class JNI/reflection dari .aar
-keep class go.** { *; }
-keep class libXray.** { *; }
-keepclassmembers class libXray.** { *; }
-keepattributes *Annotation*,Signature,InnerClasses,EnclosingMethod
