package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/s3cy/autoebiten/internal/cli"
	"github.com/s3cy/autoebiten/internal/rpc"
)

var (
	// Flags
	pidFlag         int
	keyFlag         string
	inputActionFlag string
	mouseActionFlag string
	durationTicks   int64
	xFlag           int
	yFlag           int
	wheelXFlag      float64
	wheelYFlag      float64
	buttonFlag      string
	outputFlag      string
	scriptFlag      string
	inlineFlag      string
	asyncFlag       bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "autoebiten",
		Short: "CLI tool for automating Ebitengine games",
		Long: `autoebiten is a CLI tool that enables AI agents to automate Ebitengine games.
It communicates with games via JSON-RPC over a Unix socket.

If --pid is not provided, autoebiten automatically detects a running game instance.
If multiple games are running, use --pid to specify the target.`,
		PersistentPreRunE: persistentPreRunRootCommand,
	}
	rootCmd.PersistentFlags().IntVarP(&pidFlag, "pid", "p", 0, "Target game process PID (auto-detected if not specified)")

	// input command
	inputCmd := &cobra.Command{
		Use:   "input",
		Short: "Send keyboard input to the game",
		Long: `Inject keyboard input into the game.

Actions:
  press   - Press and immediately release the key
  release - Release a held key
  hold    - Press and hold for duration_ticks`,
		RunE: runInputCommand,
	}
	inputCmd.Flags().StringVarP(&keyFlag, "key", "k", "", "Key name (e.g., KeyA, KeySpace, KeyArrowUp)")
	inputCmd.Flags().StringVarP(&inputActionFlag, "action", "a", "hold", "Action: press, release, or hold")
	inputCmd.Flags().Int64VarP(&durationTicks, "duration_ticks", "d", 6, "Duration in ticks for hold action")
	inputCmd.MarkFlagRequired("key")

	// mouse command
	mouseCmd := &cobra.Command{
		Use:   "mouse",
		Short: "Send mouse input to the game",
		Long: `Inject mouse input into the game.

Actions:
  position - Move cursor to (x, y) coordinates
  press    - Press mouse button at current position
  release  - Release mouse button
  hold     - Press and hold for duration_ticks`,
		RunE: runMouseCommand,
	}
	mouseCmd.Flags().StringVarP(&mouseActionFlag, "action", "a", "position", "Action: position, press, release, or hold")
	mouseCmd.Flags().IntVarP(&xFlag, "x", "x", 0, "X coordinate")
	mouseCmd.Flags().IntVarP(&yFlag, "y", "y", 0, "Y coordinate")
	mouseCmd.Flags().StringVarP(&buttonFlag, "button", "b", "", "Mouse button (e.g., MouseButtonLeft, MouseButtonRight)")
	mouseCmd.Flags().Int64VarP(&durationTicks, "duration_ticks", "d", 6, "Duration in ticks for hold action")

	// wheel command
	wheelCmd := &cobra.Command{
		Use:   "wheel",
		Short: "Send wheel input to the game",
		RunE:  runWheelCommand,
	}
	wheelCmd.Flags().Float64VarP(&wheelXFlag, "x", "x", 0, "Horizontal scroll (negative=left, positive=right)")
	wheelCmd.Flags().Float64VarP(&wheelYFlag, "y", "y", 0, "Vertical scroll (negative=down, positive=up)")

	// screenshot command
	screenshotCmd := &cobra.Command{
		Use:   "screenshot",
		Short: "Take a screenshot of the game",
		Long: `Capture the game window and save to a file or output as base64.
If --output is not specified, a timestamped filename is generated.`,
		RunE: runScreenshotCommand,
	}
	screenshotCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file path (optional)")
	screenshotCmd.Flags().BoolVarP(&asyncFlag, "async", "a", false, "Async mode: return immediately without waiting for capture")

	// run command
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a script file",
		Long: `Execute a JSON script for automated game control.

Scripts support: input, mouse, wheel, screenshot, delay, and repeat commands.
Use --script for file path or --inline for JSON string.`,
		RunE: runScriptCommand,
	}
	runCmd.Flags().StringVarP(&scriptFlag, "script", "s", "", "Path to script file")
	runCmd.Flags().StringVar(&inlineFlag, "inline", "", "Inline JSON script string")

	// ping command
	pingCmd := &cobra.Command{
		Use:   "ping",
		Short: "Check if game is running",
		RunE:  runPingCommand,
	}

	// keys command
	keysCmd := &cobra.Command{
		Use:   "keys",
		Short: "List all available key names",
		RunE:  runKeysCommand,
	}

	// mouse_buttons command
	mouseButtonsCmd := &cobra.Command{
		Use:   "mouse_buttons",
		Short: "List all available mouse button names",
		RunE:  runMouseButtonsCommand,
	}

	rootCmd.AddCommand(inputCmd)
	rootCmd.AddCommand(mouseCmd)
	rootCmd.AddCommand(wheelCmd)
	rootCmd.AddCommand(screenshotCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(pingCmd)
	rootCmd.AddCommand(keysCmd)
	rootCmd.AddCommand(mouseButtonsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func persistentPreRunRootCommand(cmd *cobra.Command, args []string) error {
	if pidFlag == 0 {
		if err := cli.EnsureTargetPID(); err != nil {
			return err
		}
	} else {
		rpc.SetTargetPID(pidFlag)
	}
	return nil
}

func runInputCommand(cmd *cobra.Command, args []string) error {
	executor := cli.NewCommandExecutor()
	return executor.RunInputCommand(keyFlag, inputActionFlag, durationTicks)
}

func runMouseCommand(cmd *cobra.Command, args []string) error {
	executor := cli.NewCommandExecutor()
	return executor.RunMouseCommand(mouseActionFlag, xFlag, yFlag, buttonFlag, durationTicks)
}

func runWheelCommand(cmd *cobra.Command, args []string) error {
	executor := cli.NewCommandExecutor()
	return executor.RunWheelCommand(wheelXFlag, wheelYFlag)
}

func runScreenshotCommand(cmd *cobra.Command, args []string) error {
	executor := cli.NewCommandExecutor()
	return executor.RunScreenshotCommand(outputFlag, asyncFlag)
}

func runScriptCommand(cmd *cobra.Command, args []string) error {
	executor := cli.NewCommandExecutor()

	var input string
	switch {
	case inlineFlag != "":
		input = inlineFlag
	case scriptFlag != "":
		input = scriptFlag
	default:
		return fmt.Errorf("either --script or --inline must be provided")
	}

	return executor.RunScriptCommand(input, scriptFlag != "")
}

func runPingCommand(cmd *cobra.Command, args []string) error {
	executor := cli.NewCommandExecutor()
	return executor.RunPingCommand()
}

func runKeysCommand(cmd *cobra.Command, args []string) error {
	executor := cli.NewCommandExecutor()
	return executor.ListKeysCommand()
}

func runMouseButtonsCommand(cmd *cobra.Command, args []string) error {
	executor := cli.NewCommandExecutor()
	return executor.ListMouseButtonsCommand()
}
