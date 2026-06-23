plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
    id("org.jetbrains.kotlin.plugin.compose")
}

android {
    namespace = "com.example.barkdroid"
    compileSdk = 35

    defaultConfig {
        applicationId = "com.example.barkdroid"
        minSdk = 26
        targetSdk = 35
        versionCode = 1
        versionName = "1.0.0"
    }

    // Read config from config.properties (user-provided)
    val configFile = rootProject.file("config.properties")
    if (configFile.exists()) {
        val props = java.util.Properties()
        props.load(configFile.inputStream())

        defaultConfig {
            // Let the user override applicationId
            props.getProperty("APPLICATION_ID")?.let { applicationId = it }
            // JPush AppKey (safe to embed in APK — only for registration, not sending)
            buildConfigField("String", "JPUSH_APP_KEY", "\"${props.getProperty("JPUSH_APP_KEY", "")}\"")
            buildConfigField("String", "SERVER_BASE_URL", "\"${props.getProperty("SERVER_BASE_URL", "http://10.0.2.2:8080")}\"")
            buildConfigField("String", "SERVER_TOKEN", "\"${props.getProperty("SERVER_TOKEN", "")}\"")
        }
    } else {
        defaultConfig {
            buildConfigField("String", "JPUSH_APP_KEY", "\"\"")
            buildConfigField("String", "SERVER_BASE_URL", "\"http://10.0.2.2:8080\"")
            buildConfigField("String", "SERVER_TOKEN", "\"\"")
        }
    }

    buildTypes {
        release {
            isMinifyEnabled = true
            proguardFiles(getDefaultProguardFile("proguard-android-optimize.txt"), "proguard-rules.pro")
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
        buildConfig = true
    }
}

dependencies {
    // Compose
    implementation(platform("androidx.compose:compose-bom:2024.06.00"))
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.material3:material3")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.activity:activity-compose:1.9.0")

    // JPush SDK
    implementation("cn.jiguang.sdk:jpush:5.4.0")
    implementation("cn.jiguang.sdk:jcore:4.3.0")

    // Network
    implementation("com.squareup.okhttp3:okhttp:4.12.0")

    // Lifecycle
    implementation("androidx.lifecycle:lifecycle-runtime-ktx:2.8.2")
}
