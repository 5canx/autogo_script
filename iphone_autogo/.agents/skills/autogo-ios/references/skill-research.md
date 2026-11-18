# Open-Source Skill Research

Repository metadata and Skill files checked on 2026-11-18. Maintenance dates below are the GitHub `pushed_at` values observed in the public API at research time; they are not a guarantee of release quality. No code or dependencies from these repositories are vendored.

| Project | Skill files actually reviewed | License / last pushed | Fit and what to reuse |
|---|---|---|---|
| [appium/skills](https://github.com/appium/skills) | `skills/setup-xcuitest/SKILL.md`; tree includes `setup-xcuitest`, `xcuitest-real-device-config`, and `appium-troubleshooting` | Apache-2.0 / 2026-10-22 | Maintained official Appium skills. Its XCUITest/WDA setup is a different runtime from AutoGo. Borrow task routing, platform prerequisites, troubleshooting and real-device readiness checks only; do not add Appium to AutoGo scripts. |
| [LambdaTest/agent-skills](https://github.com/LambdaTest/agent-skills) | `appium-skill/SKILL.md`, `appium-skill/reference/ios-specific.md` | MIT / 2026-11-17 | Detailed Appium and iOS capability guidance, with modular references. Useful for trigger wording, collecting environment facts, and locator/wait/report discipline. LambdaTest cloud/session examples do not describe AutoGo. |
| [TerminalSkills/skills](https://github.com/TerminalSkills/skills) | `skills/appium/SKILL.md`, `skills/ios-debugger-agent/SKILL.md` | Apache-2.0 / 2026-11-08 | Large multi-tool library. The debugger Skill uses evidence-first screen/log collection and chooses between command-line and available MCP tools. Reuse its route-selection and evidence/report structure; simulator/Xcode commands are not AutoGo iPhone runtime instructions. |
| [magnus919/agent-skills](https://github.com/magnus919/agent-skills) | `agent-skills/SKILL.md`, `binary-analysis/SKILL.md` | MIT / 2026-11-16 | General skill-authoring and binary-analysis material, not an AutoGo/iOS automation integration. The authoring Skill emphasizes real project grounding, conditional references, accurate metadata, and validation. Borrow progressive disclosure/reference organization only; avoid pulling in binary-analysis workflows or tooling for normal script requests. |
| [owen800q/mobile-reverse-skill](https://github.com/owen800q/mobile-reverse-skill) | Could not read a `SKILL.md`: GitHub repository API returned 404 at research time | Unknown | The supplied repository path is unavailable under that name, so license, maintenance, and content cannot be verified. Do not copy or attribute techniques to it. Re-check if the user provides a corrected URL. |

## Compatibility Decision

AutoGo runs Go code in its own iOS host and exposes screen/app/device APIs there. Appium uses a WebDriver server and platform drivers such as XCUITest/WDA; Xcode debugger Skills assume an Xcode-built app and simulator/device tooling. Those projects are useful examples of agent workflow and evidence handling, but their capabilities, commands, packages, and runtime assumptions must not be copied into AutoGo code. The Skill therefore remains dependency-free and project-first.

Selected primary references:

- [Appium XCUITest setup Skill](https://github.com/appium/skills/blob/main/skills/setup-xcuitest/SKILL.md)
- [LambdaTest Appium Skill](https://github.com/LambdaTest/agent-skills/blob/main/appium-skill/SKILL.md) and [iOS reference](https://github.com/LambdaTest/agent-skills/blob/main/appium-skill/reference/ios-specific.md)
- [TerminalSkills Appium](https://github.com/TerminalSkills/skills/blob/main/skills/appium/SKILL.md) and [iOS debugger agent](https://github.com/TerminalSkills/skills/blob/main/skills/ios-debugger-agent/SKILL.md)
- [Magnus agent-skills authoring Skill](https://github.com/magnus919/agent-skills/blob/main/agent-skills/SKILL.md) and [binary-analysis Skill](https://github.com/magnus919/agent-skills/blob/main/binary-analysis/SKILL.md)
