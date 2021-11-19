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

type ConfigT struct {
	Validation ValidationT `toml:"validation"`
	Limits     LimitsT     `toml:"limits"`
	Moderation ModerationT `toml:"moderation"`
}

type ValidationT struct {
	// email required for submitting a comment?
	RequireEmail bool `toml:"require_email"`
	// name required for submitting a comment?
	RequireAuthor bool `toml:"require_author"`
	// minimum length for a comment
	CommentLengthMin int `toml:"comment_length_min"`
	// maximum length for a comment
	CommentLengthMax int `toml:"comment_length_max"`
}

type LimitsT struct {
	// how many levels of nesting to allow for comments, 0 = no nesting allowed, flat
	CommentDepthMax int `toml:"comment_depth_max"`
	// when initially loading a thread, how many levels of nesting can be queried
	InitialQueryDepthMax int `toml:"initial_query_depth_max"`
	// how many comments can be queried at once
	QueryLimitMax int `toml:"query_limit_max"`
}

type ModerationT struct {
	// whether a posted comment is not shown until approved by an administrator
	RequireApproval bool `toml:"require_approval"`
}

func GetDefaultConfig() ConfigT {
	return ConfigT{
		Validation: ValidationT{
			RequireEmail:     false,
			RequireAuthor:    true,
			CommentLengthMin: 3,
			CommentLengthMax: 20000,
		},
		Limits: LimitsT{
			CommentDepthMax:      8,
			InitialQueryDepthMax: 4,
			QueryLimitMax:        200,
		},
		Moderation: ModerationT{
			RequireApproval: true,
		},
	}
}

func ValidateConfig(config ConfigT) error {
	if config.Validation.CommentLengthMin < 1 || config.Validation.CommentLengthMin > config.Validation.CommentLengthMax {
		return fmt.Errorf("the minimum comment length (comment_length_min) must be in the range [1,comment_length_max], was %d", config.Validation.CommentLengthMin)
	}

	if config.Validation.CommentLengthMin > maxCommentLength || config.Validation.CommentLengthMax < config.Validation.CommentLengthMin {
		return fmt.Errorf("the maximum comment length (comment_length_max) must be in the range [comment_length_min,%d], was %d", maxCommentLength, config.Validation.CommentLengthMin)
	}

	if config.Limits.CommentDepthMax < 0 {
		return fmt.Errorf("the maximum comment depth (comment_depth_max) must be >= 0, was %d", config.Limits.CommentDepthMax)
	}

	if config.Limits.CommentDepthMax > warnMaxCommentDepth {
		zap.L().Sugar().Warnf("the maximum comment depth (comment_depth_max) is > %d, this might degrade performance, was %d", warnMaxCommentDepth, config.Limits.CommentDepthMax)
	}

	if config.Limits.InitialQueryDepthMax < 0 || config.Limits.InitialQueryDepthMax > config.Limits.CommentDepthMax {
		return fmt.Errorf("the maximum initial query depth (initial_query_depth_max) must be in the range [0,comment_depth_max], was %d", config.Limits.InitialQueryDepthMax)
	}

	if config.Limits.QueryLimitMax < 1 {
		return fmt.Errorf("the maximum query limit (query_limit_max) must be > 0, was %d", config.Limits.QueryLimitMax)
	}

	if config.Limits.QueryLimitMax > warnMaxQueryLimit {
		zap.L().Sugar().Warnf("the maximum query limit (query_limit_max) is > %d, this might degrade performance, was %d", warnMaxQueryLimit, config.Limits.QueryLimitMax)
	}

	return nil
}

func ReadConfig() (ConfigT, error) {
	config := GetDefaultConfig()

	if _, err := os.Stat(configFilename); !errors.Is(err, os.ErrNotExist) {
		// file exists
		_, err := toml.DecodeFile(configFilename, &config)
		if err != nil {
			return ConfigT{}, fmt.Errorf("unable to read or parse gomment configuration file: %w", err)
		}
	} else {
		zap.L().Sugar().Infof("configuration file '%s' not found, using default configuration", configFilename)
	}

	err := ValidateConfig(config)
	if err != nil {
		return ConfigT{}, fmt.Errorf("configuration validation failed: %w", err)
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	err = encoder.Encode(config)
	if err != nil {
		return ConfigT{}, fmt.Errorf("failed to serialize validated config: %w", err)
	}

	fmt.Println(string(buf.Bytes()))

	return config, nil
}
