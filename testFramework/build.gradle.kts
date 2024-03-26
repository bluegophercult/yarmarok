plugins {
    kotlin("jvm") version "1.9.10"
    application
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(kotlin("test"))
    implementation("io.rest-assured:rest-assured:4.3.3")
    implementation("org.springframework:spring-web:6.0.12")
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin:2.15.+")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jsr310:2.3.0-rc1")
    implementation("org.assertj:assertj-core:3.11.1")
    testImplementation("org.assertj:assertj-core:3.11.1")
}

tasks.test {
    useJUnitPlatform()
}

task<Test>("itTest") {
    useJUnitPlatform()

    systemProperties["junit.jupiter.execution.parallel.enabled"] = true
    systemProperties["junit.jupiter.execution.parallel.config.strategy"] = "fixed"
    systemProperties["junit.jupiter.execution.parallel.config.fixed.parallelism"] = 8
    systemProperties["junit.jupiter.execution.parallel.mode.default"] = "concurrent"
}