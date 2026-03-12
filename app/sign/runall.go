package sign

import (
	"fmt"
	"time"

	"github.com/Dasongzi1366/AutoGo/app"
	"github.com/Dasongzi1366/AutoGo/motion"
	"github.com/Dasongzi1366/AutoGo/utils"
)

type appTask struct {
	name string
	pkg  string
	run  func() bool
}

func RunAllApps() {
	tasks := []appTask{
		{
			name: "一汽奥迪",
			pkg:  DefaultAudiPkg,
			run: func() bool {
				return Run("")
			},
		},
		{
			name: "德邦快递",
			pkg:  DefaultDepponPkg,
			run: func() bool {
				return RunDeppon("", "")
			},
		},
		{
			name: "华住会",
			pkg:  DefaultHuazhuPkg,
			run: func() bool {
				return RunHuazhu("")
			},
		},
		{
			name: "中国联通",
			pkg:  DefaultLTPkg,
			run: func() bool {
				return RunLT("")
			},
		},
	}

	toast("[start] 顺序执行 APP 任务")
	for i, t := range tasks {
		toast(fmt.Sprintf("[task %d/%d] %s", i+1, len(tasks), t.name))
		ok, timedOut := runWithTimeout(10*time.Second, t.run)
		if timedOut {
			toast(fmt.Sprintf("[skip] %s 卡住超时(10s)，跳过", t.name))
			bestEffortStop(t.pkg)
			continue
		}
		if ok {
			toast(fmt.Sprintf("[done] %s", t.name))
		} else {
			toast(fmt.Sprintf("[warn] %s 返回失败，继续下一个", t.name))
			bestEffortStop(t.pkg)
		}
		// 给系统一点喘息时间，避免连续启动过快
		time.Sleep(800 * time.Millisecond)
	}
	toast("[done] 全部任务执行结束")
}

func toast(message string) {
	fmt.Println(message)
	utils.Toast(message)
}

func runWithTimeout(timeout time.Duration, fn func() bool) (ok bool, timedOut bool) {
	if timeout <= 0 {
		return fn(), false
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- false
			}
		}()
		done <- fn()
	}()

	select {
	case ok = <-done:
		return ok, false
	case <-time.After(timeout):
		return false, true
	}
}

func bestEffortStop(pkg string) {
	if pkg != "" {
		app.ForceStop(pkg)
	}
	motion.Home(0)
	// 稍微等一下让界面稳定
	time.Sleep(600 * time.Millisecond)
}
