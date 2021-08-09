package x.x.@appname@
import android.app.Service
import android.annotation.SuppressLint
import android.content.ComponentName
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.os.Message
import android.util.Log
import android.webkit.WebChromeClient
import android.webkit.WebView
import android.webkit.WebViewClient
import android.widget.TextView
import android.widget.Toast
import android.view.View
import androidx.appcompat.app.AppCompatActivity
import org.json.JSONObject
import java.io.BufferedReader
import java.io.File
import java.io.InputStreamReader
import java.util.concurrent.Executors
import android.os.IBinder

class MainActivity : Service() {
    override fun onStartCommand(intent: Intent, flags: Int, startId: Int): Int {
        Log.i("@appname@", "onStartCommand")
        Toast.makeText(
            applicationContext, 
            "This is a Service running in Background",
            Toast.LENGTH_SHORT
        ).show()
        startForeground()
        rund()
        return START_STICKY
    }
    override fun onTaskRemoved(rootIntent: Intent) {
        Log.d("@appname@", "onTaskRemoved")
        val restartServiceIntent = Intent(applicationContext, this.javaClass)
        restartServiceIntent.setPackage("x.x.@appname@")
        startService(restartServiceIntent)
        super.onTaskRemoved(rootIntent)
    }
    override fun onBind(intent: Intent): IBinder? {
        // TODO: Return the communication channel to the service.
        throw UnsupportedOperationException("Not yet implemented")
    }
    private fun exec(command: String, params: String): String {
        try {
            Log.d("@appname@", "Running command: " + command + " " + params)
            val process = ProcessBuilder()
                    .directory(File(filesDir.parentFile!!, "lib"))
                    .command(command, params)
                    .redirectErrorStream(true)
                    .start()
            val reader = BufferedReader(
                    InputStreamReader(process.inputStream)
            )
            var line = reader.readLine()
            while(line != null) {
                Log.i("@appname@", "command:$line")
                if (line.startsWith("goprod:")) {
                    val task = JSONObject(line.substring(7))
                    Log.i("@appname@", "goprod action: $task")
                    val type = task.getString("type")
                    val data = task.getJSONObject("data")
                    when (type) {
                        "android-intent" -> {
                            val uri = Uri.parse(data.getString("uri"))
                            val goprodintent = Intent(Intent.ACTION_VIEW)
                            goprodintent.setPackage(data.getString("package"))
                            goprodintent.data = uri
                            if (data.getBoolean("extraused")) {
                                val extra = data.getJSONArray("extra")
                                for (i in 0 until extra.length()) {
                                    val item = extra.getJSONObject(i)
                                    goprodintent.putExtra(item.getString("key"), item.getString("value"))
                                }
                            }
                            if (data.getBoolean("customcomponent")) {
                                val compo = data.getJSONObject("component")
                                goprodintent.component = ComponentName(compo.getString("pkg"), compo.getString("cls"))
                            }
                            startActivity(goprodintent)
                            Log.d("@appname", "activity started!")
                        }
                        "toast" -> {
                            Toast.makeText(this,  data.getJSONObject("data").getString("text"), Toast.LENGTH_SHORT).show()
                        }
                        else -> {
                            Log.i("@appname@", "Invalid type '$type' provided.")
                        }
                    }
                }
                line = reader.readLine()
            }
            Log.d("@appname@", "Process terminated")
            reader.close()
            process.waitFor()
            return "end"
        } catch (e: Exception) {
            Log.d("@appname@", e.message ?: "IOException")
            return e.message ?: "IOException"
        }
    }
    private fun rund() {
        val myPool = Executors.newFixedThreadPool(5)
        myPool.submit { exec(getApplicationInfo().nativeLibraryDir + "/libbin.so", "") }
    }
}