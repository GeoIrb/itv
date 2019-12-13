package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/GeoIrb/itv/models"
)

func TestDo(t *testing.T) {
	testReq := models.ClientRequest{
		Method: "GET",
		URL:    "https://swapi.co/api/planets/3/",
		Headers: map[string]string{
			"User-Agent": "itv-test",
		},
		Body: "1234",
	}

	res, err := testReq.Do(10 * time.Second)
	if err != nil {
		t.Errorf("TestDo %v", err)
	}

	fmt.Println(res)
}
