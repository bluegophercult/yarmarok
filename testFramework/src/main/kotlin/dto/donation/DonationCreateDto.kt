package dto.donation

data class DonationCreateDto(
    val amount: Int,
    val participantId: String
)