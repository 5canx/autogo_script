# Project Crosswalk

This Skill was researched against two different local trees. Treat them as separate project layouts; never assume files or helpers from one exist in the other.

## `iphone_autogo`

Location supplied in the request: `/Users/scan/Desktop/Code/autogo代码/iphone_autogo/`.

- Minimal Go module `app`, Go 1.26, with `replace github.com/Dasongzi1366/AutoGo => ./AutoGo`.
- Root `main.go` is a standalone API smoke/demo entry point. `resources/META-INF/ios.toml` sets `appPackage="com.auto.go"`, `appName="AutoGo"`, and `autoRun=false`.
- The current root `main.go` performs swipe and click actions and saves a screenshot. Do not execute it against a live device as a harmless compile check; compile only unless the user explicitly requests that device test.
- `AutoGo/` has package surfaces for `app`, `device`, `files`, `uiacc`, `https`, `ime`, `motion`, `images`, `opencv`, `dotocr`, `ppocr`, `yolo`, `storages`, `system`, `utils`, `apkctl`, `console`, `hud`, and `imgui`.
- Most device-facing functions return zero values or do nothing in the checked-in Go source. Native iOS dylibs are included under `resources/libs/ios`, but their source and complete build/deploy workflow are not included. Do not report a successful local build as proof that an API worked on a phone.
- The directory is part of an untracked parent-repository tree and includes large IPA extraction/build artifacts. Keep Skill files isolated under `.agents/skills/autogo-ios`; do not edit or repackage those artifacts as part of script generation.
- Existing `dump_file.py`/`dump_main.py`/`dump_sections.py` and paired Frida JS scripts attach to `com.SafetyTraining` and export executable bytes, segments, or section metadata into `ipa_out/`. `netcap.py` attaches to that app and expects `/tmp/netcap.js`, then prints request URLs and bodies. These are specialized analysis prototypes, not general AutoGo APIs; require explicit authorization for that app/data and keep captured traffic/binaries private. Do not automatically run them or copy their target-specific behavior into a new Skill example.

## `awesomeProject1`

Location supplied by the user: `/Users/scan/GolandProjects/awesomeProject1/`.

- Go module `app` also replaces `github.com/Dasongzi1366/AutoGo` with `./AutoGo`; its checked-in API layer is a local stub for syntax/API checking, while real behavior requires the AutoGo iOS host.
- `main.go` is the sole task registry. Add business tasks under `appscript/` (native apps) or `minicode/` (WeChat mini-programs), expose `RunXxx()` functions, then register them in `main.go`. Shared helpers belong in `common/` and should be reused.
- `common/wxinit.go` supplies OCR helpers, `DumpScreen`, coordinate scaling (`Sx`, `Sy`, `Rect`), mini-program navigation, and node-tree helpers. `common/pagefeat.go` supplies page-state matching and bounded wait helpers. `common/runner.go` adds task timeout/stall handling and stop events.
- `appscript/debug_dump.go` includes `RunScreenRecorder` (manual device navigation with OCR-based page-change detection and full screen text/center coordinates in logs) and a pre-programmed flow dump. The recorder detects page changes; it does not record the user's touch coordinates.
- `mcp-iphone-mirror/` is a separate macOS MCP server. Its tools include `get_status`, `screenshot` (device-sized PNG returned as an image), `ocr` (local PaddleOCR boxes), `click`, `swipe`, `scroll`, `type_text`, `set_clipboard`, and `pin`. Setup requires an already-open iPhone Mirroring window and macOS Accessibility/Screen Recording permissions. The MCP README documents default 1170x2532 coordinates while AutoGo scripts use 1284x2778.
- Repository `CLAUDE.md` marks the mirror-driven coordinate interaction workflow deprecated and says to keep one AutoGo coordinate system. This Skill uses mirror MCP for read-only live screenshots/OCR when available, but uses AutoGo `DumpScreen`/OCR for runtime behavior and never initiates live clicks just to generate a script.
- Known project limitation: `CLAUDE.md` documents `uiacc` crashes on WeChat mini-program pages and some edit menus. Reuse its OCR fallback there; use `uiacc` only on verified native pages.
- Existing scripts include JD/Taobao native-app tasks, mini-program scripts, control-panel examples, coordinate helpers, and the MCP server. Inspect relevant files before writing so the Skill extends rather than duplicates existing flows.

## Choosing a Screen-Reading Route

| Need | First choice | What it returns | Constraint |
|---|---|---|---|
| Generate a script that can recognize visible text on-device | `ppocr` through an existing helper or direct AutoGo call | Labels, confidence, box/center coordinates in logs | Requires AutoGo host/model support; local API stubs return no real detections. |
| Capture pixels from the AutoGo device context | `images.CaptureScreen` | `*image.NRGBA` screenshot | Does not itself deliver the image to the AI; the script must save or otherwise return it through an existing authorized workflow. |
| Dump current text and coordinates from `awesomeProject1` | `common.DumpScreen` / `RunScreenRecorder` | Ordered OCR text and centers in host logs | Runs in the AutoGo host; recorder is manual page observation and has a task timeout. |
| Let the agent visually inspect a live iPhone mirror | Configured read-only `iphone-mirror` MCP | PNG image and optionally OCR box JSON | Requires macOS permissions; screenshot coordinates differ from the AutoGo script baseline. |

The project includes Codex and Claude Code configuration for the reference MCP server; details and permission limits are in [MCP setup](mcp-setup.md). This only starts a server; it does not grant macOS Accessibility/Screen Recording permission or connect the phone.

## Verified Read-Only Example

`examples/screen-observe/main.go` shows `images.CaptureScreen` plus `ppocr` OCR without tapping or changing app state. Use `common.DumpScreen()` instead when the active project already provides that helper.
