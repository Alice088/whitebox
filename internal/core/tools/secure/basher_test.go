package secure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"whitebox/internal/paths"
)

type Basher struct {
	Blacklist []*regexp.Regexp
	Whitelist []*regexp.Regexp
	enabled   bool
}

func Command(cmd string) error {
	if !Basherx.enabled {
		return nil
	}

	cmd = strings.TrimSpace(cmd)

	if cmd == "" {
		return fmt.Errorf("empty command")
	}

	// 1. проверка whitelist
	allowed := false
	for _, r := range Basherx.Whitelist {
		if r.MatchString(cmd) {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("command not allowed")
	}

	// 2. проверка blacklist
	for _, r := range Basherx.Blacklist {
		if r.MatchString(cmd) {
			return fmt.Errorf("command blocked: unsecure pattern: %s", r.String())
		}
	}

	return nil
}

type RawBasher struct {
	Blacklist []string `json:"blacklist"`
	Whitelist []string `json:"whitelist"`
}

func (rb RawBasher) compile(patterns []string) ([]*regexp.Regexp, error) {
	var result []*regexp.Regexp

	for _, p := range patterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

var Basherx Basher

func init() {
	raw, err := os.ReadFile(filepath.Join(paths.CommandsDir, "rules.json"))
	if err != nil {
		Basherx.enabled = false
		return
	}
	Basherx.enabled = true

	var rules RawBasher
	err = json.Unmarshal(raw, &rules)
	if err != nil {
		panic("failed to parse commands/rules: " + err.Error())
	}

	Basherx.Whitelist, err = rules.compile(rules.Whitelist)
	if err != nil {
		panic("failed to compile whitelist rules.json: " + err.Error())
	}
	Basherx.Blacklist, err = rules.compile(rules.Blacklist)
	if err != nil {
		panic("failed to compile blacklist rules.json: " + err.Error())
	}
}
