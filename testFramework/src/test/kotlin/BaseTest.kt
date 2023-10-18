import org.junit.jupiter.api.AfterAll
import org.junit.jupiter.api.BeforeAll
import java.io.File
import org.assertj.core.api.Assertions.assertThat


abstract class BaseTest {
    // TODO try to rewrite killing process and uncomment code
//    companion object {
//        private lateinit var connection: Process
//
//        @JvmStatic
//        @BeforeAll
//        fun init() {
//            connection = ProcessBuilder(
//                "go",
//                "test",
//                "-tags",
//                "local",
//                "-timeout",
//                "0",
//                "-count=1",
//                "-v",
//                "./testinfra/local/run_test.go"
//            )
//                .directory(File("./../"))
//                .redirectOutput(ProcessBuilder.Redirect.INHERIT)
//                .redirectError(ProcessBuilder.Redirect.INHERIT)
//                .start()
//            Thread.sleep(15_000)
//
//        }
//
//        @JvmStatic
//        @AfterAll
//        fun after() {
//            assertThat(utils.Utils.isPortInUse(8080)).isTrue()
//            assertThat(utils.Utils.isPortInUse(8081)).isTrue()
//            connection.destroy()
//            Thread.sleep(5_000)
//
//        }
//    }
}