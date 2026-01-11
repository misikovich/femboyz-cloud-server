package uidgenerator

import "testing"

func TestGenerate(t *testing.T) {

	testSample := "87253RTYQJ"
	if !Validate(testSample) {
		t.Errorf("Validate test sample is not valid")
	}

	for i := 0; i < 100000; i++ {
		uid := Generate()
		if len(uid) != 10 {
			t.Errorf("uid length is not 10")
		}
		if !Validate(uid) {
			t.Errorf("uid is not valid")
		}
	}
}
