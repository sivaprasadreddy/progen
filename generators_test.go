package main

import (
	"fmt"
	minimalgo "github.com/sivaprasadreddy/progen/generators/minimal-go"
	minimaljava "github.com/sivaprasadreddy/progen/generators/minimal-java"
	springboot "github.com/sivaprasadreddy/progen/generators/spring-boot"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestGenerateMinimalJavaMavenApp(t *testing.T) {
	pc := minimaljava.ProjectConfig{
		AppName:     "my-minimal-java-mvn-app",
		GroupID:     "com.sivalabs",
		ArtifactID:  "my-minimal-java-mvn-app",
		AppVersion:  "1.0",
		BasePackage: "com.sivalabs.myapp",
		BuildTool:   "Maven",
	}
	err := minimaljava.GenerateProject(pc)
	assert.Nil(t, err)
	cmd := exec.Command("/bin/sh", "-c", "cd my-minimal-java-mvn-app; ./mvnw test;")
	err = cmd.Run()
	fmt.Println("error:", err)
	assert.Nil(t, err)

	//cleanup
	cmd = exec.Command("/bin/sh", "-c", "rm -rf my-minimal-java-mvn-app;")
	err = cmd.Run()
	assert.Nil(t, err)
}

func TestGenerateMinimalJavaGradleApp(t *testing.T) {
	pc := minimaljava.ProjectConfig{
		AppName:     "my-minimal-java-gradle-app",
		GroupID:     "com.sivalabs",
		ArtifactID:  "my-minimal-java-gradle-app",
		AppVersion:  "1.0",
		BasePackage: "com.sivalabs.myapp",
		BuildTool:   "Gradle",
	}
	err := minimaljava.GenerateProject(pc)
	assert.Nil(t, err)
	cmd := exec.Command("/bin/sh", "-c", "cd my-minimal-java-gradle-app; ./gradlew build;")
	err = cmd.Run()
	fmt.Println("error:", err)
	assert.Nil(t, err)

	//cleanup
	cmd = exec.Command("/bin/sh", "-c", "rm -rf my-minimal-java-gradle-app;")
	err = cmd.Run()
	assert.Nil(t, err)
}

func TestGenerateSpringBootMavenApp(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppName:         "my-spring-boot-mvn-app",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-mvn-app",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Maven",
		DbType:          "postgresql",
		DbMigrationTool: "flyway",
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	cmd := exec.Command("/bin/sh", "-c", "cd my-spring-boot-mvn-app; ./mvnw test;")
	err = cmd.Run()
	fmt.Println("error:", err)
	assert.Nil(t, err)

	//cleanup
	cmd = exec.Command("/bin/sh", "-c", "rm -rf my-spring-boot-mvn-app;")
	err = cmd.Run()
	assert.Nil(t, err)
}

func TestGenerateSpringBootGradleApp(t *testing.T) {
	pc := springboot.ProjectConfig{
		AppName:         "my-spring-boot-gradle-app",
		GroupID:         "com.sivalabs",
		ArtifactID:      "my-spring-boot-gradle-app",
		AppVersion:      "1.0",
		BasePackage:     "com.sivalabs.myapp",
		BuildTool:       "Gradle",
		DbType:          "postgresql",
		DbMigrationTool: "flyway",
	}
	err := springboot.GenerateProject(pc)
	assert.Nil(t, err)
	cmd := exec.Command("/bin/sh", "-c", "cd my-spring-boot-gradle-app; ./gradlew test;")
	err = cmd.Run()
	fmt.Println("error:", err)
	assert.Nil(t, err)

	//cleanup
	cmd = exec.Command("/bin/sh", "-c", "rm -rf my-spring-boot-gradle-app;")
	err = cmd.Run()
	assert.Nil(t, err)
}

func TestGenerateMinimalGoGinApp(t *testing.T) {
	pc := minimalgo.ProjectConfig{
		AppName:        "my-minimal-go-gin-app",
		ModulePath:     "github.com/sivalabs/myginapp",
		RoutingLibrary: "Gin",
	}
	err := minimalgo.GenerateProject(pc)
	assert.Nil(t, err)
	cmd := exec.Command("/bin/sh", "-c", "cd my-minimal-go-gin-app; go build;")
	err = cmd.Run()
	fmt.Println("error:", err)
	assert.Nil(t, err)

	//cleanup
	cmd = exec.Command("/bin/sh", "-c", "rm -rf my-minimal-go-gin-app;")
	err = cmd.Run()
	assert.Nil(t, err)
}

func TestGenerateMinimalGoChiApp(t *testing.T) {
	pc := minimalgo.ProjectConfig{
		AppName:        "my-minimal-go-chi-app",
		ModulePath:     "github.com/sivalabs/mychiapp",
		RoutingLibrary: "Chi",
	}
	err := minimalgo.GenerateProject(pc)
	assert.Nil(t, err)
	cmd := exec.Command("/bin/sh", "-c", "cd my-minimal-go-chi-app; go build;")
	err = cmd.Run()
	fmt.Println("error:", err)
	assert.Nil(t, err)

	//cleanup
	cmd = exec.Command("/bin/sh", "-c", "rm -rf my-minimal-go-chi-app;")
	err = cmd.Run()
	assert.Nil(t, err)
}
