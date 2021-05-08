package x.x.@appname@
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
import androidx.appcompat.app.AppCompatActivity
import org.json.JSONObject
import java.io.BufferedReader
import java.io.File
import java.io.InputStreamReader
import java.util.concurrent.Executors


class MainActivity : AppCompatActivity() {
    private lateinit var webview: WebView
    @SuppressLint("SetJavaScriptEnabled")
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        supportActionBar?.hide()
        webview = findViewById(R.id.webview)
        webview.settings.javaScriptEnabled = true
        webview.settings.domStorageEnabled = true
        webview.settings.loadWithOverviewMode = true
        webview.settings.useWideViewPort = true
        webview.settings.builtInZoomControls = true
        webview.settings.displayZoomControls = false
        webview.settings.setSupportZoom(true)
        webview.settings.defaultTextEncodingName = "utf-8"
        webview.webChromeClient = WebChromeClient()
        webview.webViewClient = WebViewClient()
        webview.getSettings().setSupportMultipleWindows(true)
        webview.setWebChromeClient(object : WebChromeClient() {
            override fun onCreateWindow(
                    view: WebView,
                    dialog: Boolean,
                    userGesture: Boolean,
                    resultMsg: Message
            ): Boolean {
                val result = view.hitTestResult
                val data = result.extra
                val context: Context = view.context
                val browserIntent = Intent(Intent.ACTION_VIEW, Uri.parse(data))
                context.startActivity(browserIntent)
                return false
            }
        })
        webview.loadUrl("data:text/html,<meta http-equiv=\"Refresh\" content=\"1; url='http://127.0.0.1:@port@'\" />")
        rund()
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