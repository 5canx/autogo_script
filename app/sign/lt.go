package sign

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/motion"
	"github.com/Dasongzi1366/AutoGo/ppocr"
	"github.com/Dasongzi1366/AutoGo/uiacc"
)

// 中国联通
const DefaultLTPkg = "com.sinovatech.unicom.ui"

const patternLtOpened = "977,182,950f0f,984,182,971614,991,182,971614,998,182,b54e3f,977,189,950f0f,984,189,be645d,991,189,971513,998,189,f0d3c3,977,196,d3928a,984,196,d08774,991,196,d39288,998,196,ae4440,977,203,ffebce,984,203,feebc6,991,203,ffe9cd,998,203,ffe5c5,1005,161,ffa417,1012,161,f8eadf,1019,161,fe9615,1026,161,feead2,1005,168,ff9716,1012,168,e1a678,1019,168,fe9315,1026,168,fab15c,1005,175,fffdfb,1012,175,fffcf6,1019,175,fdf6ee,1026,175,e8263d,1005,182,ac423f,1012,182,efcdbe,1019,182,fef2e5,1026,182,e8263e"

// 已签到页面特征
const patternLtSigned = "982,185,ff490e,975,156,ffa22e,982,156,ff9a1e,989,156,ff991c,996,156,ff991c,975,163,ff7e09,982,163,ff7d06,989,163,ff7d06,996,163,ff7c06,975,170,ff8e3d,982,170,ff6708,989,170,fc9251,996,170,ffc5ad,975,177,ff540a,982,177,ff4b09,989,177,f84213,996,177,ffb498,987,197,ff2e1a,994,197,ffe3c8,1001,197,ff2e1a,1008,197,ffe3c8,987,204,ff2237,994,204,ffdcba,1001,204,ff2338,1008,204,ffd4b5,987,211,fe1c4c,994,211,fe1d4f,1001,211,fe1d4e,1008,211,fe1b47,987,218,f33b4f,994,218,f33b50,1001,218,f33b50,1008,218,f23d4d"

// 我的页-切换账号入口按钮特征
const patternLtSwitchAccount = "867,177,424242,874,177,f8f8f8,881,177,f8f8f8,888,177,f8f8f8,867,184,333333,874,184,f8f8f8,881,184,f8f8f8,888,184,808080,867,191,f8f8f8,874,191,333333,881,191,f8f8f8,888,191,f8f8f8,867,198,f8f8f8,874,198,333333,881,198,f8f8f8,888,198,f8f8f8,906,153,b0b0b0,913,153,f8f8f8,920,153,f8f8f8,927,153,f8f8f8,906,160,f8f8f8,913,160,6a6a6a,920,160,f8f8f8,927,160,f8f8f8,906,167,f8f8f8,913,167,626262,920,167,f8f8f8,927,167,f8f8f8,906,174,f8f8f8,913,174,f8f8f8,920,174,333333,927,174,f8f8f8"

const (
	patternLtAdClose = "508,1470,ffffff,533,1482,ffffff,545,1472,666666"
	// 广告 OCR 区域
	ltAdOcrX1 = 131
	ltAdOcrY1 = 1086
	ltAdOcrX2 = 955
	ltAdOcrY2 = 1156

	// 打开联通后：当前账号区域
	ltAccountX1 = 821
	ltAccountY1 = 478
	ltAccountX2 = 1005
	ltAccountY2 = 529

	// 底栏“我的”区域（用于切换账号前回到我的页）
	ltMyTabX1 = 927
	ltMyTabY1 = 2309
	ltMyTabX2 = 1023
	ltMyTabY2 = 2362
)

