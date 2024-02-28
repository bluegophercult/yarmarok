package utils

import java.text.SimpleDateFormat
import java.util.*

object Logger {
    private val dateFormat = SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault())

    fun info(message: String) {
        log("INFO", message)
    }

    fun error(message: String) {
        log("ERROR", message)
    }

    private fun log(level: String, message: String) {
        val timestamp = dateFormat.format(Date())
        println("[$timestamp][$level] $message")
    }
}