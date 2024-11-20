package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	workingDir     = "/"
	stackCache     = make(map[uintptr]*logContext)
	stackCacheLock sync.RWMutex
)

type logContext struct {
	FuncName  string
	Line      int
	ShortPath string
	FullPath  string
	FileName  string
	CallTime  time.Time
}

// SpecifyContext TODO
func SpecifyContext(skip int) (*logContext, error) {
	callTime := time.Now()
	if skip < 0 {
		err := fmt.Errorf("can not skip negative stack frames")
		return nil, err
	}
	caller, err := extractCallerInfo(skip + 2)
	if err != nil {
		return nil, err
	}
	ctx := new(logContext)
	*ctx = *caller
	ctx.CallTime = callTime
	return ctx, nil
}

func extractCallerInfo(skip int) (*logContext, error) {
	var stack [1]uintptr
	if runtime.Callers(skip+1, stack[:]) != 1 {
		return nil, fmt.Errorf("error  during runtime.Callers")
	}
	pc := stack[0]

	// do we have a cache entry?
	stackCacheLock.RLock()
	ctx, ok := stackCache[pc]
	stackCacheLock.RUnlock()
	if ok {
		return ctx, nil
	}

	// look up the details of the given caller
	funcInfo := runtime.FuncForPC(pc)
	if funcInfo == nil {
		return nil, fmt.Errorf("error during runtime.FuncForPC")
	}

	var shortPath string
	fullPath, line := funcInfo.FileLine(pc)
	if strings.HasPrefix(fullPath, workingDir) {
		shortPath = fullPath[len(workingDir):]
	} else {
		shortPath = fullPath
	}
	funcName := funcInfo.Name()
	if strings.HasPrefix(funcName, workingDir) {
		funcName = funcName[len(workingDir):]
	}

	_, fileFullPath, lineTemp, ok := runtime.Caller(skip)
	if !ok {
		lineTemp = 0
	} else {
		if lineTemp > 0 {
			line = lineTemp
			fullPath = fileFullPath
		}
	}

	ctx = &logContext{
		FuncName:  funcName,
		Line:      line,
		ShortPath: shortPath,
		FullPath:  fullPath,
		FileName:  filepath.Base(fullPath),
	}

	// save the details in the cache; note that it's possible we might
	// have written an entry into the map in between the test above and
	// this section, but the behaviour is still correct
	stackCacheLock.Lock()
	stackCache[pc] = ctx
	stackCacheLock.Unlock()
	return ctx, nil
}
