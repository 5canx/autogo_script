# Live Screen Access

The AutoGo Skill can inspect a live screen when the agent client has a screen tool. This project wires the existing `awesomeProject1/mcp-iphone-mirror` server into Codex and Claude Code with read-only tools. The server runs from source with `go run`, so later changes to that MCP project are picked up without rebuilding a binary.

## What the Agent Can Do

- `get_status`: check whether the iPhone Mirroring window is available and report display geometry.
- `screenshot`: return a device-sized PNG as image content for visual inspection.
- `ocr`: return text labels and bounding boxes/centers in mirror device coordinates.

The agent can use that evidence to describe the page, identify labels, plan an AutoGo script, and compare before/after screenshots when the user explicitly asks for a live test. The MCP server also implements touch/input tools, but the project config keeps them unavailable in Codex and denies them in Claude Code. Production scripts still run in the AutoGo host.

## Requirements and Limits

- The iPhone must already be connected to and visible in macOS iPhone Mirroring.
- The user must grant Screen Recording and Accessibility permissions to the terminal/client process as macOS requests. This configuration does not grant either permission.
- `ocr` uses the separate PaddleOCR service at `127.0.0.1:9527`; without it, use the image returned by `screenshot` or AutoGo's own OCR.
- The MCP screenshot defaults to 1170x2532 device coordinates; `awesomeProject1` AutoGo scripts use a 1284x2778 baseline. Do not copy MCP OCR coordinates directly into a script; convert through the project's verified `common.Sx/Sy` path and validate on the AutoGo device.
- Claude Code will ask to trust the project MCP configuration on first use. Codex must trust the project before it loads `.codex/config.toml`.
- If `awesomeProject1` or the Go executable moves, update the paths in `.mcp.json` and `.codex/config.toml`. The provided paths use this machine's current locations.

## If Setup Fails

1. Confirm `go version` is at least 1.25 and run `go test ./...` from `/Users/scan/GolandProjects/awesomeProject1/mcp-iphone-mirror`.
2. Open iPhone Mirroring and reconnect the phone.
3. Grant Screen Recording for screenshots and Accessibility for window lookup, then restart the agent.
4. Start the optional OCR server from `mcp-iphone-mirror/PaddleOCRSvr.py` only when OCR boxes are needed. The image screenshot does not require that service.
