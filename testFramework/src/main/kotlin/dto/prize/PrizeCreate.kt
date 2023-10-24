package dto.prize

data class PrizeCreate(
    val name: String,
    val ticketCost: Int,
    val description: String
)