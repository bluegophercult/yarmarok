package raffle

import BaseTest
import api.controller.RaffleController
import dto.raffle.RaffleCreateDto
import org.apache.commons.lang3.RandomStringUtils
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import steps.RaffleSteps
import java.util.*

class RaffleUpdateTests : BaseTest() {
    @Test
    fun `update raffle with valid data - success`() {
        val raffle = RaffleSteps.createRaffle()

        val newRaffleDto = RaffleCreateDto(
            name = RandomStringUtils.randomAlphanumeric(5),
            note = RandomStringUtils.randomAlphanumeric(5)
        )

        RaffleController.updateRaffle(
            raffleId = raffle.id,
            entity = newRaffleDto
        )

        val raffleInDb = RaffleController.getRaffles().firstOrNull { it.id == raffle.id }
        assertThat(raffleInDb).isNotNull

        assertThat(raffleInDb!!.name).isEqualTo(newRaffleDto.name)
        assertThat(raffleInDb.note).isEqualTo(newRaffleDto.note)
    }

    @Test
    fun `update raffle with valid raffle id - should fail`() {
        val newRaffleDto = RaffleCreateDto(
            name = RandomStringUtils.randomAlphanumeric(5),
            note = RandomStringUtils.randomAlphanumeric(5)
        )

        val raffleId = UUID.randomUUID().toString()
        val response = RaffleController.updateRaffle(
            raffleId = raffleId,
            entity = newRaffleDto
        ).extract()

        assertThat(response.asString()).contains("get raffle: item not found")

        val raffleInDb = RaffleController.getRaffles().firstOrNull { it.id == raffleId }
        assertThat(raffleInDb).isNull()
    }

    @Test
    fun `update raffle with invalid name - should fail`() {
        val raffle = RaffleSteps.createRaffle()

        val newRaffleDto = RaffleCreateDto(
            name = RandomStringUtils.random(5),
            note = RandomStringUtils.randomAlphanumeric(5)
        )

        val response = RaffleController.updateRaffle(
            raffleId = raffle.id,
            entity = newRaffleDto
        ).extract()

        assertThat(response.asString()).contains("Key: 'RaffleRequest.Name' Error:Field validation for 'Name' failed on the 'charsValidation' tag")

        val raffleInDb = RaffleController.getRaffles().first { it.id == raffle.id }

        assertThat(raffleInDb.name).isNotEqualTo(newRaffleDto.name)
        assertThat(raffleInDb.note).isNotEqualTo(newRaffleDto.note)
    }

    @Test
    fun `update raffle with invalid note - should fail`() {
        val raffle = RaffleSteps.createRaffle()

        val newRaffleDto = RaffleCreateDto(
            name = RandomStringUtils.randomAlphanumeric(5),
            note = RandomStringUtils.random(5)
        )

        val response = RaffleController.updateRaffle(
            raffleId = raffle.id,
            entity = newRaffleDto
        ).extract()

        assertThat(response.asString()).contains("Key: 'RaffleRequest.Note' Error:Field validation for 'Note' failed on the 'charsValidation' tag")

        val raffleInDb = RaffleController.getRaffles().first { it.id == raffle.id }

        assertThat(raffleInDb.name).isNotEqualTo(newRaffleDto.name)
        assertThat(raffleInDb.note).isNotEqualTo(newRaffleDto.note)
    }

    @Test
    fun `update raffle with invalid name & note - should fail`() {
        val raffle = RaffleSteps.createRaffle()

        val newRaffleDto = RaffleCreateDto(
            name = RandomStringUtils.random(5),
            note = RandomStringUtils.random(5)
        )

        val response = RaffleController.updateRaffle(
            raffleId = raffle.id,
            entity = newRaffleDto
        ).extract()

        assertThat(response.asString()).contains(
            "Key: 'RaffleRequest.Name' Error:Field validation for 'Name' failed on the 'charsValidation' tag\n" +
                    "Key: 'RaffleRequest.Note' Error:Field validation for 'Note' failed on the 'charsValidation' tag"
        )

        val raffleInDb = RaffleController.getRaffles().first { it.id == raffle.id }

        assertThat(raffleInDb.name).isNotEqualTo(newRaffleDto.name)
        assertThat(raffleInDb.note).isNotEqualTo(newRaffleDto.note)
    }
}