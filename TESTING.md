# Testing Guide for Wayne

This document describes the testing strategy and infrastructure for the wayne library.

## Platform Testing Matrix

Wayne is a cross-platform GUI library targeting Windows, macOS, Android, and iOS. Each platform has different testing requirements and capabilities.

### Current Testing Status

| Platform | Automated Tests | CI/CD | Manual Testing Required |
|----------|----------------|-------|-------------------------|
| **Windows** | ✅ Yes | ✅ GitHub Actions (windows-latest) | ⚠️ Visual/Integration |
| **macOS** | ✅ Yes | ✅ GitHub Actions (macos-latest) | ⚠️ Visual/Integration |
| **Android** | ⚠️ Limited | 🔄 Planned | ✅ Required |
| **iOS** | ⚠️ Limited | 🔄 Planned | ✅ Required |

## Running Tests

### Prerequisites

All platforms require:
- **CGO enabled**: Tests need CGO for Ebitengine's graphics backend
- **Go 1.21+**: Minimum Go version supported by Ebitengine v2
- **Platform-specific graphics libraries**: Automatically handled by Ebitengine dependencies

### On Windows

```bash
# Set CGO_ENABLED and run tests
set CGO_ENABLED=1
go test -v -race ./...

# Or in PowerShell
$env:CGO_ENABLED=1
go test -v -race ./...
```

Tests will execute:
- General test files: `*_test.go` (with platform build tags)
- Windows-specific: `app_windows_test.go`

### On macOS

```bash
# Set CGO_ENABLED and run tests
export CGO_ENABLED=1
go test -v -race ./...
```

Tests will execute:
- General test files: `*_test.go` (with platform build tags)
- macOS-specific: `app_darwin_test.go`

### On Android

Android testing requires gomobile and an Android emulator or physical device:

```bash
# Install gomobile
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init

# Run tests on Android device/emulator
# Note: Requires adb and a running emulator or connected device
gomobile test -target=android ./...
```

Platform-specific tests: `app_android_test.go`

**Current Status**: Android tests are defined but require manual execution with proper gomobile setup. CI automation is planned.

### On iOS

iOS testing requires gomobile and iOS Simulator or physical device:

```bash
# Install gomobile (macOS only)
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init

# Run tests on iOS simulator
gomobile test -target=ios ./...
```

Platform-specific tests: `app_ios_test.go`

**Current Status**: iOS tests are defined but require manual execution with proper gomobile setup. CI automation is planned.

### On Linux

**Linux is not a supported platform.** All source files have build tags `//go:build windows || darwin || android || ios` which explicitly exclude Linux and other Unix systems.

Running `go test ./...` on Linux will report "no packages to test" because no files match the build constraints.

## Test Organization

### General Tests

These test files have platform build tags and run on all supported platforms:
- `app_test.go` - Application lifecycle and window management
- `concretewidgets_test.go` - Widget behavior and properties
- `dispatcher_test.go` - Event dispatching and hit-testing
- `event_test.go` - Event creation and properties
- `game_test.go` - Game loop and input translation
- `layout_test.go` - Layout algorithms and constraints
- `resource_test.go` - Resource loading and management
- `theme_propagation_test.go` - Theme inheritance

### Platform-Specific Smoke Tests

These test files verify basic functionality on specific platforms:
- `app_windows_test.go` - Windows-specific smoke tests
- `app_darwin_test.go` - macOS-specific smoke tests
- `app_android_test.go` - Android-specific smoke tests
- `app_ios_test.go` - iOS-specific smoke tests

Smoke tests verify:
1. App and window creation
2. Widget instantiation
3. Basic widget tree construction
4. Canvas availability

## Continuous Integration

### GitHub Actions

Automated testing runs on:
- **Windows**: `windows-latest` runner - See `.github/workflows/test.yml`
- **macOS**: `macos-latest` runner - See `.github/workflows/test.yml`
- **Android**: `macos-latest` with Android SDK - See `.github/workflows/test-android.yml`
- **iOS**: `macos-latest` with Xcode - See `.github/workflows/test-ios.yml`

### Android CI Details

The Android CI workflow (`.github/workflows/test-android.yml`):
1. Installs JDK 17 and Android SDK via `android-actions/setup-android`
2. Installs Android NDK 25.2.9519653
3. Installs gomobile and initializes it
4. Builds a verification AAR using `gomobile bind`

