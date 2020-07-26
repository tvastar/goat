package goat_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tvastar/goat"
	"golang.org/x/oauth2"
)

func TestAll(t *testing.T) {
	s := suite{}
	t.Run("Consent", s.TestConsent)
	t.Run("AccessToken", s.TestAccessToken)
	t.Run("AccessTokenRefresh", s.TestAccessTokenRefresh)
	t.Run("RefreshToken", s.TestRefreshToken)
}

type suite struct{}

func (s suite) TestConsent(t *testing.T) {
	// startServers creates a fake provider with the provided paths
	// as well as an oauth server to test against with credentials
	// already configured.
	goats, oauths := s.startServers(goat.Paths{
		Consent:         "/fakeprov/consent",
		Code:            "/fakeprov/code",
		SetRefreshToken: "/fakeprov/refreshtoken",
		GetAccessToken:  "/fakeprov/accesstoken",
	})
	defer goats.Close()
	defer oauths.Close()

	// we create a fake resorce to redirect to after completion of
	// consent flow.
	fakeResource := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("consent flow succeeded"))
	}))
	defer fakeResource.Close()

	// fetch consent URL and ensure result is "consent flow succeeded".
	v := url.Values{"redirect_url": {fakeResource.URL}}
	body, err := s.fetchBody(http.Get(goats.URL + "/fakeprov/consent?" + v.Encode()))

	if body != "consent flow succeeded" || err != nil {
		t.Fatal("redirect didn't go as planned?", body, err)
	}
}

func (s suite) TestAccessToken(t *testing.T) {
	// startServers creates a fake provider with the provided paths
	// as well as an oauth server to test against with credentials
	// already configured.
	goats, oauths := s.startServers(goat.Paths{
		Consent:         "/fakeprov/consent",
		Code:            "/fakeprov/code",
		SetRefreshToken: "/fakeprov/refreshtoken",
		GetAccessToken:  "/fakeprov/accesstoken",
	})
	defer goats.Close()
	defer oauths.Close()

	// we create a fake resorce to redirect to after completion of
	// consent flow.
	fakeResource := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("consent flow succeeded"))
	}))
	defer fakeResource.Close()

	// fetch consent URL and ensure result is "consent flow succeeded".
	v := url.Values{"redirect_url": {fakeResource.URL}}
	body, err := s.fetchBody(http.Get(goats.URL + "/fakeprov/consent?" + v.Encode()))

	if body != "consent flow succeeded" || err != nil {
		t.Fatal("redirect didn't go as planned?", body, err)
	}

	authResource := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if oauths.isAuthorized(r) {
			w.Write([]byte("Authorized!"))
		}
	}))
	defer authResource.Close()

	token, err := s.fetchToken(goats.URL + "/fakeprov/accesstoken")
	if err != nil {
		t.Fatal("Could not get token", err)
	}

	client := (&oauth2.Config{}).Client(context.Background(), token)
	body, err = s.fetchBody(client.Get(authResource.URL))
	if err != nil || body != "Authorized!" {
		t.Fatal("Unexpected resouruce fetch", body, err)
	}
}

func (s suite) TestAccessTokenRefresh(t *testing.T) {
	// startServers creates a fake provider with the provided paths
	// as well as an oauth server to test against with credentials
	// already configured.
	goats, oauths := s.startServers(goat.Paths{
		Consent:         "/fakeprov/consent",
		Code:            "/fakeprov/code",
		SetRefreshToken: "/fakeprov/refreshtoken",
		GetAccessToken:  "/fakeprov/accesstoken",
	})
	defer goats.Close()
	defer oauths.Close()

	// we create a fake resorce to redirect to after completion of
	// consent flow.
	fakeResource := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("consent flow succeeded"))
	}))
	defer fakeResource.Close()

	// fetch consent URL and ensure result is "consent flow succeeded".
	v := url.Values{"redirect_url": {fakeResource.URL}}
	body, err := s.fetchBody(http.Get(goats.URL + "/fakeprov/consent?" + v.Encode()))

	if body != "consent flow succeeded" || err != nil {
		t.Fatal("redirect didn't go as planned?", body, err)
	}

	authResource := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if oauths.isAuthorized(r) {
			w.Write([]byte("Authorized!"))
		}
	}))
	defer authResource.Close()

	for kk := 0; kk < 2; kk++ {
		token, err := s.fetchToken(goats.URL + "/fakeprov/accesstoken")
		if err != nil {
			t.Fatal("Could not get token", err)
		}

		if !token.Valid() {
			t.Fatal("Got expired token", err)
		}

		client := (&oauth2.Config{}).Client(context.Background(), token)
		body, err = s.fetchBody(client.Get(authResource.URL))
		if err != nil || body != "Authorized!" {
			t.Fatal("Unexpected resouruce fetch", body, err)
		}

		if kk != 0 {
			continue
		}
		time.Sleep(time.Until(token.Expiry))
		if token.Valid() {
			t.Fatal("invalid test -- token must be expired!", time.Until(token.Expiry))
		}
	}
}

