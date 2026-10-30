#!/usr/bin/env bash
# Checks Maven Central / Gradle Plugin Portal / nodejs.org / npm registry for the
# latest versions of the dependencies declared in generators/springboot/versions.go.
# Prints a TSV table: CONSTANT  CURRENT  LATEST  STATUS
set -uo pipefail

ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
VERSIONS_FILE="$ROOT/generators/springboot/versions.go"

MAVEN="https://repo1.maven.org/maven2"
GRADLE="https://plugins.gradle.org/m2"

# constant | repo | groupId:artifactId (or plugin id) | allow-prerelease-suffix
DEPS=$(cat <<'LIST'
springBootVersion|maven|org.springframework.boot:spring-boot|
gradleDependencyMgmtVersion|gradle|io.spring.dependency-management|
gradleGitPropertiesVersion|gradle|com.gorylenko.gradle-git-properties|
spotlessGradlePluginVersion|gradle|com.diffplug.spotless|
spotlessMavenPluginVersion|maven|com.diffplug.spotless:spotless-maven-plugin|
springModulithVersion|maven|org.springframework.modulith:spring-modulith-bom|
springCloudAWSVersion|maven|io.awspring.cloud:spring-cloud-aws-dependencies|
springBootFlociTcVersion|maven|io.floci:spring-boot-testcontainers-floci|
springdocOpenapiVersion|maven|org.springdoc:springdoc-openapi-starter-webmvc-ui|
tcMailpitVersion|maven|ch.martinelli.oss:testcontainers-mailpit|
fontAwesomeVersion|maven|org.webjars:font-awesome|
htmxVersion|maven|org.webjars.npm:htmx.org|
htmxSpringBootThymeleafVersion|maven|io.github.wimdeblauwe:htmx-spring-boot-thymeleaf|
datasourceMicrometerSpringBootVersion|maven|net.ttddyy.observation:datasource-micrometer-spring-boot|
jjwtVersion|maven|io.jsonwebtoken:jjwt-api|
bootUIVersion|maven|com.julien-dubois.bootui:bootui-spring-boot-starter|
archunitVersion|maven|com.tngtech.archunit:archunit-junit5|
taikaiVersion|maven|com.enofex:taikai|
palantirJavaFormatVersion|maven|com.palantir.javaformat:palantir-java-format|
testcontainersJooqCodegenPluginVersion|maven|dev.sivalabs:testcontainers-jooq-codegen-maven-plugin|
jacocoMavenPluginVersion|maven|org.jacoco:jacoco-maven-plugin|
frontendMavenPluginVersion|maven|com.github.eirslett:frontend-maven-plugin|
gradleNodePluginVersion|gradle|com.github.node-gradle.node|
mavenAntrunPluginVersion|maven|org.apache.maven.plugins:maven-antrun-plugin|
opentelemetryLogbackAppenderVersion|maven|io.opentelemetry.instrumentation:opentelemetry-logback-appender-1.0|alpha
opentelemetryAPIIncubatorVersion|maven|io.opentelemetry:opentelemetry-api-incubator|alpha
LIST
)

current_version() {
  grep -E "^[[:space:]]*$1[[:space:]]*=" "$VERSIONS_FILE" | sed -E 's/.*=[[:space:]]*"?([^"]*)"?.*/\1/'
}

metadata_url() {
  local repo=$1 coord=$2
  if [[ $repo == gradle ]]; then
    echo "$GRADLE/${coord//.//}/$coord.gradle.plugin/maven-metadata.xml"
  else
    local g=${coord%%:*} a=${coord#*:}
    echo "$MAVEN/${g//.//}/$a/maven-metadata.xml"
  fi
}

# Latest stable version from a maven-metadata.xml. Pre-releases (M, RC, alpha,
# beta, SNAPSHOT, CR, ...) are skipped unless the artifact only ships with a
# known suffix (e.g. OpenTelemetry "-alpha").
latest_from_metadata() {
  local url=$1 allow=$2
  curl -fsSL --max-time 20 "$url" \
    | grep -oE '<version>[^<]+</version>' | sed -E 's#</?version>##g' \
    | if [[ -n $allow ]]; then grep -E "^[0-9.]+-$allow$"; else grep -E '^[0-9]+(\.[0-9]+)*$'; fi \
    | sort -V | tail -1
}

status_of() {
  local cur=$1 latest=$2
  if [[ -z $latest ]]; then echo "LOOKUP_FAILED"; return; fi
  if [[ $cur == "$latest" ]]; then echo "up-to-date"; return; fi
  local newest; newest=$(printf '%s\n%s\n' "$cur" "$latest" | sort -V | tail -1)
  if [[ $newest != "$latest" ]]; then echo "current-is-newer"; return; fi
  if [[ ${cur%%.*} != "${latest%%.*}" ]]; then echo "MAJOR"; else echo "update"; fi
}

printf 'CONSTANT\tCURRENT\tLATEST\tSTATUS\n'
while IFS='|' read -r const repo coord allow; do
  [[ -z $const ]] && continue
  cur=$(current_version "$const")
  latest=$(latest_from_metadata "$(metadata_url "$repo" "$coord")" "$allow")
  printf '%s\t%s\t%s\t%s\n' "$const" "$cur" "$latest" "$(status_of "$cur" "$latest")"
done <<< "$DEPS"

# Node.js: latest LTS within the same major line as the current value.
cur=$(current_version nodeVersion)
latest=$(curl -fsSL --max-time 20 https://nodejs.org/dist/index.json \
  | grep -oE '"version":"v[0-9.]+"[^}]*"lts":"[A-Za-z]+"' \
  | grep -oE 'v[0-9.]+' | sed 's/^v//' | grep -E "^${cur%%.*}\." | sort -V | tail -1)
printf '%s\t%s\t%s\t%s\n' nodeVersion "$cur" "$latest" "$(status_of "$cur" "$latest")"

# npm
cur=$(current_version npmVersion)
latest=$(curl -fsSL --max-time 20 https://registry.npmjs.org/-/package/npm/dist-tags | grep -oE '"latest":"[^"]+"' | cut -d'"' -f4)
printf '%s\t%s\t%s\t%s\n' npmVersion "$cur" "$latest" "$(status_of "$cur" "$latest")"
