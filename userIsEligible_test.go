package main

import (
	"errors"
	"testing"
)

func TestUserIsEligible(t *testing.T) {
	var tests = []struct {
		email       string
		password    string
		age         int
		expectedErr error
	}{
		{
			email:       "test@example.com",
			password:    "12345",
			age:         18,
			expectedErr: nil,
		},
		{
			email:       "",
			password:    "12345",
			age:         18,
			expectedErr: errors.New("email can't be empty"),
		},
		{
			email:       "test@example.com",
			password:    "",
			age:         18,
			expectedErr: errors.New("password can't be empty"),
		},
		{
			email:       "test@example.com",
			password:    "12345",
			age:         16,
			expectedErr: errors.New("age 16 is less than 18, must be at least 18"),
		},
	}

	// loop over tests and try them out
	for _, test := range tests {
		err := userIsEligible(test.email, test.password, test.age)
		errGottenString := ""
		errExpectedString := ""

		if err != nil {
			errGottenString = err.Error()
		}

		if test.expectedErr != nil {
			errExpectedString = test.expectedErr.Error()
		}

		if errGottenString != errExpectedString {
			t.Errorf("userIsEligible(%s, %s, %d) = %s, expected %s", test.email, test.password, test.age, errGottenString, errExpectedString)
		}
	}

}
