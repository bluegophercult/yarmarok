package utils

object ApiUtil {
    fun <T : Any> getListBody(response: String, toValueType: Class<T>): List<T> {
        return JsonUtil.createObjectMapper().readValue(
            response,
            JsonUtil.createObjectMapper().typeFactory.constructCollectionType(
                MutableList::class.java,
                toValueType
            )
        )
    }
}