// RunLT 启动中国联通并检测指定特征点（用于确认已进入目标页面）。
func RunLT(pkg string) bool {
	if strings.TrimSpace(pkg) == "" {
		pkg = DefaultLTPkg
	}

	fmt.Printf("[start] 打开中国联通 APP，包名: %s\n", pkg)
	if !EnsureForeground(pkg) {
		fmt.Println("[fail] 无法切到联通前台")
		return false
	}

	// 给启动页一点时间
	time.Sleep(8 * time.Second)

	processed := map[string]bool{}
	noStateStreak := 0
	for round := 0; round < 10; round++ {
		signed, opened := waitLtState(12, 800*time.Millisecond)
		if !signed && !opened {
			noStateStreak++
			fmt.Println("[warn] 未检测到联通页面特征，尝试返回恢复")
			recoverLtToHome(noStateStreak)
			if noStateStreak >= 3 {
				// 连续多次无法识别，结束本次流程（避免死循环）
				return round > 0
			}
			// 继续下一轮重试
			round--
			continue
		}
		noStateStreak = 0

		acc := strings.TrimSpace(ocrRegionText(ltAccountX1, ltAccountY1, ltAccountX2, ltAccountY2))
		if acc != "" {
			fmt.Printf("[info] 当前账号: %s\n", acc)
			processed[acc] = true
		} else {
			fmt.Println("[info] 当前账号: (未识别)")
		}

		if signed {
			fmt.Println("[info] 当前账号已签到")
		} else {
			if clickLtSignIn() {
				handleAdAfterSign()
			}
		}

		// 切换到下一个账号；若无可切换则结束
		if !gotoMyAndSwitchNext(processed) {
			return true
		}
		// 切换后等待页面稳定
		time.Sleep(6 * time.Second)
	}

	return true
}

func recoverLtToHome(streak int) {
	// 尝试通过返回键把页面带回主界面/首页
	// streak 越大，多返回几次
	backs := 1
	if streak >= 2 {
		backs = 2
	}
	for i := 0; i < backs; i++ {
		motion.Back(0)
		time.Sleep(1200 * time.Millisecond)
	}
	// 再点一次首页 tab（如果无障碍可用）
	clickHomeTab()
	time.Sleep(2000 * time.Millisecond)
}

func waitLtState(tries int, interval time.Duration) (signed bool, opened bool) {
	for i := 0; i < tries; i++ {
		if images.DetectsMultiColors(patternLtSigned, 0.9, 0) {
			return true, false
		}
		if images.DetectsMultiColors(patternLtOpened, 0.9, 0) {
			return false, true
		}
		time.Sleep(interval)
	}
	return false, false
}

func gotoMyAndSwitchNext(processed map[string]bool) bool {
	clickMyTabByOcr()
	if !clickSwitchAccountEntry() {
		return false
	}
	return switchAccountNext(processed)
}

func clickLtSignIn() bool {
	// 优先用无障碍查找“签到”相关按钮
	ui := uiacc.New(0)
	if ui != nil {
		defer ui.Release()
		for _, key := range []string{"签到", "立即签到", "去签到"} {
			if obj := ui.TextContains(key).Clickable(true).FindOnce(); obj != nil {
				if obj.Click() {
					fmt.Printf("[info] 已点击联通签到入口: %s\n", key)
					time.Sleep(2000 * time.Millisecond)
					return true
				}
			}
		}
	}

	fmt.Println("[warn] 未找到联通签到入口（需要提供签到按钮坐标/模板）")
	return false
}

