package utils

import java.net.Socket
import java.net.SocketException

object Utils {
    fun isPortInUse( port: Int): Boolean {
        var result = false
        try {
            Socket("localhost", port).close()
            result = true
        } catch (_: SocketException) {
        }
        return result
    }
}