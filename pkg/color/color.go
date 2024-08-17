package color

// Foreground colors - provides constants for foreground (text) colors that can
// be used with ANSI escape codes to change the color of terminal output.
const (
	Reset  = "\033[0m"  // Reset color to default
	Black  = "\033[30m" // Black
	Red    = "\033[31m" // Red
	Green  = "\033[32m" // Green
	Yellow = "\033[33m" // Yellow
	Blue   = "\033[34m" // Blue
	Purple = "\033[35m" // Purple
	Cyan   = "\033[36m" // Cyan
	White  = "\033[37m" // White

)

// Make - adds the specified color to the text and resets the color at the end
// to prevent subsequent text from being colored.
func Make(text string, colorCode string) string {
	return colorCode + text + Reset
}
