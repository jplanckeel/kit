package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Function to decompress based on file suffix
func DecompressFile(filename string) error {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Determine compression type based on suffix
	if strings.HasSuffix(filename, ".zip") {
		// Open a zip reader
		zipReader, err := zip.OpenReader(filename)
		if err != nil {
			return err
		}
		defer zipReader.Close()

		// Iterate through files in the zip archive
		for _, file := range zipReader.File {
			// Open each file in the archive
			zipFile, err := file.Open()
			if err != nil {
				return err
			}
			defer zipFile.Close()

			// Create a new file to write decompressed content
			outFile, err := os.Create(filepath.Join(filepath.Dir(filename), file.Name))
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Copy decompressed content to the new file
			if _, err = io.Copy(outFile, zipFile); err != nil {
				return err
			}
		}
	} else if strings.HasSuffix(filename, ".tar") || strings.HasSuffix(filename, ".gz") || strings.HasSuffix(filename, ".tgz") {
		// Open the tar file
		tarFile, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer tarFile.Close()

		// Set up reader for the tar file
		var reader io.Reader = tarFile
		if strings.HasSuffix(filename, ".gz") || strings.HasSuffix(filename, ".tgz") {
			gzReader, err := gzip.NewReader(tarFile)
			if err != nil {
				return err
			}
			defer gzReader.Close()
			reader = gzReader
		}

		tarReader := tar.NewReader(reader)

		// Iterate through files in the tar archive
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			// Create a new file to write decompressed content
			outFile, err := os.Create(filepath.Join(filepath.Dir(filename), header.Name))
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Copy decompressed content to the new file
			if _, err = io.Copy(outFile, tarReader); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("Unsupported file format for decompression")
	}

	fmt.Println("Decompression completed successfully")
	return nil
}