func (s suite) TestRefreshToken(t *testing.T) {
	// startServers creates a fake provider with the provided paths
	// as well as an oauth server to test against with credentials
	// already configured.
	goats, oauths := s.startServers(goat.Paths{
		Consent:         "/fakeprov/consent",
		Code:            "/fakeprov/code",
		SetRefreshToken: "/fakeprov/refreshtoken",
		GetAccessToken:  "/fakeprov/accesstoken",
	})
	defer goats.Close()
	defer oauths.Close()

	v := url.Values{"refresh_token": {oauths.tokens[0].RefreshToken}}
	body, err := s.fetchBody(http.Get(goats.URL + "/fakeprov/refreshtoken?" + v.Encode()))
	if err != nil || body != "{}\n" {
		t.Fatal("Unexpected refreshtoken response", body, err)
	}

	authResource := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if oauths.isAuthorized(r) {
			w.Write([]byte("Authorized!"))
		}
	}))
	defer authResource.Close()

	token, err := s.fetchToken(goats.URL + "/fakeprov/accesstoken")
	if err != nil {
		t.Fatal("Could not get token", err)
	}

	client := (&oauth2.Config{}).Client(context.Background(), token)
	body, err = s.fetchBody(client.Get(authResource.URL))
	if err != nil || body != "Authorized!" {
		t.Fatal("Unexpected resouruce fetch", body, err)
	}
}

func (s suite) fetchToken(url string) (*oauth2.Token, error) {
	body, err := s.fetchBody(http.Get(url))
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	return &token, json.Unmarshal([]byte(body), &token)
}

func (suite) fetchBody(resp *http.Response, err error) (string, error) {
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func (suite) startServers(paths goat.Paths) (*httptest.Server, *oauthServer) {
	oauths := &oauthServer{}
	oauths.Server = httptest.NewServer(oauths)
	h := goat.Handler{
		Providers: []*goat.Provider{{
			Name: "fakeprov",
			Config: oauth2.Config{
				ClientID:     "fake client ID",
				ClientSecret: "fake client secret",
				Endpoint: oauth2.Endpoint{
					AuthURL:   oauths.URL + "/auth",
					TokenURL:  oauths.URL + "/token",
					AuthStyle: oauth2.AuthStyleInParams,
				},
				RedirectURL: "", // filled later
				Scopes:      []string{"scope1", "scope2"},
			},
			Paths: paths,
		}},
		AuthenticatedUser: func(r *http.Request) (string, error) {
			return "foo", nil // no authentication in test
		},
		Tokens:   &memstore{items: map[[2]string][]byte{}},
		Sessions: &memstore{items: map[[2]string][]byte{}},
		Secrets:  &secrets{},
	}
	goats := httptest.NewServer(h)
	h.Providers[0].Config.RedirectURL = goats.URL + "/fakeprov/code"
	oauths.init(h.Providers[0].Config)
	return goats, oauths
}

type oauthServer struct {
	*httptest.Server
	clientID, secret, redirectURL, scopes string

	codes  []string
	tokens []oauth2.Token
}

func (s *oauthServer) init(c oauth2.Config) {
	s.clientID = c.ClientID
	s.secret = c.ClientSecret
	s.redirectURL = c.RedirectURL
	s.scopes = strings.Join(c.Scopes, " ")
	s.codes = []string{"my code"}
	s.tokens = []oauth2.Token{
		{AccessToken: "access token!", RefreshToken: "refresh token!"},
		{AccessToken: "access token2!", RefreshToken: "refresh token2!"},
	}
}

func (s *oauthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth":
		s.serveAuth(w, r)
	case "/token":
		s.serveToken(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *oauthServer) serveAuth(w http.ResponseWriter, r *http.Request) {
	authType := r.FormValue("response_type")
	clientID := r.FormValue("client_id")
	redirectURL := r.FormValue("redirect_uri")
	scopes := r.FormValue("scope")

	// validate
	if authType != "code" || s.clientID != clientID ||
		s.redirectURL != redirectURL || s.scopes != scopes {
		http.Error(w, "invalid args", http.StatusBadRequest)
		return
	}

	v := url.Values{"code": {s.codes[0]}, "state": {r.FormValue("state")}}
	http.Redirect(w, r, redirectURL+"?"+v.Encode(), http.StatusTemporaryRedirect)
}

func (s *oauthServer) serveToken(w http.ResponseWriter, r *http.Request) {
	data, _ := httputil.DumpRequest(r, true)
	badRequest := fmt.Sprint("bad request: " + string(data))

	clientID := r.FormValue("client_id")
	secret := r.FormValue("client_secret")
	if clientID != s.clientID || secret != s.secret {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	switch r.FormValue("grant_type") {
	case "authorization_code":
		if r.FormValue("code") != s.codes[0] {
			http.Error(w, badRequest, http.StatusBadRequest)
			return
		}
		s.codes = s.codes[1:]
	case "refresh_token":
		if r.FormValue("refresh_token") != s.tokens[0].RefreshToken {
			http.Error(w, badRequest, http.StatusBadRequest)
			return
		}
		s.tokens = s.tokens[1:]
	default:
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"access_token":  s.tokens[0].AccessToken,
		"token_type":    "Bearer",
		"refresh_token": s.tokens[0].RefreshToken,
		"expires_in":    11, // 11s > oauth2 grace period of 10s
	}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "marshal json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *oauthServer) isAuthorized(r *http.Request) bool {
	return r.Header.Get("Authorization") == "Bearer "+s.tokens[0].AccessToken
}

type memstore struct {
	items map[[2]string][]byte
}

func (m *memstore) Get(_ context.Context, provider, key string) ([]byte, error) {
	if v, ok := m.items[[2]string{provider, key}]; ok {
		return v, nil
	}
	return nil, os.ErrNotExist
}

func (m *memstore) Set(_ context.Context, provider, key string, data []byte) error {
	m.items[[2]string{provider, key}] = data
	return nil
}

type secrets struct{}

func (s *secrets) Encrypt(_ context.Context, data []byte) ([]byte, error) {
	return data, nil
}

func (s *secrets) Decrypt(_ context.Context, data []byte) ([]byte, error) {
	return data, nil
}
