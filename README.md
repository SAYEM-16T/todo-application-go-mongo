
# 📝 Todo App — Go + MongoDB (Dockerized)

A **superfast, minimal** Todo web app built with **Go**, **MongoDB**, and a clean **HTML + CSS** frontend (no JavaScript).
The whole stack is containerized with **Docker Compose**, so anyone can run it with a single command.

---

## ✨ Features

- Auth: **Register / Login / Logout** (JWT session cookie)
- **User-scoped** todos (each user sees only their own)
- Add / **Done/Undo** / Delete tasks
- Pure **HTML + CSS** (server-rendered, no JS)
- **MongoDB** persistence
- **Dockerized** stack (app + database)

---
## 📸 Screenshots (full flow)

### 1) Landing
![Landing](images/Screenshot%20from%202025-10-21%2016-00-09.png)

### 2) Register
![Register](images/Screenshot%20from%202025-10-21%2016-00-45.png)

### 3) App — first visit (empty state)
![App empty](images/Screenshot%20from%202025-10-21%2016-01-38.png)

### 4) App — first task added
![First task added](images/Screenshot%20from%202025-10-21%2016-01-48.png)

### 5) App — typing second task
![Typing second task](images/Screenshot%20from%202025-10-21%2016-03-36.png)

### 6) App — two tasks shown
![Two tasks](images/Screenshot%20from%202025-10-21%2016-03-45.png)

### 7) Login
![Login](images/Screenshot%20from%202025-10-21%2016-04-22.png)




---

## 🧱 Tech Stack

- **Go** (chi router, MongoDB driver, bcrypt, JWT)
- **MongoDB**
- **HTML + CSS** (no JS)
- **Docker & Docker Compose**

---

## 🗂️ Project Structure

```
.
├── backend/
│   ├── go.mod, go.sum
│   ├── main.go
│   ├── handlers/      # auth + todo handlers
│   ├── middleware/    # auth middleware
│   ├── models/        # User, Todo models
│   └── utils/         # DB, password, session (JWT), server-side HTML render
├── frontend/
│   ├── index.html
│   ├── login.html
│   ├── register.html
│   └── styles/main.css
├── images/            # screenshots used in README
├── Dockerfile
└── docker-compose.yml

````

---

## 🚀 Quick Start (Docker)

```bash
git clone https://github.com/<your-username>/todo-app.git
cd todo-app
docker compose up -d --build
````

Open: **[http://localhost:8080](http://localhost:8080)**

Check logs:

```bash
docker compose ps
docker compose logs -f app
docker compose logs -f mongo
```

Shut down:

```bash
docker compose down        # stop containers
docker compose down -v     # stop + remove DB volume (data will be deleted)
```

> Note: Docker Compose v2 ignores the `version:` field; you can remove it to silence the warning.

---

## ⚙️ Environment

Provided via `docker-compose.yml`:

```env
MONGODB_URI=mongodb://mongo:27017
DB_NAME=todoapp
JWT_SECRET=super-secure-change-me
ADDR=:8080
```

**Never commit real secrets.** If you run without Docker, create `backend/.env` and keep it out of Git.

---

## 🔌 Endpoints (Server-Rendered HTML)

| Method | Path                | Description                |
| -----: | ------------------- | -------------------------- |
|    GET | `/`                 | Landing                    |
|    GET | `/login`            | Login page                 |
|   POST | `/login`            | Login + set session cookie |
|    GET | `/register`         | Register page              |
|   POST | `/register`         | Create user + auto-login   |
|   POST | `/logout`           | Clear session              |
|    GET | `/app`              | Todo dashboard (auth)      |
|   POST | `/todo`             | Add todo (auth)            |
|   POST | `/todo/{id}/toggle` | Done / Undo (auth)         |
|   POST | `/todo/{id}/delete` | Delete (auth)              |

**Request flow:** HTML form → POST → server updates DB → **302 redirect** → HTML rendered again.

---

<!-- ## 🧪 Smoke Test (cURL)

```bash
# 1) Register (auto-login) and save cookies
curl -s -L -c cookie.txt \
  -d "email=test$(date +%s)@ex.com&password=secret123" \
  http://localhost:8080/register >/dev/null

# 2) Add a todo
curl -s -L -b cookie.txt \
  -d "title=Building a CI/CD pipeline with Jenkins" \
  http://localhost:8080/todo >/dev/null

# 3) See dashboard HTML
curl -s -b cookie.txt http://localhost:8080/app | head -n 20
```

---

## 🛠️ Local Development (without Docker)

```bash
# Quick Mongo via Docker
docker run -d --name mongo -p 27017:27017 mongo:7

# Backend
cd backend
cp .env.example .env     # set a strong JWT_SECRET
go mod tidy
go run .
# http://localhost:8080
``` -->

---

## 🧯 Troubleshooting

* **Port 27017 already in use**
  Remove `ports:` from the `mongo` service (so it’s internal only), or map another host port: `27018:27017`.

* **Port 8080 already in use**
  Map another host port for the app: `ports: ["9090:8080"]` and open `http://localhost:9090`, or set `ADDR=:9090`.

* **Mongo connect error**
  Ensure `MONGODB_URI=mongodb://mongo:27017` (service name `mongo` is resolvable inside the Compose network). Check `docker compose logs -f mongo`.

* **JWT invalid after config change**
  Changing `JWT_SECRET` invalidates existing sessions. Clear cookies or login again.

<!-- ---

## 🔒 Production Notes

* Use HTTPS and set cookie `Secure` flag (already `HttpOnly` + `SameSite=Lax`).
* Consider CSRF tokens for state-changing POSTs.
* Add rate limiting and request logging as needed.
* Rotating `JWT_SECRET` forces re-login for all users.
 -->
