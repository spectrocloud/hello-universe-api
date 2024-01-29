// Copyright (c) HashiCorp, Inc.

package internal

import "testing"

func TestValidateToken(t *testing.T) {

	token1 := "Bearer " + AnnoymousToken
	token2 := ""
	token3 := "Bearer DDEB2204-66C0-5CFD-A351-092EF208ADF5"

	actual := ValidateToken(token1)
	if actual != true {
		t.Errorf("Invalid token = %v; want true", actual)
	}

	actual = ValidateToken(token2)
	if actual != false {
		t.Errorf("Invalid token  = %v; want false", actual)
	}

	actual = ValidateToken(token3)
	if actual != false {
		t.Errorf("Invalid token  = %v; want false", actual)
	}
}
