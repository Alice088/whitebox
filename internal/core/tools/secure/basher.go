package secure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"whitebox/internal/paths"
)

type Basher struct {
	Blacklist []*regexp.Regexp
	Whitelist []*regexp.Regexp
}

func Command(cmd string) error {
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
			return fmt.Errorf("command blocked: unsecure patter: %s", r.String())
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
		return
	}

	var rules RawBasher
	err = json.Unmarshal(raw, &rules)
	if err != nil {
		panic("failed to parse commands/rules: " + err.Error())
	}

	Basherx.Whitelist, err = rules.compile(rules.Whitelist)
	Basherx.Blacklist, err = rules.compile(rules.Blacklist)
}
