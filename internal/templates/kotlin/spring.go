package kotlin

const (
	KotlinSpringBuildGradle = `plugins {
	id 'org.springframework.boot' version '3.2.0'
	id 'io.spring.dependency-management' version '1.1.4'
	id 'org.jetbrains.kotlin.jvm' version '1.9.20'
	id 'org.jetbrains.kotlin.plugin.spring' version '1.9.20'
}

group = 'com.example'
version = '0.0.1-SNAPSHOT'

java {
	sourceCompatibility = '17'
}

repositories {
	mavenCentral()
}

dependencies {
	implementation 'org.springframework.boot:spring-boot-starter-web'
	implementation 'com.fasterxml.jackson.module:jackson-module-kotlin'
	implementation 'org.jetbrains.kotlin:kotlin-reflect'
	testImplementation 'org.springframework.boot:spring-boot-starter-test'
}`

	KotlinSpringApplication = `package com.example.{{.ProjectName}}

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class Application

fun main(args: Array<String>) {
	runApplication<Application>(*args)
}`

	KotlinSpringController = `package com.example.{{.ProjectName}}.controller

import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController

@RestController
class HomeController {

    @GetMapping("/")
    fun home(): Map<String, String> {
        return mapOf("message" to "Hello World from {{.ProjectName}}!")
    }

    @GetMapping("/health")
    fun health(): Map<String, String> {
        return mapOf("status" to "ok")
    }
}`
)
