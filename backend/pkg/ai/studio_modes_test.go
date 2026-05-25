package ai

import "testing"

func TestParseStudioMode(t *testing.T) {
	m, ok := ParseStudioMode("")
	if !ok || m != ModeErase {
		t.Fatalf("empty -> erase, got %v %v", m, ok)
	}
	m, ok = ParseStudioMode("upscale")
	if !ok || m != ModeUpscale {
		t.Fatalf("upscale, got %v %v", m, ok)
	}
	_, ok = ParseStudioMode("invalid")
	if ok {
		t.Fatal("expected invalid mode")
	}
}

func TestStudioModeRequiresMask(t *testing.T) {
	if !StudioModeRequiresMask(ModeErase) {
		t.Fatal("erase requires mask")
	}
	if StudioModeRequiresMask(ModeUpscale) {
		t.Fatal("upscale should not require mask")
	}
}
