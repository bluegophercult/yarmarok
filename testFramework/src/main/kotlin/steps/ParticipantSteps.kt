package steps

import api.controller.ParticipantController
import dto.ResponseId
import dto.participant.ParticipantCreateDto
import org.apache.commons.lang3.RandomStringUtils
import utils.Utils

object ParticipantSteps {
    fun createParticipant(raffleId: String): ResponseId {
        return ParticipantController.addParticipant(
            raffleId,
            ParticipantCreateDto(
                name = RandomStringUtils.randomAlphanumeric(5),
                phone = Utils.generatePhoneNumber(),
                note = RandomStringUtils.randomAlphanumeric(5)
            )
        )
    }
}