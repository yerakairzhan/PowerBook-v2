# Bot Project Setup Guide

## 📌 Prerequisites
Before starting, ensure you have the following installed:
- [Go](https://go.dev/dl/) (latest version)
- [PostgreSQL](https://www.postgresql.org/download/)
- [sqlc](https://sqlc.dev)
- [Docker](https://www.docker.com/) (optional for database setup)

---

## 🚀 Setup Instructions

### 1️⃣ Create API Token in BotFather
1. Open Telegram and search for `@BotFather`.
2. Send `/newbot` and follow the instructions.
3. Save the generated **API token**, as you'll need it in the project.

---

### 2️⃣ Prepare Database
You can set up the database **manually** or **with Docker**.

#### **🔹 Option 1: Manual Setup**
1. Start PostgreSQL and create a database:
   ```sh
   psql -U postgres -c "CREATE DATABASE bot_db;"
   ```
2. Adjust ports if needed in your `.env` file.

#### **🔹 Option 2: Docker Setup**
To start PostgreSQL in a Docker container:
```sh
docker run --name bot-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=yourpassword -e POSTGRES_DB=bot_db -p 5432:5432 -d postgres
```

---

### 3️⃣ Create Tables
Run the SQL schema from `/db/tables/up.sql`:
```sh
psql -U postgres -d bot_db -f db/tables/up.sql
```

---

### 4️⃣ Configure Environment Variables
Create a `.env` file in your project root and add the following:
```env
BOT_TOKEN=your_telegram_bot_token
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_HOST=localhost
DB_PORT=5432
DB_NAME=bot_db
```
Load environment variables:
```sh
export $(grep -v '^#' .env | xargs)
```

---

### 5️⃣ Generate SQL Queries
Run `sqlc` to generate Go code for database interactions:
```sh
sqlc generate
```

---

### 6️⃣ Build & Run the Project
Finally, build and run the Golang project:
```sh
go build -o bot_main main.go
./bot_main
```

---

## ✅ You're All Set!
Your bot should now be running. 🚀 If you encounter any issues, check your database connection and API token settings.

---

### **🔹 Key Additions:**
✔ **Environment Variables (`.env`)** for easy configuration  
✔ **Single-file format** with all instructions  
✔ **Copy-paste friendly commands** for quick setup

Let me know if you need further improvements! 🚀

