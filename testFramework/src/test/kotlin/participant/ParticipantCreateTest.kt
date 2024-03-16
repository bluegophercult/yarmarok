package participant

import BaseTest
import api.controller.ParticipantController
import dto.participant.ParticipantCreateDto
import org.apache.commons.lang3.RandomStringUtils
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.springframework.http.HttpStatus
import steps.RaffleSteps
import java.util.*

class ParticipantCreateTest : BaseTest() {
    @Test
    fun `create valid participant - success`() {
        val raffle = RaffleSteps.createRaffle()

        val participantId = ParticipantController.addParticipant(
            raffle.id,
            ParticipantCreateDto("part", "+380983946652", "note")
        )
        val participant = ParticipantController.getParticipant(raffle.id).firstOrNull { it.id == participantId.id }
        assertThat(participant).isNotNull
    }

    @Test
    fun `create participant with invalid raffle Id - should fail`() {
        val raffleId = UUID.randomUUID().toString()
        val response = ParticipantController.addParticipantWithoutValidation(
            raffleId,
            ParticipantCreateDto(RandomStringUtils.random(5), "+380983946652", "note")
        ).extract()

        assertThat(response.statusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR.value())
        assertThat(response.asString()).contains("Key: 'RaffleRequest.Name' Error:Field validation for 'Name' failed on the 'charsValidation' tag")
    }

    @Test
    fun `create participant without name - should fail`() {
        val raffle = RaffleSteps.createRaffle()

        val response = ParticipantController.addParticipantWithoutValidation(
            raffle.id,
            ParticipantCreateDto(RandomStringUtils.random(5), "+380983946652", "note")
        ).extract()

        assertThat(response.statusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR.value())
        assertThat(response.asString()).contains("Key: 'RaffleRequest.Name' Error:Field validation for 'Name' failed on the 'charsValidation' tag")
    }
}