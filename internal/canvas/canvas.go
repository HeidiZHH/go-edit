package canvas

import (
	"fmt"
	"log"
	"syscall"

	"golang.org/x/term"
)

var (
	rowEnd = []byte("\r\n\x1b[K")
	prefix = []byte("~ ")
)

type Canvas struct {
	width    int
	height   int
	termFd   int
	buffer   []byte
	oldState *term.State
}

func New() (*Canvas, error) {
	termFd := int(syscall.Stdout)
	oldState, err := term.MakeRaw(termFd)
	if err != nil {
		return nil, err
	}
	// it uses syscall unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	h, w, err := term.GetSize(termFd)
	if err != nil {
		return nil, err
	}
	canvas := &Canvas{
		termFd:   termFd,
		oldState: oldState,
		width:    w,
		height:   h,
		buffer:   make([]byte, w*h),
	}
	return canvas, nil
}

// Size returns the width and height of the canvas
func (s *Canvas) Size() (int, int) {
	return s.width, s.height
}

func (s *Canvas) Set(x, y int, c byte) {
	s.buffer[y*s.width+x] = c
}

func (s *Canvas) Fill(c byte) {
	for i := range s.buffer {
		s.buffer[i] = c
	}
}

func (s *Canvas) Write(p []byte) (int, error) {
	line := append(prefix, p...)
	for i := 0; i < (s.width - len(line) - len(rowEnd)); i++ {
		line = append(line, ' ')
	}
	line = append(line, rowEnd...)

	s.buffer = line
	log.Printf("writing lines: %q", line)
	return len(p), nil
}

func (s *Canvas) Draw() error {
	if _, err := syscall.Write(s.termFd, s.buffer); err != nil {
		return err
	}
	return nil
}

func (s *Canvas) MoveCursor(x, y int) error {
	_, err := syscall.Write(s.termFd, []byte(fmt.Sprintf("\x1b[%d;%dH", y, x)))
	return err
}

func (s *Canvas) Clear() error {
	_, err := syscall.Write(s.termFd, []byte("\x1b[?25l\x1b[H"))
	return err
}

func (s *Canvas) Close() {
	term.Restore(s.termFd, s.oldState)
}
