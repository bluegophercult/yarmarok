package dto.prize

data class PrizeCreateDto(
    val name: String,
    val ticketCost: Int,
    val description: String
)