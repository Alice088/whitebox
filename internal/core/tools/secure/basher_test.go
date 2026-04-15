package secure

import (
	"regexp"
	"testing"
)

func setupStrictBasher() {
	Basherx = Basher{
		enabled: true,
		Whitelist: []*regexp.Regexp{
			regexp.MustCompile(`^ls(\s|$)`),
			regexp.MustCompile(`^echo(\s|$)`),
			regexp.MustCompile(`^mkdir(\s+[a-zA-Z0-9_./-]+)+$`),
		},
		Blacklist: []*regexp.Regexp{
			regexp.MustCompile(`rm\s+-`),
			regexp.MustCompile(`;`),
			regexp.MustCompile(`&`),
			regexp.MustCompile(`\|`),
			regexp.MustCompile(`\$\(`),
			regexp.MustCompile("`"),
			regexp.MustCompile(`\.\./`),
			regexp.MustCompile(`^/`),
		},
	}
}

func TestAllowedCommand(t *testing.T) {
	setupStrictBasher()

	if err := Command("ls -la"); err != nil {
		t.Fatalf("expected allowed, got error: %v", err)
	}
}

func TestRejectUnknownCommand(t *testing.T) {
	setupStrictBasher()

	if err := Command("pwd"); err == nil {
		t.Fatalf("expected rejection for non-whitelisted command")
	}
}

func TestRejectSemicolonInjection(t *testing.T) {
	setupStrictBasher()

	err := Command("ls; rm -rf /")
	if err == nil {
		t.Fatalf("semicolon injection passed")
	}
}

func TestRejectAndOperator(t *testing.T) {
	setupStrictBasher()

	err := Command("ls && echo hi")
	if err == nil {
		t.Fatalf("&& bypass passed")
	}
}

func TestRejectPipe(t *testing.T) {
	setupStrictBasher()

	err := Command("ls | cat")
	if err == nil {
		t.Fatalf("pipe bypass passed")
	}
}

func TestRejectSubshell(t *testing.T) {
	setupStrictBasher()

	err := Command("echo $(whoami)")
	if err == nil {
		t.Fatalf("$() bypass passed")
	}
}

func TestRejectBackticks(t *testing.T) {
	setupStrictBasher()

	err := Command("echo `whoami`")
	if err == nil {
		t.Fatalf("backtick bypass passed")
	}
}

func TestRejectTraversal(t *testing.T) {
	setupStrictBasher()

	err := Command("cat ../secret.txt")
	if err == nil {
		t.Fatalf("path traversal passed")
	}
}

func TestRejectAbsolutePath(t *testing.T) {
	setupStrictBasher()

	err := Command("cat /etc/passwd")
	if err == nil {
		t.Fatalf("absolute path passed")
	}
}

func TestUnicodeBypass(t *testing.T) {
	setupStrictBasher()

	// похожий на ; символ
	err := Command("ls ؛ rm -rf /")
	if err == nil {
		t.Fatalf("unicode bypass passed")
	}
}

func TestWeirdSpacing(t *testing.T) {
	setupStrictBasher()

	err := Command("ls    -la")
	if err != nil {
		t.Fatalf("valid command broken: %v", err)
	}
}

func TestEmptyCommand(t *testing.T) {
	setupStrictBasher()

	if err := Command(""); err == nil {
		t.Fatalf("empty command should fail")
	}
}

func TestDisabledMode(t *testing.T) {
	Basherx = Basher{enabled: false}

	if err := Command("rm -rf /"); err != nil {
		t.Fatalf("disabled mode should allow everything")
	}
}

func TestTraversalHidden(t *testing.T) {
	setupStrictBasher()

	err := Command("mkdir a/../../etc")
	if err == nil {
		t.Fatalf("path traversal via mkdir passed")
	}
}

func TestEchoAbuse(t *testing.T) {
	setupStrictBasher()

	err := Command("echo rm -rf /")
	if err == nil {
		t.Fatalf("expected echo with dangerous content to be blocked")
	}
}
