# TrackerAppBackend

TrackerAppBackend is a RESTful API designed to help users track their job applications. It provides functionalities for creating, updating, deleting, and listing job postings. Additionally, it allows users to update the status of job postings and move them between different columns.

## Features

- **User Authentication:** Secure user login using JWT-based authentication.
- **Job Management:** Users can create, update, delete, and list job postings.
- **Status Update:** Update the status of job postings (wishlist, applied, interview).
- **Job Movement:** Move job postings between different columns.

## Setup

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/username/job-tracker-api.git
   cd job-tracker-api
   ```

2. **Install Dependencies:**

   ```bash
   go mod tidy
   ```

3. **Configure the Database:**

   Update the database connection settings in the `config/db.go` file.

4. **Start the Server:**

   ```bash
   go run main.go
   ```

## API Usage

### Authentication

- **Register:** `POST /api/auth/register`
- **Login:** `POST /api/auth/login`

### Job Postings

- **List Jobs:** `GET /api/jobs`
- **Create Job:** `POST /api/jobs`
- **Update Job:** `PUT /api/jobs/:id`
- **Delete Job:** `DELETE /api/jobs/:id`
- **Update Job Status:** `PATCH /api/jobs/:id/status`
- **Move Job:** `PATCH /api/jobs/:id/move`

## CORS Configuration

This API includes CORS (Cross-Origin Resource Sharing) configuration to allow requests from different origins, which is particularly important when your frontend application runs on a different port.

### CORS Middleware

In the `routes/jobs.go` file, the CORS configuration is set up as follows:

```go
jobs.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }

    c.Next()
})
```

This configuration allows requests from `http://localhost:5173` and sets the necessary headers. It also handles OPTIONS requests efficiently.

## Contributing

If you wish to contribute, please submit a pull request or open an issue.

## License

This project is licensed under the MIT License.

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
