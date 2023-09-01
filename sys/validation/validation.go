package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Validator is an interface for validating request struct values.
type Validator interface {
	Check(v any) error
}

// Config is the validation configuration.
type Config struct {
	validate   *validator.Validate
	translator ut.Translator
}

// Ensure that Config implements the Checker interface.
var _ Validator = (*Config)(nil)

// New returns a new validation Config.
func New() *Config {
	validate := validator.New()

	// Create a translator for english so the error messages are
	// more human-readable than technical.
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the english error messages for use.
	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(err)
	}

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Config{
		validate:   validate,
		translator: translator,
	}
}

// Check validates the request struct value.
func (c *Config) Check(val any) error {
	if err := c.validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return fmt.Errorf("validation.check: %w", err)
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Err:   verror.Translate(c.translator),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}
