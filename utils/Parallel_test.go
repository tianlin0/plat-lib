package utils_test

import (
	"fmt"
	"testing"

	"github.com/tianlin0/plat-lib/httputil"
	"github.com/tianlin0/plat-lib/utils"
)

func TestParallel(t *testing.T) {

	pp := httputil.GetLogId()
	fmt.Println(pp)

	return

	filePath := "/Users/tianlin0/go/code/odp-demo/faas"
	// fileName := "/Users/tianlin0/go/code/odp-demo/faas.tar.gz"

	data, err := utils.GetGzFileFromDirectory(filePath, "")
	fmt.Print(data, err)
}
func TestUUID(t *testing.T) {

}
