package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfigByHostname(t *testing.T) {
	testConfig := Config{
		Credentials: map[string]GitlabConfig{"git.example.local": {
			BaseUrl: "https://git.example.local/api/v4",
			Token:   "glpat-testabcdef",
		}},
	}
	t.Run("should return config if hostname exists", func(t *testing.T) {
		gitlabConfig, err := GetConfigByHostname("git.example.local", testConfig)
		assert.Nil(t, err)
		assert.Equal(t, gitlabConfig.BaseUrl, "https://git.example.local/api/v4")
		assert.Equal(t, gitlabConfig.Token, "glpat-testabcdef")
	})
	t.Run("should return error if hostname doesn't exist", func(t *testing.T) {
		gitlabConfig, err := GetConfigByHostname("not-exist.local", testConfig)
		assert.Nil(t, gitlabConfig)
		assert.NotNil(t, err)
	})
}
