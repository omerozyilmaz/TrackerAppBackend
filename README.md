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

- **Register:** `POST /register` or `POST /api/register`
- **Login:** `POST /login` or `POST /api/login`

### Job Postings

- **List Jobs:** `GET /jobs` or `GET /api/jobs`
- **Create Job:** `POST /jobs` or `POST /api/jobs`
- **Update Job:** `PUT /jobs/:id` or `PUT /api/jobs/:id`
- **Delete Job:** `DELETE /jobs/:id` or `DELETE /api/jobs/:id`
- **Update Job Status:** `PATCH /jobs/:id/status` or `PATCH /api/jobs/:id/status`
- **Move Job:** `PATCH /jobs/:id/move` or `PATCH /api/jobs/:id/move`

## URL Structure

The API supports URL structures both with and without the `/api` prefix. For example, the following two endpoints serve the same function:

- `POST /jobs`
- `POST /api/jobs`

This ensures compatibility with different clients and provides flexibility in URL structure.

## CORS Configuration

This API includes CORS (Cross-Origin Resource Sharing) configuration to allow requests from different origins, which is particularly important when your frontend application runs on a different domain.

### CORS Middleware

In the `middleware/cors.go` file, the CORS configuration is set up as follows:

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "https://tracker-app-amber.vercel.app")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

This configuration allows requests from `https://tracker-app-amber.vercel.app` and sets the necessary headers. It also handles OPTIONS requests efficiently.

## Deployment

This API is deployed on the Railway platform and can be accessed via the following URL:

```
https://trackerappbackend-production.up.railway.app
```

## Contributing

If you wish to contribute, please submit a pull request or open an issue.

## License

This project is licensed under the MIT License.

## Troubleshooting

### 307 Redirect Error

You may encounter 307 redirect errors related to trailing slashes in the API. To solve this issue, routes have been added that support URLs both with and without trailing slashes:

```go
// Support both URL formats
jobs.GET("", controllers.GetJobs)      // for /jobs
jobs.GET("/", controllers.GetJobs)     // for /jobs/
```

Additionally, Gin's automatic redirects have been disabled in the `main.go` file with the following settings:

```go
r.RedirectTrailingSlash = false
r.RedirectFixedPath = false
```

These changes ensure that the API correctly handles URLs both with and without trailing slashes.
