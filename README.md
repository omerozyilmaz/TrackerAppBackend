# TrackerAppBackend

TrackerAppBackend Codes
yep

## CORS Error and Solution

### Encountered Error

A CORS (Cross-Origin Resource Sharing) error was encountered when making requests to the API from the frontend. The following error messages were observed in the terminal:

### Explanation of the Solution

- `Access-Control-Allow-Origin: "*"` - Allows requests from all origins. In a production environment, it should be restricted to specific domains for security.
- `Access-Control-Allow-Credentials: "true"` - Allows requests that include authentication credentials.
- `Access-Control-Allow-Headers` - Specifies the allowed HTTP headers.
- `Access-Control-Allow-Methods` - Specifies the allowed HTTP methods.
- Special handling for OPTIONS requests - Responds to preflight requests with a 204 (No Content) status code, indicating to the browser that it is safe to proceed with the actual request.

These changes ensure that the frontend application can access the API smoothly.

### Security Note

In a production environment, it is more secure to replace the `Access-Control-Allow-Origin` value with the actual domain of your frontend application instead of `*`:

```go
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://example.com"},
    AllowCredentials: true,
    AllowHeaders: []string{"Content-Type", "Authorization"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    MaxAge: 24 * time.Hour,
}))
```

This configuration will allow your frontend application to access the API.
