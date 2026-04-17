package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestApprovalService_List(t *testing.T) {
	svc := NewApprovalService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/p1/approvals" {
			t.Errorf("expected /posts/p1/approvals, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.ApprovalListResponse{
			Approvals: []models.Approval{{ID: "a1", Status: "Requested"}},
		}), nil
	}))

	resp, err := svc.List("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Approvals) != 1 || resp.Approvals[0].ID != "a1" {
		t.Errorf("unexpected approvals: %+v", resp.Approvals)
	}
}

func TestApprovalService_Create(t *testing.T) {
	svc := NewApprovalService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/approvals" {
			t.Errorf("expected /approvals, got %s", r.URL.Path)
		}
		var body models.CreateApprovalsRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.PostID != "p1" || len(body.Approvals) != 1 {
			t.Errorf("unexpected body: %+v", body)
		}
		return jsonResponse(t, http.StatusOK, map[string]interface{}{"approvals": []map[string]string{{"id": "a1"}}}), nil
	}))

	_, err := svc.Create(models.CreateApprovalsRequest{
		PostID: "p1",
		Approvals: []models.ApprovalRequestInput{
			{UserID: "u1", IsBlocking: true},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApprovalService_Delete(t *testing.T) {
	svc := NewApprovalService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/approvals/a1" {
			t.Errorf("expected /approvals/a1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		// /approvals/{id} DELETE returns `{"deletedApproval": Approval}`.
		return jsonResponse(t, http.StatusOK, map[string]interface{}{
			"deletedApproval": models.Approval{ID: "a1", Status: "Approved"},
		}), nil
	}))

	data, err := svc.Delete("a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got struct {
		DeletedApproval models.Approval `json:"deletedApproval"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("parse delete body: %v", err)
	}
	if got.DeletedApproval.ID != "a1" {
		t.Errorf("expected deleted approval id=a1; got %+v", got.DeletedApproval)
	}
}
