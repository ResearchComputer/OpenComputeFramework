# Installation

## Prerequisites

*   **Go**: Version 1.23 or higher.
*   **Make**: For running build commands.
*   **Git**: For version control.

## Building from Source

1.  Clone the repository:

    ```bash
    git clone https://github.com/your-org/opencomputeframework.git
    cd opencomputeframework/src
    ```

2.  Build the binary using `make`:

    ```bash
    make build
    ```

    This will create the `ocfcore` binary in the `build/` directory.

3.  (Optional) Install dependencies if needed:

    ```bash
    make build-deps
    ```

## Cross-Compilation

To build for ARM64 architecture:

```bash
make arm
```

## Debug Build

To build with debugging capabilities:

```bash
make build-debug
```

## Verify Installation

After building, you can verify the installation by checking the version:

```bash
./build/ocfcore version
```
