# wayne
API-compatible cross-platform GUI for Windows, macOS, Android, iOS

## Platform Support

**wayne** provides cross-platform GUI widgets using [Ebitengine](https://ebitengine.org/) as the rendering backend. It maintains API compatibility with [opd-ai/wain](https://github.com/opd-ai/wain), enabling source-level migration for applications targeting non-Linux platforms.

### Supported Platforms

- ✅ **Windows** - Fully tested via CI
- ✅ **macOS** - Fully tested via CI  
- ✅ **Android** - Tested via CI (emulator)
- ✅ **iOS** - Tested via CI (simulator)
- ❌ **Linux/BSD** - Not supported. For Linux, use [opd-ai/wain](https://github.com/opd-ai/wain) instead

### Testing Status

| Platform | Build | Unit Tests | Integration Tests | CI Status |
|----------|-------|------------|-------------------|-----------|
| Windows  | ✅ Yes | ✅ Yes     | ⚠️ Partial       | ✅ Automated |
| macOS    | ✅ Yes | ✅ Yes     | ⚠️ Partial       | ✅ Automated |
| Android  | ✅ Yes | ✅ Yes     | ⚠️ Partial       | ✅ Automated |
| iOS      | ✅ Yes | ✅ Yes     | ⚠️ Partial       | ✅ Automated |

**Note**: Android tests run on an x86_64 emulator (API 29) and iOS tests run on an iPhone simulator. For device-specific testing, use `gomobile test` with appropriate device configurations.

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
