package golden

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// Delimiter between section name and the line itself.
const Delimiter = "::"

// Golden file section identifiers.
const (
	SecComment   = Delimiter
	SecReqMethod = "ReqMethod"
	SecReqPath   = "ReqPath"
	SecReqQuery  = "ReqQuery"
	SecRspCode   = "RspCode"
	SecHeader    = "Header"
	SecBody      = "Body"
)

// Parsing modifiers.
const (
	ModNone  = ""  // No modifiers set for given section.
	ModMerge = "+" // Concatenate lines in the section without new lines.
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

// Golden represents golden file.
type Golden struct {
	sections []*Section // Golden file sections.
	t        T          // Test manager.
}

// NewGolden reads golden file at pth and creates new instance of Golden.
func NewGolden(t T, rdr io.Reader) *Golden {
	t.Helper()

	gld := &Golden{
		t: t,
	}

	var err error
	if gld.sections, err = parse(rdr); err != nil {
		t.Fatal(err)
	}

	return gld
}

// Section returns golden file section or nil if section is not present.
func (gld *Golden) Section(id string) *Section {
	for _, sec := range gld.sections {
		if sec.ID() == id {
			return sec
		}
	}
	return nil
}

// SectionCount returns number of sections in the golden file.
func (gld *Golden) SectionCount() int {
	return len(gld.sections)
}

// parse parses reader r representing golden file.
func parse(r io.Reader) ([]*Section, error) {
	scn := bufio.NewScanner(r)

	header := true
	var last *Section
	var sections []*Section
	for scn.Scan() {
		line := scn.Text()

		// Skip empty lines unless reading body section.
		if last != nil && line == "" && last.ID() != SecBody {
			continue
		}

		sec := NewSection(line)
		if sec != nil {
			if sec.ID() == SecComment && !header {
				return nil, errors.New("did not expect comment")
			}
			header = sec.ID() == SecComment

			if last != nil && last.ID() == sec.ID() {
				last.Add(line)
				continue
			}

			sections = append(sections, sec)
			last = sec
			continue
		}

		// Add line to last seen section.
		if last != nil {
			last.Add(line)
		}
	}

	if err := scn.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

// isComment returns true if line is a comment.
func isComment(line string) bool { return strings.HasPrefix(line, Delimiter) }
