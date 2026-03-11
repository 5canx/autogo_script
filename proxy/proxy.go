package proxy

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/motion"
	"github.com/Dasongzi1366/AutoGo/uiacc"
)

// Run 执行完整流程：导航 -> 检测/添加目标应用 -> 配置域名（可选）-> 返回首页点击抓包按钮。
// 返回值表示流程是否成功。
func Run(target, host string) bool {
	if target == "" {
		target = "京东"
	}

	fmt.Printf("[start] 检测是否存在 %s\n", target)
	ui := uiacc.New(0)
	if ui == nil {
		fmt.Println("[error] 无法创建无障碍实例")
		return false
	}
	defer ui.Release()

	preNavigate(ui)

	if ensureApp(ui, target) {
		if host != "" {
			fmt.Printf("[start] 配置抓包域名: %s\n", host)
			if !configureHost(ui, host) {
				fmt.Println("[warn] 配置抓包域名失败")
				return false
			}
		}
		backToContent(ui)
		clickHomeService(ui)
		return true
	}

	return false
}

// ensureApp 收集列表并在缺失时尝试添加/点击目标应用
func ensureApp(ui *uiacc.Uiacc, name string) bool {
	names := collectAppNames(ui)
	fmt.Printf("[info] 收集到 %d 个应用\n", len(names))
	if containsName(names, name) {
		fmt.Printf("[info] 已存在 %s，无需添加\n", name)
		return true
	}

	fmt.Printf("[info] 未找到 %s，开始添加\n", name)
	if addApp(ui, name) {
		fmt.Printf("[info] 添加/点击 %s 成功\n", name)
		return true
	}

	fmt.Printf("[warn] 未能找到或添加 %s\n", name)
	return false
}

func configureHost(ui *uiacc.Uiacc, host string) bool {
	// 尝试先返回上一级，再开始查找
	fmt.Println("[nav] 先执行返回，确保在主设置页")
	motion.Back(0)
	time.Sleep(500 * time.Millisecond)

	// 进入抓包设置，尝试多次滑动查找
	setting := findGrabSetting(ui)
	if setting == nil {
		fmt.Println("[warn] 未找到 抓包设置")
		return false
	}
	if !clickNodeWithAncestors(setting, 3) {
		fmt.Println("[warn] 无法点击 抓包设置")
		return false
	}
	fmt.Println("[info] 已进入抓包设置")
	time.Sleep(500 * time.Millisecond)

	// 进入 Host 白名单，支持滑动重试
	if !openHostWhitelist(ui) {
		fmt.Println("[warn] 未能进入 Host 白名单")
		return false
	}
	fmt.Println("[info] 已进入 Host 白名单")
	time.Sleep(500 * time.Millisecond)

	// 如已存在域名则返回
	if hasDomain(ui, host) {
		fmt.Printf("[info] 已存在域名 %s，无需添加\n", host)
		return true
	}

	// 添加域名
	if addDomain(ui, host) {
		return true
	}

	fmt.Printf("[warn] 未能添加域名 %s\n", host)
	return false
}

func hasDomain(ui *uiacc.Uiacc, host string) bool {
	for _, o := range ui.Id("text1").Find() {
		if t := o.GetText(); t != "" && strings.Contains(t, host) {
			fmt.Printf("[info] 在列表中找到域名: %s\n", t)
			return true
		}
	}
	for _, o := range ui.Id("content").Find() {
		if t := o.GetText(); t != "" && strings.Contains(t, host) {
			fmt.Printf("[info] 在 content 节点中找到域名: %s\n", t)
			return true
		}
	}
	fmt.Printf("[info] 未找到域名 %s，准备添加\n", host)
	return false
}

