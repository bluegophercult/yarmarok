package participant

import BaseTest
import api.controller.ParticipantController
import dto.participant.ParticipantCreateDto
import org.apache.commons.lang3.RandomStringUtils
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.http.HttpStatus
import steps.RaffleSteps
import utils.Utils.generatePhoneNumber
import java.util.*

class ParticipantCreateTests : BaseTest() {
    @Test
    fun `create valid participant - success`() {
        val raffle = RaffleSteps.createRaffle()

        val participantId = ParticipantController.addParticipant(
            raffle.id,
            ParticipantCreateDto(
                name = RandomStringUtils.randomAlphanumeric(5),
                phone = generatePhoneNumber(),
                note = RandomStringUtils.randomAlphanumeric(5)
            )
        )

        val participant = ParticipantController.getParticipant(raffle.id).firstOrNull { it.id == participantId.id }
        assertThat(participant).isNotNull
    }

    @Test
    fun `create participant with invalid raffle Id - should fail`() {
        val raffleId = UUID.randomUUID().toString()

        val response = ParticipantController.addParticipantWithoutValidation(
            raffleId,
            ParticipantCreateDto(
                name = RandomStringUtils.random(5),
                phone = generatePhoneNumber(),
                note = RandomStringUtils.randomAlphanumeric(5)
            )
        ).extract()

        assertThat(response.statusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR.value())
        assertThat(response.asString()).contains("Key: 'ParticipantRequest.Name' Error:Field validation for 'Name' failed on the 'charsValidation' tag")
    }

    @Test
    fun `create participant with invalid note - should fail`() {
        val raffle = RaffleSteps.createRaffle()

        val response = ParticipantController.addParticipantWithoutValidation(
            raffle.id,
            ParticipantCreateDto(
                name = RandomStringUtils.randomAlphanumeric(5),
                phone = generatePhoneNumber(),
                note = RandomStringUtils.random(5)
            )
        ).extract()

        assertThat(response.statusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR.value())
        assertThat(response.asString()).contains("Key: 'ParticipantRequest.Note' Error:Field validation for 'Note' failed on the 'charsValidation' tag")
    }

    @Test
    fun `create participant with invalid phone - should fail`() {
        val raffle = RaffleSteps.createRaffle()

        val response = ParticipantController.addParticipantWithoutValidation(
            raffle.id,
            ParticipantCreateDto(
                name = RandomStringUtils.randomAlphanumeric(5),
                phone = "",
                note = RandomStringUtils.randomAlphanumeric(5)
            )
        ).extract()

        assertThat(response.statusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR.value())
        assertThat(response.asString()).contains("Key: 'ParticipantRequest.Phone' Error:Field validation for 'Phone' failed on the 'required' tag")
    }
}