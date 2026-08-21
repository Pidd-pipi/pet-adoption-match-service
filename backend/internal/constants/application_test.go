package constants

import "testing"

func TestConfirmedRejectedTransition(t *testing.T) {
	for _, s := range NextApplicationStatuses(AppStatusConfirmed) {
		if s == AppStatusRejected {
			return
		}
	}
	t.Fatalf("confirmed status should allow a rejected transition")
}
