package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Kong/go-pdk/test"
)

func TestCustomPlugin(t *testing.T) {

	env, err := test.New(t, test.Request{
		Method:  "GET",
		Url:     "http://localhost:8000/auth/verify",
		Headers: map[string][]string{"Authorization": {"Bearer TEST_TOKEN"}},
	})

	if err != nil {
		t.Fatal(err)
	}

	// env.DoHttps(&Config{})

	// t.Run("Sent correct access token must return status code 200 ", func(t *testing.T) {
	// 	env.ClientReq.Headers = map[string][]string{"Authorization": {"TEST_TOKEN"}}
	// 	env.DoHttps(&Config{
	// 		URL: "http://authserver:80/auth/verify",
	// 	})
	// 	if env.ClientRes.Status != 200 {
	// 		t.Error("Must return status code 200")
	// 	}
	// })

	t.Run("Sent empty string must return status code 401", func(t *testing.T) {
		env.ClientReq.Headers = map[string][]string{"Authorization": {"TEST_TOKEN2"}}
		env.DoHttps(
			&Config{
				URL:            "http://authserver:80/auth/verify",
				ConnectTimeout: 1,
				SendTimeout:    1,
				ReadTimeout:    1,
			},
		)
		fmt.Println(env.ClientRes.Status, string(env.ClientRes.Body))
		if env.ClientRes.Status != http.StatusUnauthorized {
			t.Error("Must return status code 401")
		}
	})
}
