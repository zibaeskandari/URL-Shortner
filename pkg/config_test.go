package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_DefaultsValuesShouldLoadAndYamlShouldOverrideWhenAvailable(t *testing.T) {
	config, err := LoadConfig("../configs/test.yaml")
	assert.Nil(t, err)
	assert.NotEmpty(t, config)

	assert.Equal(t, "abcdefg", config.ShortCode.Chars)
	assert.Equal(t, "1234", config.ShortCode.Digits)
	assert.Equal(t, 5, config.ShortCode.MinLength)
	assert.Equal(t, 12, config.ShortCode.MaxLength)
}

func TestLoadConfig_OSEnvVariableShouldOverrideYamlWhenAvailable(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("APP_SHORT_CODE_MIN_LENGTH", "123")
	_ = os.Setenv("APP_SHORT_CODE_MAX_LENGTH", "456")
	_ = os.Setenv("APP_SHORT_CODE_CHARS", "ABCD")
	_ = os.Unsetenv("APP_SHORT_CODE_DIGITS")
	defer func() {
		_ = os.Unsetenv("APP_SHORT_CODE_MIN_LENGTH")
		_ = os.Unsetenv("APP_SHORT_CODE_MAX_LENGTH")
		_ = os.Unsetenv("APP_SHORT_CODE_CHARS")
	}()

	config, err := LoadConfig("../configs/test.yaml")
	assert.Nil(t, err)
	assert.NotEmpty(t, config)

	assert.Equal(t, "ABCD", config.ShortCode.Chars)
	assert.Equal(t, "1234", config.ShortCode.Digits)
	assert.Equal(t, 123, config.ShortCode.MinLength)
	assert.Equal(t, 456, config.ShortCode.MaxLength)
}
