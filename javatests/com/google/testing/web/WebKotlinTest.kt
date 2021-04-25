package com.google.testing.web

import com.google.testing.web.WebTest
import org.openqa.selenium.WebDriver

import org.junit.runner.RunWith
import org.junit.runners.JUnit4
import org.junit.Test as test


@RunWith(JUnit4::class)
class WebKotlinTest {
  @test
  fun testWebDriverFromKotlin() {
    val wt = WebTest()
    val driver = wt.newWebDriverSession()
    driver.get(wt.HTTPAddress().resolve("/healthz").toString())
  }
}
