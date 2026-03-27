# Wayne Examples

This directory contains example applications demonstrating wayne's features.

## Platform Requirements

All examples require building on a supported platform:
- **Windows**: Requires CGO and a C compiler (MinGW-w64 or MSVC)
- **macOS**: Requires CGO and Xcode Command Line Tools
- **Android/iOS**: Requires gomobile (see TESTING.md for setup)

## Examples

### hello

Minimal application with a button and label.

**Demonstrates:**
- App and Window creation
- Basic Panel layout
- Button click handling
- Label text updates

```bash
cd examples/hello
CGO_ENABLED=1 go build -o hello .
./hello
```

### form

Form application with text inputs and validation.

**Demonstrates:**
- Form layout with labels and text inputs
- Row containers for horizontal layout
- Reading and validating text input values
- Status feedback to users

```bash
cd examples/form
CGO_ENABLED=1 go build -o form .
./form
```

### scrollview

Scrollable list with many items.

**Demonstrates:**
- ScrollView for content larger than viewport
- Dynamic item generation
- Scroll offset control and feedback
- Buttons to control scroll position

```bash
cd examples/scrollview
CGO_ENABLED=1 go build -o scrollview .
./scrollview
```

## Building All Examples

```bash
# From the wayne repository root
cd examples

# Build all (Windows)
for d in */; do (cd "$d" && CGO_ENABLED=1 go build .); done

# Build all (macOS/Linux shell)
for d in */; do (cd "$d" && CGO_ENABLED=1 go build .); done
```

## Notes

- Examples use percentage-based sizing (0-100) for widget dimensions
- The window coordinate system has (0,0) at the top-left
- Widgets inherit theme from their parent containers unless overridden
- On high-DPI displays, set `Theme.Scale` for proper scaling
