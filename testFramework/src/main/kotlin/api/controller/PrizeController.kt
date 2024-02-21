package api.controller

import api.BaseApi
import dto.ResponseId
import dto.prize.PrizeCreateDto
import dto.prize.PrizeDto
import org.springframework.http.HttpStatus

object PrizeController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun createPrize(raffleId: String, prize: PrizeCreateDto): ResponseId {
        return post("/api/raffles/$raffleId/prizes", prize)
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(ResponseId::class.java)
    }

    fun getPrizes(raffleId: String): List<PrizeDto> {
        return get("/api/raffles/$raffleId/prizes")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", PrizeDto::class.java)
    }

    fun getPrize(raffleId: String, prizeId: String): PrizeDto {
        return get("/api/raffles/$raffleId/prizes/$prizeId")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(PrizeDto::class.java)
    }

    fun updatePrize(raffleId: String, prizeId: String, prize: PrizeCreateDto) {
        put("/api/raffles/$raffleId/prizes/$prizeId", prize)
            .then()
            .statusCode(HttpStatus.OK.value())
    }

    fun deletePrize(raffleId: String, prizeId: String) {
        delete("/api/raffles/$raffleId/prizes/$prizeId")
            .then()
            .statusCode(HttpStatus.OK.value())
    }
}