func addDomain(ui *uiacc.Uiacc, host string) bool {
	// 先点击右上角添加按钮
	if addBtn := ui.Id("menu_add").Clickable(true).FindOnce(); addBtn != nil {
		if clickNodeWithAncestors(addBtn, 2) {
			fmt.Println("[info] 已点击 menu_add 打开添加对话框")
			time.Sleep(400 * time.Millisecond)
		} else {
			fmt.Println("[warn] 点击 menu_add 失败")
		}
	} else {
		fmt.Println("[warn] 未找到 menu_add，尝试直接寻找输入框")
	}

	input := ui.Id("settings_add_host").Editable(true).FindOnce()
	if input == nil {
		input = ui.ClassName("android.widget.EditText").FindOnce()
	}
	if input == nil {
		fmt.Println("[warn] 未找到添加域名的输入框")
		return false
	}
	if !input.SetText(host) {
		fmt.Println("[warn] 设置域名文本失败")
		return false
	}
	fmt.Printf("[info] 已填入域名: %s\n", host)
	time.Sleep(300 * time.Millisecond)

	if btn := ui.Id("menu_complete").FindOnce(); btn != nil {
		if clickNodeWithAncestors(btn, 2) {
			fmt.Println("[info] 已点击 menu_complete")
			time.Sleep(500 * time.Millisecond)
			return hasDomain(ui, host)
		}
	}
	for _, label := range []string{"完成", "确定", "添加", "保存", "OK"} {
		if b := ui.Text(label).Clickable(true).FindOnce(); b != nil {
			if clickNodeWithAncestors(b, 2) {
				fmt.Printf("[info] 已点击确认按钮: %s\n", label)
				time.Sleep(500 * time.Millisecond)
				return hasDomain(ui, host)
			}
		}
	}

	fmt.Println("[warn] 未找到确认按钮，添加可能未提交")
	return hasDomain(ui, host)
}

