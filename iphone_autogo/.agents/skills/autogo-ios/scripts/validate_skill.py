#!/usr/bin/env python3
"""Check Skill metadata, local references, and AutoGo example imports."""

import argparse
import re
import sys
from pathlib import Path


LOCAL_LINK = re.compile(r"(?<!!)\[[^\]]+\]\(([^)]+)\)")
AUTO_GO_IMPORT = re.compile(r'"github\.com/Dasongzi1366/AutoGo/([A-Za-z0-9_./-]+)"')


def validate(skill_dir, project_root=None):
    skill_dir = Path(skill_dir).resolve()
    project_root = Path(project_root).resolve() if project_root else skill_dir.parents[2]
    skill_file = skill_dir / "SKILL.md"
    errors = []
    if not skill_file.is_file():
        return [f"missing {skill_file}"]

    text = skill_file.read_text(encoding="utf-8")
    if not text.startswith("---\n"):
        return ["SKILL.md must start with YAML frontmatter"]
    end = text.find("\n---\n", 4)
    if end < 0:
        return ["SKILL.md frontmatter is not closed"]
    frontmatter = text[4:end]
    fields = {}
    for line in frontmatter.splitlines():
        key, separator, value = line.partition(":")
        if separator:
            fields[key.strip()] = value.strip().strip("'\"")
    for key in ("name", "description"):
        if not fields.get(key):
            errors.append(f"frontmatter requires {key}")
    if fields.get("name") and fields["name"] != skill_dir.name:
        errors.append(f"frontmatter name {fields['name']!r} must match folder {skill_dir.name!r}")

    for markdown in skill_dir.rglob("*.md"):
        markdown_text = markdown.read_text(encoding="utf-8")
        for target in LOCAL_LINK.findall(markdown_text):
            target = target.split("#", 1)[0]
            if not target or ":" in target or target.startswith("/"):
                continue
            if not (markdown.parent / target).is_file():
                errors.append(f"{markdown.relative_to(skill_dir)} links to missing file: {target}")

    package_root = project_root / "AutoGo"
    for example in (skill_dir / "examples").rglob("*.go"):
        source = example.read_text(encoding="utf-8")
        for package in AUTO_GO_IMPORT.findall(source):
            if not (package_root / package).is_dir():
                errors.append(f"{example.relative_to(skill_dir)} imports missing AutoGo/{package}")
    return errors


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--skill-dir", type=Path, default=Path(__file__).parents[1])
    parser.add_argument("--project-root", type=Path)
    args = parser.parse_args()
    errors = validate(args.skill_dir, args.project_root)
    if errors:
        print("Skill validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print("Skill metadata, local references, and AutoGo example imports are valid.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
