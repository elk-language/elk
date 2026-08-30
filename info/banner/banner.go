// Package banner contains code that renders an info banner for Elk
package banner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/elk-language/elk/info"
	"github.com/elk-language/elk/vm"
)

func Display(customInfo ...string) {
	fmt.Println(Render())
}

func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if home, err := os.UserHomeDir(); err == nil {
		if rel, err := filepath.Rel(home, cwd); err == nil && !strings.HasPrefix(rel, "..") {
			return filepath.Join("~", rel)
		}
	}

	return cwd
}

func Render(customInfo ...string) string {
	var buff strings.Builder
	var infoSlice []string
	for titleLine := range strings.Lines(info.Title) {
		titleLine = strings.TrimRight(titleLine, "\r\n")
		infoSlice = append(infoSlice, titleLine)
	}
	infoSlice = append(infoSlice, "")
	infoSlice = append(infoSlice, fmt.Sprintf("Mode: %s", info.CurrentMode.Description()))
	infoSlice = append(infoSlice, fmt.Sprintf("Version: %s", info.Version))
	infoSlice = append(infoSlice, fmt.Sprintf("OS: %s/%s, %d threads", runtime.GOOS, runtime.GOARCH, runtime.NumCPU()))
	infoSlice = append(infoSlice, fmt.Sprintf("Global Thread Pool: %d", vm.DEFAULT_THREAD_POOL_SIZE))
	infoSlice = append(infoSlice, fmt.Sprintf("Dir: %s", cwd()))
	infoSlice = append(infoSlice, fmt.Sprintf("Time: %s", time.Now().Format(time.RFC1123)))
	infoSlice = slices.Concat(infoSlice, customInfo)

	buff.WriteByte('\n')
	i := 0
	for logoLine := range strings.Lines(info.Logo) {
		infoIndex := i - 1
		logoLine = strings.TrimRight(logoLine, "\r\n")
		buff.WriteString(logoLine)
		buff.WriteString("  ")
		if infoIndex >= 0 && infoIndex < len(infoSlice) {
			buff.WriteString("  ")
			buff.WriteString(infoSlice[infoIndex])
		}
		buff.WriteRune('\n')
		i++
	}

	return buff.String()
}
