package api.controller

import api.BaseApi
import dto.ResponseId
import dto.participant.ParticipantCreateDto
import dto.participant.ParticipantDto
import io.restassured.response.ValidatableResponse
import org.springframework.http.HttpStatus

object ParticipantController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun addParticipant(raffleId: String, participantCreate: ParticipantCreateDto): ResponseId {
        return addParticipantWithoutValidation(raffleId, participantCreate)
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
        updateParticipantWithoutValidation(raffleId, participantId, participant)
            .statusCode(HttpStatus.OK.value())
    }

    fun deleteParticipant(raffleId: String, participantId: String) {
        deleteParticipantWithoutValidation(raffleId, participantId)
            .statusCode(HttpStatus.OK.value())
    }

    fun addParticipantWithoutValidation(
        raffleId: String,
        participantCreate: ParticipantCreateDto
    ): ValidatableResponse {
        return post("/api/raffles/$raffleId/participants", participantCreate)
            .then()
    }

    fun deleteParticipantWithoutValidation(raffleId: String, participantId: String): ValidatableResponse {
        return delete("/api/raffles/$raffleId/participants/$participantId")
            .then()
    }

    fun updateParticipantWithoutValidation(
        raffleId: String,
        participantId: String,
        participant: ParticipantCreateDto
    ): ValidatableResponse {
        return put("/api/raffles/$raffleId/participants/$participantId", participant)
            .then()
    }
}