# wayne
wain-like API for irrelevant platforms

## Platform Support

**wayne** provides cross-platform GUI widgets using [Ebitengine](https://ebitengine.org/) as the rendering backend.

### Supported Platforms

- ✅ **Windows** - Fully tested via CI
- ✅ **macOS** - Fully tested via CI  
- ⚠️ **Android** - Supported but requires manual testing (CI in progress)
- ⚠️ **iOS** - Supported but requires manual testing (CI in progress)

### Testing Status

| Platform | Build | Unit Tests | Integration Tests | CI Status |
|----------|-------|------------|-------------------|-----------|
| Windows  | ✅ Yes | ✅ Yes     | ⚠️ Partial       | ✅ Automated |
| macOS    | ✅ Yes | ✅ Yes     | ⚠️ Partial       | ✅ Automated |
| Android  | ✅ Yes | ⚠️ Manual  | ⚠️ Manual        | 🔄 Planned |
| iOS      | ✅ Yes | ⚠️ Manual  | ⚠️ Manual        | 🔄 Planned |

**Note**: Android and iOS testing requires physical devices or authorized emulators. These platforms are supported by the codebase (build tags, API surface) but automated testing infrastructure is still in development.

### Running Tests

Tests require CGO and platform-specific graphics contexts:

```bash
# On Windows
CGO_ENABLED=1 go test -v ./...

# On macOS  
CGO_ENABLED=1 go test -v ./...

# On Linux (not supported - build tags exclude Linux)
# Tests will not match any packages
```

For Android and iOS, tests must be run using `gomobile test` with appropriate device configurations (see Ebitengine documentation).

## Development

```bash
# Install dependencies
go mod download

# Run tests (on supported platforms)
CGO_ENABLED=1 go test -race ./...

# Run static analysis
go vet ./...
```
