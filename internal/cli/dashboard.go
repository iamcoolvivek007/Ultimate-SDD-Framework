package cli

import (
	"fmt"
	"os/exec"
	"runtime"

	"ultimate-sdd-framework/web"

	"github.com/spf13/cobra"
)

func NewDashboardCmd() *cobra.Command {
	var port int
	var noBrowser bool

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "🌐 Open the web dashboard",
		Long: `Start the Viki web dashboard for a visual, beginner-friendly experience.

The dashboard provides:
• Visual workflow pipeline (Idea → Spec → Plan → Tasks → Code)
• Click-based navigation instead of CLI commands
• Live AI chat sidebar
• Agent selector (21+ AI personas)
• Real-time progress tracking

Perfect for beginners and visual thinkers!`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Open browser unless --no-browser flag
			if !noBrowser {
				url := fmt.Sprintf("http://localhost:%d", port)
				go openBrowser(url)
			}

			// Start web server
			server := web.NewServer(port)
			return server.Start()
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 3000, "Port to run dashboard on")
	cmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Don't open browser automatically")

	return cmd
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		fmt.Printf("Open %s in your browser\n", url)
		return
	}

	exec.Command(cmd, args...).Start()
}
