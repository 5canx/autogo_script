---
name: autogo-ios
description: Use when the user asks in natural language to create, change, debug, or inspect an iOS automation script using AutoGo; read iPhone screen text, screenshots, or UI elements; calibrate screen coordinates; or connect an existing AutoGo task to a project. Prefer the repository's existing helpers and verify every API against local source and the linked AutoGo docs.
---

# AutoGo iOS Script Builder

Turn a user's goal into a small, testable Go script that follows the AutoGo APIs and the current repository's conventions. The user should be able to describe the result they want; ask a question only when the target app, expected outcome, or a consequential action cannot be inferred from the project.

## Start Here

1. Identify the active project root, `go.mod`, `main.go`, existing `CLAUDE.md` or `AGENTS.md`, and relevant helpers. Do not assume the current directory is the project root.
2. Read [project crosswalk](references/project-crosswalk.md) to determine whether this is the minimal `iphone_autogo` template or the richer `awesomeProject1` layout.
3. Read [official API map](references/autogo-ios-api.md) for the relevant packages, parameters, platform limits, and source page. Read [open-source Skill review](references/skill-research.md) only when changing this Skill or comparing automation stacks.
4. Search existing scripts and callers before adding behavior. Reuse existing task registration, OCR, page detection, coordinate conversion, logging, wait, and stop helpers.

## When Not to Use

- For a Swift/SwiftUI app feature or native XCUITest suite with no AutoGo code, use the corresponding iOS development/testing workflow.
- For an Appium/WebDriver setup request, use an Appium Skill. Do not translate Appium capabilities into AutoGo APIs.
- For generic iOS reverse engineering or security testing unrelated to writing an authorized AutoGo script, use an applicable security-analysis workflow instead.

## Build a Script From a Request

1. Convert the request into observable steps: starting app/page, target control or text, action, expected next state, timeout, and stop condition. Keep the requested scope; do not add unrelated navigation or data collection.
2. Choose the strongest locator that is valid on the current page:
   - Reuse an existing semantic helper or page-state helper first.
   - Use `uiacc` for a native accessibility tree only where the project has verified it is safe.
   - Use OCR for visible text and coordinates when the tree is unavailable or unsafe.
   - Use colors or image templates only when text and accessibility data are insufficient.
   - Use raw coordinates only as a documented, device-calibrated fallback.
3. Make the flow state-aware: observe the current page, do one action, wait for a visible or semantic postcondition, then continue. Prefer bounded condition polling to long fixed sleeps. Log each meaningful transition and failure reason.
4. Keep the script in the project's existing package and task structure. In `awesomeProject1`, business tasks live in `appscript/` or `minicode/`, and `main.go` is the only task registry. In the minimal template, follow its existing entry point rather than inventing a second runner.
5. Run `gofmt` and the narrowest useful Go checks. A local build against this repository's stub AutoGo module proves syntax and API-shape compatibility only; report real iPhone execution as unverified unless it was run in the AutoGo host on an authorized device.

## Read Screen Information

The project configures `iphone-mirror` as a read-only MCP for current-page observation. When the user asks about the live screen, call `get_status`, then `screenshot`; use `ocr` for labels and centers. If the MCP is unavailable, ask for an authorized capture or use an AutoGo dump. A Skill does not grant OS permissions by itself.

- For an on-device text/coordinate dump, prefer `common.DumpScreen(...)` when the current project provides it.
- For direct AutoGo APIs, use `images.CaptureScreen(...)` for pixels and `ppocr.New("v5").Ocr(...)` for recognized labels and centers. Save or return only the minimum evidence needed; screenshots and OCR can contain personal information.
- For development-time visual observation, use the configured `iphone-mirror` MCP. The reference project marks its coordinate-driven interaction loop deprecated, so use it for read-only screenshots/OCR, not as the production runner.
- Never mix coordinate spaces. The reference AutoGo scripts use a 1284x2778 base with `common.Sx/Sy`; the mirror MCP defaults to 1170x2532. Convert and verify coordinates before producing device actions.

## Safety and Verification

- Work only with apps/devices the user is authorized to test. Do not bypass login, CAPTCHA, anti-abuse checks, rate limits, paywalls, or other access controls.
- Do not put credentials, session tokens, personal records, or raw screen dumps in source control or logs. Use placeholders or the project's existing secure configuration.
- Do not automate purchases, sending messages, deleting data, account changes, or other consequential actions without explicit user intent. Prefer read-only observation first.
- Do not execute device operations during code generation unless the user requested a live test and the target device/tool is clearly identified.
- Report separately what was checked locally, what requires the AutoGo host, and what was actually observed on-device.

## Reference Files

- [Official iOS APIs and runtime](references/autogo-ios-api.md)
- [Project layouts and reusable code](references/project-crosswalk.md)
- [Open-source Skill research](references/skill-research.md)
- [One-time screen MCP setup](references/mcp-setup.md)
- [Read-only screen observation example](examples/screen-observe/main.go)
- [Skill validator](scripts/validate_skill.py)
- [Trigger evaluation prompts](tests/trigger-cases.md)
