package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ConfigIn struct {
	Focused struct {
		BorderColor   *uint   `json:"border_color"`
		TextColor     *uint   `json:"text_color"`
		BoxTitleColor *uint   `json:"box_title_color"`
		BorderType    *string `json:"border_type"`
	} `json:"focused"`

	Unfocused struct {
		BorderColor   *uint   `json:"border_color"`
		TextColor     *uint   `json:"text_color"`
		BoxTitleColor *uint   `json:"box_title_color"`
		BorderType    *string `json:"border_type"`
	} `json:"unfocused"`

	ErrorTextColor *uint `json:"error_text_color"`
}

type Config struct {
	Focused        FocusedOpts
	Unfocused      UnfocusedOpts
	ErrorTextColor lipgloss.ANSIColor
}

type FocusedOpts struct {
	BorderColor   lipgloss.ANSIColor
	TextColor     lipgloss.ANSIColor
	BoxTitleColor lipgloss.ANSIColor
	BorderType    lipgloss.Border
}

type UnfocusedOpts struct {
	BorderColor   lipgloss.ANSIColor
	TextColor     lipgloss.ANSIColor
	BoxTitleColor lipgloss.ANSIColor
	BorderType    lipgloss.Border
}

var (
	defaultFocusedColor   = lipgloss.ANSIColor(4) // blue
	defaultUnfocusedColor = lipgloss.ANSIColor(7) // white
	defaultErrorColor     = lipgloss.ANSIColor(1) // red

	defaultConfig = Config{
		Focused: FocusedOpts{
			BorderColor:   defaultFocusedColor,
			TextColor:     defaultFocusedColor,
			BoxTitleColor: defaultFocusedColor,
			BorderType:    lipgloss.DoubleBorder(),
		},
		Unfocused: UnfocusedOpts{
			BorderColor:   defaultUnfocusedColor,
			TextColor:     defaultUnfocusedColor,
			BoxTitleColor: defaultUnfocusedColor,
			BorderType:    lipgloss.NormalBorder(),
		},
		ErrorTextColor: defaultErrorColor,
	}
)

const (
	cfgDirName  = "yata"
	cfgFileName = "config.json"
)

func GetConfig() (*Config, error) {
	uc, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("getting user config directory path: %w", err)
	}
	path := filepath.Join(uc, cfgDirName, cfgFileName)
	cfg, err := readConfig(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	return configInToConfig(cfg), nil
}

func configInToConfig(in *ConfigIn) *Config {
	dc := defaultConfig
	if in == nil {
		return &dc
	}
	return &Config{
		Focused: FocusedOpts{
			BorderColor:   uintPtrToColor(in.Focused.BorderColor, dc.Focused.BorderColor),
			TextColor:     uintPtrToColor(in.Focused.TextColor, dc.Focused.TextColor),
			BoxTitleColor: uintPtrToColor(in.Focused.BoxTitleColor, dc.Focused.BoxTitleColor),
			BorderType:    strPtrToBorder(in.Focused.BorderType, dc.Focused.BorderType),
		},
		Unfocused: UnfocusedOpts{
			BorderColor:   uintPtrToColor(in.Unfocused.BorderColor, dc.Unfocused.BorderColor),
			TextColor:     uintPtrToColor(in.Unfocused.TextColor, dc.Unfocused.TextColor),
			BoxTitleColor: uintPtrToColor(in.Unfocused.BoxTitleColor, dc.Unfocused.BoxTitleColor),
			BorderType:    strPtrToBorder(in.Unfocused.BorderType, dc.Unfocused.BorderType),
		},
		ErrorTextColor: uintPtrToColor(in.ErrorTextColor, dc.ErrorTextColor),
	}
}

func readConfig(path string) (*ConfigIn, error) {
	cfg := &ConfigIn{}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("checking for config file: %w", err)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := json.Unmarshal(file, cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling json: %w", err)
	}

	return cfg, nil
}

func uintPtrToColor(i *uint, defColor lipgloss.ANSIColor) lipgloss.ANSIColor {
	if i != nil && *i <= 255 {
		return lipgloss.ANSIColor(*i)
	}

	return defColor
}

func strPtrToBorder(s *string, defBorder lipgloss.Border) lipgloss.Border {
	if s == nil {
		return defBorder
	}

	switch strings.ToLower(*s) {
	case "normal":
		return lipgloss.NormalBorder()
	case "double":
		return lipgloss.DoubleBorder()
	case "rounded", "round":
		return lipgloss.RoundedBorder()
	case "thick", "thicc":
		return lipgloss.ThickBorder()
	default:
		return defBorder
	}
}
