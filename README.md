# biodata
Track basic biodata

# Features
- Env secret for passphrase. See [this article](https://michael.stapelberg.ch/posts/2025-08-24-secret-management-with-sops-nix/)
- Set db as path
- Nix module

# Improvements
- Tests (unit, integration) - dry it up
- Install shellcheck
- Build a shell linter
- Install prettier
- Build an html linter
- index page: Add entry link is too big
- Error handling
  - ties into logging (better types too).
  - Frontend too
- Favicon
- Security (csp nonce, tokens, csrf)
- Too much logging (log at lowest level only or whatever?)
- Feature: Signing out
- All the TODOs
- Clean up Waist, Weight, BP string formatting to use %g
