package cli

import (
	"fmt"
	"io/fs"

	"github.com/spf13/cobra"
	"github.com/strongo/buildinfo"
	"github.com/strongo/cli-helpers/skillsync"
	skillscmd "github.com/strongo/cli-helpers/skillsync/cobracmd"
	"github.com/strongo/cli-helpers/skillsync/githubrelease"

	"github.com/specscore/codegrapher/agentplugin"
)

const (
	codeGrapherSkillsPluginVersion = "0.0.0"
	codeGrapherSkillsUnknownSource = "0000000000000000000000000000000000000000"
)

var (
	codeGrapherSkillsCLI    = skillsync.Identity{Publisher: "code-grapher", Name: "codegrapher"}
	codeGrapherSkillsPlugin = skillsync.PluginIdentity{Publisher: "code-grapher", Name: "codegrapher"}
)

// newSkillsCmd provides the installed CLI with its exact embedded skill
// bundle. Normal sync is offline and immutable; only --newer-compatible opts
// into the release resolver provided by strongo/cli-helpers.
func newSkillsCmd() *cobra.Command {
	cfg, cfgErr := newSkillsSyncConfig()
	cmd := skillscmd.New(cfg, skillscmd.CommandOptions{
		Short: "Install CodeGrapher Agent Skills into a harness skills directory",
		Resolver: skillsync.ReleaseResolver{
			Source:         githubrelease.Source{},
			CurrentVersion: cfg.CurrentVersion,
		},
	})
	addJSONShortcut(cmd.Commands()[0])
	cmd.Long = `Install CodeGrapher's immutable, CLI-matched Agent Skills.

By default, sync reads the skill bundle embedded in this installed binary and
does not need a source checkout or network access. Use --newer-compatible only
to explicitly select a newer compatible published bundle.

With no target, every present Claude, Cursor, and Codex harness is synced. Use
--harness to select a harness or --dir to select an explicit skills directory.`
	if cfgErr != nil {
		cmd.RunE = func(*cobra.Command, []string) error {
			return fmt.Errorf("prepare embedded CodeGrapher skills: %w", cfgErr)
		}
	}
	return cmd
}

func newSkillsSyncConfig() (skillsync.Config, error) {
	source, err := fs.Sub(agentplugin.SkillsFS, "skills")
	if err != nil {
		return skillsync.Config{}, err
	}
	digest, err := skillsync.Digest(source)
	if err != nil {
		return skillsync.Config{}, err
	}
	build := buildinfo.Get("codegrapher")
	revision := build.Commit
	if len(revision) != 40 {
		revision = codeGrapherSkillsUnknownSource
	}
	version := build.Version
	if _, err := skillsync.CompareVersions(version, version); err != nil {
		version = codeGrapherSkillsPluginVersion
	}
	descriptor := skillsync.BundleDescriptor{
		Plugin: codeGrapherSkillsPlugin,
		Source: skillsync.Source{
			Repository: "github.com/code-grapher/codegrapher",
			Path:       "agentplugin/skills",
			Revision:   revision,
			Version:    version,
			Digest:     digest,
		},
	}
	if err := skillsync.ValidateDescriptor(descriptor); err != nil {
		return skillsync.Config{}, fmt.Errorf("validate embedded CodeGrapher skills descriptor: %w", err)
	}
	bundle, err := skillsync.EmbeddedBundle(descriptor, source)
	if err != nil {
		return skillsync.Config{}, fmt.Errorf("bind embedded CodeGrapher skills: %w", err)
	}
	return skillsync.Config{
		CLI:            codeGrapherSkillsCLI,
		CurrentVersion: build.Version,
		Bundles:        []skillsync.Bundle{bundle},
	}, nil
}

// addJSONShortcut keeps CodeGrapher's established --format=json / --json
// convention while using shared Cobra adapters that expose only --format.
func addJSONShortcut(cmd *cobra.Command) {
	var jsonOut bool
	cmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "Output as JSON (shortcut for --format=json)")
	original := cmd.PreRunE
	cmd.PreRunE = func(command *cobra.Command, args []string) error {
		if jsonOut {
			if err := command.Flags().Set("format", "json"); err != nil {
				return err
			}
		}
		if original != nil {
			return original(command, args)
		}
		return nil
	}
}
