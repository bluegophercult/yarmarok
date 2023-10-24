package api.controller

import api.BaseApi
import dto.ResponseId
import dto.raffle.RaffleCreate
import dto.raffle.RaffleExportResult
import dto.raffle.Raffle
import org.springframework.http.HttpStatus

object RaffleController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun getRaffles(): List<Raffle> {
        return get("/api/raffles")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", Raffle::class.java)
    }

    fun createRaffle(entity: RaffleCreate): ResponseId {
        return post("/api/raffles", entity)
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(ResponseId::class.java)
    }

    fun updateRaffle(raffleId: String, entity: RaffleCreate) {
        put("/api/raffles/$raffleId", entity)
            .then()
            .statusCode(HttpStatus.OK.value())
    }

    fun deleteRaffle(raffleId: String) {
        delete("/api/raffles/$raffleId")
            .then()
            .statusCode(HttpStatus.OK.value())
    }

    fun downloadRaffles(raffleId: String): RaffleExportResult {
        return get("/api/raffles/$raffleId/download-xlsx")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(RaffleExportResult::class.java)
    }
}