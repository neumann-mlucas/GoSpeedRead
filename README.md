### GoSpeedRead

GoSpeedRead is a toy project that functions as a speed reading GUI application. It reads from the system clipboard, adjusts words per minute, navigates through text, highlights the ORP (Optimal Recognition Point) focal letter, previews neighboring words, and accepts command-line flags. Developed using the [fyne framework](https://fyne.io/), it runs on Linux, Windows, or Mac and is compiled into a small, statically linked binary.

#### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/neumann-mlucas/GoSpeedRead.git
   cd GoSpeedRead
   ```

2. **Build the Project:**
   ```sh
   go build .
   ```

3. **Move the executable to your PATH** (or run `go install .`).

#### Usage

```
Usage of ./GoSpeedRead:
  -WPM int
        Word per minute (default 300)
  -fontsize int
        Center font size (0 = auto from height)
  -height int
        The height of the window (default 200)
  -theme string
        Theme: dark or light (default "dark")
  -width int
        The width of the window (default 800)
```

#### Keybindings

| Key            | Action                          |
| -------------- | ------------------------------- |
| `space`        | play / pause                    |
| `r`            | restart at first word           |
| `v`            | reload clipboard                |
| `h` / `←`      | jump back 5 words               |
| `l` / `→`      | jump forward 5 words            |
| `j` / `↓`      | decrease WPM by 10              |
| `k` / `↑`      | increase WPM by 10              |
| `?`            | toggle keybind help overlay     |

#### Window Manager Integration

**sxhkd** (X11):
```
super + r
    GoSpeedRead
```

**bspwm** (X11) window rules:
```
bspc rule -a SpeedRead state=floating border=off center=on
```

**Hyprland (0.55+, Lua config):**
```lua
hl.bind("SUPER + R", hl.dsp.exec_cmd("GoSpeedRead", {
    float = true, center = true, size = {800, 200},
}))
```

**Hyprland (pre-0.55, hyprlang):**
```
bind = SUPER, R, exec, GoSpeedRead
windowrulev2 = float, class:^(GoSpeedRead|SpeedRead)$
windowrulev2 = center, class:^(GoSpeedRead|SpeedRead)$
windowrulev2 = size 800 200, class:^(GoSpeedRead|SpeedRead)$
```

#### Contributing

Contributions welcome. Fork the repository and submit a pull request.

#### License

Distributed under the MIT License. See LICENSE for more information.
