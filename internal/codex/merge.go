//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bytes"
	"encoding/json"
	"strings"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
)

// parseManifest splits a hooks.json document into its top-level
// keys and the per-event matcher groups under "hooks".
//
// Parameters:
//   - data: JSON document
//
// Returns:
//   - map[string]json.RawMessage: top-level keys (including "hooks")
//   - map[string][]json.RawMessage: event name → matcher groups
//   - error: non-nil when the document or its "hooks" value does
//     not have the expected shape
func parseManifest(data []byte) (
	map[string]json.RawMessage, map[string][]json.RawMessage, error,
) {
	top := map[string]json.RawMessage{}
	if topErr := json.Unmarshal(data, &top); topErr != nil {
		return nil, nil, topErr
	}
	events := map[string][]json.RawMessage{}
	if raw, hasHooks := top[cfgCodex.KeyHooks]; hasHooks {
		if hooksErr := json.Unmarshal(raw, &events); hooksErr != nil {
			return nil, nil, hooksErr
		}
	}
	return top, events, nil
}

// foreignGroups returns the matcher groups that are not
// ctx-managed, in their original order.
//
// Parameters:
//   - groups: matcher groups of one event
//
// Returns:
//   - []json.RawMessage: groups to preserve (nil when none)
func foreignGroups(groups []json.RawMessage) []json.RawMessage {
	var kept []json.RawMessage
	for _, group := range groups {
		if !managed(group) {
			kept = append(kept, group)
		}
	}
	return kept
}

// managed reports whether a matcher group is ctx-managed: it has
// at least one handler and every handler command starts with the
// git-root anchor. Unparseable groups are treated as foreign so
// they survive the merge.
//
// Parameters:
//   - group: one matcher group
//
// Returns:
//   - bool: true when every handler is a ctx hook
func managed(group json.RawMessage) bool {
	fields := map[string]json.RawMessage{}
	if groupErr := json.Unmarshal(group, &fields); groupErr != nil {
		return false
	}
	var handlers []map[string]json.RawMessage
	if handlersErr := json.Unmarshal(
		fields[cfgCodex.KeyHandlers], &handlers,
	); handlersErr != nil || len(handlers) == 0 {
		return false
	}
	for _, h := range handlers {
		var command string
		if cmdErr := json.Unmarshal(
			h[cfgCodex.KeyCommand], &command,
		); cmdErr != nil {
			return false
		}
		current := strings.HasPrefix(command, cfgCodex.HookCommandPrefix)
		legacy := strings.HasPrefix(
			command, cfgCodex.LegacyHookCommandPrefix,
		)
		if !current && !legacy {
			return false
		}
	}
	return true
}

// encode marshals a value without HTML escaping, with the given
// indent, and a trailing newline (the json.Encoder contract).
//
// Parameters:
//   - v: value to encode
//   - indent: indent string (empty for compact output)
//
// Returns:
//   - []byte: encoded JSON
//   - error: non-nil on marshal failure
func encode(v any, indent string) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", indent)
	if encodeErr := encoder.Encode(v); encodeErr != nil {
		return nil, encodeErr
	}
	return buf.Bytes(), nil
}

// equivalent reports whether two JSON documents carry the same
// data regardless of key order and whitespace.
//
// Parameters:
//   - a: first document
//   - b: second document
//
// Returns:
//   - bool: true when both decode to the same value
//   - error: non-nil when either document does not parse
func equivalent(a, b []byte) (bool, error) {
	var av, bv any
	if aErr := json.Unmarshal(a, &av); aErr != nil {
		return false, aErr
	}
	if bErr := json.Unmarshal(b, &bv); bErr != nil {
		return false, bErr
	}
	ac, acErr := json.Marshal(av)
	if acErr != nil {
		return false, acErr
	}
	bc, bcErr := json.Marshal(bv)
	if bcErr != nil {
		return false, bcErr
	}
	return bytes.Equal(ac, bc), nil
}
