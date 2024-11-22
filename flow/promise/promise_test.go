package promise_test

import (
	"fmt"
	"github.com/tianlin0/plat-lib/flow/promise"
	"net/http"
	"testing"
)

func TestPromise(t *testing.T) {
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		resp, err := http.Get("http://gobyexample.com")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			resolve(resp.Status)
		} else {
			reject(fmt.Errorf("ERRR"))
		}
	})

	p.Then(func(value interface{}) {
		fmt.Println("Success", value)
	}).Catch(func(err error) {
		fmt.Println("Error", err)
	})

	p.Wait()
}
