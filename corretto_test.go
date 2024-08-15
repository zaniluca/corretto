package corretto

import (
	"fmt"
	"testing"
)

func TestValidationOpts(t *testing.T) {
	t.Run("custom message for validation", func(t *testing.T) {
		tests := []struct {
			name          string
			customMessage string
			expectedError string
		}{
			{
				name:          "default message",
				customMessage: "",
				expectedError: "Field1 must be at least 10",
			},
			{
				name:          "without placeholders",
				customMessage: "Field1 is supposed to be a minimum of 10",
				expectedError: "Field1 is supposed to be a minimum of 10",
			},
			{
				name:          "with placeholders",
				customMessage: "%v is supposed to be a minimum of %v",
				expectedError: "Field1 is supposed to be a minimum of 10",
			},
			{
				name:          "with extra placeholders",
				customMessage: "%v is supposed to be a minimum of %v and %v",
				expectedError: "Field1 is supposed to be a minimum of 10 and %!v(MISSING)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				schema := Schema{
					"Field1": Field().Number().Min(10, tt.customMessage),
				}

				err := schema.Parse(&struct{ Field1 int }{Field1: 5})
				if err != nil && err.Error() != tt.expectedError {
					t.Errorf("Parse() should have returned custom message")
					t.Errorf("expected: %s, got: %s", tt.expectedError, err.Error())
				}
			})
		}
	})
}

func TestKitchenSink(t *testing.T) {
	type user struct {
		FirstName string
		LastName  string
		Status    string
		Age       int
		BirthDate string
		Email     string
		Children  []user
		Hobbies   []string
	}

	uniqueHobbies := func(ctx Context) error {
		hobbies := ctx.(user).Hobbies
		unique := make(map[string]struct{})

		for i := 0; i < len(hobbies); i++ {
			hobby := hobbies[i]
			if _, ok := unique[hobby]; ok {
				return fmt.Errorf("Hobbies must be unique")
			}
			unique[hobby] = struct{}{}
		}
		return nil
	}

	userSchema := Schema{
		"FirstName": Field("Name").String().MinLength(3).MaxLength(20),
		"LastName":  Field().String().MinLength(3).MaxLength(20),
		"Status":    Field().String().OneOf([]string{"active", "inactive"}),
		"Age":       Field().Number().NonNegative(),
		"BirthDate": Field().String().NonEmpty(),
		"Email":     Field().String().Email(),
		"Hobbies":   Field().Array().Of(Field().String().NonEmpty()).MaxLength(5).Test(uniqueHobbies),
	}

	hasParentLastName := func(ctx Context) error {
		parent, ok := ctx.(user)
		if !ok {
			return fmt.Errorf("expected context to be of type user, got %T", ctx)
		}

		for i := 0; i < len(parent.Children); i++ {
			child := parent.Children[i]
			if child.LastName != parent.LastName {
				return fmt.Errorf("Children must have the same last name as the parent")
			}
		}

		return nil
	}

	childSchema := Schema{
		"FirstName": Field("Name").String().MinLength(3).MaxLength(20),
		"LastName":  Field().String().MinLength(3).MaxLength(20),
		"Age":       Field().Number().NonNegative(),
		"BirthDate": Field().String().NonEmpty(),
	}
	userSchema.Concat(Schema{
		"Children": Field().Array().Of(Field().Schema(childSchema)).Test(hasParentLastName),
	})

	tests := []struct {
		name        string
		user        user
		expectError bool
		expectedMsg string
	}{
		{
			name: "valid user",
			user: user{
				FirstName: "John",
				LastName:  "Doe",
				Status:    "active",
				Age:       30,
				BirthDate: "1990-01-01",
				Email:     "jhon@doe.com",
				Hobbies:   []string{"reading", "running"},
				Children: []user{
					{
						FirstName: "Jane",
						LastName:  "Doe",
						Age:       10,
						BirthDate: "2010-01-01",
					},
				},
			},
			expectError: false,
		},
		{
			name: "duplicate hobbies",
			user: user{
				FirstName: "John",
				LastName:  "Doe",
				Status:    "active",
				Age:       30,
				BirthDate: "1990-01-01",
				Email:     "jhon@doe.com",
				Hobbies:   []string{"reading", "reading"},
				Children: []user{
					{
						FirstName: "Jane",
						LastName:  "Doe",
						Age:       10,
						BirthDate: "2010-01-01",
					},
				},
			},
			expectError: true,
			expectedMsg: "Hobbies must be unique",
		},
		{
			name: "children with different last name",
			user: user{
				FirstName: "John",
				LastName:  "Doe",
				Status:    "active",
				Age:       30,
				BirthDate: "1990-01-01",
				Email:     "jhon@doe.com",
				Hobbies:   []string{"reading", "running"},
				Children: []user{
					{
						FirstName: "James",
						LastName:  "Smith",
						Age:       12,
						BirthDate: "2012-01-01",
					},
					{
						FirstName: "Jane",
						LastName:  "Doe",
						Age:       10,
						BirthDate: "2010-01-01",
					},
				},
			},
			expectError: true,
			expectedMsg: "Children must have the same last name as the parent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userSchema.Parse(tt.user)
			if tt.expectError && err == nil {
				t.Errorf("Parse() should have returned an error")
			}

			if tt.expectError && err != nil && err.Error() != tt.expectedMsg {
				t.Errorf("Parse() should have returned an error with the expected message")
				t.Errorf("expected: %s, got: %s", tt.expectedMsg, err.Error())
			}
		})
	}
}