func handleAdAfterSign() {
	// 给弹窗一点时间
	time.Sleep(2000 * time.Millisecond)

	// OCR 识别并打印
	if text := ocrLtAdText(); strings.TrimSpace(text) != "" {
		fmt.Printf("[info] 广告 OCR: %s\n", text)
	} else {
		fmt.Println("[info] 广告 OCR: (空)")
	}

	// 关闭广告按钮（静默轮询，不刷日志）
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if images.DetectsMultiColors(patternLtAdClose, 0.9, 0) {
			motion.Click(528, 1476, 0, 0)
			time.Sleep(1200 * time.Millisecond)
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func clickMyTabByOcr() {
	// 轮询识别“我的”，存在则点击（日志不刷屏）
	deadline := time.Now().Add(12 * time.Second)
	centerX := (ltMyTabX1 + ltMyTabX2) / 2
	centerY := (ltMyTabY1 + ltMyTabY2) / 2
	for time.Now().Before(deadline) {
		t := strings.TrimSpace(ocrRegionText(ltMyTabX1, ltMyTabY1, ltMyTabX2, ltMyTabY2))
		if strings.Contains(t, "我的") {
			motion.Click(centerX, centerY, 0, 0)
			time.Sleep(1200 * time.Millisecond)
			fmt.Println("[info] 已点击底栏'我的'")
			return
		}
		time.Sleep(600 * time.Millisecond)
	}
}

func clickSwitchAccountEntry() bool {
	deadline := time.Now().Add(10 * time.Second)
	// 模板覆盖区域大致中心点 (867..927, 153..198)
	centerX, centerY := 897, 175
	for time.Now().Before(deadline) {
		if images.DetectsMultiColors(patternLtSwitchAccount, 0.9, 0) {
			motion.Click(centerX, centerY, 0, 0)
			time.Sleep(1200 * time.Millisecond)
			fmt.Println("[info] 已点击切换账号入口")
			return true
		}
		time.Sleep(400 * time.Millisecond)
	}
	return false
}

func switchAccountNext(processed map[string]bool) bool {
	// 下滑 0.5 秒
	motion.Swipe(540, 1750, 540, 800, 500, 0, 0)
	time.Sleep(800 * time.Millisecond)

	ui := uiacc.New(0)
	if ui == nil {
		return false
	}
	defer ui.Release()

	if obj := ui.Id("user_switch_button").Clickable(true).FindOnce(); obj != nil {
		if obj.Click() {
			time.Sleep(1 * time.Second)
		}
	}

	items := ui.Id("login_bingding_item_layout").Clickable(true).Find()
	if len(items) == 0 {
		return false
	}
	titles := ui.Id("login_bingding_item_title").ClassName("android.widget.TextView").Find()
	max := len(items)
	if len(titles) < max {
		max = len(titles)
	}

	// 优先按标题去重，逐个切换
	for i := 0; i < max; i++ {
		title := strings.TrimSpace(titles[i].GetText())
		if title == "" {
			continue
		}
		if processed[title] {
			continue
		}
		fmt.Printf("[info] 切换目标: %s\n", title)
		processed[title] = true
		items[i].Click()
		time.Sleep(1200 * time.Millisecond)
		backAndGoHome()
		return true
	}

	// 如果标题不可用，则最多切一次第一个 item
	if len(titles) == 0 && items[0] != nil {
		items[0].Click()
		time.Sleep(1200 * time.Millisecond)
		backAndGoHome()
		return true
	}

	return false
}

func backAndGoHome() {
	// 切换账号后返回一层，再点首页 Tab
	motion.Back(0)
	time.Sleep(1200 * time.Millisecond)
	clickHomeTab()
	time.Sleep(2000 * time.Millisecond)
	handleHomeAdIfPresent()
}

func clickHomeTab() {
	ui := uiacc.New(0)
	if ui == nil {
		return
	}
	defer ui.Release()
	if obj := ui.Id("main_tab_home_layout").Clickable(true).FindOnce(); obj != nil {
		obj.Click()
	}
}

func handleHomeAdIfPresent() {
	ui := uiacc.New(0)
	if ui == nil {
		return
	}
	defer ui.Release()

	// 判断是否存在广告标识
	if ui.Bounds(87, 1593, 153, 1641).Text("广告").FindOnce() == nil {
		return
	}

	// 关闭按钮（按给定 index+bounds 定位）
	if btn := ui.Index(2).Bounds(498, 1767, 582, 1848).Clickable(true).FindOnce(); btn != nil {
		if btn.Click() {
			time.Sleep(1200 * time.Millisecond)
			fmt.Println("[info] 已关闭首页广告")
		}
	}
}

func ocrLtAdText() string {
	return ocrRegionText(ltAdOcrX1, ltAdOcrY1, ltAdOcrX2, ltAdOcrY2)
}

func ocrRegionText(x1, y1, x2, y2 int) string {
	ocr := ppocr.New("v5")
	if ocr == nil {
		return ""
	}
	res := ocr.Ocr(x1, y1, x2, y2, "", 0)
	if len(res) == 0 {
		return ""
	}
	parts := make([]string, 0, len(res))
	for _, r := range res {
		if s := strings.TrimSpace(r.Label); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, " ")
}
