package raffle

import BaseTest
import api.controller.RaffleController
import dto.raffle.RaffleCreateDto
import org.apache.commons.lang3.RandomStringUtils
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.http.HttpStatus

class RaffleCreateTests : BaseTest() {

    @Test
    fun `create raffle with valid data - success`() {
        val raffleDto = RaffleCreateDto(
            name = RandomStringUtils.randomAlphanumeric(5),
            note = RandomStringUtils.randomAlphanumeric(5)
        )
        val raffleId = RaffleController.createRaffle(raffleDto)
        val result = RaffleController.getRaffles().firstOrNull { it.id == raffleId.id }

        assertThat(result).isNotNull
        assertThat(result!!.name).isEqualTo(raffleDto.name)
        assertThat(result.note).isEqualTo(raffleDto.note)
    }

    @Test
    fun `create raffle without alphanumeric name - should fail`() {
        val raffleDto = RaffleCreateDto(
            name = RandomStringUtils.random(5),
            note = RandomStringUtils.randomAlphanumeric(5)
        )

        val response = RaffleController.createRaffleWithoutValidation(raffleDto).extract()
        assertThat(response.statusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR.value())
        assertThat(response.asString()).contains("Key: 'RaffleRequest.Name' Error:Field validation for 'Name' failed on the 'charsValidation' tag")
    }

    @Test
    fun `create raffle without alphanumeric note - should fail`() {
        val raffleDto = RaffleCreateDto(
            name = RandomStringUtils.randomAlphanumeric(5),
            note = RandomStringUtils.random(5)
        )

        val response = RaffleController.createRaffleWithoutValidation(raffleDto).extract()
        assertThat(response.statusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR.value())
        assertThat(response.asString()).contains("Key: 'RaffleRequest.Note' Error:Field validation for 'Note' failed on the 'charsValidation' tag")
    }

    @Test
    fun `create raffle without alphanumeric name and note - should fail`() {
        val raffleDto = RaffleCreateDto(
            name = RandomStringUtils.random(5),
            note = RandomStringUtils.random(5)
        )

        val response = RaffleController.createRaffleWithoutValidation(raffleDto).extract()
        assertThat(response.statusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR.value())
        assertThat(response.asString()).contains("Key: 'RaffleRequest.Name' Error:Field validation for 'Name' failed on the 'charsValidation' tag\n" +
                "Key: 'RaffleRequest.Note' Error:Field validation for 'Note' failed on the 'charsValidation' tag")
    }
}