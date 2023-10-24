package api.controller

import api.BaseApi
import dto.ResponseId
import dto.participant.ParticipantCreate
import dto.participant.Participant
import org.springframework.http.HttpStatus

object ParticipantController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun addParticipant(raffleId: String, participantCreate: ParticipantCreate): ResponseId {
        return post("/api/raffles/$raffleId/participants", participantCreate)
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(ResponseId::class.java)
    }

    fun getParticipant(raffleId: String): List<Participant> {
        return get("/api/raffles/$raffleId/participants")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", Participant::class.java)
    }

    fun updateParticipant(raffleId: String, participantId: String, participant: ParticipantCreate) {
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