# AutoGo iOS API and Runtime Map

Research snapshot: 2026-11-18. AutoGo's docs are a Docsify SPA. The public iOS Markdown routes were read directly after resolving the sidebar; the site currently lists 18 iOS API pages. Re-check the linked docs and changelog when behavior or packaging requirements may have changed. The changelog page's latest entry in this snapshot is dated 2026-09-29.

## Runtime Model

AutoGo is a Go automation framework packaged and run through its IDE plugin/host workflow. For iOS, the project is initialized for iOS, connected and deployed through that workflow, and delivered as an IPA. This is not an Appium WebDriver/XCUITest runtime and does not imply that an arbitrary Mac can control an iPhone over USB using these Go packages.

The iOS API calls execute in the AutoGo host environment. The local `AutoGo/` module in these projects contains Go API-shaped stubs; native `.dylib` resources are present, but the bridge/runtime source and full host build pipeline are not. A local Go build checks source/API shape, not device effects. In `iphone_autogo/resources/META-INF/ios.toml`, the observed fields are `appPackage`, `appName`, and `autoRun`; these are project packaging metadata, not Go API parameters.

Communication surfaces documented by the project:

- `https` performs ordinary HTTP GET/POST/multipart requests from a script.
- `apkctl` bridges Go to JavaScript inside the host app. A `contextID` preserves JS scope; JS receives `objc`, `app`, `scriptUi`, `go`, and `print`; `go.send(...)` calls the Go callback registered by `SetCallback`.
- `apkctl.RegEvent` listens for floating-ball pause/resume/stop events. The docs explicitly say HTTP `/task` start/stop/pause/resume does not emit these events.
- IDE device connection, remote debugging, IPA installation, and screenshot service details belong to the plugin/host workflow. The changelog notes remote debugging and an iOS screenshot port change from 8989 to 8820 in version 1.16.0; do not infer a general HTTP control API from that port note.
- The changelog added `app.Install` support for `.deb` on jailbroken devices and `app.GetList` visibility of TrollStore-installed apps in 1.13.11 (2026-05-14). Earlier iOS 1.10.0 notes describe initial TrollStore installation constraints, including iOS 14.0-17.0 at that time; treat that as historical, not a current supported-version guarantee. In 1.16.0 (2026-07-11), TrollStore script installation no longer required enabling Developer Mode. Confirm current plugin/host requirements before giving installation instructions.
- Version 1.14.4 (2026-06-12) added remote debugging and required a newer AutoGo IDE plugin; version 1.17.8 (2026-09-28) says the remote compile port changed but does not give the replacement port. Do not hard-code an inferred port.
- The separate `mcp-iphone-mirror` server in `awesomeProject1` controls a visible macOS iPhone Mirroring window using MCP stdio, RobotGo/AppleScript, and optionally local PaddleOCR. It is a development aid, not the AutoGo device runtime.

## iOS Package Index

| Package | Documented capability and notable parameters |
|---|---|
| `app` | Query/launch/stop apps by Bundle ID; enumerate apps; read bundle/data paths, name/version/icon; install IPA by path; clear/uninstall; open URL/scheme. `AppInfo` includes `PackageName`, `AppName`, `VersionName`, `VersionCode`, `IsSystemApp`. |
| `device` | Model/release/serial fields; display width/height/scale/rotation; battery/status; screen on/unlocked; brightness/memory; VPN/Wi-Fi state; wake/keep-awake; local IP/Wi-Fi MAC; reboot. |
| `files` | File/dir checks, create/read/write/append, paths, directory listing/removal, plist read/write/get/set. Use project-local paths and avoid logging private file contents. |
| `uiacc` | Chainable selectors for text/description/id/class/package and bounds; clickable/selected/enabled/scrollable/editable/checked/password/index; find/wait/click; `UiObject` exposes bounds and properties. `Rect` coordinates are physical pixels. |
| `https` | `Get(url, timeoutMs)`, `Post(url, data, headers, timeoutMs)`, `PostMultipart(url, fileName, fileData, timeoutMs)`. Timeout `0` means no configured timeout; `Post` defaults Content-Type to JSON if absent. |
| `ime` | `InputText`, `GetClipText`, `SetClipText`. |
| `motion` | Touch down/move/up, click, long-click, swipe/swipe2; Home/Recents/volume/key actions. Gesture duration is milliseconds; finger ID is 0-9. Confirm signatures in the local module before calling. |
| `images` | Capture screen/region, pixel/color matching, multi-color matching, image read/write/base64/bytes, clip/resize/rotate/grayscale/threshold transforms. For `CaptureScreen`, `x2`/`y2` equal to 0 select max width/height. Similarity is documented as 0.1-1.0 for color matching. |
| `opencv` | Template matching (`FindImage`, image-based variants, all matches); transparency uses the template's top-left RGB as transparent color when enabled; similarity 0.1-1.0. |
| `dotocr` | Dot-matrix OCR/find-string using threshold, similarity, JSON/text mode, and named dictionaries. Screen functions use region coordinates; image/base64/path variants are also documented. |
| `ppocr` | Paddle OCR `New("v2"|"v5")`, screen/image/base64/path OCR. Results include label, confidence, box and center coordinates; optional color filter. Close instances after use. |
| `yolo` | NCNN object detection (`v5`/`v8`) with CPU thread count, parameter/model paths and comma-separated labels; screen/image/base64/path inputs. |
| `storages` | String key/value operations grouped by table: get/put/remove/contains/get-all/clear. |
| `system` | Process PID, CPU/memory usage, restart self. |
| `utils` | Toast/alert/input alert, execute binary, random/sleep and scalar conversions. |
| `apkctl` | Persistent JS `Eval(contextID, script)`, JS-to-Go `SetCallback`, lifecycle `RegEvent`; Objective-C/UIKit calls through the documented `objc` bridge and script UI lifecycle. |
| `console`, `hud` | Host floating console and text HUD creation/configuration/show/hide/destroy. |
| `imgui` | Immediate-mode UI drawing APIs; see the official iOS page before using generated bindings. |

