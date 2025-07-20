package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/logger"
	"github.com/rsmacapinlac/pomodux/internal/timer"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show recent timer session(s)",
	Long:  `Show recent completed timer sessions with filtering and export options.`,
	RunE:  runHistory,
}

var (
	historyJSON   bool
	historyCSV    bool
	historyLimit  int
	historyType   string
	historyDate   string
	historyStats  bool
	historyExport string
)

func init() {
	historyCmd.Flags().BoolVar(&historyJSON, "json", false, "Output history as JSON")
	historyCmd.Flags().BoolVar(&historyCSV, "csv", false, "Output history as CSV")
	historyCmd.Flags().IntVar(&historyLimit, "limit", 10, "Number of sessions to show")
	historyCmd.Flags().StringVar(&historyType, "type", "", "Filter by session type (work, break, long-break)")
	historyCmd.Flags().StringVar(&historyDate, "date", "", "Filter by date (YYYY-MM-DD)")
	historyCmd.Flags().BoolVar(&historyStats, "stats", false, "Show session statistics")
	historyCmd.Flags().StringVar(&historyExport, "export", "", "Export to file (specify path)")
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	historyManager, err := timer.NewHistoryManager()
	if err != nil {
		return fmt.Errorf("failed to create history manager: %w", err)
	}

	sessions, err := historyManager.GetRecentSessions(100) // Get more sessions for filtering
	if err != nil {
		return fmt.Errorf("failed to get session history: %w", err)
	}

	// Apply filters
	filteredSessions := filterSessions(sessions, historyType, historyDate)

	// Limit results
	if historyLimit > 0 && len(filteredSessions) > historyLimit {
		filteredSessions = filteredSessions[:historyLimit]
	}

	// Handle export
	if historyExport != "" {
		return exportHistory(filteredSessions, historyExport, historyJSON, historyCSV)
	}

	// Show statistics if requested
	if historyStats {
		showStatistics(filteredSessions)
		return nil
	}

	// Handle output format
	if historyJSON {
		return outputHistoryJSON(filteredSessions)
	}

	if historyCSV {
		return outputHistoryCSV(filteredSessions)
	}

	// Default text output
	return outputHistoryText(filteredSessions)
}

func filterSessions(sessions []timer.SessionRecord, sessionType, date string) []timer.SessionRecord {
	var filtered []timer.SessionRecord

	for _, session := range sessions {
		// Filter by session type
		if sessionType != "" && string(session.Type) != sessionType {
			continue
		}

		// Filter by date
		if date != "" {
			sessionDate := session.StartTime.Format("2006-01-02")
			if sessionDate != date {
				continue
			}
		}

		filtered = append(filtered, session)
	}

	return filtered
}

func showStatistics(sessions []timer.SessionRecord) {
	if len(sessions) == 0 {
		fmt.Println("No sessions found for statistics.")
		return
	}

	var totalWorkTime, totalBreakTime time.Duration
	var workSessions, breakSessions, longBreakSessions int
	var completedSessions int

	for _, session := range sessions {
		actualDuration := session.EndTime.Sub(session.StartTime)

		switch session.Type {
		case timer.SessionTypeWork:
			workSessions++
			totalWorkTime += actualDuration
		case timer.SessionTypeBreak:
			breakSessions++
			totalBreakTime += actualDuration
		case timer.SessionTypeLongBreak:
			longBreakSessions++
			totalBreakTime += actualDuration
		}

		if session.Completed {
			completedSessions++
		}
	}

	totalSessions := len(sessions)
	completionRate := float64(completedSessions) / float64(totalSessions) * 100

	fmt.Printf("Session Statistics:\n")
	fmt.Printf("==================\n")
	fmt.Printf("Total Sessions:     %d\n", totalSessions)
	fmt.Printf("Completed Sessions: %d (%.1f%%)\n", completedSessions, completionRate)
	fmt.Printf("\nSession Types:\n")
	fmt.Printf("  Work Sessions:    %d (%.1f%%)\n", workSessions, float64(workSessions)/float64(totalSessions)*100)
	fmt.Printf("  Break Sessions:   %d (%.1f%%)\n", breakSessions, float64(breakSessions)/float64(totalSessions)*100)
	fmt.Printf("  Long Break Sessions: %d (%.1f%%)\n", longBreakSessions, float64(longBreakSessions)/float64(totalSessions)*100)
	fmt.Printf("\nTime Totals:\n")
	fmt.Printf("  Total Work Time:  %s\n", formatDuration(totalWorkTime))
	fmt.Printf("  Total Break Time: %s\n", formatDuration(totalBreakTime))
	fmt.Printf("  Total Time:       %s\n", formatDuration(totalWorkTime+totalBreakTime))

	if workSessions > 0 {
		avgWorkTime := totalWorkTime / time.Duration(workSessions)
		fmt.Printf("  Average Work Session: %s\n", formatDuration(avgWorkTime))
	}
}

