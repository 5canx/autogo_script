import importlib.util
import sys
import tempfile
import unittest
from pathlib import Path


sys.dont_write_bytecode = True
SCRIPT = Path(__file__).parents[1] / "scripts" / "validate_skill.py"
SPEC = importlib.util.spec_from_file_location("validate_skill", SCRIPT)
MODULE = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(MODULE)


class ValidateSkillTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name)
        self.skill = self.root / ".agents" / "skills" / "autogo-ios"
        self.project = self.root
        self.skill.mkdir(parents=True)

    def tearDown(self):
        self.temp.cleanup()

    def write(self, path, content):
        target = self.skill / path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content, encoding="utf-8")
        return target

    def test_accepts_valid_frontmatter_and_reference(self):
        self.write("SKILL.md", "---\nname: autogo-ios\ndescription: AutoGo iOS tasks\n---\n[API](references/api.md)\n")
        self.write("references/api.md", "# API\n")
        self.assertEqual(MODULE.validate(self.skill, self.project), [])

    def test_reports_missing_required_frontmatter(self):
        self.write("SKILL.md", "---\nname: autogo-ios\n---\n")
        errors = MODULE.validate(self.skill, self.project)
        self.assertTrue(any("description" in error for error in errors))

    def test_reports_missing_local_reference(self):
        self.write("SKILL.md", "---\nname: autogo-ios\ndescription: AutoGo iOS tasks\n---\n[API](references/missing.md)\n")
        errors = MODULE.validate(self.skill, self.project)
        self.assertTrue(any("references/missing.md" in error for error in errors))

    def test_reports_missing_relative_link_inside_reference(self):
        self.write("SKILL.md", "---\nname: autogo-ios\ndescription: AutoGo iOS tasks\n---\n")
        self.write("references/api.md", "[MCP](mcp-setup.md)\n")
        errors = MODULE.validate(self.skill, self.project)
        self.assertTrue(any("mcp-setup.md" in error for error in errors))

    def test_reports_unknown_project_api_package_in_examples(self):
        self.write("SKILL.md", "---\nname: autogo-ios\ndescription: AutoGo iOS tasks\n---\n")
        self.write("examples/main.go", 'package main\nimport "github.com/Dasongzi1366/AutoGo/unknown"\n')
        errors = MODULE.validate(self.skill, self.project)
        self.assertTrue(any("AutoGo/unknown" in error for error in errors))


if __name__ == "__main__":
    unittest.main()
