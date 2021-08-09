package x.x.@appname@
import android.content.Context
import android.content.Intent
import android.util.Log
import android.content.BroadcastReceiver
import android.os.Build 
class MainReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context, intent: Intent) {

        if (Intent.ACTION_BOOT_COMPLETED == intent.action) {
            Log.d("@appname", "Heyoo~~")
            val intent = Intent(context, MainActivity::class.java)
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                context.startForegroundService(intent)
            } else {
                context.startService(intent)
            }
        } else {
            Log.d("@appname", "Some weird thing happened")
        }
    }
}
