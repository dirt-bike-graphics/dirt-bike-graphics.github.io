package fitment

import "testing"

func TestFitmentLookup(t *testing.T) {
	if len(Bikes()) != 17 {
		t.Fatalf("expected 17 reference bikes, got %d", len(Bikes()))
	}
	kit := FindKit("crf 250 r", "2026")
	if kit == nil || kit.ID != "honda-crf450r-2025" {
		t.Fatalf("unexpected kit: %+v", kit)
	}
	if !SameKit("CRF450R", "2025", "CRF250R", "2026") {
		t.Fatal("expected CRF450R 2025 and CRF250R 2026 to share plastics")
	}
	if len(SharedFits("YZ250F", "2025")) != 11 {
		t.Fatalf("unexpected YZ250F shared fits count")
	}
}
