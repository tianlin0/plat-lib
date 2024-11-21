// Package retry 重试器
package retry

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type retry struct {
	attemptCount int             //最大尝试次数
	interval     time.Duration   //间隔时间
	errCallFun   ErrCallbackFunc //执行错误的方法
}

type Executable func(context.Context) (interface{}, error)

/*ErrCallbackFunc 回调函数
  nowAttemptCount 当前尝试次数
  remainCount 剩余次数
  err 错误
*/
type ErrCallbackFunc func(err error) error

// New 创建一个重试器
func New() *retry {
	r := &retry{
		attemptCount: 1,
		interval:     5 * time.Second,
	}
	return r
}

// WithInterval 设置间隔时间
func (r *retry) WithInterval(interval time.Duration) *retry {
	r.interval = interval
	return r
}

// WithAttemptCount 设置最大尝试次数, 0为不限次数
func (r *retry) WithAttemptCount(attemptCount int) *retry {
	r.attemptCount = attemptCount
	return r
}

// WithErrCallback 设置错误回调函数, 每次执行时有任何错误都会报告给该函数
func (r *retry) WithErrCallback(errFun ErrCallbackFunc) *retry {
	r.errCallFun = errFun
	return r
}

// Do 执行一个函数
func (r *retry) Do(parentCtx context.Context, f Executable, valuePtr ...interface{}) error {
	if len(valuePtr) > 0 && valuePtr != nil {
		if valuePtr[0] != nil {
			rf := reflect.ValueOf(valuePtr[0])
			if rf.Type().Kind() != reflect.Ptr {
				return fmt.Errorf("valuePtr parameter is not a pointer")
			}
		}
	}

	var retData interface{}
	var err error
	if parentCtx != nil {
		retData, err = r.doRetryWithCtx(parentCtx, f)
	} else {
		retData, err = r.doRetry(f)
	}

	if len(valuePtr) > 0 && valuePtr != nil {
		if valuePtr[0] != nil {
			rf := reflect.ValueOf(valuePtr[0])
			if rf.Elem().CanSet() {
				fv := reflect.ValueOf(retData)
				isSet := false
				if fv.Kind() == reflect.Ptr && fv.Type() == rf.Type() {
					if fv.Elem().IsValid() {
						isSet = true
						rf.Elem().Set(fv.Elem())
					}
				} else {
					if fv.IsValid() {
						isSet = true
						rf.Elem().Set(fv)
					}
				}
				if !isSet {
					if err == nil {
						err = fmt.Errorf("call Return: call of reflect.Value.Set on zero Value")
					}
				}
			}
		}
	}
	return err
}

// DoCtx 执行一个函数
func (r *retry) doRetryWithCtx(parentCtx context.Context, fn Executable) (interface{}, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	nowAttemptCount := 0
	fail := make(chan error, 1)
	success := make(chan interface{}, 1)

	for {
		go func() {
			val, err := fn(ctx)
			if err != nil {
				fail <- err
				return
			}
			success <- val
		}()

		select {
		//
		case <-parentCtx.Done():
			return nil, parentCtx.Err()

		case err := <-fail:
			if parentCtxErr := parentCtx.Err(); parentCtxErr != nil {
				return nil, parentCtxErr
			}

			nowError := fmt.Errorf("Max retries exceeded (%v).\n", nowAttemptCount)

			if r.errCallFun != nil {
				errError := r.errCallFun(err)
				if errError != nil {
					return nil, nowError
				}
			}

			nowAttemptCount++

			if nowAttemptCount >= r.attemptCount {
				return nil, nowError
			}

			if r.interval > 0 {
				time.Sleep(r.interval)
			}

		case val := <-success:
			return val, nil
		}
	}
}

// doRetry 执行一个函数
func (r *retry) doRetry(f Executable) (val interface{}, err error) {
	nowAttemptCount := 0
	parentCtx := context.Background()
	for {
		nowAttemptCount++

		val, err = f(parentCtx)
		if err == nil {
			return
		}
		if r.errCallFun != nil {
			errError := r.errCallFun(err)
			if errError != nil {
				//表示这个是致命错误，不用重试了
				return
			}
		}

		if nowAttemptCount >= r.attemptCount {
			return
		}

		if r.interval > 0 {
			time.Sleep(r.interval)
		}
	}
}
