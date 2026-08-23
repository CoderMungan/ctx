//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/codex"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates all project-local Codex integration files.
//
// Writes .codex/hooks.json (create or merge), .codex/config.toml
// (create or append the [mcp_servers.ctx] table), AGENTS.md
// (marker merge), and .agents/skills/<name>/SKILL.md for every
// embedded Codex skill. When the ctx Codex plugin is enabled in
// the user's config.toml, only AGENTS.md is deployed.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if the hooks manifest cannot be written
//     (other errors are warned but do not halt deployment)
func Deploy(cmd *cobra.Command) error {
	if codex.PluginEnabled(codex.Home()) {
		writeSetup.InfoCodexPluginActive(cmd)
		deployAgents(cmd)
		writeSetup.InfoCodexSummaryPlugin(cmd)
		return nil
	}

	if hooksErr := deployHooks(cmd); hooksErr != nil {
		return hooksErr
	}

	if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
		writeErr.WarnFile(cmd, cfgSetup.MCPConfigPathCodex, mcpErr)
	}

	deployAgents(cmd)

	if skillErr := deploySkills(cmd); skillErr != nil {
		writeErr.WarnFile(cmd, cfgSetup.SkillsPathCodex, skillErr)
	}

	writeSetup.InfoCodexSummary(cmd)
	return nil
}
