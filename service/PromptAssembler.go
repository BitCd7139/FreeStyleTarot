package service

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/storage"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	backgroundPrompts    map[string][]byte
	backgroundPromptOnce sync.Once
)

const (
	prefix = "prompt_"
	suffix = ".md"
)

// InputsAssembler is kept for compatibility; use RunTarotReading for the two-agent flow.
func InputsAssembler(predict request.Predict) (systemMsg string, userMsg string, err error) {
	materials, err := BuildSpreadMaterials(predict.Cards)
	if err != nil {
		return "", "", err
	}
	return SharedSystemPrompt(predict.Model), AgentAUserPrompt(predict, materials), nil
}

func initBackgroundPrompts() {
	backgroundPromptOnce.Do(func() {
		backgroundPrompts = make(map[string][]byte)

		entries, err := storage.Assets.ReadDir(".")
		if err != nil {
			zap.S().Errorf("CRITICAL: Failed to read assets directory: %v", err)
			return
		}

		zap.S().Debug("Starting to scan background prompts...")

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filename := entry.Name()

			if strings.HasPrefix(filename, prefix) && strings.HasSuffix(filename, suffix) {
				name := strings.TrimPrefix(filename, prefix)
				name = strings.TrimSuffix(name, suffix)

				data, err := storage.Assets.ReadFile(filename)
				if err != nil {
					zap.S().Errorf("Failed to read file %s: %v", filename, err)
					continue
				}

				backgroundPrompts[name] = data
				zap.S().Debugw("Loaded background prompt", "name", name, "file", filename)
			}
		}

		zap.S().Infof("Initialization complete. Loaded %d persona prompts.", len(backgroundPrompts))
	})
}

func getBackgroundPrompt(name string) []byte {
	initBackgroundPrompts()

	if data, ok := backgroundPrompts[name]; ok {
		return data
	}

	zap.S().Warnw("Persona prompt not found", "name", name)
	return nil
}
