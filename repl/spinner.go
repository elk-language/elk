package repl

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

// ANSI sequences that let the spinner redraw its own line.
const (
	eraseLine = "\r\x1b[K"
	cursorUp  = "\x1b[A"
)

const (
	spinnerColor = "#27F57D"
	elapsedColor = "240"
)

// A one line spinner that stays at the bottom of the terminal.
//
// Output written to the spinner is printed above the spinner line,
// so text printed by the executed program pushes the spinner down
// instead of being overwritten by the animation.
//
// The spinner does nothing when its output is not a terminal.
type Spinner struct {
	out     io.Writer
	frames  []string
	fps     time.Duration
	style   lipgloss.Style
	dimmed  lipgloss.Style
	enabled bool

	mutex sync.Mutex
	// Bytes written after the last newline.
	// The terminal has no complete line to keep them on yet,
	// so they are redrawn above the spinner on every frame.
	partialLine []byte
	message     string
	frame       int
	startTime   time.Time
	visible     bool
	stopped     bool

	stop     chan struct{}
	stopOnce sync.Once
	animated sync.WaitGroup
}

// Create a spinner that draws on the given file.
//
// Call this before Capture, because the color profile
// is detected from the given file.
func NewSpinner(out *os.File, message string) *Spinner {
	renderer := lipgloss.NewRenderer(out)

	return &Spinner{
		out:     out,
		frames:  spinner.Dot.Frames,
		fps:     spinner.Dot.FPS,
		style:   renderer.NewStyle().Foreground(lipgloss.Color(spinnerColor)),
		dimmed:  renderer.NewStyle().Foreground(lipgloss.Color(elapsedColor)),
		enabled: isatty.IsTerminal(out.Fd()),
		message: message,
		stop:    make(chan struct{}),
	}
}

// Draw the spinner and start the animation.
func (s *Spinner) Start() {
	if !s.enabled {
		return
	}

	s.mutex.Lock()
	s.startTime = time.Now()
	s.draw()
	s.mutex.Unlock()

	s.animated.Add(1)
	go s.animate()
}

// Stop the animation and erase the spinner line.
// Output that the spinner still holds back is flushed.
func (s *Spinner) Stop() {
	if !s.enabled {
		return
	}

	s.stopOnce.Do(func() { close(s.stop) })
	s.animated.Wait()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.erase()
	s.stopped = true
	s.flushPartialLine()
}

// Replace the text that follows the spinner.
func (s *Spinner) SetMessage(message string) {
	if !s.enabled {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.erase()
	s.message = message
	s.draw()
}

// Print the given bytes above the spinner line.
func (s *Spinner) Write(p []byte) (int, error) {
	if !s.enabled {
		return s.out.Write(p)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.erase()
	s.partialLine = append(s.partialLine, p...)
	if end := bytes.LastIndexByte(s.partialLine, '\n') + 1; end > 0 {
		_, err := s.out.Write(s.partialLine[:end])
		s.partialLine = append(s.partialLine[:0], s.partialLine[end:]...)
		if err != nil {
			return 0, err
		}
	}
	s.draw()

	return len(p), nil
}

func (s *Spinner) Run(fn func(SpinnerCaptureData) error) error {
	capture, err := s.Capture()
	if err != nil {
		return err
	}
	defer func() {
		capture.Restore()
		s.Stop()
	}()

	s.Start()
	return fn(capture)
}

func (s *Spinner) RunWithRestore(fn func(SpinnerCaptureData) error, restore func(SpinnerCaptureData)) error {
	capture, err := s.Capture()
	if err != nil {
		return err
	}
	defer func() {
		restore(capture)
		s.Stop()
	}()

	s.Start()
	return fn(capture)
}

type SpinnerCaptureData struct {
	Restore        func()
	OriginalStdout io.Writer
	OriginalStderr io.Writer
	NewStdout      io.Writer
	NewStderr      io.Writer
}

// Redirect os.Stdout and os.Stderr through the spinner and
// return a function that restores them.
//
// The redirection uses a pipe, so it also covers child processes
// that inherit the standard streams.
// While it is in place, stderr is written to stdout.
//
// Call the returned function before Stop.
func (s *Spinner) Capture() (captureData SpinnerCaptureData, err error) {
	if !s.enabled {
		return captureData, nil
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		return captureData, err
	}

	stdout, stderr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = writer, writer

	drained := make(chan struct{})
	go func() {
		defer close(drained)
		io.Copy(s, reader)
	}()

	return SpinnerCaptureData{
		Restore: func() {
			os.Stdout, os.Stderr = stdout, stderr
			writer.Close()

			// A grandchild process can hold the write end open,
			// which delays EOF past the lifetime of the command.
			select {
			case <-drained:
			case <-time.After(time.Second):
			}
			reader.Close()
		},
		OriginalStdout: stdout,
		OriginalStderr: stderr,
		NewStdout:      writer,
		NewStderr:      writer,
	}, nil
}

// Write the spinner line. The cursor ends up on that line.
func (s *Spinner) draw() {
	if s.visible || s.stopped {
		return
	}

	if len(s.partialLine) > 0 {
		s.out.Write(s.partialLine)
		io.WriteString(s.out, "\n")
	}

	fmt.Fprintf(
		s.out,
		"%s%s %s",
		s.style.Render(s.frames[s.frame]),
		s.message,
		s.dimmed.Render(fmt.Sprintf("(%.1fs)", time.Since(s.startTime).Seconds())),
	)
	s.visible = true
}

// Remove the lines written by draw, leaving the cursor where it started.
// It has to run with the same partial line that draw used.
func (s *Spinner) erase() {
	if !s.visible {
		return
	}

	io.WriteString(s.out, eraseLine)
	if len(s.partialLine) > 0 {
		io.WriteString(s.out, cursorUp+eraseLine)
	}
	s.visible = false
}

func (s *Spinner) flushPartialLine() {
	if len(s.partialLine) == 0 {
		return
	}

	s.out.Write(s.partialLine)
	s.partialLine = s.partialLine[:0]
}

func (s *Spinner) animate() {
	defer s.animated.Done()

	ticker := time.NewTicker(s.fps)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.mutex.Lock()
			s.erase()
			s.frame = (s.frame + 1) % len(s.frames)
			s.draw()
			s.mutex.Unlock()
		}
	}
}
