# Coding Assignment – Static Site Hosting API

Joshua Arldt's coding assignment submission.

---

## API Features

| Feature             | Route                      | Method |
|---------------------|----------------------------|--------|
| Deploy ZIP          | `/deploy`                  | POST   |
| List Sites          | `/sites`                   | GET    |
| Delete Site         | `/sites/{name}`            | DELETE |
| Serve Static Files  | `/sites/{site}/{path...}`  | GET    |
| View Deploy Logs    | `/logs/{site}`             | GET    |

---

## Getting Started

This project includes a `Makefile` for convenience.

### Prerequisites

- Go 1.24.4 or higher

### Dependencies

  ```bash
  go mod tidy
  ```

### Build the Project

Compile the Go project and output the binary to `./bin/main`:

```bash
make build
```

### Run the Server

Builds the project and runs the server:

```bash
make run
```

---

The API will be available at [http://localhost:8080](http://localhost:8080)

---

### Run Tests

Runs all tests in the project:

```bash
make test
```

### Clean Up

Cleans build artifacts and removes deployments and temp folders:

```bash
make clean
```

This will:

- Remove the compiled binary
- Delete the `deployments/` folder
- Run `go clean` internally

## API Endpoints

### `POST /deploy`

Deploy a ZIP archive of your static site.

**Form fields:**

| Field     | Type   | Description               |
|-----------|--------|---------------------------|
| siteName  | string | Unique name for your site |
| zipFile   | file   | The `.zip` file to upload |

**Example:**

```bash
curl -X POST -F "siteName=mytestsite" -F "zipFile=@./site.zip" http://localhost:8080/deploy
```

---

### `GET /sites`

List all deployed sites.

**Response:**

```json
[
  {"name":"mytestsite","deployed_at":"2025-08-03T15:54:19Z"},
  {"name":"complexsite","deployed_at":"2025-08-03T15:51:19Z"}
]
```

---

### `DELETE /delete/{siteName}`

Delete a deployed site and its files.

**Example:**

```bash
curl -X DELETE http://localhost:8080/delete/mytestsite
```

---

### `GET /sites/{siteName}`

Serve static content from a deployed site.

---

### `GET /logs/{siteName}`

Fetch the deployment log history for a site.

**Example:**

```bash
curl http://localhost:8080/logs/mytestsite
```

**Response:**

```json
[
  {
    "site_name":"mytestsite",
    "timestamp":"2025-08-03T15:54:19Z",
    "ip_address":"127.0.0.1",
    "user_agent":"curl/8.7.1"
  }
]
```

---

## Deployment Folder Structure

Deployed sites are extracted to:

`deployments/{siteName}/`

---

## Security Notes

- Path traversal protection (`..`, slashes, zip-slip)
- All file access sandboxed under the `deployments/` directory

---

## Additional Notes

See NOTES.md for commentary as well as additional features, improvements, and ideas that I would have liked to explore or implement if I had more time on this assignment.
