package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"
)

// fakeRepo is an in-memory repository for testing.
type fakeRepo struct {
	workspaces map[string]*Workspace
	accounts   map[string]*Account
	members    map[string]*Member
	invites    map[string]*Invite
	projects   map[string]*Project
	milestones map[string]*Milestone
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		workspaces: make(map[string]*Workspace),
		accounts:   make(map[string]*Account),
		members:    make(map[string]*Member),
		invites:    make(map[string]*Invite),
		projects:   make(map[string]*Project),
		milestones: make(map[string]*Milestone),
	}
}

func (f *fakeRepo) DB() *sqlx.DB                             { return nil }
func (f *fakeRepo) SetTx(tx *sqlx.Tx)                        {}
func (f *fakeRepo) StoreWorkspace(x *Workspace)              { f.workspaces[x.ID] = x }
func (f *fakeRepo) GetWorkspace(id string) (*Workspace, error) {
	if w, ok := f.workspaces[id]; ok {
		return w, nil
	}
	return nil, errNotFound
}
func (f *fakeRepo) GetWorkspaceByName(name string) (*Workspace, error) {
	for _, w := range f.workspaces {
		if w.Name == name {
			return w, nil
		}
	}
	return nil, errNotFound
}
func (f *fakeRepo) GetWorkspacesByAccount(id string) ([]*Workspace, error) {
	var ws []*Workspace
	for _, m := range f.members {
		if m.AccountID == id {
			if w, ok := f.workspaces[m.WorkspaceID]; ok {
				ws = append(ws, w)
			}
		}
	}
	return ws, nil
}
func (f *fakeRepo) DeleteWorkspace(id string) { delete(f.workspaces, id) }

func (f *fakeRepo) GetAccount(id string) (*Account, error) {
	if a, ok := f.accounts[id]; ok {
		return a, nil
	}
	return nil, errNotFound
}
func (f *fakeRepo) GetAccountByEmail(email string) (*Account, error) {
	for _, a := range f.accounts {
		if a.Email == email {
			return a, nil
		}
	}
	return nil, errNotFound
}
func (f *fakeRepo) GetAccountByConfirmationKey(key string) (*Account, error) { return nil, errNotFound }
func (f *fakeRepo) GetAccountByPasswordKey(key string) (*Account, error)     { return nil, errNotFound }
func (f *fakeRepo) FindAccountsByWorkspace(id string) ([]*Account, error)   { return nil, nil }
func (f *fakeRepo) StoreAccount(x *Account)                                 { f.accounts[x.ID] = x }
func (f *fakeRepo) DeleteAccount(id string)                                 { delete(f.accounts, id) }

func (f *fakeRepo) StoreMember(x *Member) { f.members[x.ID] = x }
func (f *fakeRepo) GetMember(workspaceID, id string) (*Member, error) {
	if m, ok := f.members[id]; ok && m.WorkspaceID == workspaceID {
		return m, nil
	}
	return nil, errNotFound
}
func (f *fakeRepo) GetMemberByAccountAndWorkspace(accountID, workspaceID string) (*Member, error) {
	for _, m := range f.members {
		if m.AccountID == accountID && m.WorkspaceID == workspaceID {
			return m, nil
		}
	}
	return nil, errNotFound
}
func (f *fakeRepo) GetMembersByAccount(id string) ([]*Member, error) {
	var ms []*Member
	for _, m := range f.members {
		if m.AccountID == id {
			ms = append(ms, m)
		}
	}
	return ms, nil
}
func (f *fakeRepo) GetMemberByEmail(workspaceID, email string) (*Member, error) { return nil, errNotFound }
func (f *fakeRepo) FindMembersByWorkspace(id string) ([]*Member, error) {
	var ms []*Member
	for _, m := range f.members {
		if m.WorkspaceID == id {
			ms = append(ms, m)
		}
	}
	return ms, nil
}
func (f *fakeRepo) DeleteMember(wsid, id string) { delete(f.members, id) }

