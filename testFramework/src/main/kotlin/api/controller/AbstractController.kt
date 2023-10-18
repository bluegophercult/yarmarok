package api.controller

import io.restassured.RestAssured
import io.restassured.response.Response
import io.restassured.specification.RequestSpecification

abstract class AbstractController(var requestSpecification: RequestSpecification) {
    fun get(url: String): Response {
        return RestAssured.given(requestSpecification)
            .with()
            .get(url)
    }

    fun <T : Any> post(url: String, body: T): Response {
        return RestAssured.given(requestSpecification)
            .body(body)
            .with()
            .post(url)

    }

    fun <T : Any> put(url: String, body: T): Response {
        return RestAssured.given(requestSpecification)
            .body(body)
            .with()
            .put(url)

    }

    fun delete(url: String): Response {
        return RestAssured.given(requestSpecification)
            .with()
            .delete(url)

    }
}