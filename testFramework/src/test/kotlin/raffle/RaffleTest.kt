package raffle

import BaseTest
import api.controller.RaffleController
import dto.raffle.RaffleCreateDto
import org.apache.commons.lang3.RandomStringUtils
import org.junit.jupiter.api.Test
import org.assertj.core.api.Assertions.assertThat

class RaffleTest : BaseTest() {

    @Test
    fun `create raffle`() {
        val raffleDto = RaffleCreateDto(RandomStringUtils.randomAlphanumeric(5), RandomStringUtils.randomAlphanumeric(5))
        val raffleId = RaffleController.createRaffle(raffleDto)
        val result = RaffleController.getRaffles().firstOrNull { it.id == raffleId.id }

        assertThat(result).isNotNull
        assertThat(result!!.name).isEqualTo(raffleDto.name)
        assertThat(result.note).isEqualTo(raffleDto.note)
    }

    @Test
    fun `create & delete raffle`() {
        val raffleDto = RaffleCreateDto(RandomStringUtils.random(5), RandomStringUtils.random(5))
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