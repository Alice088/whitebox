package secure

import (
	"os"
	"path/filepath"
	"testing"
	"whitebox/internal/paths"
)

func setupWorkspace(t *testing.T) string {
	dir, err := os.MkdirTemp("", "secure_test")
	if err != nil {
		t.Fatal(err)
	}
	paths.WorkspaceDir = dir
	return dir
}

func TestPath_Empty(t *testing.T) {
	setupWorkspace(t)

	_, err := Path("")
	if err == nil {
		t.Fatalf("expected error for empty path")
	}
}

func TestPath_Absolute(t *testing.T) {
	setupWorkspace(t)

	_, err := Path("/etc/passwd")
	if err == nil {
		t.Fatalf("absolute path should be rejected")
	}
}

func TestPath_Valid(t *testing.T) {
	dir := setupWorkspace(t)

	p, err := Path("file.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "file.txt")
	if p != expected {
		t.Fatalf("expected %s, got %s", expected, p)
	}
}

func TestPath_Traversal(t *testing.T) {
	setupWorkspace(t)

	_, err := Path("../secret.txt")
	if err == nil {
		t.Fatalf("path traversal should be rejected")
	}
}

func TestPath_NestedTraversal(t *testing.T) {
	setupWorkspace(t)

	_, err := Path("a/../../etc/passwd")
	if err == nil {
		t.Fatalf("nested traversal should be rejected")
	}
}

func TestPath_CleanTraversal(t *testing.T) {
	setupWorkspace(t)

	_, err := Path("a/../b/../../etc")
	if err == nil {
		t.Fatalf("clean traversal should be rejected")
	}
}

func TestPath_SymlinkEscape(t *testing.T) {
	dir := setupWorkspace(t)

	outsideDir, err := os.MkdirTemp("", "outside")
	if err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(outsideDir, "file.txt")
	if err := os.WriteFile(target, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	link := filepath.Join(dir, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	_, err = Path("link")
	if err == nil {
		t.Fatalf("symlink escape should be rejected")
	}
}

func TestPath_SymlinkInside(t *testing.T) {
	dir := setupWorkspace(t)

	target := filepath.Join(dir, "real.txt")
	if err := os.WriteFile(target, []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}

	link := filepath.Join(dir, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	_, err := Path("link")
	if err != nil {
		t.Fatalf("valid symlink inside workspace should pass: %v", err)
	}
}

func TestPath_NonExistent(t *testing.T) {
	setupWorkspace(t)

	_, err := Path("newfile.txt")
	if err != nil {
		t.Fatalf("non-existent file should be allowed: %v", err)
	}
}

func TestPath_Dot(t *testing.T) {
	dir := setupWorkspace(t)

	p, err := Path(".")
	if err != nil {
		t.Fatalf("dot path should pass: %v", err)
	}

	if p != dir {
		t.Fatalf("expected workspace dir, got %s", p)
	}
}

func TestPath_DoubleSlash(t *testing.T) {
	dir := setupWorkspace(t)

	p, err := Path("a//b//c.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "a/b/c.txt")
	if p != expected {
		t.Fatalf("expected %s, got %s", expected, p)
	}
}

func setupMemory(t *testing.T) string {
	dir, err := os.MkdirTemp("", "memory_test")
	if err != nil {
		t.Fatal(err)
	}
	paths.MemoriesDir = dir
	return dir
}
func TestPath_MemoryWrite(t *testing.T) {
	setupWorkspace(t)
	setupMemory(t)

	p, err := Path("memory/user/name.txt")
	if err != nil {
		t.Fatalf("memory path should be allowed: %v", err)
	}

	expected := filepath.Join(paths.MemoriesDir, "user/name.txt")

	// важно: должен резолвиться в memory, а не workspace
	if p != expected {
		t.Fatalf("expected %s, got %s", expected, p)
	}
}
