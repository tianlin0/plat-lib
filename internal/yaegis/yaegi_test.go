package yaegis

import (
	"bytes"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"os/exec"
)

import "testing"

func TestEval(t *testing.T) {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)

	_, err := i.Eval(`import "fmt"`)
	if err != nil {
		panic(err)
	}

	_, err = i.Eval(`fmt.Println("Hello Yaegi")`)
	if err != nil {
		panic(err)
	}
}

func TestTScript(t *testing.T) {
	tScriptCode := `
        print("Hello, Tscript!");
    `

	// 将 Tscript 代码传递给解释器
	cmd := exec.Command("ts-interpreter")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(tScriptCode))
	}()

	// 获取解释器的输出
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	// 打印输出结果
	fmt.Println(out.String())
}
