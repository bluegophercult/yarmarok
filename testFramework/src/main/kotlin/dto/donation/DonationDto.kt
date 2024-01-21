package dto.donation

data class DonationDto(
    val id: String,
    val prizeId: String,
    val participantId: String,
    val amount: Int,
    val ticketsNumber: Int,
    val createdAt: String
)