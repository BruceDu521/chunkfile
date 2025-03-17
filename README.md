# ChunkFile - File Splitting and Merging Tool

[English](README.md) | [简体中文](README_CN.md)

A command-line tool written in Go for splitting large files into smaller chunks and merging them back together. Particularly useful for uploading files larger than 1GB to cloud storage services with size limitations, and then reassembling them on another device.

## Features

- Split large files into smaller chunks of specified size
- Merge split chunks back into the original file
- Dynamically adjust chunk file naming to support any number of chunks
- Customize chunk size with various units (B, KB, MB, GB)
- Option to automatically clean up chunk files after successful merge
- Support for both relative and absolute paths
- Cross-platform compatibility (Windows, Linux, macOS)

## Installation

### Requirements

- Go 1.18 or higher
- The project uses Go 1.23 toolchain (automatically selected if available)

### Using Go Install (Recommended)

If you have Go installed on your system, you can install ChunkFile directly using:

```bash
go install github.com/BruceDu521/chunkfile/cmd/chunkfile@latest
```

Make sure your Go bin directory is in your PATH:
- Windows: `%USERPROFILE%\go\bin`
- Linux/macOS: `~/go/bin`

### Manual Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/BruceDu521/chunkfile.git
   cd chunkfile
   ```

2. Build the executable:
   ```bash
   go build -o chunkfile ./cmd/chunkfile
   ```

3. Move the executable to a directory in your PATH or use it from the current location.

## Usage

ChunkFile provides two main commands: `split` and `merge`.

### Splitting Files

```bash
chunkfile split --path <file_path> [--size <size>] [--unit <unit>]
```

Or using short flags:

```bash
chunkfile split -p <file_path> [-s <size>] [-u <unit>]
```

Parameters:
- `--path, -p`: Path to the file to split (required)
- `--size, -s`: Size of each chunk (default: 400)
- `--unit, -u`: Size unit (B, KB, MB, GB, case-insensitive, default: MB)

Example, splitting a large file into 500MB chunks:

```bash
chunkfile split --path "large_file.zip" --size 500 --unit MB
```

Or splitting into 1GB chunks:

```bash
chunkfile split -p "large_file.zip" -s 1 -u GB
```

This will generate a series of chunk files like: `large_file.zip.chunk.0001`, `large_file.zip.chunk.0002`, etc.

### Merging Files

```bash
chunkfile merge --path <chunk_file_prefix> [--clear]
```

Or using short flags:

```bash
chunkfile merge -p <chunk_file_prefix> [-c]
```

Parameters:
- `--path, -p`: Prefix of the chunk files (required)
- `--clear, -c`: Delete chunk files after successful merge (optional)

Example, merging previously split files:

```bash
chunkfile merge --path "large_file.zip"
```

Or merging and cleaning up chunk files:

```bash
chunkfile merge -p "large_file.zip" -c
```

This will find all files matching the pattern `large_file.zip.chunk.*`, sort them correctly, and merge them to create the original file `large_file.zip`.

## Notes

- When merging, the program automatically finds and sorts all matching chunk files
- Chunk files use `.chunk.XXXX` as suffix, where XXXX is a sequence number starting from 0001
- The number of digits in the chunk suffix is dynamically determined based on the total number of chunks
- Ensure you have enough disk space for storing chunk files or the merged file
- The program automatically converts relative paths to absolute paths for processing
- Compatible with Windows, Linux, and macOS

## Getting Help

For more information on available commands and options:

```bash
chunkfile --help
chunkfile split --help
chunkfile merge --help
```

## License

[MIT License](LICENSE) 