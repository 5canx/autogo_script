package sign

import (
	"fmt"
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

const (
	DefaultHuazhuPkg     = "com.htinns"
	DefaultHuazhuMainAct = "com.htinns/com.huazhu.main.RnMainActivity"
)

// RunHuazhu: 启动华住会 APP（默认包名），按既定坐标顺序点击
// 1) (957,2284)-(977,2304) -> 中心 (967,2294)
// 2) (532,656)-(552,676)   -> 中心 (542,666)
// 3) (699,616)-(719,636)   -> 中心 (709,626)
// 4) (520,1792)-(540,1812) -> 中心 (530,1802)
func RunHuazhu(pkg string) bool {
	if pkg == "" {
		pkg = DefaultHuazhuPkg
	}

	fmt.Printf("[start] 华住自动点击，包名: %s\n", pkg)
	EnsureForeground(pkg)
	time.Sleep(3000 * time.Millisecond)

	steps := [][2]int{
		{967, 2294},
		{542, 666},
		{709, 626},
		{530, 1802},
	}

	for i, p := range steps {
		motion.Click(p[0], p[1], 0, 0)
		fmt.Printf("[info] 点击步骤 %d 坐标 (%d,%d)\n", i+1, p[0], p[1])
		time.Sleep(5 * time.Second)
	}

	fmt.Println("[done] 华住点击流程完成")
	return true
}