func (f *fakeRepo) StoreInvite(x *Invite)           { f.invites[x.ID] = x }
func (f *fakeRepo) DeleteInvite(wsid, id string)     { delete(f.invites, id) }
func (f *fakeRepo) GetInviteByCode(code string) (*Invite, error) { return nil, errNotFound }
func (f *fakeRepo) GetInviteByEmail(wsid, email string) (*Invite, error) { return nil, errNotFound }
func (f *fakeRepo) GetInvite(workspaceID, id string) (*Invite, error)   { return nil, errNotFound }
func (f *fakeRepo) FindInvitesByWorkspace(wsid string) ([]*Invite, error) { return nil, nil }

func (f *fakeRepo) GetProjectByExternalLink(link string) (*Project, error) { return nil, errNotFound }
func (f *fakeRepo) GetProject(workspaceID, projectID string) (*Project, error) {
	if p, ok := f.projects[projectID]; ok && p.WorkspaceID == workspaceID {
		return p, nil
	}
	return nil, errNotFound
}
func (f *fakeRepo) FindProjectsByWorkspace(workspaceID string) ([]*Project, error) {
	var ps []*Project
	for _, p := range f.projects {
		if p.WorkspaceID == workspaceID {
			ps = append(ps, p)
		}
	}
	return ps, nil
}
func (f *fakeRepo) StoreProject(x *Project)                                { f.projects[x.ID] = x }
func (f *fakeRepo) DeleteProject(workspaceID, projectID string)            { delete(f.projects, projectID) }
func (f *fakeRepo) GetMilestone(workspaceID, milestoneID string) (*Milestone, error) { return nil, errNotFound }
func (f *fakeRepo) FindMilestonesByProject(workspaceID, projectID string) ([]*Milestone, error) { return nil, nil }
func (f *fakeRepo) StoreMilestone(x *Milestone)                            {}
func (f *fakeRepo) DeleteMilestone(workspaceID, milestoneID string)        {}
func (f *fakeRepo) GetWorkflow(workspaceID, workflowID string) (*Workflow, error) { return nil, errNotFound }
func (f *fakeRepo) FindWorkflowsByProject(workspaceID, projectID string) ([]*Workflow, error) { return nil, nil }
func (f *fakeRepo) StoreWorkflow(x *Workflow)                              {}
func (f *fakeRepo) DeleteWorkflow(workspaceID, workflowID string)          {}
func (f *fakeRepo) GetSubWorkflow(workspaceID, subWorkflowID string) (*SubWorkflow, error) { return nil, errNotFound }
func (f *fakeRepo) FindSubWorkflowsByProject(workspaceID, projectID string) ([]*SubWorkflow, error) { return nil, nil }
func (f *fakeRepo) FindSubWorkflowsByWorkflow(workspaceID, workflowID string) ([]*SubWorkflow, error) { return nil, nil }
func (f *fakeRepo) StoreSubWorkflow(x *SubWorkflow)                        {}
func (f *fakeRepo) DeleteSubWorkflow(workspaceID, workflowID string)       {}
func (f *fakeRepo) GetFeature(workspaceID, featureID string) (*Feature, error) { return nil, errNotFound }
func (f *fakeRepo) FindFeaturesByProject(workspaceID, projectID string) ([]*Feature, error) { return nil, nil }
func (f *fakeRepo) FindFeaturesByMilestoneAndSubWorkflow(workspaceID, mid, swid string) ([]*Feature, error) { return nil, nil }
func (f *fakeRepo) StoreFeature(x *Feature)                                {}
func (f *fakeRepo) DeleteFeature(workspaceID, workflowID string)           {}
func (f *fakeRepo) GetFeatureComment(workspaceID, ID string) (*FeatureComment, error) { return nil, errNotFound }
func (f *fakeRepo) FindFeatureCommentsByProject(workspaceID, projectID string) ([]*FeatureComment, error) { return nil, nil }
func (f *fakeRepo) StoreFeatureComment(x *FeatureComment)                  {}
func (f *fakeRepo) DeleteFeatureComment(workspaceID, commentID string)     {}
func (f *fakeRepo) GetFeatureCommentOwner(workspaceID, ID string) (*FeatureCommentOwner, error) { return nil, errNotFound }
func (f *fakeRepo) FindFeatureCommentOwnersByProject(workspaceID, projectID string) ([]*FeatureCommentOwner, error) { return nil, nil }
func (f *fakeRepo) StoreFeatureCommentOwner(x *FeatureCommentOwner)        {}
func (f *fakeRepo) GetFeatureCommentOwnerByFeatureComment(workspaceID, ID string) (*FeatureCommentOwner, error) { return nil, errNotFound }
func (f *fakeRepo) GetPersona(workspaceID, ID string) (*Persona, error)   { return nil, errNotFound }
func (f *fakeRepo) FindPersonasByProject(workspaceID, projectID string) ([]*Persona, error) { return nil, nil }
func (f *fakeRepo) StorePersona(x *Persona)                                {}
func (f *fakeRepo) DeletePersona(workspaceID, id string)                   {}
func (f *fakeRepo) GetWorkflowPersona(workspaceID, ID string) (*WorkflowPersona, error) { return nil, errNotFound }
func (f *fakeRepo) FindWorkflowPersonasByProject(workspaceID, projectID string) ([]*WorkflowPersona, error) { return nil, nil }
func (f *fakeRepo) StoreWorkflowPersona(x *WorkflowPersona)                {}
func (f *fakeRepo) DeleteWorkflowPersona(workspaceID, id string)           {}

