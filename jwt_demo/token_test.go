package myjwt

import (
	"fmt"
	"testing"
)

func TestGetToken(t *testing.T) {
	tokenStr, err := CreateToken()
	if err != nil {
		t.Fatalf("get token fail, %v\n", err)
	}
	fmt.Printf("token is : %v\n", tokenStr)
}

func TestParseToken(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJyZWdpc3Rlci1zZXJ2ZXIiLCJleHAiOjE2NDY2NDM4ODEsIm5iZiI6MTY0NjY0Mzg2MSwiaWF0IjoxNjQ2NjQzODYxfQ.HFeBSqWPUkikdUPAM8DAt2r7vGvffprCfJkza07fi4Q"
	err := ParseToken(tokenStr)
	if err != nil {
		t.Fatalf("the err is: %v\n", err)
	}
}
