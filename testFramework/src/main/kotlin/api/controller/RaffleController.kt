package api.controller

import api.BaseApi
import dto.ResponseId
import dto.raffle.RaffleCreateDto
import dto.raffle.RaffleExportResultDto
import dto.raffle.RaffleDto
import org.springframework.http.HttpStatus

object RaffleController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun getRaffles(): List<RaffleDto> {
        return get("/api/raffles")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", RaffleDto::class.java)
    }

    fun createRaffle(entity: RaffleCreateDto): ResponseId {
        return post("/api/raffles", entity)
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(ResponseId::class.java)
    }

    fun updateRaffle(raffleId: String, entity: RaffleCreateDto) {
        put("/api/raffles/$raffleId", entity)
            .then()
            .statusCode(HttpStatus.OK.value())
    }

    fun deleteRaffle(raffleId: String) {
        delete("/api/raffles/$raffleId")
            .then()
            .statusCode(HttpStatus.OK.value())
    }

    fun downloadRaffles(raffleId: String): RaffleExportResultDto {
        return get("/api/raffles/$raffleId/download-xlsx")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(RaffleExportResultDto::class.java)
    }
}