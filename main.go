package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	defaultChunkSize = 400 // Default chunk size in MB
	chunkSuffix      = ".chunk."
)

// Define file size units
const (
	B  = 1
	KB = 1024
	MB = 1024 * 1024
	GB = 1024 * 1024 * 1024
)

var (
	// Global parameters
	filePath    string
	chunkSize   int64
	sizeUnit    string
	clearChunks bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "chunkfile",
	Short: "A tool for splitting and merging large files",
	Long: `chunkfile is a command-line tool for splitting large files into smaller chunks
and merging them back together. It's mainly used for handling large files
that need to be transferred over network or stored in limited space.`,
}

// Parse and validate size unit
func parseUnit(unit string) (int64, error) {
	switch strings.ToUpper(unit) {
	case "B":
		return B, nil
	case "KB":
		return KB, nil
	case "MB":
		return MB, nil
	case "GB":
		return GB, nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s, supported units are: B, KB, MB, GB", unit)
	}
}

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split a file into smaller chunks",
	Long: `Split a large file into multiple smaller chunks.
You can specify the size and unit for each chunk, default is 400MB.
Supported units are: B, KB, MB, GB (case-insensitive).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if filePath == "" {
			return fmt.Errorf("please specify the file path to split")
		}

		// Parse unit
		unitMultiplier, err := parseUnit(sizeUnit)
		if err != nil {
			return err
		}

		// Calculate final chunk size (in bytes)
		finalChunkSize := chunkSize * unitMultiplier

		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}
		return splitFile(absPath, finalChunkSize)
	},
}

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge previously split chunks",
	Long: `Merge previously split file chunks back into a complete file.
You can choose to delete the chunk files after successful merge.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if filePath == "" {
			return fmt.Errorf("please specify the chunk file prefix")
		}
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}
		return mergeFiles(absPath, clearChunks)
	},
}

func init() {
	// Add commands
	rootCmd.AddCommand(splitCmd)
	rootCmd.AddCommand(mergeCmd)

	// Split command flags
	splitCmd.Flags().StringVarP(&filePath, "path", "p", "", "path to the file to split")
	splitCmd.Flags().Int64VarP(&chunkSize, "size", "s", defaultChunkSize, "size of each chunk")
	splitCmd.Flags().StringVarP(&sizeUnit, "unit", "u", "MB", "size unit (B, KB, MB, GB)")

	// Merge command flags
	mergeCmd.Flags().StringVarP(&filePath, "path", "p", "", "prefix of chunk files")
	mergeCmd.Flags().BoolVarP(&clearChunks, "clear", "c", false, "delete chunk files after successful merge")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Calculate required number of digits
func calculateDigits(n int64) int {
	if n <= 0 {
		return 1
	}
	return int(math.Floor(math.Log10(float64(n)))) + 1
}

// Format file size with appropriate unit
func formatSize(size int64) string {
	if size < KB {
		return fmt.Sprintf("%d B", size)
	} else if size < MB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	} else if size < GB {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}
}

// splitFile splits a file into multiple chunks
func splitFile(filePath string, chunkSize int64) error {
	// Open source file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Calculate number of chunks
	totalChunks := (fileInfo.Size() + chunkSize - 1) / chunkSize // Round up

	// Calculate required digits for numbering
	digits := calculateDigits(totalChunks)
	// Create number format string, e.g.: %03d, %04d, %05d
	numFmt := fmt.Sprintf("%%0%dd", digits)

	fmt.Printf("File size: %s\n", formatSize(fileInfo.Size()))
	fmt.Printf("Chunk size: %s\n", formatSize(chunkSize))
	fmt.Printf("Total chunks: %d\n", totalChunks)

	// Create buffer
	buffer := make([]byte, 1024*1024) // 1MB buffer

	// Process chunks
	for i := int64(0); i < totalChunks; i++ {
		// Create chunk file with dynamic digit numbering
		chunkFileName := fmt.Sprintf("%s%s"+numFmt, filePath, chunkSuffix, i+1)
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			return fmt.Errorf("failed to create chunk file: %v", err)
		}
		defer chunkFile.Close()

		// Calculate bytes to write for current chunk
		bytesLeft := chunkSize
		if i == totalChunks-1 {
			// Last chunk might be smaller than chunkSize
			bytesLeft = fileInfo.Size() - i*chunkSize
		}

		// Write data
		var totalWritten int64
		for totalWritten < bytesLeft {
			// Calculate bytes to read
			toRead := bytesLeft - totalWritten
			if toRead > int64(len(buffer)) {
				toRead = int64(len(buffer))
			}

			// Read data
			n, err := file.Read(buffer[:toRead])
			if err != nil && err != io.EOF {
				return fmt.Errorf("failed to read file: %v", err)
			}
			if n == 0 {
				break
			}

			// Write data
			written, err := chunkFile.Write(buffer[:n])
			if err != nil {
				return fmt.Errorf("failed to write chunk file: %v", err)
			}
			totalWritten += int64(written)
		}

		fmt.Printf("Created chunk file: %s (%s)\n", chunkFileName, formatSize(totalWritten))
	}

	fmt.Println("File splitting completed")
	return nil
}

