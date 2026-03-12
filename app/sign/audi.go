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

const (
	patternSmallBtn     = ""
	patternSignedColors = "799,1013,ffffff,930,1010,fcf9f6"
	patternUnselectedMy = "954,2304,999999,977,2246,d9d9d9"
	signRegionX1        = 791
	signRegionY1        = 1067
	signRegionX2        = 943
	signRegionY2        = 1118
)

// 默认一汽奥迪包名，支持通过 Run 传入覆盖。
// 实测前台窗口为 com.timanetworks.android.frame.audisuper.release。
const DefaultAudiPkg = "com.timanetworks.android.frame.audisuper.release"

func Run(pkg string) bool {
	if strings.TrimSpace(pkg) == "" {
		pkg = DefaultAudiPkg
	}

	fmt.Printf("[start] 打开一汽奥迪 APP，包名: %s\n", pkg)
	if !EnsureForeground(pkg) {
		fmt.Println("[fail] 无法切到奥迪前台")
		return false
	}

	if launchApp(pkg) {
		fmt.Println("[success] 已请求启动奥迪 APP")
		time.Sleep(10 * time.Second)
		if ensureMyPage() {
			fmt.Println("[success] 已切换到'我的'页面")
			if signInOnce() {
				fmt.Println("[info] 已尝试点击签到并读取积分区域")
			}
		} else {
			fmt.Println("[warn] 未能确认切到'我的'页面")
		}
		return true
	}
	fmt.Println("[fail] 启动奥迪 APP 失败，请检查包名或是否已安装")
	return false
}

func launchApp(pkg string) bool {
	app.Enable(pkg)
	if app.Launch(pkg, 0) {
		fmt.Println("[info] 已通过 app.Launch 启动")
		return true
	}
	fmt.Println("[warn] app.Launch 启动失败，请检查包名是否正确且已安装")
	return false
}

// 页面识别与点击“我的”
func ensureMyPage() bool {
	// 尝试 4 次：若未选中则点击底栏右侧区域
	for i := 0; i < 4; i++ {
		if detectMySelected() {
			return true
		}
		fmt.Printf("[info] 未检测到'我的'选中状态，第 %d 次点击底栏\n", i+1)
		clickBottomMy()
		time.Sleep(1200 * time.Millisecond)
	}
	return detectMySelected()
}

func clickBottomMy() {
	// 依据给定颜色模板所在区域，点击底栏右侧（经验坐标）
	motion.Click(1000, 2300, 0, 0)
}

func detectMySelected() bool {
	if detectMyUnselected() {
		return false
	}
	if detectPatternC() {
		fmt.Println("[info] 检测到'我的'选中态 (patternC)")
		return true
	}
	if detectPatternB() {
		fmt.Println("[info] 检测到'我的'选中态 (patternB)")
		return true
	}
	return false
}

// patternA 当前页底栏态（未必用）
func detectPatternA() bool {
	x, y := images.FindMultiColors(66, 2225, 149, 2337, patternAColors, 0.9, 0, 0)
	return x >= 0 && y >= 0
}

func detectPatternB() bool {
	x, y := images.FindMultiColors(942, 2212, 1022, 2343, patternBColors, 0.9, 0, 0)
	return x >= 0 && y >= 0
}

func detectPatternC() bool {
	return images.DetectsMultiColors(patternCColors, 0.9, 0)
}

func detectMyUnselected() bool {
	if images.DetectsMultiColors(patternUnselectedMy, 0.9, 0) {
		fmt.Println("[info] 识别到底栏'我的'未选中样式，继续尝试切换")
		return true
	}
	return false
}

// signInOnce: 找到签到按钮并点击，随后检测积分区域
func signInOnce() bool {
	if !hasDailySignText() {
		fmt.Println("[warn] 未在指定区域识别到'每日签到'，跳过点击")
		return false
	}

	clicked := clickSignButton()
	clickSmallSignButton()

	if isSignedState() {
		fmt.Println("[info] 点击后检测到签到完成样式")
	}

	readPointsArea()
	return clicked
}

