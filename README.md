# biodata
Track basic biodata

# Bugs
## Session Tracking
This does not seem to be working correctly. Mobile sessions have persisted long after they should

What could be happening here:
- The session table should be keyed off the session id, not user id
- That does not explain why the session.ID check is failing

## SQL linting issues
Rename keys on new tables and push data over. Drop old tables

## HTML linting issues
Maybe just fix with `npx prettier . --write`

# Features
- API get routes
  - Will need encoding
- Signing out
- Favicon
- Charts

# Improvements
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
## Tests
- Tests (unit, integration) - dry it up
  - Need test runners too. Good examples abound here
## Static analysis
- Runner for static check and error check
## Integration Tests
- Try playwright?
