# biodata

Track basic biodata

# Features

- Favicon
- Charts

- TODO

# Improvements

## Security

- CSRF
- CSP (policy, nonce)

## Integration Tests

- Try playwright?

## Structure logging

- Experiment with `log/slog`

## Logging

- Custom handler to log out response status code. Claude recommends...

```
  type statusRecorder struct {
      http.ResponseWriter
      statusCode int
  }

  func (r *statusRecorder) WriteHeader(statusCode int) {
      r.statusCode = statusCode
      r.ResponseWriter.WriteHeader(statusCode)
  }

  func LoggingWrapper(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          start := time.Now()
          recorder := &statusRecorder{
              ResponseWriter: w,
              statusCode: 200, // default to 200 if WriteHeader is never called
          }
          next.ServeHTTP(recorder, r)
          log.Printf("%s %s %s %d %s", r.Method, r.Host, r.URL.Path,
              recorder.statusCode, time.Since(start))
      })
  }
```
