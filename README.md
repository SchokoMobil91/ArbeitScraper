# ArbeitScraper - Dockerized Setup Guide

This project consists of three services: **PostgreSQL**, **Go backend**, and **React frontend**, all running through **Docker Compose**.

The guide explains:
- Project structure
- First-time build instructions
- How to start/stop services after the initial build
- Useful Docker commands

---

## Project Structure
```
ArbeitScraper/
│
├── backend/        # Go service
│   ├── main.go
│   ├── database/
│   ├── scraper/
│   ├── go.mod
│   └── Dockerfile
│
├── frontend/       # React + Vite service
│   ├── src/
│   ├── package.json
│   └── Dockerfile
│
├── docker-compose.yml
└── README.md
```

---

## Docker Compose Overview
The `docker-compose.yml` contains three services:

### **1. PostgreSQL Database**
- Port: **5432**
- Persists data using named volume `dbdata`
- Health check ensures backend waits until DB is ready

### **2. Go Backend**
- Built from `backend/Dockerfile`
- Chromium installed to support scraping
- Exposed on port **8080**
- Connects to database via:
  ```env
  DATABASE_URL=user=postgres password=123456 dbname=arbeitdb host=db port=5432 sslmode=disable
  ```

### **3. React Frontend (Vite)**
- Built from `frontend/Dockerfile`
- Exposed on port **5173**

---

## First-Time Setup (Initial Build)
Run these commands after cloning the repository for the first time.

### **1. Ensure Docker Desktop is running**
Docker must be open before executing any commands.

### **2. Clean old containers (optional but recommended)**
```bash
docker compose down
docker container prune -f
```

### **3. Build everything for the first time**
```bash
docker compose up --build
```

This will:
- Pull base images
- Build backend (Go)
- Build frontend (Node)
- Start PostgreSQL
- Start backend → waits for DB healthcheck
- Start frontend

If successful, services will be available at:
- **Frontend:** http://localhost:5173
- **Backend:** http://localhost:8080
- **Database:** localhost:5432

---

## Starting Services After First Build
When images are already built, run:

```bash
docker compose up
```

No `--build` needed unless backend or frontend Dockerfiles are modified.

### To run in background:
```bash
docker compose up -d
```

### To stop the stack:
```bash
docker compose down
```

---

## Rebuilding Only One Service
Rebuild backend only:
```bash
docker compose build backend
```

Rebuild frontend only:
```bash
docker compose build frontend
```

Then restart:
```bash
docker compose up
```

---

## Useful Development Commands
### View running containers
```bash
docker ps
```

### View logs
Backend logs:
```bash
docker logs -f arbeit-backend
```
Frontend logs:
```bash
docker logs -f arbeit-frontend
```
Database logs:
```bash
docker logs -f arbeitscraper-postgres
```

### Connect to Postgres inside container
```bash
docker exec -it arbeitscraper-postgres psql -U postgres -d arbeitdb
```

---

## Reset Everything (Fresh Start)
Delete all database data and rebuild from scratch:

```bash
docker compose down -v
```

This removes the volume `dbdata`.

Then rebuild:
```bash
docker compose up --build
```

---

## Summary
- Run `docker compose up --build` only the first time or after major code changes.
- After that, start everything with `docker compose up`.
- All services run automatically: DB → backend → frontend.

