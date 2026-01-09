# File Sharing Service

A simple file-sharing service built with Go and Cloudflare R2. Upload files and get shareable download links.

## What It Does

- Upload files through a web interface
- Store files in Cloudflare R2 cloud storage
- Generate unique download links for each file
- Download links expire after 15 minutes for security

## Prerequisites

- Go 1.25 or higher
- Cloudflare R2 account (free tier available)

## Setup Instructions

### 1. Clone and Install

```bash
git clone https://github.com/heydeepakch/file-sharing-service.git
cd file-sharing-service
go mod download
```

### 2. Create R2 Bucket

1. Go to https://dash.cloudflare.com/
2. Navigate to R2
3. Click "Create Bucket" and name it (e.g., `my-file-sharing`)
4. Go to "Manage R2 API Tokens"
5. Create a new API token with Edit permissions
6. Save these three values:
   - Access Key ID
   - Secret Access Key
   - Endpoint URL (format: `https://[account-id].r2.cloudflarestorage.com`)

### 3. Configure Environment

Create a `.env` file in the project root:

```env
R2_ACCESS_KEY=your_access_key_id
R2_SECRET_KEY=your_secret_access_key
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
R2_BUCKET_NAME=my-file-sharing
PORT=8080
BASE_URL=http://localhost:8080
```

Replace the placeholder values with your actual R2 credentials.

### 4. Run the Server

```bash
go run .
```

The server will start on http://localhost:8080

### 5. Upload Files

Open your browser and go to http://localhost:8080

## Supported File Types

- PDF (.pdf)
- Images (.jpg, .jpeg, .png, .gif)
- Archives (.zip)

Maximum file size: 100 MB

## How It Works

**Upload Process:**

1. User uploads a file via web form
2. File is validated for type and size
3. File is uploaded to Cloudflare R2
4. Metadata is saved to local database (db.json)
5. User receives a download link

**Download Process:**

1. User clicks download link
2. Server looks up file in database
3. Server generates a temporary R2 URL (valid for 15 minutes)
4. User is redirected to download directly from R2

## Project Structure

```
file-sharing-service/
├── main.go           # Main application and HTTP handlers
├── storage.go        # Database functions (JSON file storage)
├── cloud.go          # R2 connection configuration
├── static/
│   └── index.html    # Upload form
├── .env              # Environment variables (create this)
├── db.json           # File metadata (auto-generated)
└── go.mod            # Go dependencies
```

## Environment Variables

| Variable       | Description                                         |
| -------------- | --------------------------------------------------- |
| R2_ACCESS_KEY  | Your R2 access key ID                               |
| R2_SECRET_KEY  | Your R2 secret access key                           |
| R2_ENDPOINT    | Your R2 endpoint URL                                |
| R2_BUCKET_NAME | Name of your R2 bucket                              |
| PORT           | Server port (default: 8080)                         |
| BASE_URL       | Base URL for links (default: http://localhost:8080) |

## API Endpoints

**GET /**

- Returns upload form

**POST /upload**

- Upload a file
- Form field: `file`
- Returns: Download link

**GET /download/{id}**

- Download a file
- Redirects to temporary R2 URL
