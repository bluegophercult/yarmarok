package controller.api

import application.GoApplication
import io.restassured.builder.RequestSpecBuilder
import io.restassured.filter.log.RequestLoggingFilter
import io.restassured.filter.log.ResponseLoggingFilter
import io.restassured.http.ContentType
import io.restassured.specification.RequestSpecification

object BaseApi {
    var requestSpecification: RequestSpecification =
        RequestSpecBuilder()
            .addFilter(RequestLoggingFilter())
            .addFilter(ResponseLoggingFilter())
            .setBaseUri("http://" + GoApplication.getHost())
            .setPort(GoApplication.getPort())
            .setContentType(ContentType.JSON)
            .build()
}