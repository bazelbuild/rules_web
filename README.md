# Bazel Web Testing Rules

[![Build Status](https://travis-ci.org/bazelbuild/rules_webtesting.svg?branch=master)](https://travis-ci.org/bazelbuild/rules_webtesting)
[![Build status](https://badge.buildkite.com/d9c3974f925876394ca9d3e00670c0950b6f7ebf325412def7.svg?branch=master)](https://buildkite.com/bazel/rules-webtesting-saucelabs)

Bazel rules and supporting code to allow testing against a browser with
WebDriver.

## Configure your Bazel project

For all languages, you need to add the following to your WORKSPACE file:

```bzl
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_webtesting",
    sha256 = "<version-specific-sha>",
    urls = [
        "https://github.com/bazelbuild/rules_webtesting/releases/download/<version>/rules_webtesting.tar.gz",
    ],
)

load("@io_bazel_rules_webtesting//web:repositories.bzl", "web_test_repositories")

web_test_repositories()
```

We maintain predefined versioned browsers in
`@io_bazel_rules_webtesting//web/versioned:browsers-<version>.bzl` that can be
referenced from `@io_bazel_rules_webtesting//browsers`. Note: these versions in
these files are not locked in until the corresponding release has been created.

To use, add the following to your WORKSPACE file:

```bzl
load("@io_bazel_rules_webtesting//web/versioned:browsers-<version>.bzl", "browser_repositories")

browser_repositories(chromium=True, firefox=True)
```

Then you should add the appropriate dependencies depending on what language you
are writing your tests in:

### Java

```bzl
load("@io_bazel_rules_webtesting//web:java_repositories.bzl", "java_repositories")

java_repositories()
```

#### Using `rules_jvm_external`?

If you're using `rules_jvm_external` to manage your dependencies, you can add
the required artifacts directly to your `maven_install` instead of using
`java_repositories`.

```bzl
load("@io_bazel_rules_webtesting//web:java_repositories.bzl", "RULES_WEBTESTING_ARTIFACTS")

maven_install(
    artifacts = [
      # Your artifacts
    ] + RULES_WEBTESTING_ARTIFACTS,
    # Enabling compatability support is required.
    generate_compat_repositories = True,
)
```

### Kotlin

Follow the directions at [`bazelbuild/rules_kotlin`](https://github.com/bazelbuild/rules_kotlin) to setup the Kotlin toolchain. At the time of this writing, that is:

```bzl
# add the `bazelbuild/rules_kotlin` repo to your WORKSPACE, and then...

load("@io_bazel_rules_kotlin//kotlin:kotlin.bzl", "kotlin_repositories", "kt_register_toolchains")
kotlin_repositories()
kt_register_toolchains()
```

### Scala

```bzl
load("@io_bazel_rules_webtesting//web:java_repositories.bzl", "java_repositories")

java_repositories()

http_archive(
    name = "io_bazel_rules_scala",
    sha256 = "6c69597f373a01989b9f7119bd5d28cffc9cc35d44d1f6440c409d8ef420057d",
    strip_prefix = "rules_scala-da5ba6d97a1abdadef89d509b30a9dcfde7ffbe4",
    urls = [
        "https://github.com/bazelbuild/rules_scala/archive/da5ba6d97a1abdadef89d509b30a9dcfde7ffbe4.tar.gz",
    ],
)

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()
```

### Python

```bzl
load("@io_bazel_rules_webtesting//web:py_repositories.bzl", "py_repositories")

py_repositories()
```

### Go

```bzl
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "b7a62250a3a73277ade0ce306d22f122365b513f5402222403e507f2f997d421",
    urls = [
        "https://github.com/bazelbuild/rules_go/releases/download/0.16.3/rules_go-0.16.3.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "6e875ab4b6bf64a38c352887760f21203ab054676d9c1b274963907e0768740d",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.15.0/bazel-gazelle-0.15.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("@io_bazel_rules_webtesting//web:go_repositories.bzl", "go_repositories")

go_repositories()
```

## Write your tests

Write your test in the language of your choice, but use our provided Browser API
to get an instance of WebDriver.

### Example Java Test

```java
import com.google.testing.web.WebTest;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import org.openqa.selenium.WebDriver;

@RunWith(JUnit4.class)
public class BrowserTest {
  private WebDriver driver;

  @Before public void createDriver() {
    driver = new WebTest().newWebDriverSession();
  }

  @After public void quitDriver() {
    try {
      driver.quit();
     } finally {
      driver = null;
     }
   }

  // your tests here
}
```

### Example Kotlin Test

```kotlin
package javatests.kotlinsample

import com.google.testing.web.WebTest
import org.openqa.selenium.WebDriver

import org.junit.Test as test
import org.junit.After as after
import org.junit.Before as before


class SomeTest {
  private var driver: WebDriver? = null

  @before
  fun createDriver() {
    driver = WebTest().newWebDriverSession()
  }

  @after
  fun quitDriver() {
    try {
      driver!!.quit()
    } finally {
      driver = null
    }
  }

  @test
  fun testWebDriverFromKotlin() {
    val wt = WebTest()
    driver?.get(wt.HTTPAddress().resolve("/healthz").toString())
  }
}
```

### Example Python Test

```python
import unittest
from testing.web import webtest


class BrowserTest(unittest.TestCase):
  def setUp(self):
    self.driver = webtest.new_webdriver_session()

  def tearDown(self):
    try:
      self.driver.quit()
    finally:
      self.driver = None

  # Your tests here

if __name__ == "__main__":
  unittest.main()
```

### Example Go Test

```go
import (
    "testing"

    "github.com/tebeka/selenium"
    "github.com/bazelbuild/rules_webtesting/go/webtest"
)

func TestWebApp(t *testing.T) {
    wd, err := webtest.NewWebDriverSession(selenium.Capabilities{})
    if err != nil {
        t.Fatal(err)
    }

    // your test here

    if err := wd.Quit(); err != nil {
        t.Logf("Error quitting webdriver: %v", err)
    }
}
```

### BUILD file

In your BUILD files, load the correct language specific build rule and create a
test target using it:

```bzl
load("@io_bazel_rules_webtesting//web:py.bzl", "py_web_test_suite")

py_web_test_suite(
    name = "browser_test",
    srcs = ["browser_test.py"],
    browsers = [
        "@io_bazel_rules_webtesting//browsers:chromium-local",
    ],
    local = True,
    deps = ["@io_bazel_rules_webtesting//testing/web"],
)
```
