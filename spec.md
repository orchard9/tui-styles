# Terminal CSS Specification

## Overview
This specification defines the API and behavior for a terminal styling library. The design strictly follows the chaining/builder pattern found in `lipgloss`. The core primitive is a `Style` definition that is immutable and returns new instances upon modification.

## 1. Core Types

### 1.1. Color
Represents a terminal color.
- **`Color(string)`**: Supports hex strings (`"#000000"`), ANSI color names, or distinct ANSI color codes.
- **`AdaptiveColor{Light: string, Dark: string}`**: Automatically selects the color based on the terminal's background detection (if available) or defaults.

### 1.2. Style
The fundamental building block. It holds the state for all styling attributes (colors, spacing, borders, etc.).
- **`NewStyle()`**: Returns a fresh, empty `Style`.

## 2. Style API
All methods below are defined on the `Style` type and must return a `Style` (fluent interface).

### 2.1. Text Attributes
- **`Bold(bool)`**: Sets the bold attribute.
- **`Italic(bool)`**: Sets the italic attribute.
- **`Underline(bool)`**: Sets the underline attribute.
- **`Strikethrough(bool)`**: Sets the strikethrough attribute.
- **`Faint(bool)`**: Sets the faint attribute.
- **`Blink(bool)`**: Sets the blink attribute.
- **`Reverse(bool)`**: Reverses the foreground and background colors.

### 2.2. Coloring
- **`Foreground(Color)`**: Sets the text color.
- **`Background(Color)`**: Sets the background color.
- **`SetString(string)`**: Sets the string to be rendered if not passed explicitly to `Render`.

### 2.3. Sizing & Alignment
- **`Width(int)`**: Sets the width of the container. Content matches alignment within this width.
- **`Height(int)`**: Sets the height of the container.
- **`Align(Position)`**: Sets horizontal alignment (`Left`, `Center`, `Right`).
- **`AlignVertical(Position)`**: Sets vertical alignment (`Top`, `Center`, `Bottom`).
- **`MaxWidth(int)`**: Limits the maximum width.

### 2.4. Spacing (Padding & Margin)
Both Padding and Margin support 1, 2, or 4 arguments (CSS shorthand style).
- **`Padding(top, right, bottom, left int)`**: Adds internal space inside the border.
- **`PaddingTop(int)`, `PaddingRight(int)`, `PaddingBottom(int)`, `PaddingLeft(int)`**
- **`Margin(top, right, bottom, left int)`**: Adds external space outside the border.
- **`MarginTop(int)`, `MarginRight(int)`, `MarginBottom(int)`, `MarginLeft(int)`**

### 2.5. Borders
- **`Border(BorderType, top, right, bottom, left bool)`**: Sets the border style and visibility per side.
- **`BorderForeground(Color)`**: Sets the color of the border.
- **`BorderBackground(Color)`**: Sets the background color of the border.
- **`BorderTop(bool)`, `BorderRight(bool)`, `BorderBottom(bool)`, `BorderLeft(bool)`**: Toggles specific sides.

**Standard Border Types:**
- `NormalBorder()`
- `RoundedBorder()`
- `BlockBorder()`
- `OuterHalfBlockBorder()`
- `InnerHalfBlockBorder()`
- `ThickBorder()`
- `DoubleBorder()`
- `HiddenBorder()`

## 3. Rendering
- **`Render(string) string`**: Applies the style to the provided string and returns the ANSI-escaped string.
- **`String() string`**: Renders the internal string set via `SetString` (if any).

## 4. Layout Utilities
These are package-level functions for composition.
- **`JoinHorizontal(Position, ...string) string`**: Joins multi-line strings horizontally with specified alignment (Top, Center, Bottom).
- **`JoinVertical(Position, ...string) string`**: Joins multi-line strings vertically with specified alignment (Left, Center, Right).
- **`Place(width, height int, hPos, vPos Position, content string) string`**: Places content within a given box size at specific coordinates.

## 5. Implementation Requirements
1.  **Immutability**: Calling `MyStyle.Foreground(...)` must *not* mutate `MyStyle`. It must return a new copy with the change applied.
2.  **Whitespace Handling**: Padding and Margins must render as actual spaces (or background-colored spaces).
3.  **Newline Handling**: Styles applied to multi-line strings must apply correctly to each line (especially background colors and borders).
4.  **ANSI Awareness**: The renderer must calculate string width ignoring ANSI escape codes (using a library like `mattn/go-runewidth` or similar logic).

## 6. Example Usage
```go
var style = NewStyle().
    Bold(true).
    Foreground(Color("#FAFAFA")).
    Background(Color("#7D56F4")).
    PaddingTop(2).
    PaddingLeft(4).
    Width(22)

fmt.Println(style.Render("Hello, World!"))
```
