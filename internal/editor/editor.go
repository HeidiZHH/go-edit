package editor

import (
	"bufio"
	"fmt"
	"os"

	"gihub.com/heidizhh/go-edit/internal/canvas"
)

type Editor struct {
	in      *bufio.Reader
	screen  *canvas.Canvas
	cursorX int
	curosrY int
}

func New() (*Editor, error) {
	screen, err := canvas.New()
	if err != nil {
		return nil, err
	}
	return &Editor{
		screen:  screen,
		in:      bufio.NewReader(os.Stdin),
		cursorX: 0,
		curosrY: 0,
	}, nil
}

func (e *Editor) readKey() (keyStroke, error) {
	b, err := e.in.ReadByte()
	if err != nil {
		return "", err
	}
	if b == '\x1b' {
		b, err := e.in.ReadByte()
		if err != nil {
			return "", err
		}
		if b == '[' {
			b, err := e.in.ReadByte()
			if err != nil {
				return "", err
			}
			switch b {
			case 'A':
				return keyUp, nil
			case 'B':
				return keyDown, nil
			case 'C':
				return keyRight, nil
			case 'D':
				return keyLeft, nil
			default:
				return keyStroke(fmt.Sprintf("\x1b[%s", b)), nil
			}
		}
		return keyEscape, nil
	}

	return keyStroke(b), nil
}

func (e *Editor) Run() error {
	e.screen.Clear()
	h, w := e.screen.Size()
	e.screen.Write([]byte(fmt.Sprint("Welcome. the terminal size is ", w, "x", h)))
	e.screen.Draw()
	var key string
	var err error
	for {
		key, err = e.readKey()
		if err != nil {
			return err
		}

		switch key {
		case "ESCAPE":
			fmt.Println("escape")
			return nil
		case "r":
			e.screen.MoveCursor(e.cursorX, e.curosrY)
		default:
			message := fmt.Sprintf("the char %q was hit", key)
			_, err := e.screen.Write([]byte(message))
			if err != nil {
				return err
			}
			e.screen.Draw()
		}
	}
}

func (e *Editor) Close() {
	e.screen.Clear()
	e.screen.Close()
}
