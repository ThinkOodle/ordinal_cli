package cmd

import (
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
	if out[0].UserID != "u1" || !out[0].IsBlocking {
		t.Errorf("unexpected first entry: %+v", out[0])
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
