package goat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
	"os"
)

// Paths contains the paths for a single provider.
type Paths struct {
	// Consent specifies the path for consent requests. These get
	// redirected to the aouth2 provider's auth URL
	//
	// Requests to the consent URL require a redirect_url
	// parameter that is used to redirect to after completion of
	// the full consent flow.
	Consent string

	// Code specifies the path configured as the RedirectURL with
	// the oauth provider.  This is expected to contain the code
	// that can be exchanged for the auth token.
	Code string

	// SetRefreshToken specifies the path where SetRefreshToken
	// requests happen. The refresh_token parameter (query or form
	// value) should hold the actual refresh token.  This is first
	// validated and so must be a real refresh token.
	SetRefreshToken string

	// GetAccessToken specifies the path where the acess token can
	// be fetched from. The response is a JSON with valid
	// refresh_token field or an empth hash (indicating no
	// credentials available).
	GetAccessToken string
}

// Provider specifies a specific oauth provider + associated config.
type Provider struct {
	// Name of the provider.  Used with TokenStore & SessionStore
	Name string

	// Config should specify endpoints and such
	oauth2.Config

	// Paths should specify the endpoints for the Goat
	// server. This should be unique and not shared between
	// providers.
	Paths Paths

	// AuthURLParams specifies additional auth url options.
	// Offline access is included automatically.
	AuthURLParams map[string]string
}

func (p *Provider) authCodeOptions() []oauth2.AuthCodeOption {
	opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
	for k, v := range p.AuthURLParams {
		opts = append(opts, oauth2.SetAuthURLParam(k, v))
	}
	return opts
}

// TokenStore is the interface a token store must implement.  It is a
// simple key value store.
type TokenStore interface {
	Get(ctx context.Context, provider, user string) ([]byte, error)
	Set(ctx context.Context, provider, user string, token []byte) error
}

// SessionStore is the interface a session store must implement.  It
// is a simple key value store.
//
// Note:
//
// The session store should expire the values quickly (1m is
// recommended).
//
// The sesion store should also delete the keys once they are
// fetched.  That is, all Get calls should effectively also delete
// immediately.
type SessionStore interface {
	Get(ctx context.Context, provider, nonce string) ([]byte, error)
	Set(ctx context.Context, provider, nonce string, session []byte) error
}

// EncrypterDecrypter is the interface to be implemented by an
// encryption engine.
type EncrypterDecrypter interface {
	Encrypt(ctx context.Context, data []byte) ([]byte, error)
	Decrypt(ctx context.Context, data []byte) ([]byte, error)
}

// Handler is a http.Handler that serves all the http requests.
type Handler struct {
	// Providers contains the list of providers.
	Providers []*Provider

	// AuthenticatedUser is a function that checks that the
	// request is authenticated and gets the identity of the
	// requester.  All Goat requests must be authenticated.
	AuthenticatedUser func(r *http.Request) (string, error)

	// Tokens is the token storage.  All the tokens are encrypted
	// before-hand and so this does not require special handling
	// for security.
	Tokens TokenStore

	// Sessions is used to take a random nonce to prevent CSRF
	// attacks.  It is a transient store.
	Sessions SessionStore

	// Secrets provides the encryption/description support.
	Secrets EncrypterDecrypter
}

// ServeHTTP implements http.Handler.  If better control over error
// handling and logging is required, please directly call Handle
// which does not write any responses other than redirect.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err, status := h.Handle(r)
	if status == http.StatusTemporaryRedirect {
		http.Redirect(w, r, body.(string), http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		log.Println("Error", err)
	}
	if body != nil {
		w.Header().Add("Content-type", "application/json")
	} else if err != nil {
		w.Header().Add("Content-type", "text/plain")
		body = err.Error()
	}

	w.WriteHeader(status)
	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Println("Error", err)
		}
	}
}

