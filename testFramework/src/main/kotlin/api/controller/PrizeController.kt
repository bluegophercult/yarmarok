package api.controller

import api.BaseApi
import dto.ResponseId
import dto.prize.PrizeCreate
import dto.prize.Prize
import org.springframework.http.HttpStatus

object PrizeController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun createPrize(raffleId: String, prize: PrizeCreate): ResponseId {
        return post("/api/raffles/$raffleId/prizes", prize)
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(ResponseId::class.java)
    }

    fun getAllPrizes(raffleId: String): List<Prize> {
        return get("/api/raffles/$raffleId/prizes")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", Prize::class.java)
    }

    fun getPrize(raffleId: String, prizeId: String): Prize {
        return get("/api/raffles/$raffleId/prizes/$prizeId")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(Prize::class.java)
    }

    fun updatePrize(raffleId: String, prizeId: String, prize: PrizeCreate) {
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