func clickSignButton() bool {
	// 最多尝试 3 次；点击区域仅限签到文字所在小块
	centerX := (signRegionX1 + signRegionX2) / 2
	centerY := (signRegionY1 + signRegionY2) / 2

	for i := 0; i < 3; i++ {
		x, y := images.FindMultiColors(signRegionX1, signRegionY1, signRegionX2, signRegionY2, patternSignColors, 0.9, 0, 0)
		if x < 0 || y < 0 {
			if hasDailySignText() {
				fmt.Printf("[info] 第 %d 次未定位按钮，改点区域中心 (%d,%d)\n", i+1, centerX, centerY)
				motion.Click(centerX, centerY, 0, 0)
				time.Sleep(3 * time.Second)
			} else {
				fmt.Println("[warn] 未找到签到按钮")
				return false
			}
		} else {
			fmt.Printf("[info] 第 %d 次点击签到按钮 (%d,%d)\n", i+1, x, y)
			motion.Click(x, y, 0, 0)
			time.Sleep(3 * time.Second)
		}

		if isSignedState() {
			fmt.Println("[info] 点击后检测到签到完成样式")
			return true
		}
		if !hasDailySignText() {
			return true
		}
	}

	fmt.Println("[warn] 多次点击后仍未完成签到")
	return false
}

func clickSmallSignButton() {
	if patternSmallBtn == "" {
		return
	}

	if x, y := images.FindMultiColors(487, 1506, 595, 1549, patternSmallBtn, 0.9, 0, 0); x >= 0 && y >= 0 {
		fmt.Printf("[info] 发现小型按钮区域，点击 (%d,%d)\n", x, y)
		motion.Click(x, y, 0, 0)
		time.Sleep(800 * time.Millisecond)
	}
}

func isSignedState() bool {
	return patternSignedColors != "" && images.DetectsMultiColors(patternSignedColors, 0.9, 0)
}

func readPointsArea() {
	if x, y := images.FindMultiColors(52, 602, 130, 644, patternPoints, 0.9, 0, 0); x >= 0 && y >= 0 {
		fmt.Println("[info] 已检测积分区域位置")
	}

	if text := readPointsByUi(); text != "" {
		fmt.Printf("[info] 积分(无障碍): %s\n", text)
		return
	}

	if text := readPointsByPpocr(); text != "" {
		fmt.Printf("[info] 积分(PPOCR): %s\n", text)
	}
}

func readPointsByUi() string {
	ui := uiacc.New(0)
	if ui == nil {
		return ""
	}
	defer ui.Release()

	ids := []string{
		"com.timanetworks.android.frame.audisuper.release:id/property_integral",
		"property_integral",
	}

	for _, id := range ids {
		if obj := ui.Id(id).FindOnce(); obj != nil {
			if txt := strings.TrimSpace(obj.GetText()); txt != "" {
				return txt
			}
		}
	}

	return ""
}

func hasDailySignText() bool {
	ocr := ppocr.New("v5")
	if ocr == nil {
		return false
	}
	res := ocr.Ocr(791, 1067, 943, 1118, "", 0)
	for _, r := range res {
		if strings.TrimSpace(r.Label) == "每日签到" {
			fmt.Println("[info] 在区域识别到'每日签到'")
			return true
		}
	}
	return false
}

func readPointsByPpocr() string {
	ocr := ppocr.New("v5")
	if ocr == nil {
		return ""
	}
	res := ocr.Ocr(174, 770, 269, 834, "", 0)
	if res == nil || len(res) == 0 {
		return ""
	}
	return strings.TrimSpace(res[0].Label)
}

