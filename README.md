# remont

Backend for a renovation company: projects, construction stages, workers, documents and progress photos in one place. Clients follow their renovation online, managers and foremen run the work.

![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7)
![SQLite](https://img.shields.io/badge/SQLite-GORM-003B57?logo=sqlite&logoColor=white)
[![ci](https://github.com/HowkaCoder/remont/actions/workflows/ci.yml/badge.svg)](https://github.com/HowkaCoder/remont/actions/workflows/ci.yml)

## Features

- Projects with members: client, manager, foreman (`prorab`) and workers
- Stages of work (`states`) with assigned workers and repair details
- Document and photo folders per project, file upload and download
- Project characteristics (area, rooms and so on)
- Role based access control: every route is protected by a permission, roles and permissions are stored in the database and can be changed at runtime

## Access control

```
user ──< user_roles >── role ──< role_permissions >── permission
```

Each request goes through two middlewares: `Authenticate` checks the JWT, `Require("create_document")` checks that one of the user's roles has the permission.

| Role | Can do |
|---|---|
| `admin` | everything, including users, roles and permissions |
| `manager`, `prorab` | manage projects, stages, documents, photos |
| `worker` | projects, photos, own stages |
| `client` | read their project, documents and photos |

Roles and permissions are seeded on startup and the seeding is idempotent. Public registration creates a `client`.

## Architecture

```
main.go                  config, seeding, routes
internal/
  config/                environment variables
  app.go                 database, migrations, seeding
  app/
    entity/              GORM models
    repository/          database access
    usecase/             business logic
    handler/             HTTP handlers, auth and upload helpers
```

## Getting started

```bash
cp .env.example .env      # set JWT_SECRET and admin credentials
go run .
```

or

```bash
docker compose up --build
```

Set `SEED_DEMO=true` and `DEMO_PASSWORD` to create demo users `worker@demo.local`, `manager@demo.local`, `client@demo.local`, `prorab@demo.local`.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3000` | HTTP port |
| `DB_PATH` | `database/database.db` | SQLite file |
| `UPLOADS_DIR` | `uploads` | Uploaded documents and photos |
| `PUBLIC_URL` | `http://localhost:3000` | Base URL for photo links |
| `JWT_SECRET` | — | Required |
| `ADMIN_EMAIL`, `ADMIN_PASSWORD` | — | Creates an admin on first start |
| `SEED_DEMO`, `DEMO_PASSWORD` | `false` | Demo users |

## API

| Method | Path | Permission |
|---|---|---|
| `POST` | `/register`, `/login` | public |
| `GET` | `/get-profile` | any user |
| `GET`, `DELETE` | `/users[/:id]` | `manage_users` |
| `GET`, `POST` | `/roles`, `/permissions`, `/roles/assign`, `/users/assign-role` | `manage_roles` |
| `GET`, `POST`, `PATCH`, `DELETE` | `/api/projects[/:id]` | `*_project` |
| `GET` | `/api/projects/users/:id` | `get_project_by_id` |
| `GET`, `POST`, `PATCH`, `DELETE` | `/api/docs[/:id]`, `/api/docs/folder/:folderID` | `*_document` |
| `GET`, `POST`, `PATCH`, `DELETE` | `/api/photos[/:id]`, `/api/photos/folder/:folderID` | `*_photo` |
| `GET`, `POST`, `PATCH`, `DELETE` | `/api/document-folders`, `/api/photo-folder` | `*_folder` |
| `GET`, `POST`, `PATCH`, `DELETE` | `/api/states`, `/api/details`, `/api/chars` | `*_state`, `*_detail`, `*_char` |

Files are uploaded as `multipart/form-data` with fields `file`, `projectID`, `folderID` and `name` / `title`. Stored files get random names, so the original file name can't be used to write outside the uploads folder.

## Tests

```bash
go test ./...
```

HTTP tests cover login, registration, role checks, hidden password hashes and document upload and removal.
