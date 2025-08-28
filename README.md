# Ben's Pomodoro (`pomodoro-do-ben`)

![Ben's Pomodoro](./pomodoro-do-ben.png)

A simple yet powerful Pomodoro timer for your Linux desktop, designed to keep you focused and mindful. Written in Go with a clean GTK-based graphical interface.

---

## ✨ Features

*   **Classic Pomodoro Timer:** Boost your productivity with the classic focus/break cycle.
*   **Clean Desktop GUI:** A simple and intuitive graphical interface built with Go and GTK.
*   **Desktop Notifications:** Stay informed about your Pomodoro status without leaving your workflow.
*   **Audio Cues:** Sound notifications to signal the start of focus or break periods.
*   **Internationalization:** UI translated into multiple languages.
*   **Easy Installation:** A simple script to install the application and desktop entries on your system.

## 🚀 Installation

You can easily install Ben's Pomodoro using the provided shell script.

### Prerequisites

Make sure you have `go` and `imagemagick` installed. On Debian/Ubuntu-based systems, you can install them with:

```bash
sudo apt-get update
sudo apt-get install golang-go imagemagick
```

### 1. Clone the Repository

```bash
git clone https://github.com/evandrojr/pomodoro-do-ben.git
cd pomodoro-do-ben
```
*(Please replace `evandrojr` with the actual repository path)*

### 2. Build the Application

Compile the Go program to create the executable:

```bash
go build
```

### 3. Run the Installation Script

The script will copy the application binary, desktop file, and icons to your local user directories (`~/.local`).

```bash
chmod +x install.sh
./install.sh
```

## 🏃‍♀️ Usage

After installation, you can launch **Ben's Pomodoro** directly from your desktop's application menu.

Alternatively, you can run the application from your terminal:

```bash
pomodoro-do-ben
```

## 🛠️ Building from Source

If you prefer to build and run the application manually without installing it system-wide:

### 1. Prerequisites

Ensure you have the necessary development libraries installed. On Debian/Ubuntu-based systems:

```bash
sudo apt-get update
sudo apt-get install golang-go libgtk-3-dev libasound2-dev
```

### 2. Build

```bash
go build
```

### 3. Run

```bash
./pomodoro-do-ben
```

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/evandrojr/pomodoro-do-ben/issues).

## 📝 License

This project is licensed under the **[Your License Here]**. Please add a `LICENSE` file to the project.