package colorz // import "ezpkg.io/colorz"

import (
	"fmt"
	"io"
)

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	Reset Color = 0

	_N = White - Black + 1
)

type Color byte

var cachedCodes = initCodes()

func initCodes() (codes [_N]string) {
	for c := Black; c < White; c++ {
		codes[c-Black] = fmt.Sprintf("\x1b[%dm", byte(c))
	}
	return codes
}

func (c Color) String() string {
	return c.Code()
}

func (c Color) Code() string {
	if c == 0 {
		return "\x1b[0m"
	}
	if c >= Black && c <= White {
		return cachedCodes[c-Black]
	}
	return ""
}

func (c Color) Wrap(s string) string {
	return c.Code() + s + Reset.Code()
}

func (c Color) Format(s fmt.State, verb rune) {
	switch {
	case verb == 'v' && (s.Flag('+') || s.Flag('#')):
		fprintf(s, "Color(%d)", byte(c))
	case verb == 'd':
		fprintf(s, "%d", byte(c))
	default:
		writeString(s, c.Code())
	}
}

func (c Color) Sprint(args ...any) string {
	return c.Wrap(fmt.Sprint(args...))
}

func (c Color) Sprintf(format string, args ...any) string {
	s := fmt.Sprintf(format, args...)
	return c.Sprint(s)
}

func (c Color) Fprintf(w io.Writer, format string, args ...any) (n int, err error) {
	ni, err := fmt.Fprintf(w, "\x1b[%dm", byte(c))
	n += ni
	if err != nil {
		return n, err
	}
	ni, err = fmt.Fprintf(w, format, args...)
	n += ni
	if err != nil {
		return n, err
	}
	ni, err = fmt.Fprintf(w, "\x1b[0m")
	return n + ni, err
}

func (c Color) Print(args ...any) {
	fmt.Printf("\x1b[%dm", byte(c))
	fmt.Print(args...)
	fmt.Printf("\x1b[0m")
}

func (c Color) Printf(format string, args ...any) {
	fmt.Printf("\x1b[%dm", byte(c))
	fmt.Printf(format, args...)
	fmt.Printf("\x1b[0m")
}

func (c Color) Println(args ...any) {
	fmt.Printf("\x1b[%dm", byte(c))
	fmt.Println(args...)
	fmt.Printf("\x1b[0m")
}

func writeString(w fmt.State, s string) {
	_, _ = io.WriteString(w, s)
}

func fprintf(w fmt.State, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
