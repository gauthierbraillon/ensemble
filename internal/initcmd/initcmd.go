package initcmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

const hookCommand = "ensemble hook"

func WriteSettings(dir string) (bool, error) {
	settingsPath := filepath.Join(dir, ".claude", "settings.json")
	settings := map[string]interface{}{}
	if data, err := os.ReadFile(settingsPath); err == nil { // #nosec G304
		_ = json.Unmarshal(data, &settings)
	}
	if hasHook(settings) {
		return false, nil
	}
	mergeHook(settings)
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0750); err != nil {
		return false, err
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return false, err
	}
	if err := os.WriteFile(settingsPath, append(data, '\n'), 0600); err != nil { // #nosec G306
		return false, err
	}
	return true, nil
}

func EnsembleOnPath() bool {
	_, err := exec.LookPath("ensemble")
	return err == nil
}

func hasHook(settings map[string]interface{}) bool {
	hooks, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		return false
	}
	preToolUse, ok := hooks["PreToolUse"].([]interface{})
	if !ok {
		return false
	}
	for _, entry := range preToolUse {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		innerHooks, ok := entryMap["hooks"].([]interface{})
		if !ok {
			continue
		}
		for _, h := range innerHooks {
			hMap, ok := h.(map[string]interface{})
			if !ok {
				continue
			}
			if hMap["command"] == hookCommand {
				return true
			}
		}
	}
	return false
}

func mergeHook(settings map[string]interface{}) {
	hooks, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		hooks = map[string]interface{}{}
		settings["hooks"] = hooks
	}
	entry := map[string]interface{}{
		"matcher": "Write|Edit",
		"hooks": []interface{}{
			map[string]interface{}{
				"type":    "command",
				"command": hookCommand,
			},
		},
	}
	preToolUse, ok := hooks["PreToolUse"].([]interface{})
	if !ok {
		preToolUse = []interface{}{}
	}
	hooks["PreToolUse"] = append(preToolUse, entry)
}
