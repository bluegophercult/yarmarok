package participant

import BaseTest
import api.controller.ParticipantController
import dto.participant.ParticipantCreateDto
import org.apache.commons.lang3.RandomStringUtils
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import steps.ParticipantSteps
import steps.RaffleSteps
import utils.Utils

class ParticipantUpdateTests : BaseTest() {
    @Test
    fun `update participant - success`() {
        val raffle = RaffleSteps.createRaffle()
        val participantId = ParticipantSteps.createParticipant(raffle.id)

        val newParticipant = ParticipantCreateDto(
            name = RandomStringUtils.randomAlphanumeric(10),
            note = RandomStringUtils.randomAlphanumeric(10),
            phone = Utils.generatePhoneNumber()
        )

        ParticipantController.updateParticipant(
            raffleId = raffle.id,
            participantId = participantId.id,
            participant = newParticipant
        )

        val participant = ParticipantController.getParticipant(raffle.id).first { it.id == participantId.id }

        assertThat(participant.name).isEqualTo(newParticipant.name)
        assertThat(participant.note).isEqualTo(newParticipant.note)
        assertThat(participant.phone).isEqualTo(newParticipant.phone)
    }

    @Test
    fun `update participant with invalid name - should fail`() {
        val raffle = RaffleSteps.createRaffle()
        val participantId = ParticipantSteps.createParticipant(raffle.id)

        val newParticipant = ParticipantCreateDto(
            name = RandomStringUtils.random(10),
            note = RandomStringUtils.randomAlphanumeric(10),
            phone = Utils.generatePhoneNumber()
        )

        val response = ParticipantController.updateParticipantWithoutValidation(
            raffleId = raffle.id,
            participantId = participantId.id,
            participant = newParticipant
        ).extract()

        val participant = ParticipantController.getParticipant(raffle.id).first { it.id == participantId.id }

        assertThat(response.asString()).contains("Key: 'ParticipantRequest.Name' Error:Field validation for 'Name' failed on the 'charsValidation' tag")
        assertThat(participant.name).isNotEqualTo(newParticipant.name)
    }

    @Test
    fun `update participant with invalid note - should fail`() {
        val raffle = RaffleSteps.createRaffle()
        val participantId = ParticipantSteps.createParticipant(raffle.id)

        val newParticipant = ParticipantCreateDto(
            name = RandomStringUtils.randomAlphanumeric(10),
            note = RandomStringUtils.random(10),
            phone = Utils.generatePhoneNumber()
        )

        val response = ParticipantController.updateParticipantWithoutValidation(
            raffleId = raffle.id,
            participantId = participantId.id,
            participant = newParticipant
        ).extract()

        val participant = ParticipantController.getParticipant(raffle.id).first { it.id == participantId.id }

        assertThat(response.asString()).contains("Key: 'ParticipantRequest.Note' Error:Field validation for 'Note' failed on the 'charsValidation' tag")
        assertThat(participant.note).isNotEqualTo(newParticipant.note)
    }

    @Test
    fun `update participant with invalid phone - should fail`() {
        val raffle = RaffleSteps.createRaffle()
        val participantId = ParticipantSteps.createParticipant(raffle.id)

        val newParticipant = ParticipantCreateDto(
            name = RandomStringUtils.randomAlphanumeric(10),
            note = RandomStringUtils.randomAlphanumeric(10),
            phone = RandomStringUtils.randomAlphanumeric(3)
        )

        val response = ParticipantController.updateParticipantWithoutValidation(
            raffleId = raffle.id,
            participantId = participantId.id,
            participant = newParticipant
        ).extract()

        val participant = ParticipantController.getParticipant(raffle.id).first { it.id == participantId.id }

        assertThat(response.asString()).contains("ParticipantRequest.Phone' Error:Field validation for 'Phone' failed on the 'phoneValidation' tag")
        assertThat(participant.phone).isNotEqualTo(newParticipant.phone)
    }
}