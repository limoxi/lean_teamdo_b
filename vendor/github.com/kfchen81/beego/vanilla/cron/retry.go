package cron

import (
	"context"
	"fmt"
	"github.com/kfchen81/beego"
	"github.com/kfchen81/beego/vanilla/backoff"
	"github.com/kfchen81/beego/logs"
	"runtime"
	"time"
	"bytes"
)

type RetryTaskParam struct {
	NewContext func() context.Context
	GetDatas func() []interface{}
	BeforeAction func(data interface{}) error
	DoAction func(ctx context.Context, times int, data interface{}) error
	AfterActionSuccess func(data interface{}) error
	AfterActionFail func(data interface{}) error
	GetTaskDataId func(data interface{}) string
	RecordFailByPanic func(data interface{}, error string)
}

func retryTaskGorutione(maxMinutes int, taskParam *RetryTaskParam, goroutineTimes int, data interface{}) {
	defer func() {
		if err := recover(); err != nil {
			beego.Error(err)
			errMsg := fmt.Sprintf("%s", err)
			var buffer bytes.Buffer
			buffer.WriteString(fmt.Sprintf("[Unprocessed_Exception] %s\n", errMsg))
			for i := 1; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				buffer.WriteString(fmt.Sprintf("%s:%d\n", file, line))
			}
			logs.Error(buffer.String())
			
			if goroutineTimes <= 3 {
				time.Sleep(4 * time.Second)
				beego.Warn(fmt.Sprintf("[retry] restart goroutine for %d times", goroutineTimes))
				go retryTaskGorutione(maxMinutes, taskParam, goroutineTimes+1, data)
			} else {
				//需要捕捉在AfterActionFail和
				defer func() {
					if err3 := recover(); err3 != nil {
						beego.Error(err3)
					}
				}()
				
				err2 := taskParam.AfterActionFail(data)
				if err2 != nil {
					beego.Error(err2)
				}
				
				taskParam.RecordFailByPanic(data, errMsg)
			}
		}
	}()
	
	expBackoff := &backoff.ExponentialBackOff{
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          1.8,// * backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      time.Duration(maxMinutes) * time.Minute,
		Clock:               backoff.SystemClock,
	}
	expBackoff.Reset()
	
	if taskParam.BeforeAction != nil {
		err := taskParam.BeforeAction(data)
		if err != nil {
			beego.Error(err)
			return
		}
	}
	
	times := 0
	ctx := context.Background()
	err := backoff.RetryNotify(func() error {
		times += 1
		return taskParam.DoAction(ctx, times, data)
	}, expBackoff, func (err error, duration time.Duration) {
		beego.Warn(fmt.Sprintf("[push_order_payment] push '%s' fail %d times, because of : %s, next push after %v", taskParam.GetTaskDataId(data), times, err.Error(), duration))
	})
	
	if err != nil {
		beego.Error(err)
		
		err := taskParam.AfterActionFail(data)
		if err != nil {
			beego.Error(err)
		}
	} else {
		err := taskParam.AfterActionSuccess(data)
		if err != nil {
			beego.Error(err)
		}
	}
}

func StartRetryTask(maxMinutes int, taskParam *RetryTaskParam) {
	if taskParam.BeforeAction == nil {
		beego.Error("[retry] Need taskParam.BeforeAction != nil")
		return
	}
	if taskParam.DoAction == nil {
		beego.Error("[retry] Need taskParam.DoAction != nil")
		return
	}
	if taskParam.GetDatas == nil {
		beego.Error("[retry] Need taskParam.GetDatas != nil")
		return
	}
	if taskParam.AfterActionSuccess == nil {
		beego.Error("[retry] Need taskParam.AfterActionSuccess != nil")
		return
	}
	if taskParam.AfterActionFail == nil {
		beego.Error("[retry] Need taskParam.AfterActionFail != nil")
		return
	}
	if taskParam.GetTaskDataId == nil {
		beego.Error("[retry] Need taskParam.GetTaskDataId != nil")
		return
	}
	if taskParam.RecordFailByPanic == nil {
		beego.Error("[retry] Need taskParam.RecordFailByPanic != nil")
		return
	}
	
	//datas := getOrdersNeedPush()
	datas := taskParam.GetDatas()
	
	for _, data := range datas {
		go retryTaskGorutione(maxMinutes, taskParam, 1, data)
	}
}
