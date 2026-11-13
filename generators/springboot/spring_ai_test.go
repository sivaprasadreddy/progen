package springboot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeAIModelsDropsInvalidAndDuplicatesInCanonicalOrder(t *testing.T) {
	got := sanitizeAIModels([]AIModelProvider{"anthropic", "bogus", "ollama", "anthropic"})
	require.Equal(t, []AIModelProvider{AIModelOllama, AIModelAnthropic}, got)
}

func TestDefaultAIChatModelIsFirstSelectedInCanonicalOrder(t *testing.T) {
	pc := ProjectConfig{SpringAIModels: []AIModelProvider{AIModelAnthropic, AIModelOpenAI}}
	require.True(t, pc.SpringAISupport())
	require.False(t, pc.OllamaSupport())
	require.Equal(t, AIModelOpenAI, pc.DefaultAIChatModel())

	require.False(t, ProjectConfig{}.SpringAISupport())
	require.Equal(t, AIModelProvider(""), ProjectConfig{}.DefaultAIChatModel())
}
