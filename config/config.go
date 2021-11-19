package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"go.uber.org/zap"
)

const configFilename = "gomment.toml"
const maxCommentLength = 100000
const warnMaxCommentDepth = 20
const warnMaxQueryLimit = 500

type configT struct {
	Validation validationT `toml:"validation"`
	Limits     limitsT     `toml:"limits"`
	Moderation moderationT `toml:"moderation"`
}

type validationT struct {
	// email required for submitting a comment?
	RequireEmail bool `toml:"require_email"`
	// name required for submitting a comment?
	RequireName bool `toml:"require_name"`
	// minimum length for a comment
	CommentLengthMin int `toml:"comment_length_min"`
	// maximum length for a comment
	CommentLengthMax int `toml:"comment_length_max"`
}

type limitsT struct {
	// how many levels of nesting to allow for comments, 0 = no nesting allowed, flat
	MaxCommentDepth int `toml:"max_comment_depth"`
	// when initially loading a thread, how many levels of nesting can be queried
	MaxInitialQueryDepth int `toml:"max_initial_query_depth"`
	// how many comments can be queried at once
	MaxQueryLimit int `toml:"max_query_limit"`
}

type moderationT struct {
	// whether a posted comment is not shown until approved by an administrator
	RequireApproval bool `toml:"require_approval"`
}

func GetDefaultConfig() configT {
	return configT{
		Validation: validationT{
			RequireEmail:     false,
			RequireName:      true,
			CommentLengthMin: 3,
			CommentLengthMax: 20000,
		},
		Limits: limitsT{
			MaxCommentDepth:      8,
			MaxInitialQueryDepth: 4,
			MaxQueryLimit:        200,
		},
		Moderation: moderationT{
			RequireApproval: true,
		},
	}
}

func ValidateConfig(config configT) error {
	if config.Validation.CommentLengthMin < 1 || config.Validation.CommentLengthMin > config.Validation.CommentLengthMax {
		return fmt.Errorf("the minimum comment length (comment_length_min) must be in the range [1,comment_length_max], was %d", config.Validation.CommentLengthMin)
	}

	if config.Validation.CommentLengthMin > maxCommentLength || config.Validation.CommentLengthMax < config.Validation.CommentLengthMin {
		return fmt.Errorf("the maximum comment length (comment_length_max) must be in the range [comment_length_min,%d], was %d", maxCommentLength, config.Validation.CommentLengthMin)
	}

	if config.Limits.MaxCommentDepth < 0 {
		return fmt.Errorf("the maximum comment depth (max_comment_depth) must be >= 0, was %d", config.Limits.MaxCommentDepth)
	}

	if config.Limits.MaxCommentDepth > warnMaxCommentDepth {
		zap.L().Sugar().Warnf("the maximum comment depth (max_comment_depth) is > %d, this might degrade performance, was %d", warnMaxCommentDepth, config.Limits.MaxCommentDepth)
	}

	if config.Limits.MaxInitialQueryDepth < 0 || config.Limits.MaxInitialQueryDepth > config.Limits.MaxCommentDepth {
		return fmt.Errorf("the maximum initial query depth (max_initial_query_depth) must be in the range [0,max_comment_depth], was %d", config.Limits.MaxInitialQueryDepth)
	}

	if config.Limits.MaxQueryLimit < 1 {
		return fmt.Errorf("the maximum query limit (max_query_limit) must be > 0, was %d", config.Limits.MaxQueryLimit)
	}

	if config.Limits.MaxQueryLimit > warnMaxQueryLimit {
		zap.L().Sugar().Warnf("the maximum query limit (max_query_limit) is > %d, this might degrade performance, was %d", warnMaxQueryLimit, config.Limits.MaxQueryLimit)
	}

	return nil
}

func ReadConfig() (configT, error) {
	config := GetDefaultConfig()

	if _, err := os.Stat(configFilename); !errors.Is(err, os.ErrNotExist) {
		// file exists
		_, err := toml.DecodeFile(configFilename, &config)
		if err != nil {
			return configT{}, fmt.Errorf("unable to read or parse gomment configuration file: %w", err)
		}
	} else {
		zap.L().Sugar().Infof("configuration file '%s' not found, using default configuration", configFilename)
	}

	err := ValidateConfig(config)
	if err != nil {
		return configT{}, fmt.Errorf("configuration validation failed: %w", err)
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	err = encoder.Encode(config)
	if err != nil {
		return configT{}, fmt.Errorf("failed to serialize validated config: %w", err)
	}

	fmt.Println(string(buf.Bytes()))

	return config, nil
}
