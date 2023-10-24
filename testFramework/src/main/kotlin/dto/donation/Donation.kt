package dto.donation

data class Donation(
    val id: String,
    val prizeId: String,
    val participantId: String,
    val amount: Int,
    val ticketsNumber: Int,
    val createdAt: String
)