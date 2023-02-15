package sol

import (
	"fmt"
	"testing"
)

func TestCardID(t *testing.T) {
	cid := NewCardID(0, 1, 1)
	str := fmt.Sprint(cid)
	if str != "0 Club 1" {
		t.Errorf("wrong string for 0 1 1: %s", cid)
	}
	// col := cid.Color()
	// if col != BasicColors["Black"] {
	// 	t.Errorf("wrong color for %s", cid)
	// }

	cid2 := NewCardID(1, 1, 1)
	if !SameCard(cid, cid2) {
		t.Errorf("cards should be the same: %s and %s", cid, cid2)
	}
	if SameCardAndPack(cid, cid2) {
		t.Errorf("cards should NOT be the same: %s and %s", cid, cid2)
	}
}
