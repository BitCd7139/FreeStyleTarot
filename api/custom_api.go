package api

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/service"
	"context"
)

// withCustomAPI injects a CustomLLM into the request context when the user
// provided custom API config. Returns the original context if no config or
// initialization fails (caller may decide to fall back to global LLM).
func withCustomAPI(ctx context.Context, cfg *request.CustomAPIConfig) (context.Context, error) {
	if cfg == nil || cfg.APIKey == "" {
		return ctx, nil
	}
	llm, err := service.NewCustomLLM(*cfg)
	if err != nil {
		return ctx, err
	}
	return service.WithCustomLLM(ctx, llm), nil
}
