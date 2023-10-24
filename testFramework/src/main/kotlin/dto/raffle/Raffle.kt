package dto.raffle

data class Raffle(
    val id: String,
    val organizerId: String,
    val name: String,
    val note: String,
    val createdAt: String
)