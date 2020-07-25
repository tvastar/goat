// Package goat implements a native Go oauth access token server.
//
// The expected flow for this server:
//
// 1. Client requests an access token for a provider, say, "google",
// for the current user.
//
// 2. The goat server returns an access token (or an empty string if
// one isn't available)
//
// 3. If an access token is not available, the client sends the
// browser to /google/url?redirect_url=something on the goat server.
//
// 4. The goat server redirects to the consent page for the provider
// ("google").  Once consent succeeds, the server then redirects the
// browser to the redirect URL configured for the provider which
// souuld be the goat server.
//
// 5. When redirected from the provider, the goat server gets the
// access token via standard oauth2 flows. This is then saved to its
// secure storage. The goat server than redirects to the URI provided
// in step3.
//
// 6. The client gets the redirect to this page and starts with step
// 1 which shoudl now succeed.
//
//
// Session & Secure Storage
//
// The goat server requries three components for security:
//
// 1. A session storage to hold the state parameter to prevent CSRF
// attacks.
//
// 2. An token storage to store encrypted tokens.
//
// 3. An encrpyption engine to encrypt/decrypt tokens.
//
// Example implementations for these are available in
// sub-directories.
//
// An example that puts this together is available in the cmd/goat
// example.
package goat