**Note**: Full test execution with `gomobile test` requires a running Android emulator, which significantly increases CI time. The current workflow verifies the build succeeds. For comprehensive testing, use local emulators or a dedicated mobile CI service.

### iOS CI Details

The iOS CI workflow (`.github/workflows/test-ios.yml`):
1. Uses the default Xcode installation on macOS runners
2. Installs gomobile and initializes it
3. Builds a verification XCFramework using `gomobile bind`

**Note**: Full test execution with `gomobile test` requires iOS Simulator, which significantly increases CI time. The current workflow verifies the build succeeds. For comprehensive testing, use local simulators or a dedicated mobile CI service.

### Mobile Emulator/Simulator Setup (Local Development)

#### Android Emulator Setup

1. **Install Android Studio** or standalone Android SDK Command-Line Tools
2. **Install required components**:
   ```bash
   sdkmanager "platform-tools" "platforms;android-34" "system-images;android-34;google_apis;arm64-v8a"
   ```
3. **Create an AVD (Android Virtual Device)**:
   ```bash
   avdmanager create avd -n test_device -k "system-images;android-34;google_apis;arm64-v8a"
   ```
4. **Start the emulator**:
   ```bash
   emulator -avd test_device -no-window &
   adb wait-for-device
   ```
5. **Run tests**:
   ```bash
   gomobile test -target=android ./...
   ```

#### iOS Simulator Setup (macOS only)

1. **Install Xcode** from the Mac App Store
2. **Install Xcode Command Line Tools**:
   ```bash
   xcode-select --install
   ```
3. **List available simulators**:
   ```bash
   xcrun simctl list devices
   ```
4. **Boot a simulator**:
   ```bash
   xcrun simctl boot "iPhone 15 Pro"
   ```
5. **Run tests**:
   ```bash
   gomobile test -target=ios ./...
   ```

### Troubleshooting Mobile Testing

#### Android Issues

- **"NDK not found"**: Ensure `ANDROID_NDK_HOME` is set to your NDK installation path
- **"adb: device not found"**: Run `adb devices` to verify emulator is connected
- **Build fails with "API level not found"**: Install the required platform with `sdkmanager`

#### iOS Issues

- **"xcode-select: error: no developer tools"**: Install Xcode Command Line Tools
- **"Simulator not found"**: Verify simulator name with `xcrun simctl list devices`
- **Code signing errors**: Use simulator target (`-target=ios/simulator`) to avoid signing requirements

## Known Limitations

### Graphics Context Requirements

Ebitengine tests may require a graphics context, which can be problematic in headless CI environments. Current workarounds:
- Use virtual display servers (Xvfb) if needed
- Run tests in GUI-enabled CI runners
- Mock graphics initialization for unit tests where possible

### Test Coverage

Test coverage metrics may not be accurate across platforms due to:
- Build tag exclusions
- Platform-specific code paths
- CGO dependency requirements

### Manual Testing

Visual and integration testing still requires manual validation:
- Widget rendering appearance
- User interaction behavior
- Performance characteristics
- Platform-specific UI conventions

## Adding New Tests

When adding tests:

1. **Add platform build tags** to all test files:
   ```go
   //go:build windows || darwin || android || ios
   ```

2. **Enable CGO** when running tests:
   ```go
   // Use t.Skip() to skip tests that require specific platform features
   if runtime.GOOS != "windows" {
       t.Skip("Test requires Windows")
   }
   ```

3. **Document platform requirements** in test godoc:
   ```go
   // TestFooBar verifies feature X.
   // Requires: Windows/macOS for graphics context initialization.
   ```

4. **Consider platform-specific test files** for platform-unique behavior:
   - Use `//go:build windows` for Windows-only tests
   - Use `//go:build darwin` for macOS-only tests
   - Etc.

## Validation Checklist

Before marking a platform as "fully tested":

- [x] Unit tests pass with `-race` flag
- [x] `go vet ./...` reports no issues
- [x] Smoke tests verify basic app creation
- [x] Widget creation and tree building works
- [x] Event dispatch functions correctly
- [x] Layout resolution produces expected results
- [x] Resource loading (fonts, images) succeeds
- [x] Theme propagation works
- [x] Manual visual testing confirms rendering quality

## Resources

- [Ebitengine Testing Guide](https://ebitengine.org/en/documents/testing.html)
- [gomobile Documentation](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile)
- [GitHub Actions for Go](https://github.com/actions/setup-go)
