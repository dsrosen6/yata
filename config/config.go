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

type (
	ConfigIn struct {
		General struct {
			SelectedProject string  `json:"selected_project"`
			ShowHelp        HelpOpt `json:"show_help"`
		} `json:"general"`
		Style struct {
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
			ErrorText *uint `json:"error_text_color"`
		} `json:"style"`
	}

	Config struct {
		General GeneralOpts
		Style   StyleOpts
	}

	GeneralOpts struct {
		// SelectedProject is the project that will be selected when the app is opened. If not specified,
		// it will default to the most recently selected project. If specified and a project by that title
		// doesn't exist, it will be created.
		//
		// Special options which won't create or select a specific project:
		// 	"all": all filter at the top of the project list will be selected.
		// 	"most_recent": default (works, but not necessary)
		SelectedProject string

		// ShowHelp determines if the help keys are shown at the bottom of the app when opened.
		// Default is to use the most recent.
		//
		// Valid options:
		// 	"enable"
		// 	"disable"
		// 	"most_recent": default (works, but not necessary)
		ShowHelp HelpOpt
	}

	StyleOpts struct {
		Focused        FocusedOpts
		Unfocused      UnfocusedOpts
		ErrorTextColor lipgloss.ANSIColor
	}

	FocusedOpts struct {
		BorderColor   lipgloss.ANSIColor
		TextColor     lipgloss.ANSIColor
		BoxTitleColor lipgloss.ANSIColor
		BorderType    lipgloss.Border
	}

	UnfocusedOpts struct {
		BorderColor   lipgloss.ANSIColor
		TextColor     lipgloss.ANSIColor
		BoxTitleColor lipgloss.ANSIColor
		BorderType    lipgloss.Border
	}

	HelpOpt string
)

const (
	HelpOptMostRecent HelpOpt = "most_recent"
	HelpOptEnable     HelpOpt = "enable"
	HelpOptDisable    HelpOpt = "disable"
)

var (
	defaultFocusedColor    = lipgloss.ANSIColor(7) // white
	defaultUnfocusedColor  = lipgloss.ANSIColor(8) // gray
	defaultErrorColor      = lipgloss.ANSIColor(1) // red
	defaultHelpOpt         = HelpOptMostRecent
	defaultSelectedProject = "most_recent"

	defaultConfig = Config{
		General: GeneralOpts{
			SelectedProject: defaultSelectedProject,
			ShowHelp:        defaultHelpOpt,
		},
		Style: StyleOpts{
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
		},
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
		General: GeneralOpts{
			SelectedProject: selectedProject(in.General.SelectedProject, dc.General.SelectedProject),
			ShowHelp:        helpOpt(in.General.ShowHelp, dc.General.ShowHelp),
		},
		Style: StyleOpts{
			Focused: FocusedOpts{
				BorderColor:   uintPtrToColor(in.Style.Focused.BorderColor, dc.Style.Focused.BorderColor),
				TextColor:     uintPtrToColor(in.Style.Focused.TextColor, dc.Style.Focused.TextColor),
				BoxTitleColor: uintPtrToColor(in.Style.Focused.BoxTitleColor, dc.Style.Focused.BoxTitleColor),
				BorderType:    strPtrToBorder(in.Style.Focused.BorderType, dc.Style.Focused.BorderType),
			},
			Unfocused: UnfocusedOpts{
				BorderColor:   uintPtrToColor(in.Style.Unfocused.BorderColor, dc.Style.Unfocused.BorderColor),
				TextColor:     uintPtrToColor(in.Style.Unfocused.TextColor, dc.Style.Unfocused.TextColor),
				BoxTitleColor: uintPtrToColor(in.Style.Unfocused.BoxTitleColor, dc.Style.Unfocused.BoxTitleColor),
				BorderType:    strPtrToBorder(in.Style.Unfocused.BorderType, dc.Style.Unfocused.BorderType),
			},
			ErrorTextColor: uintPtrToColor(in.Style.ErrorText, dc.Style.ErrorTextColor),
		},
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

func selectedProject(in, def string) string {
	switch in {
	case "":
		return def
	default:
		return in
	}
}

func helpOpt(in, def HelpOpt) HelpOpt {
	switch in {
	case HelpOptMostRecent, HelpOptEnable, HelpOptDisable:
		return in
	default:
		return def
	}
}
