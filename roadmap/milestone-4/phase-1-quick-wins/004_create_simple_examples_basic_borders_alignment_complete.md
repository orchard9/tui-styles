## Purpose

Create runnable example programs demonstrating basic library usage to help new users get started quickly and validate API usability.

## Acceptance Criteria

- [ ] examples/ directory created with three programs
- [ ] examples/basic/main.go - simple text styling (colors, bold, italic)
- [ ] examples/borders/main.go - box rendering with different border styles
- [ ] examples/alignment/main.go - text alignment and padding
- [ ] Each example compiles and runs without errors
- [ ] Each example includes inline comments explaining the code
- [ ] Each example produces visually appealing terminal output
- [ ] README.md in examples/ directory with screenshots and descriptions

## Technical Approach

**Directory Structure**:
```
examples/
├── README.md
├── basic/
│   └── main.go
├── borders/
│   └── main.go
└── alignment/
    └── main.go
```

**Example 1: Basic Styling** (examples/basic/main.go)
- Demonstrates Style.New(), Bold(), Italic(), Underline()
- Shows foreground and background colors
- Prints styled text with different combinations
- ~30 lines of simple, readable code

**Example 2: Borders** (examples/borders/main.go)
- Demonstrates border.Box() with different styles
- Shows SingleBorder, DoubleBorder, RoundedBorder
- Renders multiple boxes with titles
- Demonstrates border colors and padding
- ~40 lines showcasing border variety

**Example 3: Alignment** (examples/alignment/main.go)
- Demonstrates text alignment (left, center, right)
- Shows padding and margin utilities
- Renders aligned text blocks
- Demonstrates width constraints
- ~35 lines showing layout capabilities

**Example README**:
- Brief description of each example
- Instructions to run: `go run examples/basic/main.go`
- Visual output (ASCII art or screenshots)
- Link to main README and godoc

**Files to Create/Modify**:
- examples/README.md (new)
- examples/basic/main.go (new)
- examples/borders/main.go (new)
- examples/alignment/main.go (new)

**Dependencies**:
- tui-styles library (local)

## Testing Strategy

**Compilation**:
- Run `go build` in each example directory
- Verify no compilation errors
- Test with `go run` to ensure runtime correctness

**Visual Validation**:
- Run each example in terminal and verify output
- Test on macOS Terminal and iTerm2
- Ensure colors render correctly
- Verify borders draw properly with Unicode characters

**Code Quality**:
- Ensure examples follow Go conventions
- Add comments explaining each step
- Keep code simple and readable (educational purpose)
- No error handling complexity (focus on library usage)

## Notes

**Example Template**:
```go
package main

import (
    "fmt"
    "github.com/yourusername/tui-styles/style"
    "github.com/yourusername/tui-styles/color"
)

func main() {
    // Create a bold red style
    errorStyle := style.New().Bold().Foreground(color.Red)
    fmt.Println(errorStyle.Render("Error: Something went wrong"))

    // Create a green success style
    successStyle := style.New().Foreground(color.Green)
    fmt.Println(successStyle.Render("Success: Operation completed"))
}
```

**Best Practices for Examples**:
- One concept per example (focus)
- Short and runnable (<50 lines)
- Include comments for clarity
- Show practical use cases
- Visual output that demonstrates value

**Visual Output Goals**:
- Basic: Colorful text with various styles
- Borders: 3-4 boxes with different border styles
- Alignment: Text blocks showing left/center/right alignment

**Future Examples** (Phase 2):
- examples/dashboard/ - complex composition
- examples/table/ - data table rendering
- examples/progress/ - progress bars and indicators

**Reference**:
- Go by Example: https://gobyexample.com/
- Cobra examples: https://github.com/spf13/cobra/tree/main/doc

