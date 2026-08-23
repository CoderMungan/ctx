//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"strings"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// tomlAssign renders the ` = ` separator between a TOML key and
// its value.
//
// Returns:
//   - string: the assignment separator
func tomlAssign() string {
	return token.Space + token.KeyValueSep + token.Space
}

// tableHeaderPresent reports whether any line of the document,
// after trimming, equals the ctx MCP table header.
//
// Parameters:
//   - data: config.toml content
//
// Returns:
//   - bool: true when the header line exists
func tableHeaderPresent(data []byte) bool {
	for _, line := range splitLines(data) {
		if strings.TrimSpace(line) == cfgCodex.TOMLHeaderMCPCtx {
			return true
		}
	}
	return false
}

// splitLines splits a document into lines without a line-length
// ceiling (bufio.Scanner caps lines at 64KB and stops scanning
// silently, which would make header detection miss everything
// after one very long line).
//
// Parameters:
//   - data: document content
//
// Returns:
//   - []string: the document's lines
func splitLines(data []byte) []string {
	return strings.Split(string(data), token.NewlineLF)
}
