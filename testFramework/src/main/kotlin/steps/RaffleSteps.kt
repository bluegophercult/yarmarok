package steps

import api.controller.RaffleController
import dto.ResponseId
import dto.raffle.RaffleCreateDto
import org.apache.commons.lang3.RandomStringUtils

object RaffleSteps {
    fun createRaffle(): ResponseId {
        val raffleDto = RaffleCreateDto(
            name = RandomStringUtils.randomAlphanumeric(5),
            note = RandomStringUtils.randomAlphanumeric(5)
        )
        return RaffleController.createRaffle(raffleDto)
    }
}