package main

import (
	"testing"
)

func TestFakeRepoProjectCRUD(t *testing.T) {
	repo := newFakeRepo()

	wsID := "ws-1"
	projectID := "proj-1"

	// Create
	p := &Project{
		WorkspaceID: wsID,
		ID:          projectID,
		Title:       "Test Project",
	}
	repo.StoreProject(p)

	// Read
	got, err := repo.GetProject(wsID, projectID)
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}
	if got.Title != "Test Project" {
		t.Errorf("expected title 'Test Project', got %q", got.Title)
	}

	// Read non-existent
	_, err = repo.GetProject(wsID, "nonexistent")
	if err == nil {
		t.Error("expected error for non-existent project")
	}

	// Find by workspace
	repo.StoreProject(&Project{WorkspaceID: wsID, ID: "proj-2", Title: "Second"})
	all, err := repo.FindProjectsByWorkspace(wsID)
	if err != nil {
		t.Fatalf("FindProjectsByWorkspace failed: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 projects, got %d", len(all))
	}

	// Update
	got.Title = "Updated Project"
	repo.StoreProject(got)
	updated, _ := repo.GetProject(wsID, projectID)
	if updated.Title != "Updated Project" {
		t.Errorf("expected 'Updated Project', got %q", updated.Title)
	}

	// Delete
	repo.DeleteProject(wsID, projectID)
	_, err = repo.GetProject(wsID, projectID)
	if err == nil {
		t.Error("expected error after delete")
	}

	// Verify count after delete
	all, _ = repo.FindProjectsByWorkspace(wsID)
	if len(all) != 1 {
		t.Errorf("expected 1 project after delete, got %d", len(all))
	}
}

func TestFakeRepoWorkspaceCRUD(t *testing.T) {
	repo := newFakeRepo()

	w := &Workspace{
		ID:                   "ws-1",
		Name:                 "testws",
		AllowExternalSharing: true,
	}
	repo.StoreWorkspace(w)

	// Read by ID
	got, err := repo.GetWorkspace("ws-1")
	if err != nil {
		t.Fatalf("GetWorkspace failed: %v", err)
	}
	if got.Name != "testws" {
		t.Errorf("expected 'testws', got %q", got.Name)
	}

	// Read by name
	got, err = repo.GetWorkspaceByName("testws")
	if err != nil {
		t.Fatalf("GetWorkspaceByName failed: %v", err)
	}

	// Read non-existent
	_, err = repo.GetWorkspaceByName("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent workspace")
	}

	// Delete
	repo.DeleteWorkspace("ws-1")
	_, err = repo.GetWorkspace("ws-1")
	if err == nil {
		t.Error("expected error after workspace delete")
	}
}

func TestFakeRepoAccountCRUD(t *testing.T) {
	repo := newFakeRepo()

	a := &Account{
		ID:    "acc-1",
		Name:  "Test User",
		Email: "test@example.com",
	}
	repo.StoreAccount(a)

	// Read by ID
	got, err := repo.GetAccount("acc-1")
	if err != nil {
		t.Fatalf("GetAccount failed: %v", err)
	}
	if got.Email != "test@example.com" {
		t.Errorf("expected 'test@example.com', got %q", got.Email)
	}

	// Read by email
	got, err = repo.GetAccountByEmail("test@example.com")
	if err != nil {
		t.Fatalf("GetAccountByEmail failed: %v", err)
	}

	// Delete
	repo.DeleteAccount("acc-1")
	_, err = repo.GetAccount("acc-1")
	if err == nil {
		t.Error("expected error after account delete")
	}
}

func TestFakeRepoMemberCRUD(t *testing.T) {
	repo := newFakeRepo()

	m := &Member{
		ID:          "mem-1",
		WorkspaceID: "ws-1",
		AccountID:   "acc-1",
		Level:       "OWNER",
	}
	repo.StoreMember(m)

	// Read
	got, err := repo.GetMember("ws-1", "mem-1")
	if err != nil {
		t.Fatalf("GetMember failed: %v", err)
	}
	if got.Level != "OWNER" {
		t.Errorf("expected OWNER, got %q", got.Level)
	}

	// Read by account + workspace
	got, err = repo.GetMemberByAccountAndWorkspace("acc-1", "ws-1")
	if err != nil {
		t.Fatalf("GetMemberByAccountAndWorkspace failed: %v", err)
	}

	// Find by workspace
	repo.StoreMember(&Member{ID: "mem-2", WorkspaceID: "ws-1", AccountID: "acc-2", Level: "VIEWER"})
	all, err := repo.FindMembersByWorkspace("ws-1")
	if err != nil {
		t.Fatalf("FindMembersByWorkspace failed: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 members, got %d", len(all))
	}

	// Delete
	repo.DeleteMember("ws-1", "mem-1")
	_, err = repo.GetMember("ws-1", "mem-1")
	if err == nil {
		t.Error("expected error after member delete")
	}
}