// mergeFiles merges chunk files back together
func mergeFiles(chunkPrefix string, clear bool) error {
	// Get all chunk files
	dir := filepath.Dir(chunkPrefix)
	if dir == "" {
		dir = "."
	}
	base := filepath.Base(chunkPrefix)

	// Read directory entries
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Filter chunk files
	var chunkFiles []string
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), base) && strings.Contains(entry.Name(), chunkSuffix) {
			chunkFiles = append(chunkFiles, filepath.Join(dir, entry.Name()))
		}
	}

	if len(chunkFiles) == 0 {
		return fmt.Errorf("no chunk files found")
	}

	// Sort chunk files
	sort.Slice(chunkFiles, func(i, j int) bool {
		// Extract sequence numbers from filenames
		iNum := extractChunkNumber(chunkFiles[i])
		jNum := extractChunkNumber(chunkFiles[j])
		return iNum < jNum
	})

	// Determine output filename (remove suffix)
	outputFileName := strings.TrimSuffix(chunkPrefix, chunkSuffix)
	if outputFileName == chunkPrefix {
		// If prefix doesn't contain suffix, use prefix as output filename
		outputFileName = chunkPrefix
	} else {
		// If prefix contains suffix, need further processing
		parts := strings.Split(chunkPrefix, chunkSuffix)
		if len(parts) > 0 {
			outputFileName = parts[0]
		}
	}

	// Create output file
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Merge files
	buffer := make([]byte, 1024*1024) // 1MB buffer
	totalSize := int64(0)

	for i, chunkFile := range chunkFiles {
		fmt.Printf("Processing chunk file %d/%d: %s\n", i+1, len(chunkFiles), chunkFile)

		// Open chunk file
		file, err := os.Open(chunkFile)
		if err != nil {
			return fmt.Errorf("failed to open chunk file %s: %w", chunkFile, err)
		}

		// Get file size
		fileInfo, err := file.Stat()
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to get chunk file info: %w", err)
		}

		chunkSize := fileInfo.Size()
		totalSize += chunkSize

		// Copy content to output file
		for {
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				file.Close()
				return fmt.Errorf("failed to read chunk file: %w", err)
			}
			if n == 0 {
				break
			}

			_, err = outputFile.Write(buffer[:n])
			if err != nil {
				file.Close()
				return fmt.Errorf("failed to write output file: %w", err)
			}
		}

		file.Close()
	}

	fmt.Printf("Created merged file: %s (%s)\n", outputFileName, formatSize(totalSize))

	if clear {
		// Clean up chunk files
		fmt.Println("Starting cleanup of chunk files...")
		for _, chunkFile := range chunkFiles {
			err := os.Remove(chunkFile)
			if err != nil {
				fmt.Printf("Warning: failed to delete chunk file %s: %v\n", chunkFile, err)
			} else {
				fmt.Printf("Deleted chunk file: %s\n", chunkFile)
			}
		}
		fmt.Println("Chunk file cleanup completed")
	}

	return nil
}

// extractChunkNumber extracts the sequence number from a chunk filename
func extractChunkNumber(filename string) int {
	// Find suffix position
	suffixIndex := strings.LastIndex(filename, chunkSuffix)
	if suffixIndex == -1 {
		return 0
	}

	// Extract number part (now supports any number of digits)
	numStr := filename[suffixIndex+len(chunkSuffix):]
	// Remove possible file extension
	if dotIndex := strings.LastIndex(numStr, "."); dotIndex != -1 {
		numStr = numStr[:dotIndex]
	}

	// Convert to integer
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}

	return num
}
