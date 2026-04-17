package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestParseApprovalsBody_Array(t *testing.T) {
	var out []models.ApprovalRequestInput
	err := parseApprovalsBody([]byte(`[{"userId":"u1","isBlocking":true},{"userId":"u2"}]`), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].UserID != "u1" || out[0].IsBlocking == nil || !*out[0].IsBlocking {
		t.Errorf("unexpected first entry: %+v", out[0])
	}
	if out[1].IsBlocking != nil {
		t.Errorf("second entry should have nil IsBlocking (key absent), got %+v", out[1])
	}
	if out[1].UserID != "u2" {
		t.Errorf("unexpected second entry: %+v", out[1])
	}
}

func TestParseApprovalsBody_ObjectWrapper(t *testing.T) {
	var out []models.ApprovalRequestInput
	err := parseApprovalsBody([]byte(`{"approvals":[{"userId":"u1","message":"please"}]}`), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].UserID != "u1" || out[0].Message != "please" {
		t.Errorf("unexpected result: %+v", out)
	}
}

func TestParseApprovalsBody_ObjectWithoutApprovalsKey(t *testing.T) {
	var out []models.ApprovalRequestInput
	err := parseApprovalsBody([]byte(`{"postId":"p1"}`), &out)
	if err == nil {
		t.Fatal("expected error for missing approvals key")
	}
	if !strings.Contains(err.Error(), "approvals") {
		t.Errorf("error should mention approvals, got: %v", err)
	}
}

func TestParseApprovalsBody_InvalidShape(t *testing.T) {
	var out []models.ApprovalRequestInput
	err := parseApprovalsBody([]byte(`"not json object or array"`), &out)
	if err == nil {
		t.Fatal("expected error for non-array/non-object body")
	}
}

// TestApprovalRequestInput_IsBlockingMarshal locks in the three-state wire
// behavior of the isBlocking field: explicit true serializes as true,
// explicit false serializes as false (not omitted), and a nil pointer is
// dropped by omitempty so the server sees "key absent" and applies its
// default. The pre-fix bool-with-omitempty collapsed explicit-false into
// the omitted case, so this test would have failed at the false check.
func TestApprovalRequestInput_IsBlockingMarshal(t *testing.T) {
	t.Run("explicit true marshals as true", func(t *testing.T) {
		v := true
		b, err := json.Marshal(models.ApprovalRequestInput{UserID: "u1", IsBlocking: &v})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), `"isBlocking":true`) {
			t.Errorf("expected isBlocking:true, got: %s", b)
		}
	})
	t.Run("explicit false marshals as false", func(t *testing.T) {
		v := false
		b, err := json.Marshal(models.ApprovalRequestInput{UserID: "u1", IsBlocking: &v})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), `"isBlocking":false`) {
			t.Errorf("expected isBlocking:false, got: %s", b)
		}
	})
	t.Run("nil pointer omits key", func(t *testing.T) {
		b, err := json.Marshal(models.ApprovalRequestInput{UserID: "u1"})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), "isBlocking") {
			t.Errorf("expected isBlocking key to be omitted, got: %s", b)
		}
	})
	t.Run("round-trips through JSON", func(t *testing.T) {
		src := `{"userId":"u1","isBlocking":false}`
		var got models.ApprovalRequestInput
		if err := json.Unmarshal([]byte(src), &got); err != nil {
			t.Fatal(err)
		}
		if got.IsBlocking == nil {
			t.Fatal("expected non-nil IsBlocking after unmarshaling explicit false")
		}
		if *got.IsBlocking {
			t.Errorf("expected false, got true")
		}
		b, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), `"isBlocking":false`) {
			t.Errorf("round-trip dropped explicit false: %s", b)
		}
	})
}

func TestParseApprovalsBody_LeadingWhitespace(t *testing.T) {
	var out []models.ApprovalRequestInput
	err := parseApprovalsBody([]byte("   \n[{\"userId\":\"u1\"}]\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].UserID != "u1" {
		t.Errorf("unexpected result: %+v", out)
	}
}
