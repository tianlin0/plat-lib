package goroutines_test

import (
	"fmt"
	"github.com/tianlin0/plat-lib/goroutines"
	"math"
	"testing"
	"time"
)

type AA struct {
	K int
}

func (m *AA) Add(a int, b int) (int, map[string]interface{}, error) {
	return a + b + m.K, map[string]interface{}{
		"aaaaa": "1111",
	}, fmt.Errorf("error is null")
}

func TestParallel(t *testing.T) {
	m := &AA{}
	cf := goroutines.NewCallFunc().SetFunc(m.Add)
	m.K = 10
	var c int
	var s map[string]string
	var ee error
	err := cf.SetParams(1, 3).GetAll(&c, &s, &ee)
	fmt.Println(ee == nil)
	fmt.Println(c, s, err)
	fmt.Println(ee)
}
func TestParallel2(t *testing.T) {
	dataList := 13
	stepNum := 3
	timeFloat := float64(dataList) / float64(stepNum)
	asynTimes := int(math.Ceil(timeFloat))

	fmt.Println(asynTimes)

	i := 0

	for j := 0; j < asynTimes; j++ {
		for ; i < dataList; i++ {
			if i >= (j+1)*stepNum {
				break
			}
			fmt.Println(j, i)
		}
	}

}
func TestParallel33(t *testing.T) {
	var err2 any
	err2 = fmt.Errorf("aaaaa")
	if interface{}(err2) != nil {
		fmt.Println(222)
	}

}

func getData(aa string) (string, error) {
	time.Sleep(2 * time.Second)
	return aa + "getData", nil
}

func TestTimeout(t *testing.T) {
	name := ""
	aa, err := goroutines.RunWithTimeout(1*time.Second, func() (string, error) {
		return getData(name)
	})

	fmt.Println(aa, err)

}
