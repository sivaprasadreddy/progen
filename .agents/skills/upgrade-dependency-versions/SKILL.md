---
name: upgrade-dependency-versions
description: Checks Maven Central, the Gradle Plugin Portal, nodejs.org and the npm registry for newer versions of the dependencies and plugins used in generated Spring Boot projects, and updates generators/springboot/versions.go. Use when asked to upgrade, bump, refresh or check dependency/plugin versions for the generated Maven (pom.xml.tmpl) or Gradle (build.gradle.tmpl) builds.
---

# Upgrade dependency versions for generated Spring Boot projects

All versions rendered into `generators/springboot/templates/pom.xml.tmpl` and
`generators/springboot/templates/build.gradle.tmpl` come from the constants in
`generators/springboot/versions.go`. This skill finds newer releases and updates
those constants. Only edit `versions.go` — never hardcode versions in templates.

## Steps

1. **Check for newer versions.** Run from the repo root:

   ```bash
   .agents/skills/upgrade-dependency-versions/scripts/check-versions.sh | column -t
   ```

   It prints `CONSTANT  CURRENT  LATEST  STATUS` for every constant mapped to an
   artifact. `STATUS` is one of:
   - `up-to-date` — nothing to do.
   - `update` — newer minor/patch release in the same major line. Apply it.
   - `MAJOR` — newer major version. **Do not apply automatically**; list these
     for the user and ask (major bumps often break Spring Boot compatibility or
     template syntax). Apply only the ones the user approves.
   - `current-is-newer` / `LOOKUP_FAILED` — report it; verify the coordinate
     manually (see mapping below) instead of guessing a version.

2. **Respect compatibility before applying.**
   - `springModulithVersion`, `springCloudAWSVersion`, `springdocOpenapiVersion`,
     `htmxSpringBootThymeleafVersion`, `datasourceMicrometerSpringBootVersion`
     and `springBootFlociTcVersion` track Spring Boot releases. If Spring Boot is
     not being bumped to a new minor, skip a minor bump of these unless their
     release notes say they support the current Boot version.
   - `opentelemetryLogbackAppenderVersion` (instrumentation) and
     `opentelemetryAPIIncubatorVersion` (SDK) are released in pairs: the
     instrumentation release N depends on a specific SDK version. Pick the
     incubator version that matches the chosen instrumentation release (check its
     POM on Maven Central) rather than blindly taking the latest of each.
   - `nodeVersion` stays on the current LTS major (script already limits this);
     `npmVersion` must be compatible with that Node version.

3. **Edit `generators/springboot/versions.go`** — change only the constant
   values; keep alignment (`gofmt` handles it via `make fmt`).

4. **Verify:**

   ```bash
   make fmt
   make test
   ```

   For a Spring Boot or other significant bump, also run the slow suite, which
   generates and builds real projects:

   ```bash
   make all_tests
   ```

   If a build fails because of an upgrade, revert that single constant and
   report it rather than changing templates to force it through (unless the
   user asks for the migration).

5. **Report** a table of `constant: old → new`, the skipped `MAJOR` updates,
   and any lookup failures.

## Constant → artifact mapping

| Constant | Repository | Coordinates |
|---|---|---|
| springBootVersion | Maven Central | org.springframework.boot:spring-boot |
| gradleDependencyMgmtVersion | Gradle Plugin Portal | io.spring.dependency-management |
| gradleGitPropertiesVersion | Gradle Plugin Portal | com.gorylenko.gradle-git-properties |
| spotlessGradlePluginVersion | Gradle Plugin Portal | com.diffplug.spotless |
| spotlessMavenPluginVersion | Maven Central | com.diffplug.spotless:spotless-maven-plugin |
| springModulithVersion | Maven Central | org.springframework.modulith:spring-modulith-bom |
| springCloudAWSVersion | Maven Central | io.awspring.cloud:spring-cloud-aws-dependencies |
| springBootFlociTcVersion | Maven Central | io.floci:spring-boot-testcontainers-floci |
| springdocOpenapiVersion | Maven Central | org.springdoc:springdoc-openapi-starter-webmvc-ui |
| tcMailpitVersion | Maven Central | ch.martinelli.oss:testcontainers-mailpit |
| fontAwesomeVersion | Maven Central | org.webjars:font-awesome |
| htmxVersion | Maven Central | org.webjars.npm:htmx.org |
| htmxSpringBootThymeleafVersion | Maven Central | io.github.wimdeblauwe:htmx-spring-boot-thymeleaf |
| datasourceMicrometerSpringBootVersion | Maven Central | net.ttddyy.observation:datasource-micrometer-spring-boot |
| jjwtVersion | Maven Central | io.jsonwebtoken:jjwt-api |
| bootUIVersion | Maven Central | com.julien-dubois.bootui:bootui-spring-boot-starter |
| archunitVersion | Maven Central | com.tngtech.archunit:archunit-junit5 |
| taikaiVersion | Maven Central | com.enofex:taikai |
| palantirJavaFormatVersion | Maven Central | com.palantir.javaformat:palantir-java-format |
| testcontainersJooqCodegenPluginVersion | Maven Central | dev.sivalabs:testcontainers-jooq-codegen-maven-plugin |
| jacocoMavenPluginVersion | Maven Central | org.jacoco:jacoco-maven-plugin |
| frontendMavenPluginVersion | Maven Central | com.github.eirslett:frontend-maven-plugin |
| gradleNodePluginVersion | Gradle Plugin Portal | com.github.node-gradle.node |
| mavenAntrunPluginVersion | Maven Central | org.apache.maven.plugins:maven-antrun-plugin |
| opentelemetryLogbackAppenderVersion | Maven Central | io.opentelemetry.instrumentation:opentelemetry-logback-appender-1.0 (`-alpha`) |
| opentelemetryAPIIncubatorVersion | Maven Central | io.opentelemetry:opentelemetry-api-incubator (`-alpha`) |
| nodeVersion | nodejs.org/dist/index.json | latest LTS in current major |
| npmVersion | registry.npmjs.org | npm `latest` dist-tag |

Manual lookup URLs:
- Maven Central: `https://repo1.maven.org/maven2/<groupId with / >/<artifactId>/maven-metadata.xml`
- Gradle plugin: `https://plugins.gradle.org/m2/<plugin id with / >/<plugin id>.gradle.plugin/maven-metadata.xml`

Not covered: `javaVersion` and the Docker image constants (`*Image`) — change
those only when the user explicitly asks.

If a new constant is added to `versions.go`, add it to the `DEPS` list in
`scripts/check-versions.sh` and to the table above.
