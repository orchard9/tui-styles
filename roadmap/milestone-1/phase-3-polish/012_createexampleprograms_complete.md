## Purpose

Create example programs demonstrating core type usage (Color, AdaptiveColor, Position, Border). These serve as practical references for library users and validate that the API is usable.

## Acceptance Criteria

- [ ] `examples/colors/` program demonstrates Color usage
- [ ] `examples/adaptive-colors/` program demonstrates AdaptiveColor
- [ ] `examples/positions/` program demonstrates Position enum
- [ ] `examples/borders/` program demonstrates all border types
- [ ] All examples compile and run successfully
- [ ] Each example includes explanatory comments
- [ ] Examples demonstrate practical use cases
- [ ] Terminal output is visually clear and colorful

## Technical Approach

**Example Programs**:

1. **`examples/colors/main.go`** - Color demonstration:
   ```go
   package main

   import (
       "fmt"
       "log"

       "github.com/orchard9/tui-styles"
   )

   func main() {
       fmt.Println("TUI Styles - Color Examples\n")

       // Hex colors
       red, err := tuistyles.NewColor("#FF0000")
       if err != nil {
           log.Fatal(err)
       }
       fmt.Printf("%sRed from hex (#FF0000)\x1b[0m\n", red.ToANSI())

       // Short hex
       blue, _ := tuistyles.NewColor("#00F")
       fmt.Printf("%sBlue from short hex (#00F)\x1b[0m\n", blue.ToANSI())

       // ANSI color names
       green, _ := tuistyles.NewColor("green")
       fmt.Printf("%sGreen from ANSI name (green)\x1b[0m\n", green.ToANSI())

       // ANSI 256-color codes
       magenta, _ := tuistyles.NewColor("201")
       fmt.Printf("%sMagenta from ANSI code (201)\x1b[0m\n", magenta.ToANSI())

       // Error handling
       fmt.Println("\nError Handling:")
       if _, err := tuistyles.NewColor("invalid"); err != nil {
           fmt.Printf("Invalid color rejected: %v\n", err)
       }
   }
   ```

2. **`examples/adaptive-colors/main.go`** - AdaptiveColor demo:
   ```go
   package main

   import (
       "fmt"
       "log"
       "os"

       "github.com/orchard9/tui-styles"
   )

   func main() {
       fmt.Println("TUI Styles - Adaptive Color Examples\n")

       // Create adaptive color
       textColor, err := tuistyles.NewAdaptiveColor("#000000", "#FFFFFF")
       if err != nil {
           log.Fatal(err)
       }

       // Display current terminal background
       termBg := os.Getenv("TERM_BACKGROUND")
       if termBg == "" {
           termBg = "unknown (defaulting to dark)"
       }
       fmt.Printf("Terminal background: %s\n\n", termBg)

       // Show adaptive text
       color := textColor.ToColor()
       fmt.Printf("%sThis text adapts to your terminal background\x1b[0m\n", color.ToANSI())

       // Demonstrate manual override
       fmt.Println("\nManual Override:")
       os.Setenv("TERM_BACKGROUND", "light")
       fmt.Printf("%sLight terminal mode\x1b[0m\n", textColor.ToColor().ToANSI())

       os.Setenv("TERM_BACKGROUND", "dark")
       fmt.Printf("%sDark terminal mode\x1b[0m\n", textColor.ToColor().ToANSI())
   }
   ```

3. **`examples/positions/main.go`** - Position enum demo:
   ```go
   package main

   import (
       "fmt"

       "github.com/orchard9/tui-styles"
   )

   func main() {
       fmt.Println("TUI Styles - Position Examples\n")

       // Horizontal positions
       fmt.Println("Horizontal Positions:")
       fmt.Printf("  Left:   %s (value: %d)\n", tuistyles.Left, tuistyles.Left)
       fmt.Printf("  Center: %s (value: %d)\n", tuistyles.Center, tuistyles.Center)
       fmt.Printf("  Right:  %s (value: %d)\n", tuistyles.Right, tuistyles.Right)

       // Vertical positions
       fmt.Println("\nVertical Positions:")
       fmt.Printf("  Top:    %s (value: %d)\n", tuistyles.Top, tuistyles.Top)
       fmt.Printf("  Center: %s (value: %d)\n", tuistyles.Center, tuistyles.Center)
       fmt.Printf("  Bottom: %s (value: %d)\n", tuistyles.Bottom, tuistyles.Bottom)

       // Type checking
       fmt.Println("\nType Checking:")
       fmt.Printf("  Left.IsHorizontal(): %v\n", tuistyles.Left.IsHorizontal())
       fmt.Printf("  Left.IsVertical():   %v\n", tuistyles.Left.IsVertical())
       fmt.Printf("  Center.IsHorizontal(): %v\n", tuistyles.Center.IsHorizontal())
       fmt.Printf("  Center.IsVertical():   %v\n", tuistyles.Center.IsVertical())
   }
   ```

