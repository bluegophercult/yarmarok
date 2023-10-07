package dto.raffle

data class RaffleGetDto(
    val id: String,
    val organizerId: String,
    val name: String,
    val note: String,
    val createdAt: String
)