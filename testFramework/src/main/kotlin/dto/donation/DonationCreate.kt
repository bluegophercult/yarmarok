package dto.donation

data class DonationCreate(
    val amount: Int,
    val participantId: String
)