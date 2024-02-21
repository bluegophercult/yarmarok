package api.controller

import api.BaseApi
import dto.ResponseId
import dto.participant.ParticipantCreateDto
import dto.participant.ParticipantDto
import org.springframework.http.HttpStatus

object ParticipantController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun addParticipant(raffleId: String, participantCreate: ParticipantCreateDto): ResponseId {
        return post("/api/raffles/$raffleId/participants", participantCreate)
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(ResponseId::class.java)
    }

    fun getParticipant(raffleId: String): List<ParticipantDto> {
        return get("/api/raffles/$raffleId/participants")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", ParticipantDto::class.java)
    }

    fun updateParticipant(raffleId: String, participantId: String, participant: ParticipantCreateDto) {
        put("/api/raffles/$raffleId/participants/$participantId", participant)
            .then()
            .statusCode(HttpStatus.OK.value())
    }

    fun deleteParticipant(raffleId: String, participantId: String) {
        delete("/api/raffles/$raffleId/participants/$participantId")
            .then()
            .statusCode(HttpStatus.OK.value())
    }
}