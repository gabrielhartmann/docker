// +build windows

package winterm

import (
	"os"

	"github.com/Sirupsen/logrus"
	. "github.com/docker/docker/pkg/term/ansiterm"
)

// ansiWriter wraps a standard output file (e.g., os.Stdout) providing ANSI sequence translation.
type ansiWriter struct {
	file           *os.File
	fd             uintptr
	infoReset      *CONSOLE_SCREEN_BUFFER_INFO
	command        []byte
	escapeSequence []byte
	inAnsiSequence bool
	parser         *AnsiParser
}

func newAnsiWriter(nFile int) *ansiWriter {
	file, fd := getStdFile(nFile)
	info, err := GetConsoleScreenBufferInfo(fd)
	if err != nil {
		return nil
	}

	parser := CreateParser(Ground, CreateWinEventHandler(fd, file))

	return &ansiWriter{
		file:           file,
		fd:             fd,
		infoReset:      info,
		command:        make([]byte, 0, ANSI_MAX_CMD_LENGTH),
		escapeSequence: []byte(KEY_ESC_CSI),
		parser:         parser,
	}
}

func (aw *ansiWriter) Fd() uintptr {
	return aw.fd
}

// Write writes len(p) bytes from p to the underlying data stream.
func (aw *ansiWriter) Write(p []byte) (total int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	logrus.Infof("Write: % x", p)
	logrus.Infof("Write: %s", string(p))
	return aw.parser.Parse(p)
}
