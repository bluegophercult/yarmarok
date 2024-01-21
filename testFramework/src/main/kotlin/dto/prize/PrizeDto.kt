package dto.prize

data class PrizeDto (
    val id: String,
    val name: String,
    val ticketCost: Int,
    val description: String,
    val createdAt: String
)