The official iOS sidebar does not list Android-only `media`, `rhino`, `plugin`, or `vdisplay`; do not silently import those capabilities into an iOS script.

## Official-Docs-to-Source Crosswalk

The paths below are in the `iphone_autogo/AutoGo/` module unless the project crosswalk says otherwise. The module files expose the package API but mostly contain placeholders rather than device implementations.

| Official page | Local source | Assessment |
|---|---|---|
| `app` | `AutoGo/app/app.go` | Same API family; implementations return empty/default values locally. |
| `device` | `AutoGo/device/device.go` | Same API family; device access is stubbed locally. |
| `files` | `AutoGo/files/files.go` | Path, file, and plist API surface; local methods are stubs. |
| `uiacc` | `AutoGo/uiacc/uiacc.go` | Selector/object API is present; local methods return nil/default values. |
| `https` | `AutoGo/https/https.go` | GET/POST/multipart surface is present; local methods are stubs. |
| `ime` | `AutoGo/ime/ime.go` | Input and clipboard API surface; local methods are stubs. |
| `motion` | `AutoGo/motion/motion.go`, `AutoGo/motion/keycode.go` | Touch/gesture/key API surface; local methods are stubs. |
| `images` | `AutoGo/images/images.go` | Capture, color, image IO and transforms; local methods are stubs. |
| `opencv` | `AutoGo/opencv/*.go` | OpenCV Go wrapper API; native behavior depends on included libraries/host. |
| `dotocr` | `AutoGo/dotocr/dotocr.go` | OCR API surface. **Drift:** official `Ocr`/`FindStr` docs mention a `displayId` argument, while this checked-in Go signature does not. Use the local signature when compiling here and verify the installed host version before relying on virtual-display arguments. |
| `ppocr` | `AutoGo/ppocr/ppocr.go` | OCR result and input API surface; local calls return empty results. |
| `yolo` | `AutoGo/yolo/yolo.go` | Detector API surface; model/runtime implementation is not present in this Go file. |
| `storages` | `AutoGo/storages/storages.go` | Table/key storage surface; local methods are stubs. |
| `system` | `AutoGo/system/system.go` | Process stats/restart surface; local methods are stubs. |
| `utils` | `AutoGo/utils/utils.go` | Utility API surface; local methods are stubs. |
| `apkctl` | `AutoGo/apkctl/apkctl.go` | Detailed JS/ObjC bridge contract and Go method shapes are present; native bridge source is absent. |
| `console`, `hud` | `AutoGo/console/console.go`, `AutoGo/hud/hud.go` | Floating display API surfaces; local device-facing methods are stubs. |
| `imgui` | `AutoGo/imgui/*.go` | Large generated CGo bindings and iOS library references; avoid pulling it into simple script tasks. |

The sibling `awesomeProject1` checkout contains higher-level code that wraps these APIs, but those helpers are not automatically available in `iphone_autogo`. Read [project crosswalk](project-crosswalk.md) before using a `common.*` helper.

## Official Sources

- [iOS app](https://autogo.cc/ios/app.md), [device](https://autogo.cc/ios/device.md), [files](https://autogo.cc/ios/files.md), [uiacc](https://autogo.cc/ios/uiacc.md)
- [iOS https](https://autogo.cc/ios/https.md), [ime](https://autogo.cc/ios/ime.md), [motion](https://autogo.cc/ios/motion.md), [images](https://autogo.cc/ios/images.md), [opencv](https://autogo.cc/ios/opencv.md)
- [iOS dotocr](https://autogo.cc/ios/dotocr.md), [storages](https://autogo.cc/ios/storages.md), [system](https://autogo.cc/ios/system.md), [utils](https://autogo.cc/ios/utils.md), [ppocr](https://autogo.cc/ios/ppocr.md), [yolo](https://autogo.cc/ios/yolo.md)
- [iOS apkctl](https://autogo.cc/ios/apkctl.md), [imgui](https://autogo.cc/ios/imgui.md), [console](https://autogo.cc/ios/console.md), [hud](https://autogo.cc/ios/hud.md)
- [Overview](https://autogo.cc/README.md), [changelog](https://autogo.cc/changelog.md), [sidebar](https://autogo.cc/_sidebar.md)
