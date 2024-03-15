package test

import (
	"context"
	"kis-flow/log"
	"testing"
)

func TestKisLogger(t *testing.T) {
	ctx := context.Background()
	log.GetLogger().InfoFX(ctx, "Test kisflow InfoFX")
	log.GetLogger().DebugFX(ctx, "Test kisflow DebugFX")
	log.GetLogger().ErrorFX(ctx, "Test kisflow ErrorFX")
	log.GetLogger().InfoF("Test kisflow InfoF")
	log.GetLogger().DebugF("Test kisflow DebugF")
	log.GetLogger().ErrorF("Test kisflow ErrorF")
}
