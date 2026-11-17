package springboot

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAxelixGeneration(t *testing.T) {
	for _, buildTool := range []BuildTool{Maven, Gradle} {
		for _, frontend := range []Frontend{FrontendNone, FrontendThymeleaf, FrontendJTE, FrontendAngular, FrontendReact} {
			for _, enabled := range []bool{false, true} {
				name := buildTool.String() + "/" + frontend.String() + "/disabled"
				if enabled {
					name = buildTool.String() + "/" + frontend.String() + "/enabled"
				}
				t.Run(name, func(t *testing.T) {
					t.Chdir(t.TempDir())
					pc := defaultProjectConfig
					pc.Versions = defaultVersions()
					pc.BuildTool = buildTool
					pc.Frontend = frontend
					pc.SetFeatures(nil)
					if enabled {
						pc.SetFeatures([]string{FeatureAxelixSupport})
					}
					require.Equal(t, enabled, pc.AxelixSupport)
					pg := projectGenerator{tmplFS: tmplsFS}
					for _, gen := range []configGenerator{NewBuildToolConfig(pg), NewSecurityConfig(pg), NewReadMeConfig(pg)} {
						require.NoError(t, gen.generate(pc))
					}
					buildFile := "pom.xml"
					if buildTool == Gradle {
						buildFile = "build.gradle"
					}
					data, err := os.ReadFile(filepath.Join(pc.AppName, buildFile))
					require.NoError(t, err)
					if buildTool == Maven {
						var doc any
						require.NoError(t, xml.Unmarshal(data, &doc))
					}
					security, err := os.ReadFile(filepath.Join(pc.AppName, "src/main/java/com/mycompany/myapp/config/SecurityConfig.java"))
					require.NoError(t, err)
					readme, err := os.ReadFile(filepath.Join(pc.AppName, "README.md"))
					require.NoError(t, err)
					if enabled {
						require.Contains(t, string(data), "axelix-spring-boot-4-starter")
						require.Contains(t, string(data), axelixVersion)
						if buildTool == Maven {
							require.Contains(t, string(data), "axelix-maven-plugin")
							require.Contains(t, string(data), "<goal>axelix-generate-project-info</goal>")
						} else {
							require.Contains(t, string(data), `id "com.axelixlabs.axelix"`)
						}
						require.Contains(t, string(security), `http.securityMatcher("/actuator/axelix-*/**")`)
						require.Contains(t, string(security), "@Order(0)")
						if pc.RestApiEnabled() {
							require.Contains(t, string(security), "registration.setEnabled(false)")
						} else {
							require.NotContains(t, string(security), "FilterRegistrationBean")
						}
						require.Contains(t, string(security), ".csrf(csrf -> csrf.disable())")
						require.Contains(t, string(readme), "## Axelix")
					} else {
						require.NotContains(t, string(data), "axelix")
						require.NotContains(t, string(security), "axelix")
						require.NotContains(t, string(readme), "Axelix")
					}
					configPath := filepath.Join(pc.AppName, ProjectConfigFile)
					require.NoError(t, WriteConfigFile(pc, configPath))
					loaded, _, err := LoadConfigFile(configPath)
					require.NoError(t, err)
					require.Equal(t, enabled, loaded.AxelixSupport)
				})
			}
		}
	}
}
