package dto.prize

data class Prize (
    val id: String,
    val name: String,
    val ticketCost: Int,
    val description: String,
    val createdAt: String
)