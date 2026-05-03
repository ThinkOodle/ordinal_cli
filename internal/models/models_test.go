package models

import (
	"encoding/json"
	"strings"
	"testing"
)

// Optional nested objects (User on Comment/Approval, SlackWebhook on
// SlackBoost, CreatedBy on Profile) must disappear from marshaled output when
// the API omitted them — not serialize as `"field": {}`. json omitempty only
// works on pointer-to-struct; a bare struct with omitempty always renders.
// The CLI re-marshals typed responses for table/csv, so a synthetic empty
// object here diverges from the raw API body.
func TestOptionalNestedObjects_OmittedWhenUnset(t *testing.T) {
	cases := []struct {
		name string
		in   interface{}
		keys []string
	}{
		{"comment", Comment{ID: "c_1", Message: "hi"}, []string{"user"}},
		{"approval", Approval{ID: "a_1", Status: "pending"}, []string{"user", "requestedBy"}},
		{"slackBoost", SlackBoost{ID: "sb_1"}, []string{"slackWebhook"}},
		{"profile", Profile{ID: "p_1", Channel: "linkedIn"}, []string{"createdBy"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := json.Marshal(c.in)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			got := string(b)
			for _, k := range c.keys {
				if strings.Contains(got, `"`+k+`":`) {
					t.Errorf("%s: expected key %q to be omitted; got %s", c.name, k, got)
				}
			}
		})
	}
}

// When the nested object IS present on the wire, marshaling must preserve it.
// This guards against over-correcting by dropping optional fields that do
// carry data.
func TestOptionalNestedObjects_PreservedWhenSet(t *testing.T) {
	comment := Comment{
		ID:      "c_1",
		Message: "hi",
		User:    &CommentUser{ID: "u_1", Email: "a@b.com"},
	}
	b, err := json.Marshal(comment)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	if !strings.Contains(got, `"user":{"id":"u_1"`) {
		t.Errorf("expected user block; got %s", got)
	}
}

// Post and idea responses carry nested, channel-specific payloads. The CLI
// keeps those as raw JSON so new channel schemas are not lost while formatting
// typed responses.
func TestChannelRawMessages_PreserveNewChannels(t *testing.T) {
	postJSON := []byte(`{
		"id":"p_1",
		"title":"Post",
		"status":"Scheduled",
		"tikTok":{"copy":"caption"},
		"youTubeShorts":{"title":"short"}
	}`)
	var post Post
	if err := json.Unmarshal(postJSON, &post); err != nil {
		t.Fatalf("unmarshal post: %v", err)
	}
	encodedPost, err := json.Marshal(post)
	if err != nil {
		t.Fatalf("marshal post: %v", err)
	}
	gotPost := string(encodedPost)
	for _, want := range []string{`"tikTok":{"copy":"caption"}`, `"youTubeShorts":{"title":"short"}`} {
		if !strings.Contains(gotPost, want) {
			t.Errorf("expected post output to preserve %s; got %s", want, gotPost)
		}
	}

	ideaJSON := []byte(`{
		"id":"i_1",
		"title":"Idea",
		"tikTok":{"copy":"caption"},
		"youTubeShorts":{"title":"short"}
	}`)
	var idea Idea
	if err := json.Unmarshal(ideaJSON, &idea); err != nil {
		t.Fatalf("unmarshal idea: %v", err)
	}
	encodedIdea, err := json.Marshal(idea)
	if err != nil {
		t.Fatalf("marshal idea: %v", err)
	}
	gotIdea := string(encodedIdea)
	for _, want := range []string{`"tikTok":{"copy":"caption"}`, `"youTubeShorts":{"title":"short"}`} {
		if !strings.Contains(gotIdea, want) {
			t.Errorf("expected idea output to preserve %s; got %s", want, gotIdea)
		}
	}
}