var errNotFound = &testError{"not found"}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

// newTestRouter creates a chi router with all middleware and handlers for testing.
func newTestRouter() chi.Router {
	config := Configuration{
		Environment:        "development",
		Mode:               "hosted",
		AppSiteURL:         "http://localhost:3000",
		DbConnectionString: "test",
		JWTSecret:          "test-secret-key-for-testing",
		Port:               "5000",
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	auth := jwtauth.New("HS256", []byte(config.JWTSecret), nil)

	r.Use(jwtauth.Verifier(auth))
	r.Use(ContextSkeleton(config))
	r.Use(testTransaction(newFakeRepo()))
	r.Use(Auth(auth))
	r.Use(User())

	r.Route("/v1/users", usersAPI)
	r.Route("/v1/link", linkAPI)
	r.Route("/v1/account", accountAPI)
	r.Route("/v1/", workspaceAPI)

	return r
}

// testTransaction creates middleware that injects a fake repo into the service.
func testTransaction(repo Repository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := GetEnv(r).Service
			s.SetRepoObject(repo)
			next.ServeHTTP(w, r)
		})
	}
}

// apiTest provides helpers for API tests.
type apiTest struct {
	server *httptest.Server
	client *http.Client
}

func newAPITest(t *testing.T) *apiTest {
	t.Helper()
	router := newTestRouter()
	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)
	return &apiTest{
		server: srv,
		client: srv.Client(),
	}
}

func (a *apiTest) postJSON(path string, body interface{}) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", a.server.URL+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return a.client.Do(req)
}

func (a *apiTest) postJSONWithToken(path, token string, body interface{}) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", a.server.URL+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return a.client.Do(req)
}

func (a *apiTest) get(path, token, workspaceID string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", a.server.URL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if workspaceID != "" {
		req.Header.Set("Workspace", workspaceID)
	}
	return a.client.Do(req)
}

func readBody(t *testing.T, r *http.Response) string {
	t.Helper()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	r.Body.Close()
	return string(b)
}

func readJSON(r *http.Response, v interface{}) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return json.Unmarshal(b, v)
}

// --- Tests ---