4. **`examples/borders/main.go`** - Border styles demo:
   ```go
   package main

   import (
       "fmt"

       "github.com/orchard9/tui-styles"
   )

   func main() {
       fmt.Println("TUI Styles - Border Examples\n")

       borders := []struct {
           name   string
           border tuistyles.Border
       }{
           {"Normal", tuistyles.NormalBorder()},
           {"Rounded", tuistyles.RoundedBorder()},
           {"Thick", tuistyles.ThickBorder()},
           {"Double", tuistyles.DoubleBorder()},
           {"Block", tuistyles.BlockBorder()},
           {"Outer Half Block", tuistyles.OuterHalfBlockBorder()},
           {"Inner Half Block", tuistyles.InnerHalfBlockBorder()},
           {"Hidden", tuistyles.HiddenBorder()},
       }

       for _, b := range borders {
           fmt.Printf("%s Border:\n", b.name)
           fmt.Printf("  %s────────%s\n", b.border.TopLeft, b.border.TopRight)
           fmt.Printf("  %s        %s\n", b.border.Left, b.border.Right)
           fmt.Printf("  %s────────%s\n\n", b.border.BottomLeft, b.border.BottomRight)
       }
   }
   ```

**Directory Structure**:
```
examples/
├── colors/
│   └── main.go
├── adaptive-colors/
│   └── main.go
├── positions/
│   └── main.go
└── borders/
    └── main.go
```

**Files to Create/Modify**:
- `examples/colors/main.go`
- `examples/adaptive-colors/main.go`
- `examples/positions/main.go`
- `examples/borders/main.go`

**Dependencies**:
- All Phase 2 tasks complete (core types implemented)

## Testing Strategy

**Manual Testing**:
```bash
# Build and run each example
cd examples/colors
go run main.go

cd ../adaptive-colors
go run main.go

cd ../positions
go run main.go

cd ../borders
go run main.go
```

**Validation Checklist**:
- [ ] All examples compile without errors
- [ ] Output is visually clear and formatted
- [ ] Colors render correctly in terminal
- [ ] Adaptive colors respond to TERM_BACKGROUND
- [ ] Border characters display correctly (Unicode support)
- [ ] Error handling examples work as expected
- [ ] Comments explain what each section demonstrates

**Terminal Testing**:
Test examples on multiple terminals:
- iTerm2 (macOS)
- Terminal.app (macOS)
- Alacritty (cross-platform)
- GNOME Terminal (Linux)
- Windows Terminal (Windows)

## Notes

**Visual Presentation**: Examples should be visually appealing and demonstrate the library's capabilities. Use colors, formatting, and clear output.

**Error Handling**: Show both successful usage and error cases. Teach users proper error handling patterns.

**Comments**: Add explanatory comments so users understand what each section does. Examples serve as tutorials.

**Keep It Simple**: Don't over-engineer examples. Show practical, real-world usage patterns.

**Reset Codes**: Always reset ANSI codes after colored output (`\x1b[0m`) to avoid affecting subsequent terminal output.

**Example Programs as Tests**: Consider running examples in CI to ensure they always compile:
```bash
# In CI pipeline
for dir in examples/*/; do
    echo "Building ${dir}..."
    (cd "$dir" && go build)
done
```

**README Integration**: Reference examples in README Quick Start section with links:
```markdown
See [examples/colors](examples/colors/main.go) for complete color usage.
```

**Future Enhancements**:
- Add `examples/advanced/` for complex use cases (Milestone 2)
- Add `examples/themes/` for pre-built color schemes
- Add benchmarking examples

**Reference**: Similar to lipgloss examples at https://github.com/charmbracelet/lipgloss/tree/master/examples




