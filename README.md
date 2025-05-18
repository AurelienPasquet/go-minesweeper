# 🧨 Minesweeper in Go (Ebiten)

A classic Minesweeper clone written in Go, using the [Ebiten](https://ebitengine.org/) 2D game engine.  
⚡ Fast, clean, and focused on gameplay accuracy and responsiveness.

---

## 📦 Features

- 🎮 Classic Minesweeper mechanics: left-click to reveal, right-click to flag.
- 🚩 Flag handling, automatic reveal of empty areas, and flood-fill propagation.
- 💥 Win/loss detection with visual feedback.
- 🔁 Press `SPACE` to instantly restart the game.
- 🖥️ Grid tile size adapts automatically to your screen resolution.

---

## 🧪 Running the Game

### 🔧 Requirements

- [Go](https://golang.org/dl/) 1.21 or later
- Ebiten (installed automatically with `go get`)

### ▶️ Run Locally

```bash
go run . -w 30 -h 16 -m 99
```

**Command-line flags:**
- `-w`: grid width (e.g., 30)
- `-h`: grid height (e.g., 16)
- `-m`: number of mines (e.g., 99)

---

## ⌨️ Controls

| Action                | Input             |
|-----------------------|------------------|
| Reveal tile           | Left click        |
| Flag/unflag tile      | Right click       |
| Restart game          | Spacebar          |

---

### Have Fun!
