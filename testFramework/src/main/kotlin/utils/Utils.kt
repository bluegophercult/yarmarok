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

    fun generatePhoneNumber():String{
        return "+380${(99999999..999999999).random()}"
    }
}