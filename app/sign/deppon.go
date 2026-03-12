package sign

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dasongzi1366/AutoGo/app"
	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/motion"
	"github.com/Dasongzi1366/AutoGo/ppocr"
	"github.com/Dasongzi1366/AutoGo/uiacc"
)

// 默认德邦快递包名与主入口 Activity
const (
	DefaultDepponPkg     = "com.deppon.dpapp"
	DefaultDepponMainAct = "com.deppon.dpapp/.MainActivity"
	patternAdCloseButton = "533,1564,ebebeb,573,1547,8c8c8c,580,1547,fbfbfb,587,1547,8c8c8c,594,1547,8c8c8c,573,1554,8c8c8c,580,1554,8c8c8c,587,1554,8c8c8c,594,1554,8c8c8c,573,1561,8c8c8c,580,1561,8c8c8c,587,1561,a9a9a9,594,1561,8c8c8c,573,1568,8c8c8c,580,1568,8c8c8c,587,1568,c5c5c5,594,1568,8c8c8c,515,1593,8c8c8c,522,1593,8c8c8c,529,1593,8c8c8c,536,1593,8c8c8c,515,1600,8c8c8c,522,1600,8c8c8c,529,1600,8c8c8c,536,1600,8c8c8c,515,1607,e5e5e5,522,1607,8c8c8c,529,1607,8c8c8c,536,1607,8c8c8c,515,1614,8c8c8c,522,1614,a4a4a4,529,1614,fbfbfb,536,1614,dedede"
	patternMyPoints      = "581,661,7c482d,588,661,bb9872,595,661,bb9872,602,661,bb9872,581,668,7c482d,588,668,956848,595,668,fff1be,602,668,f9e9b6,581,675,7c482d,588,675,fff1be,595,675,7c482d,602,675,7c482d,581,682,ae8965,588,682,fff1be,595,682,fff1be,602,682,fdedbb,652,661,f3e0ae,659,661,ddc396,666,661,ffefbb,673,661,ffefbb,652,668,dec698,659,668,7c482d,666,668,7c482d,673,668,ffefbb,652,675,fff0bc,659,675,7c482d,666,675,fbebb8,673,675,ffefbb,652,682,fff0bc,659,682,fff0bc,666,682,fff0bc,673,682,fbebb7,687,658,8b5c3e,694,658,ffefba,701,658,ddc294,708,658,ffeeb8,687,665,ac8560,694,665,ac8560,701,665,ac8560,708,665,f1dcaa,687,672,7c482d,694,672,fff0ba,701,672,7c482d,708,672,ffefb9,687,679,ffefbb,694,679,fff0bb,701,679,c8a87e,708,679,ffefb9"
	patternSignBtn       = "ffffff,-9,-11,ffffff,85,1,569eff,42,19,1a78ff"
	patternSignedDone    = "373,1167,e3a67c,258,1172,fbe1c6,521,1174,ffecd3,527,1188,c55a1c"
)

// RunDeppon 启动德邦快递 APP，尝试进入 MainActivity。
func RunDeppon(pkg, activity string) bool {
	if strings.TrimSpace(pkg) == "" {
		pkg = DefaultDepponPkg
	}
	if strings.TrimSpace(activity) == "" {
		activity = DefaultDepponMainAct
	}

	fmt.Printf("[start] 打开德邦快递 APP，包名: %s, Activity: %s\n", pkg, activity)
	app.Enable(pkg)
	app.EnableAccessibility(pkg)
	fmt.Println("[info] 已请求开启无障碍以便节点检测")

	if !EnsureForeground(pkg) {
		fmt.Println("[fail] 无法切到德邦前台")
		return false
	}

	// 先按包名启动，若需要特定 Activity，可在后续补充 Intent 方式。
	if app.Launch(pkg, 0) {
		fmt.Println("[success] 已请求启动德邦快递 APP")
		time.Sleep(5 * time.Second)

		if closeAdIfPresent() {
			fmt.Println("[info] 发现并尝试关闭广告")
		}
		if ensureMyTab() {
			fmt.Println("[success] 已切换到'我的'")
			clickMyPoints()
			clickSignInPoints()
			readMyBalance()
		}
		return true
	}

	fmt.Println("[fail] 启动德邦快递 APP 失败，请检查是否已安装")
	return false
}

func closeAdIfPresent() bool {
	if images.DetectsMultiColors(patternAdCloseButton, 0.9, 0) {
		// 点击模板区域中心附近的经验坐标
		motion.Click(555, 1585, 0, 0)
		time.Sleep(800 * time.Millisecond)
		return true
	}
	return false
}

