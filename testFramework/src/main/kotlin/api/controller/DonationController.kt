package api.controller

import api.BaseApi
import dto.ResponseId
import dto.donation.DonationCreateDto
import dto.donation.DonationDto
import org.springframework.http.HttpStatus

object DonationController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun createDonation(raffleId: String, prizeId: String, donation: DonationCreateDto): ResponseId {
        return post("/api/raffles/$raffleId/prizes/$prizeId/donations", donation)
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(ResponseId::class.java)
    }

    fun getDonations(raffleId: String, prizeId: String): List<DonationDto> {
        return get("/api/raffles/$raffleId/prizes/$prizeId/donations")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", DonationDto::class.java)
    }

    fun getDonation(raffleId: String, prizeId: String, donationId: String): DonationDto {
        return get("/api/raffles/$raffleId/prizes/$prizeId/donations/$donationId")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(DonationDto::class.java)
    }

    fun updateDonation(raffleId: String, prizeId: String, donationId: String, donation: DonationCreateDto) {
        put("/api/raffles/$raffleId/prizes/$prizeId/donations/$donationId", donation)
            .then()
            .statusCode(HttpStatus.OK.value())
    }

    fun deleteDonation(raffleId: String, prizeId: String, donationId: String) {
        delete("/api/raffles/$raffleId/prizes/$prizeId/donations/$donationId")
            .then()
            .statusCode(HttpStatus.OK.value())
    }
}