# HTMX TODO App with Go

A simple and responsive TODO app built using **HTMX** for interactivity, **Tailwind CSS** for styling, and **Go** for the backend API. Inspired by GitHub's clean design principles.

## 🛠️ Features

- ✅ Add tasks dynamically without reloading the page.
- 📱 Responsive and mobile-friendly UI.
- 🧰 Lightweight and easy to customize.

## 📦 Prerequisites

Ensure you have the following installed:

- Go (1.20+ recommended)
- Basic knowledge of Go and HTMX

## 🚀 Getting Started

1. **Clone the repository:**

   ```bash
   git clone https://github.com/000xs/htmx-todo-app.git
   cd htmx-todo-app
   ```

2. **Run the Go server:**

   ```bash
   go run main.go
   ```

3. **Open the app:**

   Navigate to `http://localhost:3000` in your browser.

## 📂 Project Structure

```
.
├── main.go        # Entry point of the Go backend
├── db
│    └── redis.go  # Redis database connection and operations
├── models
│    └── models.go # Data models (e.g., Task structure)
├── routes
│    ├── html
│    │    ├── index.html    # Main HTMX TODO app
│    │    ├── login.html    # User login page
│    │    └── register.html # User registration page
│    ├── auth.go      # Authentication routes (register/login)
│    ├── frontend.go  # HTML rendering handlers
│    ├── routes.go    # Main route handler
│    └── todo.go      # Task-related routes (add, delete)
├── UI
│    └── ui.png       # Screenshot of the app
├── modd.conf        # Modd configuration for live reload
├── color.go         # Utilities for console output
├── go.mod           # Go module configuration
├── go.sum           # Dependencies checksum
└── README.md        # Project documentation
```

 

## ✨ Customization

- **HTMX Interactions:** Modify `hx-post` and `hx-get` attributes for your API.
- **Go Logic:** Enhance task handling with database integration.
- **Styling:** Tailwind CSS is used for styling; customize classes as needed.

## 🤝 Contributing

Feel free to fork the repo and submit pull requests. All contributions are welcome!

## 📄 License

This project is licensed under the MIT License.

