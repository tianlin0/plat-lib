package curl

import (
	"fmt"
	"github.com/tianlin0/plat-lib/templates"
)

type RetryPolicy struct {
	RetryConditionFunc func(resp string) bool //条件判断的方法，处理比较复杂的问题
	RespDateType       string                 //返回的数据类型，为string和其他类型
	RetryCondition     string                 //重试的条件，满足条件true则重试
	MaxAttempts        int                    //最大重试次数
	leftAttempts       int                    //还剩下的重试次数
}

func canRetry(r *RetryPolicy, respData string) bool {
	if r == nil {
		return false
	}
	if r.leftAttempts <= 0 {
		r.leftAttempts = 0
		return false
	}

	if r.MaxAttempts > 0 {
		retOk, _ := r.onlyCheckCondition(respData)
		if retOk {
			r.leftAttempts--
			return true
		}
	}
	return false
}

func (r *RetryPolicy) onlyCheckCondition(respData string) (bool, error) {
	//方法优先
	if r.RetryConditionFunc != nil {
		return r.RetryConditionFunc(respData), nil
	}

	if r.RetryCondition == "" {
		return false, nil
	}

	if r.RespDateType == "string" {
		//TODO 需要进行字符串的判断，不如返回的是5，大于5的话，条件怎么写呢？
		if respData != "" {
			return false, nil
		}
		return true, nil
	}
	retOk, err := templates.RuleExpr(r.RetryCondition, respData)
	if err != nil {
		return true, err
	}
	return retOk, nil
}

func (r *RetryPolicy) buildRetryable() bool {
	if r.MaxAttempts <= 0 {
		return false
	}
	if r.RetryCondition == "" {
		r.RetryCondition = fmt.Sprintf("%s!=%s", defaultReturnKey, defaultReturnVal)
		return true
	}

	return false
}