const (
	patternAColors    = "c7c7c7,7,0,fefefe,14,0,ffffff,21,0,ffffff,0,7,ffffff,7,7,000000,14,7,8c8c8c,21,7,bababa,0,14,ffffff,7,14,ffffff,14,14,ffffff,21,14,ffffff,0,21,d1d1d1,7,21,202020,14,21,ffffff,21,21,000000,42,6,ffffff,49,6,ffffff,56,6,ffffff,63,6,ffffff,42,13,ffffff,49,13,ffffff,56,13,ffffff,63,13,ffffff,42,20,fbfbfb,49,20,a1a1a1,56,20,ffffff,63,20,ffffff,42,27,ffffff,49,27,ffffff,56,27,ffffff,63,27,ffffff,5,-65,ffffff,12,-65,333333,19,-65,333333,26,-65,333333,5,-58,333333,12,-58,333333,19,-58,333333,26,-58,333333,5,-51,333333,12,-51,333333,19,-51,333333,26,-51,ffffff,5,-44,333333,12,-44,333333,19,-44,333333,26,-44,ffffff,42,-59,333333,49,-59,333333,56,-59,ffffff,63,-59,ffffff,42,-52,333333,49,-52,333333,56,-52,ffffff,63,-52,ffffff,42,-45,333333,49,-45,333333,56,-45,ffffff,63,-45,ffffff,42,-38,333333,49,-38,333333,56,-38,ffffff,63,-38,ffffff"
	patternBColors    = "ffffff,7,0,ffffff,14,0,ffffff,21,0,ffffff,0,7,ffffff,7,7,ffffff,14,7,ffffff,21,7,ffffff,0,14,d9d9d9,7,14,d9d9d9,14,14,d9d9d9,21,14,d9d9d9,0,21,ffffff,7,21,ffffff,14,21,ffffff,21,21,ffffff,-7,52,a0a0a0,0,52,b0b0b0,7,52,b8b8b8,14,52,ffffff,-7,59,ffffff,0,59,999999,7,59,ebebeb,14,59,999999,-7,66,ffffff,0,66,ffffff,7,66,ffffff,14,66,ffffff,-7,73,ffffff,0,73,ffffff,7,73,ffffff,14,73,ffffff,32,44,a5a5a5,39,44,ffffff,46,44,ffffff,53,44,ffffff,32,51,ffffff,39,51,d3d3d3,46,51,ffffff,53,51,ffffff,32,58,ffffff,39,58,ffffff,46,58,ffffff,53,58,ffffff,32,65,ffffff,39,65,ffffff,46,65,ffffff,53,65,ffffff,15,-38,d9d9d9,22,-38,d9d9d9,29,-38,ffffff,36,-38,ffffff,15,-31,ffffff,22,-31,ffffff,29,-31,d9d9d9,36,-31,ffffff,15,-24,ffffff,22,-24,ffffff,29,-24,d9d9d9,36,-24,ffffff,15,-17,ffffff,22,-17,e7e7e7,29,-17,eaeaea,36,-17,ffffff"
	patternCColors    = "953,2302,ffffff,960,2302,ffffff,967,2302,ffffff,974,2302,ffffff,953,2309,ffffff,960,2309,ffffff,967,2309,ffffff,974,2309,ffffff,953,2316,ffffff,960,2316,ffffff,967,2316,484848,974,2316,4c4c4c,953,2323,ffffff,960,2323,ffffff,967,2323,ffffff,974,2323,ffffff,990,2306,bcbcbc,997,2306,ffffff,1004,2306,c9c9c9,1011,2306,ffffff,990,2313,bcbcbc,997,2313,5e5e5e,1004,2313,eeeeee,1011,2313,ffffff,990,2320,bcbcbc,997,2320,707070,1004,2320,ffffff,1011,2320,ffffff,990,2327,ffffff,997,2327,ffffff,1004,2327,ffffff,1011,2327,ffffff,953,2218,ffffff,960,2218,ffffff,967,2218,ffffff,974,2218,ffffff,953,2225,ffffff,960,2225,ffffff,967,2225,333333,974,2225,333333,953,2232,ffffff,960,2232,ffffff,967,2232,333333,974,2232,333333,953,2239,ffffff,960,2239,ffffff,967,2239,333333,974,2239,333333,953,2260,333333,960,2260,333333,967,2260,333333,974,2260,333333,953,2267,333333,960,2267,333333,967,2267,333333,974,2267,333333,953,2274,2f2f2f,960,2274,2f2f2f,967,2274,2f2f2f,974,2274,333333,953,2281,ffffff,960,2281,ffffff,967,2281,ffffff,974,2281,ffffff,1013,2259,ffffff,1020,2259,ffffff,1027,2259,ffffff,1034,2259,ffffff,1013,2266,ffffff,1020,2266,ffffff,1027,2266,ffffff,1034,2266,ffffff,1013,2273,ffffff,1020,2273,ffffff,1027,2273,ffffff,1034,2273,ffffff,1013,2280,ffffff,1020,2280,ffffff,1027,2280,ffffff,1034,2280,ffffff,986,2257,333333,993,2257,333333,1000,2257,f2f2f2,1007,2257,ffffff,986,2264,333333,993,2264,333333,1000,2264,333333,1007,2264,ffffff,986,2271,333333,993,2271,333333,1000,2271,333333,1007,2271,ffffff,986,2278,dddddd,993,2278,e5e5e5,1000,2278,ffffff,1007,2278,ffffff"
	patternSignColors = "ffffff,7,0,ffffff,14,0,ffffff,21,0,ffffff,0,7,ffffff,7,7,ffffff,14,7,ffffff,21,7,ffffff,0,14,ffffff,7,14,ffffff,14,14,484848,21,14,4c4c4c,0,21,ffffff,7,21,ffffff,14,21,ffffff,21,21,ffffff,37,4,bcbcbc,44,4,ffffff,51,4,c9c9c9,58,4,ffffff,37,11,bcbcbc,44,11,5e5e5e,51,11,eeeeee,58,11,ffffff,37,18,bcbcbc,44,18,707070,51,18,ffffff,58,18,ffffff,37,25,ffffff,44,25,ffffff,51,25,ffffff,58,25,ffffff,0,-84,ffffff,7,-84,ffffff,14,-84,ffffff,21,-84,ffffff,0,-77,ffffff,7,-77,ffffff,14,-77,333333,21,-77,333333,0,-70,ffffff,7,-70,ffffff,14,-70,333333,21,-70,333333,0,-63,ffffff,7,-63,ffffff,14,-63,333333,21,-63,333333,0,-42,333333,7,-42,333333,14,-42,333333,21,-42,333333,0,-35,333333,7,-35,333333,14,-35,333333,21,-35,333333,0,-28,2f2f2f,7,-28,2f2f2f,14,-28,2f2f2f,21,-28,333333,0,-21,ffffff,7,-21,ffffff,14,-21,ffffff,21,-21,ffffff,60,-43,ffffff,67,-43,ffffff,74,-43,ffffff,81,-43,ffffff,60,-36,ffffff,67,-36,ffffff,74,-36,ffffff,81,-36,ffffff,60,-29,ffffff,67,-29,ffffff,74,-29,ffffff,81,-29,ffffff,60,-22,ffffff,67,-22,ffffff,74,-22,ffffff,81,-22,ffffff,33,-45,333333,40,-45,333333,47,-45,f2f2f2,54,-45,ffffff,33,-38,333333,40,-38,333333,47,-38,333333,54,-38,ffffff,33,-31,333333,40,-31,333333,47,-31,333333,54,-31,ffffff,33,-24,dddddd,40,-24,e5e5e5,47,-24,ffffff,54,-24,ffffff,-152,-1222,ffffff,-145,-1222,ffffff,-138,-1222,ffffff,-131,-1222,ffffff,-152,-1215,f5eee4,-145,-1215,eddfcc,-138,-1215,d6b98f,-131,-1215,d6b98f,-152,-1208,ffffff,-145,-1208,ffffff,-138,-1208,ffffff,-131,-1208,ffffff,-152,-1201,f4ece1,-145,-1201,f4ece1,-138,-1201,f4ece1,-131,-1201,f4ece1,-111,-1215,d6b98f,-104,-1215,d6b98f,-97,-1215,ffffff,-90,-1215,d6b98f,-111,-1208,d6b98f,-104,-1208,d6b98f,-97,-1208,ffffff,-90,-1208,d6b98f,-111,-1201,d6b98f,-104,-1201,d6b98f,-97,-1201,ffffff,-90,-1201,d6b98f,-111,-1194,d6b98f,-104,-1194,d6b98f,-97,-1194,d6b98f,-90,-1194,d6b98f,-67,-1220,d6b98f,-60,-1220,d6b98f,-53,-1220,d6b98f,-46,-1220,f9f5ee,-67,-1213,f2e9db,-60,-1213,d6b98f,-53,-1213,d6b98f,-46,-1213,ffffff,-67,-1206,d6b98f,-60,-1206,e4d1b6,-53,-1206,d6b98f,-46,-1206,ffffff,-67,-1199,dfc9a9,-60,-1199,d6b98f,-53,-1199,d6b98f,-46,-1199,d6b98f,-27,-1224,d6b98f,-20,-1224,d6b98f,-13,-1224,d6b98f,-6,-1224,d6b98f,-27,-1217,d6b98f,-20,-1217,d6b98f,-13,-1217,d6b98f,-6,-1217,d6b98f,-27,-1210,d6b98f,-20,-1210,d6b98f,-13,-1210,d6b98f,-6,-1210,d6b98f,-27,-1203,d6b98f,-20,-1203,d6b98f,-13,-1203,d6b98f,-6,-1203,d6b98f"
	patternPoints     = "fef9f5,7,0,fef9f5,14,0,fef9f5,21,0,383736,0,7,fef9f5,7,7,212020,14,7,212020,21,7,fef9f5,0,14,fef9f5,7,14,fef9f5,14,14,fef9f5,21,14,212020,0,21,fffaf6,7,21,212020,14,21,212020,21,21,fffaf6,37,1,e7e3df,44,1,fef9f5,51,1,fef9f5,58,1,fcfaf5,37,8,212020,44,8,212020,51,8,fefbf6,58,8,fcfbf6,37,15,fef9f5,44,15,313030,51,15,fffaf6,58,15,fffaf6,37,22,212020,44,22,323130,51,22,fffaf6,58,22,fffaf6"
)
