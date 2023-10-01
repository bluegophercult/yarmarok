import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test

class FirstTest : BaseTest() {
    @Test
    fun `base first test`() {
        assertThat(1).isEqualTo(3 - 2)
    }
}