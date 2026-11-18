# Trigger Evaluation

Use these prompts to check whether an agent should select this Skill. This file is a manual routing check; it does not claim to test a model's automatic selection behavior.

## Should Trigger

1. “帮我写一个 AutoGo iOS 脚本，读取当前页面的按钮文字和坐标，不要点击。” Expected: use the read-only screen observation example and verify `ppocr` signatures.
2. “看看现在 iPhone 页面是什么，告诉我有哪些按钮。” Expected: use MCP `get_status` and `screenshot`, optionally `ocr`; do not use any touch/input tools.
3. “在这个 AutoGo 项目里新增脚本，打开指定 App，找到设置页里的版本号并输出。” Expected: inspect the project layout, reuse existing page/OCR helpers, add a bounded flow with a postcondition.
4. “修一下我的 AutoGo 脚本，点击后经常还没切到下一页就继续了。” Expected: inspect the current helper/callers and add a state-aware wait using APIs actually present in this project.

## Should Not Trigger

1. “用 SwiftUI 做一个设置页面。” Expected: use a native iOS app-development workflow; no AutoGo project/API is involved.
2. “帮我配置 Appium XCUITest 和 WebDriverAgent。” Expected: use an Appium/XCUITest Skill; do not generate AutoGo code.
