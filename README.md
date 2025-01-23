# Bango

Welcome to the **Bango** repository! Bango is a server application written in Go (Golang). It is designed to consume HTTP requests, process them based on the configured URL endpoints, and return the appropriate results. While Golang already provides a robust standard library for building HTTP servers, this project is a learning initiative to recreate and understand the underlying mechanics of HTTP server implementation in Go.

---

## Table of Contents

1. [About the Project](#about-the-project)
2. [Features](#features)
3. [Installation](#installation)

---

## About the Project

Bango is a server application built in Go. It accepts HTTP requests, processes them according to the configured URL endpoints, and returns the appropriate responses. The primary goal of this project is to deepen the understanding of Go's HTTP server implementation by recreating core functionalities from scratch. This is purely a learning project and not intended to replace or compete with Go's standard library.

---

## Features

- **HTTP Request Handling**: Processes incoming HTTP requests based on configured endpoints.
- **Custom Routing**: Implements custom URL routing logic.
- **Learning-Oriented**: Designed to help developers understand how HTTP servers work in Go.
- **Lightweight**: Minimalistic and easy to set up.

---

---

## Installation

To get started with Bango, follow these steps:

1. **Clone the repository**:
   ```bash
   git clone https://github.com/saravanastar/bango.git
   cd bango

2. **Build the Project**
   ```bash
   go build -o bango-server ./cmd
   ```
3. **Run the Server**
   ```bash
   ./bango-server
   ```

## Example Requests
```bash
curl -X GET http://localhost:4221/
```
**Response**
```html
All OK!
```
