### GoSpeedRead

GoSpeedRead is a toy project that functions as a speed reading GUI application. It offers basic features such as reading from the system clipboard, adjusting words per minute, navigating through text (going back or skipping), providing word previews, and accepting command-line flags, among others. Developed using the [fyne framework](https://fyne.io/) framework, it is designed to run seamlessly on Linux, Windows, or Mac. Additionally, the application is compiled into a small, statically linked binary for ease of distribution.

#### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/neumann-mlucas/GoSpeedRead.git
   cd GoSpeedRead

2. **Build the Project:**
   ```sh
   go build .
   ```

3. **Move the executable to your path**

#### Usage

You can run GoSpeedRead in the command line with cmd arguments:

```
Usage of ./GoSpeedRead:
  -WPM int
        Word per minute (default 300)
  -height int
        The height of the window (default 200)
  -width int
        The width of the window (default 800)
```

I primarily use the application with a keybind in **sxhkd**:

```
super + r
    GoSpeedRead
```

For **bspwm**, I also like to set some window rules:


```base
bspc rule -a SpeedRead state="floating" border=off center=on
```

This setup enhances my reading workflow by allowing quick and efficient access to the GoSpeedRead application

#### Contributing

Contributions are welcome! If you have suggestions for improvements or bug fixes, please feel free to fork the repository and submit a pull request.

#### License

Distributed under the MIT License. See LICENSE for more information.

