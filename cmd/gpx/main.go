package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ZeraiGR/gpx/internal/app"
	"github.com/ZeraiGR/gpx/internal/config"
	"github.com/ZeraiGR/gpx/internal/envx"
	"github.com/ZeraiGR/gpx/internal/shell"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "init":
		initCmd(os.Args[2:])
	case "list":
		listCmd(os.Args[2:])
	case "status":
		statusCmd(os.Args[2:])
	case "use":
		useCmd(os.Args[2:])
	case "set":
		setCmd(os.Args[2:])
	case "unset":
		unsetCmd(os.Args[2:])
	case "diff":
		diffCmd(os.Args[2:])
	case "apply":
		applyCmd(os.Args[2:])
	case "profile":
		profileCmd(os.Args[2:])
	case "version":
		versionCmd()
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("gpx - manage environment presets for Go workflows")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gpx init   [--force] [--config PATH]")
	fmt.Println("  gpx list   [--config PATH]")
	fmt.Println("  gpx status [--config PATH]")
	fmt.Println("  gpx use <profile> [--config PATH]")
	fmt.Println("  gpx set KEY=VALUE [KEY=VALUE ...] [--config PATH]")
	fmt.Println("  gpx unset KEY [KEY ...]")
	fmt.Println("  gpx diff <profile> [--config PATH]")
	fmt.Println("  gpx apply [--rc PATH] [--shell zsh|bash] [--dry-run] [--backup] <profile> [--config PATH]")
	fmt.Println()
	fmt.Println("Config editing:")
	fmt.Println("  gpx profile add <name> [--config PATH]")
	fmt.Println("  gpx profile rm <name> [--config PATH]")
	fmt.Println("  gpx profile rename <old> <new> [--config PATH]")
	fmt.Println("  gpx profile show <name> [--config PATH]")
	fmt.Println("  gpx profile set <name> KEY=VALUE [KEY=VALUE ...] [--config PATH]")
	fmt.Println("  gpx profile unset <name> KEY [KEY ...] [--config PATH]")
	fmt.Println()
	fmt.Println("  gpx version")
	fmt.Println()
	fmt.Println("Tips:")
	fmt.Println(`  eval "$(gpx use public)"`)
}

func resolveConfigPath(fs *flag.FlagSet) *string {
	return fs.String("config", "", "path to config file (default: ~/.config/gpx/config.json)")
}

func makeApp(cfgPath string) app.App {
	return app.App{ConfigPath: cfgPath}
}

func defaultConfigPathOrExit() string {
	p, err := config.DefaultPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	return p
}

// enforce flags-first contract: <cmd> [flags] <args>
func ensureFlagsBeforeArgs(args []string, cmdName string) {
	pos := -1
	for i, a := range args {
		if !strings.HasPrefix(a, "-") {
			pos = i
			break
		}
	}
	if pos == -1 {
		return
	}
	for j := pos + 1; j < len(args); j++ {
		if strings.HasPrefix(args[j], "-") {
			fmt.Fprintf(os.Stderr, "error: flags must come before <args> for %q. Try: gpx %s [flags] <args>\n", cmdName, cmdName)
			os.Exit(2)
		}
	}
}

func initCmd(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	force := fs.Bool("force", false, "overwrite existing config")
	cfgPath := resolveConfigPath(fs)
	_ = fs.Parse(args)

	path := *cfgPath
	if path == "" {
		path = defaultConfigPathOrExit()
	}

	res, err := app.InitConfig(path, *force)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("Config: %s (%s)\n", res.Path, res.Status)
}

func listCmd(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	cfgPath := resolveConfigPath(fs)
	_ = fs.Parse(args)

	path := *cfgPath
	if path == "" {
		path = defaultConfigPathOrExit()
	}

	a := makeApp(path)
	items, err := a.ListProfiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print(app.FormatProfiles(items))
}

func statusCmd(args []string) {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	cfgPath := resolveConfigPath(fs)
	_ = fs.Parse(args)

	path := *cfgPath
	if path == "" {
		path = defaultConfigPathOrExit()
	}

	a := makeApp(path)
	rows, err := a.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print(app.FormatStatus(rows))
}

