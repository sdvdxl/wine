package util

import (
	"fmt"
	"testing"
)

func TestIsBlank(t *testing.T) {
	if !IsBlank("") {
		t.FailNow()
	}

	if !IsBlank(`
	`) {
		t.FailNow()
	}

}

func TestHashPassword(t *testing.T) {
	fmt.Println(HashAndSalt("123456", "245accec-3c12-4642-967f-e476cef558c0"))
}
