// Package agentplugin embeds CodeGrapher's publishable Agent Skills bundle.
//
// The command package consumes this immutable snapshot so `codegrapher skills
// sync` works from an installed binary, without requiring a source checkout.
package agentplugin

import "embed"

// SkillsFS contains the complete publishable skills tree.
//
//go:embed all:skills
var SkillsFS embed.FS