func useCmd(args []string) {
	ensureFlagsBeforeArgs(args, "use")
	fs := flag.NewFlagSet("use", flag.ExitOnError)
	cfgPath := resolveConfigPath(fs)
	_ = fs.Parse(args)

	rest := fs.Args()
	if len(rest) < 1 {
		fmt.Fprintln(os.Stderr, "error: missing profile name")
		os.Exit(2)
	}
	name := rest[0]

	path := *cfgPath
	if path == "" {
		path = defaultConfigPathOrExit()
	}

	a := makeApp(path)
	lines, err := a.UseProfile(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	for _, ln := range lines {
		fmt.Println(ln)
	}
}

func setCmd(args []string) {
	ensureFlagsBeforeArgs(args, "set")
	fs := flag.NewFlagSet("set", flag.ExitOnError)
	cfgPath := resolveConfigPath(fs)
	_ = fs.Parse(args)

	tokens := fs.Args()
	if len(tokens) == 0 {
		fmt.Fprintln(os.Stderr, "error: expected at least one KEY=VALUE assignment")
		os.Exit(2)
	}

	path := *cfgPath
	if path == "" {
		path = defaultConfigPathOrExit()
	}

	a := makeApp(path)
	lines, err := a.SetVars(tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	for _, ln := range lines {
		fmt.Println(ln)
	}
}

func unsetCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: expected at least one KEY")
		os.Exit(2)
	}
	lines, err := envx.UnsetLines(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	for _, ln := range lines {
		fmt.Println(ln)
	}
}

func diffCmd(args []string) {
	ensureFlagsBeforeArgs(args, "diff")
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	cfgPath := resolveConfigPath(fs)
	_ = fs.Parse(args)

	rest := fs.Args()
	if len(rest) < 1 {
		fmt.Fprintln(os.Stderr, "error: missing profile name")
		os.Exit(2)
	}
	profile := rest[0]

	path := *cfgPath
	if path == "" {
		path = defaultConfigPathOrExit()
	}

	a := makeApp(path)
	rows, err := a.DiffProfile(profile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print(app.FormatDiff(rows))
}

func applyCmd(args []string) {
	ensureFlagsBeforeArgs(args, "apply")
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	cfgPath := resolveConfigPath(fs)
	shName := fs.String("shell", "zsh", "shell type: zsh or bash (used only when --rc is not provided)")
	rc := fs.String("rc", "", "rc file path (overrides --shell default)")
	dryRun := fs.Bool("dry-run", false, "show what would be written, but do not modify any file")
	backup := fs.Bool("backup", false, "create a backup of rc file before modifying it")
	_ = fs.Parse(args)

	rest := fs.Args()
	if len(rest) < 1 {
		fmt.Fprintln(os.Stderr, "error: missing profile name")
		os.Exit(2)
	}
	profile := rest[0]

	path := *cfgPath
	if path == "" {
		path = defaultConfigPathOrExit()
	}

	rcPath := *rc
	if rcPath == "" {
		p, err := shell.DefaultRC(*shName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		rcPath = p
	}

	a := makeApp(path)
	report, err := a.ApplyProfileToRC(profile, rcPath, shell.ApplyOptions{
		DryRun: *dryRun,
		Backup: *backup,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if !*backup {
		fmt.Println("Note: no backup was created (use --backup to enable)")
	}

	if *dryRun {
		fmt.Printf("Dry-run: would write GPX block to %s\n", report.RCPath)
		fmt.Println()
		fmt.Print(report.NewContent)
		return
	}

	fmt.Printf("Applied profile %q to %s\n", profile, report.RCPath)
	if report.BackupPath != "" {
		fmt.Printf("Backup: %s\n", report.BackupPath)
	}
	fmt.Printf("Next: source %s (or restart shell)\n", report.RCPath)
}

func profileCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: missing subcommand (add|rm|rename|show|set|unset)")
		os.Exit(2)
	}

	sub := args[0]
	rest := args[1:]

	// flags-first for all profile subcommands
	ensureFlagsBeforeArgs(rest, "profile "+sub)

	fs := flag.NewFlagSet("profile "+sub, flag.ExitOnError)
	cfgPath := resolveConfigPath(fs)
	_ = fs.Parse(rest)

	path := *cfgPath
	if path == "" {
		path = defaultConfigPathOrExit()
	}
	a := makeApp(path)

	argv := fs.Args()

	switch sub {
	case "add":
		if len(argv) < 1 {
			fmt.Fprintln(os.Stderr, "error: profile add <name>")
			os.Exit(2)
		}
		if err := a.AddProfile(argv[0]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("OK")
	case "rm":
		if len(argv) < 1 {
			fmt.Fprintln(os.Stderr, "error: profile rm <name>")
			os.Exit(2)
		}
		if err := a.RemoveProfile(argv[0]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("OK")
	case "rename":
		if len(argv) < 2 {
			fmt.Fprintln(os.Stderr, "error: profile rename <old> <new>")
			os.Exit(2)
		}
		if err := a.RenameProfile(argv[0], argv[1]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("OK")
	case "show":
		if len(argv) < 1 {
			fmt.Fprintln(os.Stderr, "error: profile show <name>")
			os.Exit(2)
		}
		vars, err := a.ShowProfile(argv[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Print(app.FormatProfileVars(argv[0], vars))
	case "set":
		if len(argv) < 2 {
			fmt.Fprintln(os.Stderr, "error: profile set <name> KEY=VALUE [KEY=VALUE ...]")
			os.Exit(2)
		}
		name := argv[0]
		tokens := argv[1:]
		if err := a.SetProfileVars(name, tokens); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("OK")
	case "unset":
		if len(argv) < 2 {
			fmt.Fprintln(os.Stderr, "error: profile unset <name> KEY [KEY ...]")
			os.Exit(2)
		}
		name := argv[0]
		keys := argv[1:]
		if err := a.UnsetProfileVars(name, keys); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("OK")
	default:
		fmt.Fprintln(os.Stderr, "error: unknown profile subcommand:", sub)
		os.Exit(2)
	}
}

func versionCmd() {
	fmt.Printf("gpx %s (commit=%s, date=%s)\n", version, commit, date)
}
