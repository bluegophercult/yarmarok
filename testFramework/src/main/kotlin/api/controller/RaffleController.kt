package api.controller

import api.BaseApi
import dto.ResponseId
import dto.raffle.RaffleCreateDto
import dto.raffle.RaffleDto
import dto.raffle.RaffleExportResultDto
import io.restassured.response.ValidatableResponse
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

    fun updateRaffle(raffleId: String, entity: RaffleCreateDto): ValidatableResponse {
        return put("/api/raffles/$raffleId", entity)
            .then()
    }

    fun deleteRaffle(raffleId: String): ValidatableResponse {
        return delete("/api/raffles/$raffleId")
            .then()
    }

    fun downloadRaffles(raffleId: String): RaffleExportResultDto {
        return get("/api/raffles/$raffleId/download-xlsx")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(RaffleExportResultDto::class.java)
    }

    fun createRaffleWithoutValidation(entity: RaffleCreateDto): ValidatableResponse {
        return post("/api/raffles", entity)
            .then()
    }
}