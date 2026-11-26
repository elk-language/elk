package lsp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lsp/headers"
	"github.com/elk-language/elk/lsp/types"
)

type lsp struct {
	logger *slog.Logger
}

func new(logger *slog.Logger) *lsp {
	return &lsp{
		logger: logger,
	}
}

func Start() {
	logfile, err := os.OpenFile(".elk-lsp-log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("LSP open log file error: %s", err))
	}

	lsp := new(
		slog.New(slog.NewTextHandler(logfile, &slog.HandlerOptions{})),
	)

	lsp.run()
}

func (l *lsp) run() {

}

const contentType = "application/vscode-jsonrpc; charset=utf-8"

func buildHeaders(contentLength int) *headers.Headers {
	return headers.New(
		ds.MakePair("Content-Type", contentType),
		ds.MakePair("Content-Length", strconv.Itoa(contentLength)),
	)
}

func encodeMessage(msg types.Type) []byte {
	encoded, err := json.Marshal(msg)
	if err != nil {
		panic(fmt.Sprintf("LSP message encoding error: %s", err))
	}

	var buf bytes.Buffer
	encodedHeaders := headers.Encode(buildHeaders(len(encoded)))
	buf.Write(encodedHeaders)

	return buf.Bytes()
}

func decodeMessageBase(msg []byte) (string, []byte, error) {
	headerMap, contentStart, err := headers.Parse(msg)
	if err != nil {
		return "", nil, fmt.Errorf("header decoding error: %w", err)
	}

	strContentLength, ok := headerMap.GetOk("Content-Length")
	if !ok {
		return "", nil, errors.New("missing Content-Length header")
	}

	contentLength, err := strconv.Atoi(strContentLength)
	if err != nil {
		return "", nil, fmt.Errorf("invalid content length: %w", err)
	}

	content := msg[contentStart : contentStart+contentLength]

	var message types.MessageWithMethod
	err = json.Unmarshal(content, &message)
	if err != nil {
		return "", nil, fmt.Errorf("message content decoding error: %w", err)
	}

	return message.Method, content, nil
}
