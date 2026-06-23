package com.example.barkdroid.push

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import cn.jpush.android.api.JPushInterface
import org.json.JSONObject

/**
 * Receives system-level push events from JPush SDK.
 * - REGISTRATION: JPush assigned a registration ID → publish to event bus
 * - MESSAGE_RECEIVED: custom message received (in-app handling)
 * - NOTIFICATION_OPENED: user tapped a notification
 */
class JPushReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {
        val bundle = intent.extras ?: return
        val action = intent.action ?: return

        Log.d(TAG, "onReceive: action=$action")

        when (action) {
            JPushInterface.ACTION_REGISTRATION_ID -> {
                val regId = bundle.getString(JPushInterface.EXTRA_REGISTRATION_ID)
                if (!regId.isNullOrEmpty()) {
                    Log.i(TAG, "JPush registration ID: $regId")
                    PushEventBus.onRegistrationId(regId)
                }
            }

            JPushInterface.ACTION_MESSAGE_RECEIVED -> {
                val message = bundle.getString(JPushInterface.EXTRA_MESSAGE)
                val extras = bundle.getString(JPushInterface.EXTRA_EXTRA)
                Log.d(TAG, "Message received: $message, extras: $extras")
            }

            JPushInterface.ACTION_NOTIFICATION_OPENED -> {
                val extras = bundle.getString(JPushInterface.EXTRA_EXTRA)
                Log.d(TAG, "Notification opened, extras: $extras")
                // If the notification carried a URL extra, open it
                extras?.let {
                    try {
                        val json = JSONObject(it)
                        val url = json.optString("url", "")
                        if (url.isNotEmpty()) {
                            val openIntent = Intent(Intent.ACTION_VIEW, android.net.Uri.parse(url))
                            openIntent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                            context.startActivity(openIntent)
                        }
                    } catch (_: Exception) {}
                }
            }
        }
    }

    companion object {
        private const val TAG = "JPushReceiver"
    }
}

/**
 * Simple event bus to decouple JPushReceiver from MainActivity.
 */
object PushEventBus {
    private val _registrationId = MutableLiveData<String>()
    val registrationId: LiveData<String> = _registrationId

    fun onRegistrationId(regId: String) {
        _registrationId.postValue(regId)
    }
}
