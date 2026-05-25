package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireMember(t *testing.T) {
	tests := []struct {
		name       string
		setMember  bool
		member     *Member
		wantStatus int
	}{
		{"no member", false, nil, 401},
		{"has member", true, &Member{Level: "VIEWER"}, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireMember()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			s := NewFeatmapService()
			if tt.setMember {
				s.SetMemberObject(tt.member)
			}
			ctx := context.WithValue(r.Context(), contextKey, &Env{Service: s})
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestRequireAdmin(t *testing.T) {
	tests := []struct {
		name       string
		member     *Member
		wantStatus int
	}{
		{"viewer", &Member{Level: "VIEWER"}, 401},
		{"editor", &Member{Level: "EDITOR"}, 401},
		{"admin", &Member{Level: "ADMIN"}, 200},
		{"owner", &Member{Level: "OWNER"}, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireAdmin()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			s := NewFeatmapService()
			s.SetMemberObject(tt.member)
			ctx := context.WithValue(r.Context(), contextKey, &Env{Service: s})
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestRequireEditor(t *testing.T) {
	tests := []struct {
		name       string
		member     *Member
		wantStatus int
	}{
		{"viewer", &Member{Level: "VIEWER"}, 401},
		{"editor", &Member{Level: "EDITOR"}, 200},
		{"admin", &Member{Level: "ADMIN"}, 200},
		{"owner", &Member{Level: "OWNER"}, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireEditor()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			s := NewFeatmapService()
			s.SetMemberObject(tt.member)
			ctx := context.WithValue(r.Context(), contextKey, &Env{Service: s})
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestRequireOwner(t *testing.T) {
	tests := []struct {
		name       string
		member     *Member
		wantStatus int
	}{
		{"admin", &Member{Level: "ADMIN"}, 401},
		{"owner", &Member{Level: "OWNER"}, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireOwner()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			s := NewFeatmapService()
			s.SetMemberObject(tt.member)
			ctx := context.WithValue(r.Context(), contextKey, &Env{Service: s})
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