// Handle does all the work for handling HTTP requests but does not
// respond.  It returns the response for the caller to respond.
//
// For redirects, the body contains the URL to be redirected to.
func (h Handler) Handle(r *http.Request) (body interface{}, err error, status int) {
	user, err := h.AuthenticatedUser(r)
	if err != nil {
		return nil, err, http.StatusUnauthorized
	}

	for _, p := range h.Providers {
		switch r.URL.Path {
		case p.Paths.Consent:
			return h.handleConsent(r, user, p)
		case p.Paths.Code:
			return h.handleCode(r, user, p)
		case p.Paths.SetRefreshToken:
			return h.handleSetRefreshToken(r, user, p)
		case p.Paths.GetAccessToken:
			return h.handleGetAccessToken(r, user, p)
		}
	}
	return nil, errors.New("not found"), http.StatusNotFound
}

func (h Handler) handleConsent(r *http.Request, user string, p *Provider) (interface{}, error, int) {
	redirectURL := r.FormValue("redirect_url")
	if redirectURL == "" {
		return nil, errors.New("missing redirect_url"), http.StatusBadRequest
	}

	state := fmt.Sprintf("%x-%x", rand.Int63(), rand.Int63())
	if err := h.Sessions.Set(r.Context(), p.Name, state, []byte(redirectURL)); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	url := p.Config.AuthCodeURL(state, p.authCodeOptions()...)
	return url, nil, http.StatusTemporaryRedirect
}

func (h Handler) handleCode(r *http.Request, user string, p *Provider) (interface{}, error, int) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	if state == "" || code == "" {
		return nil, errors.New("missing state or code"), http.StatusBadRequest
	}

	redirectURL, err := h.Sessions.Get(r.Context(), p.Name, state)
	if err != nil {
		err = errors.Wrap(err, "bad state")
		return nil, err, http.StatusBadRequest
	}

	token, err := p.Config.Exchange(r.Context(), code, p.authCodeOptions()...)
	if err != nil {
		err = errors.Wrap(err, "token exchange")
		return nil, err, http.StatusBadRequest
	}

	if err = h.saveToken(r.Context(), user, p, token); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return string(redirectURL), nil, http.StatusTemporaryRedirect
}

func (h Handler) handleSetRefreshToken(r *http.Request, user string, p *Provider) (interface{}, error, int) {
	token := &oauth2.Token{RefreshToken: r.FormValue("refresh_token")}
	if token.RefreshToken == "" {
		return nil, errors.New("missing refresh_token"), http.StatusBadRequest
	}

	tok, err := p.Config.TokenSource(r.Context(), token).Token()
	if err != nil {
		return nil, errors.New("invalid refresh_token"), http.StatusBadRequest
	}

	if err = h.saveToken(r.Context(), user, p, tok); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return map[string]string{}, nil, http.StatusOK
}

func (h Handler) handleGetAccessToken(r *http.Request, user string, p *Provider) (interface{}, error, int) {
	data, err := h.Tokens.Get(r.Context(), p.Name, user)
	if os.IsNotExist(err) {
		return &oauth2.Token{}, nil, http.StatusOK
	}
	if err != nil {
		return nil, errors.Wrap(err, "read token"), http.StatusInternalServerError
	}

	if data, err = h.Secrets.Decrypt(r.Context(), data); err != nil {
		return nil, errors.Wrap(err, "decrypt"), http.StatusInternalServerError
	}

	var token oauth2.Token
	if err = json.Unmarshal(data, &token); err != nil {
		return nil, errors.Wrap(err, "unmarshal"), http.StatusInternalServerError
	}

	if token.Valid() {
		token.RefreshToken = ""
		return &token, nil, http.StatusOK
	}

	tok, err := p.Config.TokenSource(r.Context(), &token).Token()
	if err != nil {
		// TODO: in case of a permanent error, it would be better
		// to return "", err, http.StatusOK and force the client to
		// refresh again
		return nil, errors.New("invalid refresh_token"), http.StatusInternalServerError
	}

	if err = h.saveToken(r.Context(), user, p, tok); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	tok.RefreshToken = ""
	return tok, nil, http.StatusOK
}

func (h Handler) saveToken(ctx context.Context, user string, p *Provider, token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	if data, err = h.Secrets.Encrypt(ctx, data); err != nil {
		return errors.Wrap(err, "encrypt")
	}

	if err = h.Tokens.Set(ctx, p.Name, user, data); err != nil {
		return errors.Wrap(err, "saving token")
	}
	return nil
}
