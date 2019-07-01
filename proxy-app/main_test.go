package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		main()
		wg.Done()
	}(wg)
	fmt.Println("Server running...")
}

type Response struct {
	Status       int            `json:"status,omitempty"`
	Response     string         `json:"result,omitempt"`
	ResponseText []ResponseText `json:"res,omitempty"`
}

type ResponseText struct {
	Domain string
}

func TestMain(t *testing.T) {
	main()

	cases := []struct {
		Domain string
		Output string
	}{
		{Domain: "alpha", Output: `["alpha"]`},
		{Domain: "", Output: "domain error"},
	}
	valuesToCompare := &Response{}
	client := &http.Client{}
	for _, singleCase := range cases {

		request, err := http.NewRequest("GET", "http://localhost:8080/ping", nil)
		assert.Nil(t, err)

		request.Header.Add("domain", singleCase.Domain)
		response, err := client.Do(request)
		assert.Nil(t, err)

		bytes, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bytes, valuesToCompare)
		assert.Nil(t, err)
		fmt.Println(valuesToCompare.Response)
		assert.Equal(t, singleCase.Output, valuesToCompare.Response)
	}
}
