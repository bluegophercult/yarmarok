package participant

import BaseTest
import api.controller.ParticipantController
import api.controller.RaffleController
import dto.participant.ParticipantCreateDto
import dto.raffle.RaffleCreateDto
import org.junit.jupiter.api.Test

class ParticipantTest : BaseTest() {
    @Test
    fun `create participant`() {
        val raffleDto = RaffleCreateDto("name", "haha")
        val raffleId = RaffleController.createRaffle(raffleDto)

        val participant =
            ParticipantController.addParticipant(raffleId.id, ParticipantCreateDto("part", "0983946652", "note"))
        val all = ParticipantController.getParticipant(raffleId.id)
    }
}