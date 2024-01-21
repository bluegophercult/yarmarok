package utils

import java.net.Socket
import java.net.SocketException


object Utils {
    fun isPortInUse( port: Int): Boolean {
        try {
            Socket("localhost", port).close()
        } catch (_: SocketException) {
            return false
        }
        return true
    }
}