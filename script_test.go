package script

import (
	"encoding/hex"
	"testing"
)

func TestScript(t *testing.T) {

	scripts := []string{
		"0b3c4b616e7965323032303e7588",
	}
	for line, scriptHex := range scripts {
		if len(scriptHex) == 0 {
			continue
		}
		script, err := hex.DecodeString(scriptHex)
		if err != nil {
			t.Logf("ignore line: %d", line)
			continue
		}

		pc, ok := GetLockingScriptPushDropPosition(script)

		t.Logf("script: %d, ok: %v", pc, ok)
	}
}
