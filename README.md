# biodata
Track basic biodata

# Bugs
## Session Tracking
This does not seem to be working correctly. Mobile sessions have persisted long after they should

What could be happening here:
- The session table should be keyed off the session id, not user id
- That does not explain why the session.ID check is failing

# Improvements
## Tests
- Tests (unit, integration) - dry it up
  - Need test runners too. Good examples abound here
## Linting
- Lint sql
- Install prettier
- Build an html linter
## Clean up
- All the TODOs
- Clean up Waist, Weight, BP string formatting to use %g
## Logging
- Too much logging (log at lowest level only or whatever?)
- Custom handler to log out response status code
## Error handling
- Error handling
  - ties into logging (better types too).
  - Frontend too
## Security
- Security (csp nonce, tokens, csrf)

# Features
- index page: Add entry link is too big
- Favicon
- Update the page title to biotrak
- Signing out
- API get routes
  - Will need encoding

