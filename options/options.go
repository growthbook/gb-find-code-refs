package options

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/growthbook/gb-find-code-refs/internal/validation"
)

const (
	maxProjKeyLength = 20 // Maximum project key length
)

// TODO audit
type Options struct {
	Branch              string `mapstructure:"branch"`
	Dir                 string `mapstructure:"dir" yaml:"-"`
	OutDir              string `mapstructure:"outDir"`
	Revision            string `mapstructure:"revision"`
	FlagsPath           string `mapstructure:"flagsPath"`
	ContextLines        int    `mapstructure:"contextLines"`
	Lookback            int    `mapstructure:"lookback"`
	AllowTags           bool   `mapstructure:"allowTags"`
	Debug               bool   `mapstructure:"debug"`
	IgnoreServiceErrors bool   `mapstructure:"ignoreServiceErrors"`

	// The following options can only be configured via YAML configuration

	Aliases    []Alias    `mapstructure:"aliases"`
	Delimiters Delimiters `mapstructure:"delimiters"`
}

type Delimiters struct {
	// If set to `true`, the default delimiters (single-quote, double-qoute, and backtick) will not be used unless provided as `additional` delimiters
	DisableDefaults bool     `mapstructure:"disableDefaults"`
	Additional      []string `mapstructure:"additional"`
}

func Init(flagSet *pflag.FlagSet) error {
	for _, f := range flags {
		usage := strings.ReplaceAll(f.usage, "\n", " ")
		switch value := f.defaultValue.(type) {
		case string:
			flagSet.StringP(f.name, f.short, value, usage)
		case int:
			flagSet.IntP(f.name, f.short, value, usage)
		case bool:
			flagSet.BoolP(f.name, f.short, value, usage)
		}
	}

	flagSet.VisitAll(func(f *pflag.Flag) {
		viper.BindEnv(f.Name, "GB_"+strcase.ToScreamingSnake(f.Name))
	})

	return viper.BindPFlags(flagSet)
}

func InitYAML() error {
	err := validateYAMLPreconditions()
	if err != nil {
		return err
	}
	absPath, err := validation.NormalizeAndValidatePath(viper.GetString("dir"))
	if err != nil {
		return err
	}
	viper.SetConfigName("coderefs")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(absPath, ".growthbook"))
	err = viper.ReadInConfig()
	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return err
	}
	return nil
}

// validatePreconditions ensures required flags have been set
func validateYAMLPreconditions() error {
	dir := viper.GetString("dir")
	flagsPath := viper.GetString("flagsPath")
	missingRequiredOptions := []string{}
	if dir == "" {
		missingRequiredOptions = append(missingRequiredOptions, "dir")
	}
	if flagsPath == "" {
		missingRequiredOptions = append(missingRequiredOptions, "flagsPath")
	}
	if len(missingRequiredOptions) > 0 {
		return fmt.Errorf("missing required option(s): %v", missingRequiredOptions)
	}
	return nil
}

func GetOptions() (Options, error) {
	var opts Options
	err := viper.Unmarshal(&opts)
	return opts, err
}

func GetWrapperOptions(dir string, merge func(Options) (Options, error)) (Options, error) {
	flags := pflag.CommandLine

	err := Init(flags)
	if err != nil {
		return Options{}, err
	}

	// Set precondition flags
	err = flags.Set("dir", dir)
	if err != nil {
		return Options{}, err
	}

	err = InitYAML()
	if err != nil {
		return Options{}, err
	}

	opts, err := GetOptions()
	if err != nil {
		return opts, err
	}

	return merge(opts)
}

func (o Options) ValidateRequired() error {
	missingRequiredOptions := []string{}
	if o.Dir == "" {
		missingRequiredOptions = append(missingRequiredOptions, "dir")
	}
	if len(missingRequiredOptions) > 0 {
		return fmt.Errorf("missing required option(s): %v", missingRequiredOptions)
	}

	return nil
}

// Validate ensures all options have been set to a valid value
func (o Options) Validate() error {
	if err := o.ValidateRequired(); err != nil {
		return err
	}

	maxContextLines := 5
	if o.ContextLines > maxContextLines {
		return fmt.Errorf(`invalid value %q for "contextLines": must be <= %d`, o.ContextLines, maxContextLines)
	}

	// match all non-control ASCII characters
	validDelims := regexp.MustCompile("^[\x20-\x7E]$")
	for i, d := range o.Delimiters.Additional {
		if !validDelims.MatchString(d) {
			return fmt.Errorf(`invalid value %q for "delimiters.additional[%d]": each delimiter must be a valid non-control ASCII character`, d, i)
		}
	}

	if _, err := validation.NormalizeAndValidatePath(o.Dir); err != nil {
		return fmt.Errorf(`invalid value for "dir": %+v`, err)
	}

	if o.OutDir != "" {
		if _, err := validation.NormalizeAndValidatePath(o.OutDir); err != nil {
			return fmt.Errorf(`invalid valid for "outDir": %+v`, err)
		}
	}

	for _, a := range o.Aliases {
		if err := a.IsValid(); err != nil {
			return err
		}
	}

	if o.Revision != "" && o.Branch == "" {
		return fmt.Errorf(`"branch" option is required when "revision" option is set`)
	}

	return nil
}

func projKeyValidation(projKey string) error {
	if strings.HasPrefix(projKey, "sdk-") {
		return fmt.Errorf("provided project key (%s) appears to be a LaunchDarkly SDK key", "sdk-xxxx")
	} else if strings.HasPrefix(projKey, "api-") {
		return fmt.Errorf("provided project key (%s) appears to be a LaunchDarkly API access token", "api-xxxx")
	}

	return nil
}

func (o Options) GetProjectKeys() (projects []string) {
	return projects
}