func ensureMyTab() bool {
	// 优先用无障碍 ContentDesc="我的" 点击
	ui := uiacc.New(0)
	if ui != nil {
		defer ui.Release()
		if obj := ui.Desc("我的").Clickable(true).FindOnce(); obj != nil {
			if obj.Click() {
				time.Sleep(1200 * time.Millisecond)
				return true
			}
		}
	}

	// 回退：点击底栏经验坐标（靠右）
	motion.Click(1000, 2300, 0, 0)
	time.Sleep(1200 * time.Millisecond)
	return true
}

func clickMyPoints() {
	ui := uiacc.New(0)
	if ui == nil {
		return
	}
	defer ui.Release()

	// 优先比色定位“我的积分”卡片
	if images.DetectsMultiColors(patternMyPoints, 0.9, 0) {
		motion.Click(645, 672, 0, 0) // 模板中心经验坐标
		time.Sleep(1500 * time.Millisecond)
		fmt.Println("[info] 已通过比色点击'我的积分'")
		return
	}

	// 回退：用无障碍 DescContains 点击
	objs := ui.DescContains("我的积分").Find()
	for _, o := range objs {
		if o == nil {
			continue
		}
		if o.Click() {
			fmt.Println("[info] 已点击'我的积分'")
			time.Sleep(800 * time.Millisecond)
			return
		}
	}
	fmt.Println("[info] 未找到可点击的'我的积分'")
}

func clickSignInPoints() {
	ui := uiacc.New(0)
	if ui == nil {
		return
	}
	defer ui.Release()

	for attempts := 0; attempts < 4; attempts++ {
		// 等待页面元素加载
		time.Sleep(2000 * time.Millisecond)

		if hasSignedToday(ui) {
			return
		}
		if images.DetectsMultiColors(patternSignedDone, 0.9, 0) {
			fmt.Println("[info] 比色检测到已签到提示，返回上一页")
			motion.Back(0)
			time.Sleep(800 * time.Millisecond)
			return
		}

		// 优先比色找“签到领积分”按钮
		if x, y := images.FindMultiColors(393, 1137, 516, 1199, patternSignBtn, 0.9, 0, 0); x >= 0 && y >= 0 {
			motion.Click(x, y, 0, 0)
			time.Sleep(1200 * time.Millisecond)
			return
		}

		// 收集候选
		exact := ui.Text("签到领积分").Find()
		contains := ui.TextContains("签到领积分").Find()
		descs := ui.DescContains("签到领积分").Find()

		// 兜底：TextView 同时包含“签到”、“积分”
		fallback := []*uiacc.UiObject{}
		tvs := ui.ClassName("android.widget.TextView").Find()
		for _, o := range tvs {
			if o == nil {
				continue
			}
			txt := strings.TrimSpace(o.GetText())
			if txt == "" {
				continue
			}
			if strings.Contains(txt, "签到") && strings.Contains(txt, "积分") {
				fallback = append(fallback, o)
			}
		}

		// 先用 exact，目标第2个（index 1）
		if len(exact) > 0 {
			idx := 1
			if len(exact) <= 1 {
				idx = 0
			}
			if obj := exact[idx]; obj != nil && obj.Click() {
				time.Sleep(1200 * time.Millisecond)
				return
			}
		}

		// 否则用 contains+desc 合并
		cands := append(contains, descs...)
		if len(cands) > 0 {
			idx := 1
			if len(cands) <= 1 {
				idx = 0
			}
			if obj := cands[idx]; obj != nil && obj.Click() {
				time.Sleep(1200 * time.Millisecond)
				return
			}
		}

		// 最后尝试 fallback
		if len(fallback) > 0 {
			if obj := fallback[0]; obj != nil && obj.Click() {
				time.Sleep(1200 * time.Millisecond)
				return
			}
		}
	}

	fmt.Println("[info] 点击'签到领积分'失败")
}

func readMyBalance() {
	ocr := ppocr.New("v5")
	if ocr == nil {
		return
	}
	res := ocr.Ocr(707, 637, 797, 695, "", 0)
	if len(res) == 0 {
		fmt.Println("[info] 未识别到积分余额")
		return
	}
	fmt.Printf("[info] 积分余额: %s\n", strings.TrimSpace(res[0].Label))
}

func hasSignedToday(ui *uiacc.Uiacc) bool {
	if ui == nil {
		return false
	}

	exact := ui.Text("今日已签到，请明天再来").Find()
	contains := ui.TextContains("今日已签到").Find()
	contains = append(contains, ui.TextContains("签到，请明天再来").Find()...)
	contains = append(contains, ui.TextContains("已签到").Find()...)
	contains = append(contains, ui.DescContains("已签到").Find()...)

	fmt.Printf("[info] 已签到提示匹配: exact=%d, contains=%d\n", len(exact), len(contains))

	if len(exact) > 0 {
		fmt.Printf("[info] 检测到'今日已签到，请明天再来' 文本，数量: %d\n", len(exact))
		return true
	}
	if len(contains) > 0 {
		fmt.Printf("[info] 检测到包含'已签到' 文本，数量: %d\n", len(contains))
		return true
	}
	return false
}
