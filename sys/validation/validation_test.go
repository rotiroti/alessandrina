package validation_test

import (
	"testing"

	"github.com/rotiroti/alessandrina/sys/validation"
)

func TestCheck(t *testing.T) {
	type dummyBook struct {
		Title string `validate:"required"`
		ISBN  string `validate:"required,isbn"`
		Pages int    `validate:"required,min=1"`
	}

	tests := []struct {
		name    string
		val     any
		wantErr bool
	}{
		{
			name: "valid",
			val: dummyBook{
				Title: "The Hobbit",
				ISBN:  "978-0547928227",
				Pages: 310,
			},
			wantErr: false,
		},
		{
			name: "invalid",
			val: dummyBook{
				Title: "",
				ISBN:  "isbn",
				Pages: 0,
			},
			wantErr: true,
		},
		{
			name:    "bad input",
			val:     "bad input",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validation.New()
			err := v.Check(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFieldErrors_Error(t *testing.T) {
	tests := []struct {
		name string
		fe   validation.FieldErrors
		want string
	}{
		{
			name: "empty",
			fe:   validation.FieldErrors{},
			want: "[]",
		},
		{
			name: "one",
			fe: validation.FieldErrors{
				validation.FieldError{
					Field: "Title",
					Err:   "required",
				},
			},
			want: `[{"field":"Title","error":"required"}]`,
		},
		{
			name: "two",
			fe: validation.FieldErrors{
				validation.FieldError{
					Field: "Title",
					Err:   "required",
				},
				validation.FieldError{
					Field: "ISBN",
					Err:   "required",
				},
			},
			want: `[{"field":"Title","error":"required"},{"field":"ISBN","error":"required"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fe.Error(); got != tt.want {
				t.Errorf("FieldErrors.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
