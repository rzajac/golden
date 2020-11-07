package golden

// Golden file body type.
const (
	// TypeText represents golden file text body type (default).
	TypeText = "text"

	// TypeJSON represents golden file JSON body type.
	TypeJSON = "json"
)

// T is a subset of testing.TB interface.
// It's primarily used to test golden package but can be used to implement
// custom actions to be taken on errors.
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
