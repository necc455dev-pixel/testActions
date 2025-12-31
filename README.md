# testActions

Minimal Go API that exposes `/pulls` to fetch open pull requests from GitHub using the official [`google/go-github`](https://github.com/google/go-github) client (best-practice reference).

## Run locally
Set repository and (optionally) token, then start the server:
```powershell
$env:GITHUB_OWNER="your-org"
$env:GITHUB_REPO="your-repo"
$env:GITHUB_TOKEN="ghp_xxx" # optional for public repos, recommended for higher rate limits
go run .
```

Call the endpoint:
```bash
curl http://localhost:8080/pulls | jq
```

## GitHub Actions
The workflow in `.github/workflows/blank.yml` builds and tests on push/PR to `main` using Go 1.22.
