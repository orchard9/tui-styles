## Purpose

Create a sophisticated dashboard example demonstrating complex composition, showcasing the library's capabilities for building real-world terminal UIs.

## Acceptance Criteria

- [ ] examples/dashboard/main.go created (~100-150 lines)
- [ ] Dashboard includes: title header, stats boxes, colored status indicators
- [ ] Demonstrates nested box compositions
- [ ] Shows mixed alignment (center headers, left content, right metrics)
- [ ] Uses multiple border styles for visual hierarchy
- [ ] Includes color-coded elements (success=green, error=red, warning=yellow)
- [ ] Compiles and runs with visually appealing output
- [ ] Documented with inline comments explaining composition

## Technical Approach

**Dashboard Layout**:
```
┌──────────────────────────────────────────────────┐
│              System Dashboard v1.0               │
└──────────────────────────────────────────────────┘

┏━━━━━━━━━━━┓ ┏━━━━━━━━━━━┓ ┏━━━━━━━━━━━┓
┃ CPU: 45%  ┃ ┃ Memory    ┃ ┃ Disk      ┃
┃ ✓ Normal  ┃ ┃ 8.2/16 GB ┃ ┃ 256/512GB ┃
┗━━━━━━━━━━━┛ ┗━━━━━━━━━━━┛ ┗━━━━━━━━━━━┛

╔═══════════════════════════════════════════════╗
║ Recent Activity                               ║
╟───────────────────────────────────────────────╢
║ ✓ Deployment succeeded       2 minutes ago    ║
║ ⚠ High memory usage          5 minutes ago    ║
║ ✗ Build failed               10 minutes ago   ║
╚═══════════════════════════════════════════════╝
```

**Composition Elements**:
1. Title header (centered, bold, underlined)
2. Stats boxes (3 side-by-side, different border styles)
3. Activity log (nested box, styled entries)
4. Status indicators (colored symbols: ✓✗⚠)
5. Timestamps (right-aligned, dimmed)

**Code Structure**:
```go
func main() {
    // Title
    title := renderTitle()

    // Stats boxes
    cpu := renderCPUBox()
    memory := renderMemoryBox()
    disk := renderDiskBox()
    stats := renderStatsRow(cpu, memory, disk)

    // Activity log
    activity := renderActivityLog()

    // Compose dashboard
    fmt.Println(title)
    fmt.Println()
    fmt.Println(stats)
    fmt.Println()
    fmt.Println(activity)
}
```

**Styling Strategy**:
- Headers: Bold + underline
- Success indicators: Green foreground
- Error indicators: Red foreground
- Warning indicators: Yellow foreground
- Timestamps: Dim/gray
- Boxes: Various border styles (single, double, rounded, thick)

**Files to Create/Modify**:
- examples/dashboard/main.go (new)
- examples/README.md (add dashboard description)

**Dependencies**:
- tui-styles library (all packages)

## Testing Strategy

**Visual Validation**:
- Run in macOS Terminal and verify rendering
- Test in iTerm2 for color accuracy
- Ensure alignment is consistent
- Verify borders connect properly
- Check Unicode symbols display correctly

**Code Quality**:
- Modular functions for each dashboard component
- Clear inline comments explaining composition
- No hardcoded magic numbers (use constants)
- Clean separation of concerns

**Demo Script**:
```bash
cd examples/dashboard
go run main.go
```

## Notes

**Dashboard Components**:

1. **Title Header**:
```go
func renderTitle() string {
    s := style.New().Bold().Underline().Foreground(color.Cyan)
    text := "System Dashboard v1.0"
    centered := layout.Center(text, 50)
    return border.Box(s.Render(centered), border.Single, 52, 3)
}
```

2. **Stats Box**:
```go
func renderCPUBox() string {
    content := "CPU: 45%\n✓ Normal"
    s := style.New().Foreground(color.Green)
    return border.Box(s.Render(content), border.Double, 15, 4)
}
```

3. **Activity Log**:
```go
func renderActivityLog() string {
    entries := []string{
        renderEntry("✓", "Deployment succeeded", "2 min ago", color.Green),
        renderEntry("⚠", "High memory usage", "5 min ago", color.Yellow),
        renderEntry("✗", "Build failed", "10 min ago", color.Red),
    }
    content := strings.Join(entries, "\n")
    header := style.New().Bold().Render("Recent Activity")
    return border.Box(header+"\n"+content, border.Thick, 50, 8)
}
```

**Visual Polish**:
- Consistent spacing between elements
- Color hierarchy (errors stand out)
- Proper alignment (metrics right-aligned, timestamps right-aligned)
- Border variety for visual interest

**Future Enhancements** (post-v1.0):
- examples/dashboard-live/ - real-time updating dashboard
- examples/table/ - data table rendering
- examples/progress/ - progress bars and spinners

**Reference**:
- Inspiration: htop, btop, lazygit TUIs
- Border styles: Box Drawing Unicode characters
- Color theory: Use sparingly for emphasis


