package golden

// Types of payload.
const (
	// PayloadText represents text payload (default).
	PayloadText = "text"

	// PayloadJSON represents JSON payload.
	PayloadJSON = "json"
)

// T is a subset of testing.TB interface.
type T interface {
	// Fatal is equivalent to Log followed by FailNow.
	Fatal(args ...interface{})

	// Fatalf is equivalent to Logf followed by FailNow.
	Fatalf(format string, args ...interface{})

	// Helper marks the calling function as a test helper function.
	// When printing file and line information, that function will be skipped.
	// Helper may be called simultaneously from multiple goroutines.
	Helper()
}
