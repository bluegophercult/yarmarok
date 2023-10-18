package raffle

import BaseTest
import api.controller.RaffleController
import dto.raffle.RaffleCreateDto
import org.junit.jupiter.api.Test
import org.assertj.core.api.Assertions.assertThat

class RaffleTest : BaseTest() {
    @Test
    fun `base first test`() {
        val raffleDto = RaffleCreateDto("name", "haha")
        val raffleId = RaffleController.createRaffle(raffleDto)
        var result = RaffleController.getRaffles().firstOrNull { it.id == raffleId.id }

        assertThat(result).isNotNull
        assertThat(result!!.name).isEqualTo(raffleDto.name)
        assertThat(result.note).isEqualTo(raffleDto.note)

        RaffleController.deleteRaffle(raffleId.id)

        result = RaffleController.getRaffles().firstOrNull { it.id == raffleId.id }
        assertThat(result).isNull()
    }
}