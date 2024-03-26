package utils

import java.io.FileInputStream
import java.util.*

object EnvProperties {
    private val envProps = Properties().also { it.load(FileInputStream(".src/main/resources/env.properties")) }

    fun getHost(): String {
        return envProps["host"].toString()
    }

    fun getPort(): Int {
        return envProps["port"].toString().toInt()
    }
}