func TestSignup(t *testing.T) {
	api := newAPITest(t)

	resp, err := api.postJSON("/v1/users/signup", map[string]string{
		"workspaceName": "testworkspace",
		"name":          "Test User",
		"email":         "test@example.com",
		"password":      "password123",
	})
	if err != nil {
		t.Fatalf("signup request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, readBody(t, resp))
	}

	var tokenResp TokenResponse
	if err := readJSON(resp, &tokenResp); err != nil {
		t.Fatalf("failed to decode token response: %v", err)
	}
	if tokenResp.Token == "" {
		t.Fatal("expected non-empty token")
	}
	t.Logf("signup successful, got token")
}

func TestLoginFlow(t *testing.T) {
	api := newAPITest(t)

	// Signup
	resp, err := api.postJSON("/v1/users/signup", map[string]string{
		"workspaceName": "logintest",
		"name":          "Login User",
		"email":         "login@example.com",
		"password":      "password123",
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("signup failed: %v %d", err, resp.StatusCode)
	}
	var signupResp TokenResponse
	readJSON(resp, &signupResp)

	// Login
	resp, err = api.postJSON("/v1/users/login", map[string]string{
		"email":    "login@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, readBody(t, resp))
	}

	var loginResp TokenResponse
	if err := readJSON(resp, &loginResp); err != nil {
		t.Fatalf("failed to decode login token: %v", err)
	}
	if loginResp.Token == "" {
		t.Fatal("expected non-empty login token")
	}
	t.Logf("login successful")
}

func TestUnauthenticatedAccess(t *testing.T) {
	api := newAPITest(t)

	resp, err := api.get("/v1/account/app", "", "")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
		body := readBody(t, resp)
		t.Fatalf("expected 401 or 403 for unauthenticated access, got %d: %s", resp.StatusCode, body)
	}
	t.Logf("unauthenticated access correctly denied: %d", resp.StatusCode)
}

func TestGetAppWithToken(t *testing.T) {
	api := newAPITest(t)

	// Signup to get a token
	resp, err := api.postJSON("/v1/users/signup", map[string]string{
		"workspaceName": "apptest",
		"name":          "App User",
		"email":         "appuser@example.com",
		"password":      "password123",
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("signup failed: %v %d", err, resp.StatusCode)
	}
	var tokenResp TokenResponse
	readJSON(resp, &tokenResp)

	// Fetch /v1/account/app with the token
	resp, err = api.get("/v1/account/app", tokenResp.Token, "")
	if err != nil {
		t.Fatalf("app request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, readBody(t, resp))
	}

	var appResp struct {
		Mode        string       `json:"mode"`
		Account     *Account     `json:"account"`
		Workspaces  []*Workspace `json:"workspaces"`
		Memberships []*Member    `json:"memberships"`
	}
	if err := readJSON(resp, &appResp); err != nil {
		t.Fatalf("failed to decode app response: %v", err)
	}
	if appResp.Account == nil {
		t.Fatal("expected non-nil account in app response")
	}
	if appResp.Account.Email != "appuser@example.com" {
		t.Fatalf("expected email 'appuser@example.com', got '%s'", appResp.Account.Email)
	}
	if len(appResp.Workspaces) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(appResp.Workspaces))
	}
	t.Logf("app response: mode=%s, email=%s, workspaces=%d", appResp.Mode, appResp.Account.Email, len(appResp.Workspaces))
}

func TestInvalidLogin(t *testing.T) {
	api := newAPITest(t)

	resp, err := api.postJSON("/v1/users/login", map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	body := readBody(t, resp)
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non-200 for invalid login, got 200: %s", body)
	}
	if !strings.Contains(body, "email") && !strings.Contains(body, "password") && !strings.Contains(body, "incorrect") {
		t.Logf("invalid login response: %d %s", resp.StatusCode, body)
	}
	t.Logf("invalid login correctly rejected: %d", resp.StatusCode)
}

func TestSignupValidation(t *testing.T) {
	api := newAPITest(t)

	tests := []struct {
		name    string
		payload map[string]string
	}{
		{"short password", map[string]string{"workspaceName": "vtest", "name": "Test", "email": "v@test.com", "password": "123"}},
		{"invalid email", map[string]string{"workspaceName": "vtest2", "name": "Test", "email": "not-an-email", "password": "password123"}},
		{"short name", map[string]string{"workspaceName": "vtest3", "name": "", "email": "valid@test.com", "password": "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := api.postJSON("/v1/users/signup", tt.payload)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.StatusCode == http.StatusOK {
				body := readBody(t, resp)
				t.Fatalf("expected validation failure for '%s', got 200: %s", tt.name, body)
			}
		})
	}
}
