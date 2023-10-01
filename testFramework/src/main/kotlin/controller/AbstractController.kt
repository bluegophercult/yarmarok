package controller

import io.restassured.RestAssured
import io.restassured.response.Response
import io.restassured.specification.RequestSpecification

abstract class AbstractController(var requestSpecification: RequestSpecification) {
    fun get(url: String): Response {
        return RestAssured.given(requestSpecification)
            .with()
            .get(url)
    }
}