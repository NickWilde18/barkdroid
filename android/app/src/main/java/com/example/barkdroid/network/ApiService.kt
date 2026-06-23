package com.example.barkdroid.network

import android.util.Log
import okhttp3.*
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONObject
import java.io.IOException

/**
 * Talks to the barkdroid server to register this device.
 * After registration, the server returns a Bark-compatible key
 * that the user can use to send pushes.
 */
class ApiService(
    private val baseUrl: String,
    private val token: String,
) {
    private val client = OkHttpClient()
    private val jsonMediaType = "application/json; charset=utf-8".toMediaType()

    /**
     * Register this device with the server.
     * Called when JPush assigns a registration ID.
     */
    fun register(registrationId: String, onSuccess: (deviceKey: String) -> Unit) {
        val payload = JSONObject().apply {
            put("platform", "android")
            put("push_provider", "jpush")
            put("registration_id", registrationId)
        }

        val body = payload.toString().toRequestBody(jsonMediaType)

        val requestBuilder = Request.Builder()
            .url("$baseUrl/register")
            .post(body)
            .header("Content-Type", "application/json")

        if (token.isNotEmpty()) {
            requestBuilder.header("Authorization", "Bearer $token")
        }

        client.newCall(requestBuilder.build()).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                Log.e(TAG, "Register failed: ${e.message}")
            }

            override fun onResponse(call: Call, response: Response) {
                response.body?.string()?.let { body ->
                    Log.d(TAG, "Register response: $body")
                    try {
                        val json = JSONObject(body)
                        if (json.optInt("code") == 200) {
                            val data = json.getJSONObject("data")
                            val key = data.getString("key")
                            onSuccess(key)
                        }
                    } catch (e: Exception) {
                        Log.e(TAG, "Parse register response failed: ${e.message}")
                    }
                }
            }
        })
    }

    companion object {
        private const val TAG = "ApiService"
    }
}
