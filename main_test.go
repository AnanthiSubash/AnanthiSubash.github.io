package main

import (
	//"go-postgres/router"
	//"io/ioutil"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	total := Add(1, 3)
	assert.Equal(t, 4, total, "Expecting 4")

}
