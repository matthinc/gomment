package config_test

import (
	"testing"

	"github.com/matthinc/gomment/config"
	"github.com/stretchr/testify/assert"
)

func TestValidateDefault(t *testing.T) {
	err := config.ValidateConfig(config.GetDefaultConfig())
	assert.NoError(t, err, "expected default configuration to be validated successfully")
}

func TestValidateCommentLength(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Validation.CommentLengthMin = 10
	cfg.Validation.CommentLengthMax = 100

	err := config.ValidateConfig(cfg)
	assert.NoError(t, err, "expected validation of comment length [10,100] to be successful")
}

func TestValidateCommentLengthMinNegative(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Validation.CommentLengthMin = -1

	err := config.ValidateConfig(cfg)
	assert.Error(t, err, "expected validation to fail for negative minimum comment length")
}

func TestValidateCommentLengthMaxNegative(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Validation.CommentLengthMax = -1

	err := config.ValidateConfig(cfg)
	assert.Error(t, err, "expected validation to fail for negative maximum comment length")
}

func TestValidateCommentLengthSwapped(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Validation.CommentLengthMin = 100
	cfg.Validation.CommentLengthMax = 10

	err := config.ValidateConfig(cfg)
	assert.Error(t, err, "expected validation of comment length [100,10] to fail")
}

func TestValidateMaxCommentDepth(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Limits.MaxCommentDepth = 10

	err := config.ValidateConfig(cfg)
	assert.NoError(t, err, "expected validation of comment depth 10 to be successful")
}

func TestValidateMaxCommentDepthNegative(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Limits.MaxCommentDepth = -1

	err := config.ValidateConfig(cfg)
	assert.Error(t, err, "expected validation to fail for negative maximum comment depth")
}

func TestValidateMaxInitialQueryDepthHigher(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Limits.MaxCommentDepth = 9
	cfg.Limits.MaxInitialQueryDepth = 10

	err := config.ValidateConfig(cfg)
	assert.Error(t, err, "expected validation to fail for maximum comment depth being smaller that maximum initial query depth")
}

func TestValidateMaxQueryLimitZero(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Limits.MaxQueryLimit = 0

	err := config.ValidateConfig(cfg)
	assert.Error(t, err, "expected validation to fail for negative maximum comment limit")
}
