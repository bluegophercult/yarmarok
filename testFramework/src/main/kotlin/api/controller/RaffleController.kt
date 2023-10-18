package api.controller

import api.BaseApi
import dto.ResponseId
import dto.raffle.RaffleCreateDto
import dto.raffle.RaffleGetDto
import org.springframework.http.HttpStatus

object RaffleController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun getRaffles(): List<RaffleGetDto> {
        return get("/api/raffles")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", RaffleGetDto::class.java)
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
}