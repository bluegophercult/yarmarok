package raffle

import BaseTest
import api.controller.RaffleController
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.http.HttpStatus
import steps.RaffleSteps
import java.util.*

class RaffleDeleteTests : BaseTest() {
    @Test
    fun `delete raffle - success`() {
        val raffle = RaffleSteps.createRaffle()

        var result = RaffleController.getRaffles().firstOrNull { it.id == raffle.id }
        assertThat(result).isNotNull

        val response = RaffleController.deleteRaffle(raffle.id).extract()
        assertThat(response.statusCode()).isEqualTo(HttpStatus.OK.value())

        result = RaffleController.getRaffles().firstOrNull { it.id == raffle.id }
        assertThat(result).isNull()
    }

    @Test
    fun `delete raffle with invalid id - should fail`() {
        val response = RaffleController.deleteRaffle(UUID.randomUUID().toString()).extract()

        assertThat(response.asString()).contains("deleting raffle: item not found")
    }
}