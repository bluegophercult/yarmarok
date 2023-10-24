package api.controller

import api.BaseApi
import dto.ResponseId
import dto.donation.DonationCreate
import dto.donation.Donation
import org.springframework.http.HttpStatus

object DonationController : AbstractController(requestSpecification = BaseApi.requestSpecification) {
    fun createDonation(raffleId: String, prizeId: String, donation: DonationCreate): ResponseId {
        return post("/api/raffles/$raffleId/prizes/$prizeId/donations", donation)
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(ResponseId::class.java)
    }

    fun getAllDonations(raffleId: String, prizeId: String): List<Donation> {
        return get("/api/raffles/$raffleId/prizes/$prizeId/donations")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body()
            .jsonPath().getList("items", Donation::class.java)
    }

    fun getDonation(raffleId: String, prizeId: String, donationId: String): Donation {
        return get("/api/raffles/$raffleId/prizes/$prizeId/donations/$donationId")
            .then()
            .statusCode(HttpStatus.OK.value())
            .extract().body().`as`(Donation::class.java)
    }

    fun updateDonation(raffleId: String, prizeId: String, donationId: String, donation: DonationCreate) {
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