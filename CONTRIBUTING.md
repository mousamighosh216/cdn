# Contributing to CDN Platform

Thanks for your interest in contributing! 🚀  
This project is primarily a learning-focused system design project, but improvements, ideas, and discussions are welcome.

---

## 🧰 Prerequisites

- Go 1.20+
- Basic understanding of HTTP & distributed systems
- Familiarity with Go modules

---

## 🛠️ Development Setup

```bash
git clone [https://github.com/mousamighosh216/cdn.git](https://github.com/mousamighosh216/cdn.git)
cd cdn
go mod tidy

# Run individual components:

go run edge/cmd/edge/main.go
go run origin/main.go
go run control-plane/cmd/server/main.go
go run dns-resolver/cmd/resolver/main.go

# Run using Docker-Compose:

docker-compose up --build

--- 

## 📐 Code Guidelines

- Follow standard Go formatting (gofmt)
- Keep components modular
- Avoid introducing global state
- Prefer clarity over premature optimization

--- 

## 💡 What You Can Contribute

- Bug fixes
- Performance improvements
- Architecture discussions
- Documentation improvements
- Tests

---