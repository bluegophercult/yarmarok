package participant

import BaseTest
import api.controller.ParticipantController
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import steps.ParticipantSteps
import steps.RaffleSteps
import java.util.*

class ParticipantDeleteTests : BaseTest() {
    @Test
    fun `delete participant - success`() {
        val raffle = RaffleSteps.createRaffle()
        val participantId = ParticipantSteps.createParticipant(raffle.id)

        var participant = ParticipantController.getParticipant(raffle.id).firstOrNull { it.id == participantId.id }
        assertThat(participant).isNotNull

        ParticipantController.deleteParticipant(raffle.id, participantId.id)

        participant = ParticipantController.getParticipant(raffle.id).firstOrNull { it.id == participantId.id }
        assertThat(participant).isNull()
    }

    @Test
    fun `delete participant with invalid raffle id - should fail`() {
        val raffle = RaffleSteps.createRaffle()
        val participantId = ParticipantSteps.createParticipant(raffle.id)

        val deleteInvalidRaffleId =
            ParticipantController.deleteParticipantWithoutValidation(UUID.randomUUID().toString(), participantId.id)
                .extract()

        assertThat(deleteInvalidRaffleId.asString()).contains("deleting participant: item not found")

        val response =
            ParticipantController.deleteParticipantWithoutValidation(raffle.id, UUID.randomUUID().toString())
                .extract()

        assertThat(response.asString()).contains("deleting participant: item not found")
    }
}