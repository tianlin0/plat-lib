package logs_test

import (
	"context"
	"fmt"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
	"testing"
)

func logOne(ctx context.Context) {
	logs.CtxLogger(ctx).Error("ddddd")
}

func TestToString(t *testing.T) {
	logs.DefaultLogger().Error("aaaaa")
	logs.DefaultLogger().Info("bbbb")

	ctx := context.Background()

	logger, newCtx := logs.NewCtxLogger(ctx, logs.INFO, nil)

	logger.Info("cccc")

	logOne(newCtx)
}

func TestFile(t *testing.T) {
	pp, err := utils.SpecifyContext(0)
	fmt.Println(conv.String(pp), err)
}
func TestDebug(t *testing.T) {
	logger := logs.DefaultLogger()
	logger.Debug("aaa")

	//ctx := context.Background()

	//oldGdpLogger := gdplog.Logger(ctx)
	//gdpLogger, newCtx := gdplog.NewGdpLogger(ctx, logs.DEBUG, oldGdpLogger)
	//logs.SetConfig(&logs.Config{
	//	//CtxLogger: func(ctx context.Context) logs.ILogger {
	//	//	return gdplog.GetGdpLogger(ctx)
	//	//},
	//})
	//gdpLogger.Info("gdplog.NewGdpLogger init:")

	//logOne(newCtx)

	logger.Error("ambbdfdf")
}