func findGrabSetting(ui *uiacc.Uiacc) *uiacc.UiObject {
	for i := 0; i < 8; i++ {
		setting := ui.PackageName("com.guoshi.httpcanary").Text("抓包设置").ClassName("android.widget.TextView").FindOnce()
		if setting != nil {
			return setting
		}
		fmt.Println("[info] 未找到 抓包设置，滑动重试")
		motion.Swipe(540, 1600, 540, 500, 700, 0, 0)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func openHostWhitelist(ui *uiacc.Uiacc) bool {
	for i := 0; i < 6; i++ {
		if wl := ui.Id("settings_host_whitelist").FindOnce(); wl != nil {
			if clickNodeWithAncestors(wl, 3) {
				return true
			}
		}
		fmt.Println("[info] 未找到 Host 白名单，滑动重试")
		motion.Swipe(540, 1600, 540, 500, 700, 0, 0)
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func preNavigate(ui *uiacc.Uiacc) {
	const (
		x1       = 70
		y1       = 2263
		x2       = 283
		y2       = 2324
		sim      = 0.9
		dir      = 0
		template = "ffffff,7,0,ffffff,14,0,757575,21,0,757575,0,7,ffffff,7,7,ffffff,14,7,a0a0a0,21,7,757575,0,14,ffffff,7,14,ffffff,14,14,757575,21,14,757575,0,21,757575,7,21,757575,14,21,757575,21,21,757575,172,10,f1f1f1,179,10,ffffff,186,10,333333,193,10,333333,172,17,ffffff,179,17,ffffff,186,17,333333,193,17,ffffff,172,24,333333,179,24,ffffff,186,24,333333,193,24,333333,172,31,ffffff,179,31,ffffff,186,31,ffffff,193,31,ffffff,139,-10,ffffff,146,-10,d2d2d2,153,-10,ffffff,160,-10,333333,139,-3,ffffff,146,-3,646464,153,-3,ffffff,160,-3,ffffff,139,4,333333,146,4,333333,153,4,848484,160,4,f2f2f2,139,11,ffffff,146,11,333333,153,11,ffffff,160,11,ffffff"
	)

	fmt.Println("[nav] 点击第一个可点击的 ImageButton")
	if obj := ui.ClassName("android.widget.ImageButton").Clickable(true).FindOnce(); obj != nil {
		if obj.Click() {
			fmt.Println("[nav] ImageButton 点击成功")
			time.Sleep(500 * time.Millisecond)
		} else {
			fmt.Println("[nav-warn] ImageButton 点击失败")
		}
	} else {
		fmt.Println("[nav-warn] 未找到可点击的 ImageButton")
	}

	fmt.Println("[nav] 多点找色点击")
	if x, y := images.FindMultiColors(x1, y1, x2, y2, template, sim, dir, 0); x >= 0 && y >= 0 {
		motion.Click(x, y, 0, 0)
		fmt.Printf("[nav] 找色命中并点击 (%d,%d)\n", x, y)
		time.Sleep(500 * time.Millisecond)
	} else {
		fmt.Println("[nav-warn] 找色未命中，跳过")
	}

	fmt.Println("[nav] 点击第一个可点击的 LinearLayout")
	if obj := ui.ClassName("android.widget.LinearLayout").Clickable(true).FindOnce(); obj != nil {
		if obj.Click() {
			fmt.Println("[nav] LinearLayout 点击成功")
			time.Sleep(500 * time.Millisecond)
		} else {
			fmt.Println("[nav-warn] LinearLayout 点击失败")
		}
	} else {
		fmt.Println("[nav-warn] 未找到可点击的 LinearLayout")
	}

	fmt.Println("[nav] 向上滑动进入列表")
	motion.Swipe(540, 1800, 540, 600, 400, 0, 0)
	time.Sleep(600 * time.Millisecond)
}

func backToContent(ui *uiacc.Uiacc) {
	for i := 0; i < 3; i++ {
		if c := ui.Id("content").FindOnce(); c != nil {
			if clickNodeWithAncestors(c, 3) {
				fmt.Println("[info] 已回到 content 界面")
			} else if c.ClickCenter() {
				fmt.Println("[info] 已回到 content 界面（中心点击兜底）")
			} else {
				fmt.Println("[warn] 点击 content 节点失败")
			}
			return
		}
		fmt.Println("[info] 未找到 content，按返回键重试")
		motion.Back(0)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("[warn] 尝试返回 content 界面失败")
}

func clickHomeService(ui *uiacc.Uiacc) {
	for i := 0; i < 4; i++ {
		if btn := ui.Id("home_service_button").Clickable(true).FindOnce(); btn != nil {
			if clickNodeWithAncestors(btn, 2) {
				fmt.Println("[info] 已点击 home_service_button")
				return
			}
			if btn.ClickCenter() {
				fmt.Println("[info] 已点击 home_service_button（中心点击兜底）")
				return
			}
			fmt.Println("[warn] 点击 home_service_button 失败")
			return
		}
		fmt.Println("[info] 未找到 home_service_button，按返回键重试")
		motion.Back(0)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("[warn] 尝试多次仍未找到 home_service_button")
}

func collectAppNames(ui *uiacc.Uiacc) []string {
	seen := make(map[string]struct{})
	var result []string

	list := findList(ui)
	for i := 0; i < 8; i++ {
		before := len(result)
		grabNames(ui, seen, &result)
		added := len(result) - before
		fmt.Printf("[info] 第 %d 轮抓取新增 %d，累计 %d\n", i+1, added, len(result))
		if list == nil {
			fmt.Println("[warn] 未找到列表节点，停止抓取")
			break
		}
		children := list.GetChildren()
		// 如果可见子项不足 8 个，就不再滑动
		if len(children) < 8 {
			fmt.Printf("[info] 可见子项 %d (<8)，停止滑动\n", len(children))
			break
		}
		fmt.Println("[info] 向上滑动以继续收集")
		motion.Swipe(540, 1600, 540, 500, 700, 0, 0)
		time.Sleep(700 * time.Millisecond)
		list = findList(ui)
	}
	// 收尾再抓一次
	grabNames(ui, seen, &result)
	fmt.Printf("[info] 收尾抓取完成，累计 %d\n", len(result))
	return result
}

func grabNames(ui *uiacc.Uiacc, seen map[string]struct{}, result *[]string) {
	objs := ui.Id("app_name").Find()
	for _, o := range objs {
		t := o.GetText()
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		*result = append(*result, t)
	}
}

func findList(ui *uiacc.Uiacc) *uiacc.UiObject {
	if l := ui.Id("list").FindOnce(); l != nil {
		return l
	}
	if l := ui.ClassName("android.widget.ListView").FindOnce(); l != nil {
		return l
	}
	return ui.ClassName("androidx.recyclerview.widget.RecyclerView").FindOnce()
}

func containsName(names []string, target string) bool {
	for _, n := range names {
		if n == target {
			fmt.Printf("[info] 列表中找到 %s\n", target)
			return true
		}
	}
	fmt.Printf("[info] 列表中未找到 %s\n", target)
	return false
}

func addApp(ui *uiacc.Uiacc, name string) bool {
	btn := ui.Id("menu_add").Clickable(true).FindOnce()
	if btn == nil || !btn.Click() {
		fmt.Println("[warn] 未找到或无法点击 menu_add")
		return false
	}
	fmt.Println("[info] 已点击 menu_add")
	time.Sleep(500 * time.Millisecond)

	input := ui.Id("search_src_text").Text("请输入应用名或包名").Focusable(true).Clickable(true).FindOnce()
	if input == nil {
		fmt.Println("[warn] 未找到搜索输入框")
		return false
	}

	ok := input.SetText(name)
	time.Sleep(300 * time.Millisecond)
	if !ok {
		fmt.Println("[warn] 设置搜索文本失败")
		return false
	}
	fmt.Printf("[info] 已输入 %s\n", name)

	// 等待搜索结果刷新
	time.Sleep(1200 * time.Millisecond)

	// 调试：打印 app_name 节点与 FrameLayout 子节点
	dumpAppNameNodes(ui)
	dumpFrameLayoutChildren(ui)

	// 直接查找并点击 Id=app_name 且文本为目标
	if obj := ui.Id("app_name").Text(name).FindOnce(); obj != nil {
		if clickNodeWithAncestors(obj, 3) {
			fmt.Printf("[info] 已点击搜索结果 %s (向上寻找可点击父节点)\n", name)
			return true
		}
		fmt.Printf("[warn] 点击搜索结果 %s 失败\n", name)
	}

	// 尝试在 list 节点下递归查找并点击目标
	if clickAppInList(ui, name) {
		return true
	}

	fmt.Printf("[warn] 搜索结果未找到 %s\n", name)
	return false
}

func clickAppInList(ui *uiacc.Uiacc, name string) bool {
	list := findList(ui)
	if list == nil {
		fmt.Println("[warn] 未找到 list 节点")
		return false
	}
	if findAndClickApp(list, name) {
		fmt.Printf("[info] 在 list 中点击了 %s\n", name)
		return true
	}
	return false
}

func findAndClickApp(node *uiacc.UiObject, name string) bool {
	if node == nil {
		return false
	}
	if node.GetText() == name || node.GetDesc() == name {
		return clickNodeWithAncestors(node, 3)
	}
	if node.GetId() == "app_name" && node.GetText() == name {
		return clickNodeWithAncestors(node, 3)
	}
	children := node.GetChildren()
	for _, c := range children {
		if findAndClickApp(c, name) {
			return true
		}
	}
	return false
}

func dumpAppNameNodes(ui *uiacc.Uiacc) {
	objs := ui.Id("app_name").Find()
	fmt.Printf("[debug] app_name 节点数: %d\n", len(objs))
	for i, o := range objs {
		fmt.Printf("[debug] app_name[%d]: text=%s desc=%s id=%s class=%s\n", i, o.GetText(), o.GetDesc(), o.GetId(), o.GetClassName())
	}
}

func dumpFrameLayoutChildren(ui *uiacc.Uiacc) {
	frames := ui.ClassName("android.widget.FrameLayout").Find()
	fmt.Printf("[debug] FrameLayout 数量: %d\n", len(frames))
	for idx, f := range frames {
		children := f.GetChildren()
		fmt.Printf("[debug] FrameLayout[%d] 子节点数: %d\n", idx, len(children))
		for j, c := range children {
			fmt.Printf("[debug]  - child[%d]: text=%s desc=%s id=%s class=%s\n", j, c.GetText(), c.GetDesc(), c.GetId(), c.GetClassName())
		}
	}
}

func clickNodeWithAncestors(node *uiacc.UiObject, maxUp int) bool {
	cur := node
	for i := 0; i <= maxUp && cur != nil; i++ {
		if cur.GetClickable() {
			if cur.Click() {
				return true
			}
		}
		cur = cur.GetParent()
	}
	// 兜底点击中心
	cur = node
	if cur != nil {
		return cur.ClickCenter()
	}
	return false
}