func outputHistoryJSON(sessions []timer.SessionRecord) error {
	type sessionOutput struct {
		Type           string    `json:"type"`
		Duration       string    `json:"duration"`
		StartTime      time.Time `json:"start_time"`
		EndTime        time.Time `json:"end_time"`
		ActualDuration string    `json:"actual_duration"`
		Completed      bool      `json:"completed"`
	}

	var output []sessionOutput
	for _, session := range sessions {
		output = append(output, sessionOutput{
			Type:           string(session.Type),
			Duration:       formatDuration(session.Duration),
			StartTime:      session.StartTime,
			EndTime:        session.EndTime,
			ActualDuration: formatDuration(session.EndTime.Sub(session.StartTime)),
			Completed:      session.Completed,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

func outputHistoryCSV(sessions []timer.SessionRecord) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header
	header := []string{"Type", "Duration", "Start Time", "End Time", "Actual Duration", "Completed"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, session := range sessions {
		row := []string{
			string(session.Type),
			formatDuration(session.Duration),
			session.StartTime.Format("2006-01-02 15:04:05"),
			session.EndTime.Format("2006-01-02 15:04:05"),
			formatDuration(session.EndTime.Sub(session.StartTime)),
			strconv.FormatBool(session.Completed),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func outputHistoryText(sessions []timer.SessionRecord) error {
	logger.Debug("outputHistoryText called", map[string]interface{}{"session_count": len(sessions)})
	if len(sessions) == 0 {
		fmt.Println("No session history found.")
		return nil
	}

	fmt.Printf("Recent completed sessions (%d total):\n\n", len(sessions))
	for i, session := range sessions {
		fmt.Printf("%d. %s Session\n", i+1, session.Type)
		fmt.Printf("   Duration: %s\n", session.Duration)
		fmt.Printf("   Start: %s\n", session.StartTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("   End: %s\n", session.EndTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Actual Duration: %s\n", session.EndTime.Sub(session.StartTime))
		fmt.Printf("   Completed: %t\n\n", session.Completed)
	}

	return nil
}

func exportHistory(sessions []timer.SessionRecord, filepath string, jsonFormat, csvFormat bool) error {
	// Validate file path for security
	if err := validateExportPath(filepath); err != nil {
		return fmt.Errorf("invalid export path: %w", err)
	}

	file, err := os.Create(filepath) // #nosec G304 -- filepath is validated by validateExportPath
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	if jsonFormat {
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		return enc.Encode(sessions)
	}

	if csvFormat {
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		header := []string{"Type", "Duration", "Start Time", "End Time", "Actual Duration", "Completed"}
		if err := writer.Write(header); err != nil {
			return err
		}

		// Write data
		for _, session := range sessions {
			row := []string{
				string(session.Type),
				formatDuration(session.Duration),
				session.StartTime.Format("2006-01-02 15:04:05"),
				session.EndTime.Format("2006-01-02 15:04:05"),
				formatDuration(session.EndTime.Sub(session.StartTime)),
				strconv.FormatBool(session.Completed),
			}
			if err := writer.Write(row); err != nil {
				return err
			}
		}

		return nil
	}

	// Default to text format
	return outputHistoryText(sessions)
}

func getStateDir() (string, error) {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		stateHome = filepath.Join(homeDir, ".local", "state")
	}
	return filepath.Join(stateHome, "pomodux"), nil
}

// validateExportPath validates an export file path for security
func validateExportPath(exportPath string) error {
	// Check for path traversal attempts
	if strings.Contains(exportPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Check for dangerous characters
	dangerousChars := []string{"|", "&", ";", "`", "$", "(", ")", "<", ">", "*", "?"}
	for _, char := range dangerousChars {
		if strings.Contains(exportPath, char) {
			return fmt.Errorf("dangerous character '%s' not allowed in file path", char)
		}
	}

	// For export files, we can be more permissive but still validate
	// Allow relative paths and absolute paths in user's home directory
	if filepath.IsAbs(exportPath) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory for path validation")
		}

		// Check if path is within user's home directory
		relPath, err := filepath.Rel(homeDir, exportPath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			return fmt.Errorf("export path must be within user's home directory")
		}
	}

	return nil
}
