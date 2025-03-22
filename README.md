# Quickbooks `goth` Provider

An OAuth2 Code Flow Quickbooks `goth` provider. Inspired by the [Fitbit Goth Provider](https://github.com/markbates/goth/tree/master/providers/fitbit).

See https://github.com/markbates/goth for more information.

## Getting Started

```bash
go get github.com/munchpass/gothquickbooks
```

To use the provider:

```go
// Initialize the provider
// (replace the apiCtx values with your own values)
goth.UseProviders(qb.New(apiCtx.QuickbooksClientId, apiCtx.QuickbooksSecret, apiCtx.QuickbooksRedirectUrl))

// Create your HTTP Handlers.
r.GET("/quickbooks/start", qb.OAuthStart)
r.GET("/quickbooks/callback", qb.OAuthCallback)

// You're done!
// Just go to /quickbooks to start the OAuth flow.
```
