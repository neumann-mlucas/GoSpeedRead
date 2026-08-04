### GoSpeedRead

GoSpeedRead is a toy project that functions as a speed reading GUI application. It reads from the system clipboard, adjusts words per minute, navigates through text, highlights the ORP (Optimal Recognition Point) focal letter, previews neighboring words, and accepts command-line flags. Developed using the [fyne framework](https://fyne.io/), it runs on Linux, Windows, or Mac and is compiled into a small, statically linked binary.

#### Installation

Requires Go 1.21+.

**Option A — `go install` (recommended):**
```sh
go install github.com/neumann-mlucas/GoSpeedRead@latest
```
Binary lands in `$(go env GOBIN)` (fallback `$(go env GOPATH)/bin`, usually `~/go/bin` or `~/.go/bin`). Ensure that directory is on your `PATH`.

**Option B — build from source:**
```sh
git clone https://github.com/neumann-mlucas/GoSpeedRead.git
cd GoSpeedRead
go build .
# then move ./GoSpeedRead somewhere on your PATH
```

**Runtime deps:** fyne uses OpenGL + your display server. On Linux you need `libgl`, `libx11`, `xorg-server` (X11) and/or `wayland`, `libxkbcommon` (Wayland). Most desktop distros ship these.

#### Usage

```
Usage of ./GoSpeedRead:
  -WPM int
        Word per minute (default 300)
  -fontsize int
        Center font size (0 = auto from height)
  -height int
        The height of the window (default 200)
  -width int
        The width of the window (default 800)
```

Theme follows your system preference (light/dark). On GNOME/GTK-based desktops, toggle with:
```sh
gsettings set org.gnome.desktop.interface color-scheme 'prefer-dark'   # dark
gsettings set org.gnome.desktop.interface color-scheme 'default'       # light
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

#### Hyprland Integration

Fyne registers with `app_id = "SpeedRead"` on Wayland. Hyprland-spawned processes may not inherit your login shell's `PATH` — if the bind does nothing, use the absolute binary path (e.g. `/home/you/.go/bin/GoSpeedRead`).

**Hyprland 0.55+ (Lua config at `~/.config/hypr/hyprland.lua`):**
```lua
hl.bind("SUPER + R", hl.dsp.exec_cmd("GoSpeedRead", {
    float = true, center = true, size = {800, 200},
}))
```

**Hyprland pre-0.55 (hyprlang at `~/.config/hypr/hyprland.conf`):**
```
bind = SUPER, R, exec, GoSpeedRead
windowrulev2 = float, class:^(SpeedRead)$
windowrulev2 = center, class:^(SpeedRead)$
windowrulev2 = size 800 200, class:^(SpeedRead)$
```

Check the actual `class` / `title` Hyprland sees:
```sh
hyprctl clients | grep -B1 -A6 -i speed
```

#### Xorg Integration

X11 class is `SpeedRead`. Verify with `xprop WM_CLASS` and click the window.

**sxhkd** (keybind daemon):
```
super + r
    GoSpeedRead
```

**bspwm** window rule (float + center):
```
bspc rule -a SpeedRead state=floating border=off center=on
```

#### Contributing

Contributions welcome. Fork the repository and submit a pull request.

#### License

Distributed under the MIT License. See LICENSE for more information.
