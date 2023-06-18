package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/llmuz/ijk/log"
	"github.com/llmuz/ijk/log/config"
	"github.com/llmuz/ijk/log/zapimpl"
)

type H struct {
	Key string
}

func (c *H) Levels() (lvs []log.Level) {
	return []log.Level{log.DebugLevel, log.InfoLevel, log.ErrorLevel, log.WarnLevel, log.PanicLevel}
}

func (c *H) Fire(e log.Entry) (err error) {
	err = e.AppendField(log.Any(c.Key, e.Context().Value("trace_id")))
	return err
}

var (
	pf          = flag.String("conf", "log/examples/zaplog/config.toml", "配置文件")
	cfg         config.LogConfig
	complexData = make(map[string]interface{})
)

func main() {
	flag.Parse()
	if _, err := toml.DecodeFile(*pf, &cfg); err != nil {
		panic(err)
	}

	logger, err := zapimpl.NewZapLogger(&cfg)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	// 移除日志文件
	//defer os.RemoveAll(filepath.Dir(cfg.GetFileName()))
	n := time.Now()
	ctx := context.WithValue(context.TODO(), "trace_id", md5.New().Sum([]byte(time.Now().String())))
	helper := zapimpl.NewHelper(logger, zapimpl.AddHook(&H{Key: "hello"}), zapimpl.AddHook(&H{Key: "900x"}), zapimpl.AddHook(&H{Key: "trace_id"}))

	fmt.Println("start ", n)
	helper.WithContext(ctx).Errorf("hello %s %s", "你", "好")
	for i := 0; i < 100; i++ {
		helper.WithContext(ctx).Debugf("hello %#v %#v", log.Any("hello", complexData), log.Any("now", time.Now()))
		helper.WithContext(ctx).Infof("hello %#v %#v", log.Any("hello", complexData), log.Any("now", time.Now()))
		helper.WithContext(ctx).Warnf("hello %#v %#v", log.Any("hello", complexData), log.Any("now", time.Now()))
		helper.WithContext(ctx).Errorf("hello %#v %#v", log.Any("hello", complexData), log.Any("now", time.Now()))
	}
	fmt.Println("end ", time.Now().Sub(n))
}
