import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.AfterAll
import org.junit.jupiter.api.BeforeAll
import utils.Utils
import java.io.File


abstract class BaseTest {
    companion object {
        private lateinit var connection: Process

        @JvmStatic
        @BeforeAll
        fun init() {
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
            assertThat(Utils.isPortInUse(8081)).isTrue()
        }

        @JvmStatic
        @AfterAll
        fun after() {
            println("Stopping main: ${connection.pid()}")

            connection.children().forEach {
                println("Stopping child: ${it.pid()}")
                if (it.isAlive) {
                    it.destroy()
                }
            }

            connection.destroy()
            connection.waitFor()
            println("Stopped ${connection.exitValue()}")
        }
    }
}