package participant

import BaseTest
import api.controller.ParticipantController
import dto.participant.ParticipantCreateDto
import org.junit.jupiter.api.Test
import steps.RaffleSteps

class ParticipantTest : BaseTest() {
    @Test
    fun `create participant`() {
        val raffle = RaffleSteps.createRaffle()

        val participant =
            ParticipantController.addParticipant(raffle.id, ParticipantCreateDto("part", "+380983946652", "note"))
        val all = ParticipantController.getParticipant(raffle.id)
    }
}