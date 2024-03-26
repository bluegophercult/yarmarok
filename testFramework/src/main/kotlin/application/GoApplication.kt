package application

import org.assertj.core.api.Assertions.assertThat
import utils.EnvProperties
import utils.Logger
import utils.Utils
import java.io.File

object GoApplication {
    private lateinit var connection: Process
    val HOST = EnvProperties.getHost()
    val PORT = EnvProperties.getPort()

    fun runApplication(){
        val pb = ProcessBuilder(
            "go",
            "run",
            "-tags",
            "local",
            "./testinfra/local/run.go"
        ).directory(File("./../"))
        pb.inheritIO()

        connection = pb.start()

        Thread.sleep(8_000)
        assertThat(Utils.isPortInUse(PORT)).isTrue()
    }

    fun stopApplication(){
        Logger.info("Stopping main: ${connection.pid()}")

        connection.children().forEach {
            Logger.info("Stopping child: ${it.pid()}")
            if (it.isAlive) {
                it.destroy()
            }
        }

        connection.destroy()
        connection.waitFor()
        Logger.info("Stopped ${connection.exitValue()}")
    }
}