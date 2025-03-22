# Echo Examples

## Getting Started

The following will start an echo server that implements the OAuth flows:

```bash
QB_CLIENT_ID=<YOUR_CLIENT_ID> QB_CLIENT_SECRET=<YOUR_CLIENT_SECRET> go run main.go
```

Go to `http://localhost:3000/quickbooks/start` to start the OAuth flow. You'll know it worked by checking that it correctly fetched the user info properly in the server logs (stdout).

## How does it work?

Start the OAuth flow with the `OAuthStart` handler (i.e. at `/quickbooks/start`).

Then, the redirect uri should be configured to redirect to whatever endpoint `OAuthCallback` is mapped to (i.e. `/quickbooks/callback`).
