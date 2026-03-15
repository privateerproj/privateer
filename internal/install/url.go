package install

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const defaultHTTPTimeout = 60 * time.Second

// FromURL downloads the plugin binary from downloadURL and writes it to destPath (full path to the output file).
// If the URL ends with .tar.gz or .zip, the archive is extracted and the binary inside (matching binaryName) is written to destPath.
// binaryName is the expected executable name (no path). The parent directory must already exist.
func FromURL(downloadURL, destPath, binaryName string) error {
	client := &http.Client{Timeout: defaultHTTPTimeout}
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "pvtr/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", downloadURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d for %s", resp.StatusCode, downloadURL)
	}

	urlLower := strings.ToLower(downloadURL)
	switch {
	case strings.HasSuffix(urlLower, ".zip"):
		return extractZip(resp.Body, destPath, binaryName)
	case strings.HasSuffix(urlLower, ".tar.gz"):
		return extractTarGz(resp.Body, destPath, binaryName)
	default:
		return writeRaw(resp.Body, destPath)
	}
}

func writeRaw(r io.Reader, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", destPath, err)
	}
	defer out.Close()
	written, err := out.ReadFrom(r)
	if err != nil {
		_ = os.Remove(destPath)
		return fmt.Errorf("write %s: %w", destPath, err)
	}
	if written == 0 {
		_ = os.Remove(destPath)
		return fmt.Errorf("download returned empty body")
	}
	_ = out.Chmod(0755)
	return nil
}

func extractZip(r io.Reader, destPath, binaryName string) error {
	// zip.Reader requires ReaderAt, so read into temp file or buffer
	tmpFile, err := os.CreateTemp("", "pvtr-install-*.zip")
	if err != nil {
		return fmt.Errorf("temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	if _, err := io.Copy(tmpFile, r); err != nil {
		return fmt.Errorf("read zip: %w", err)
	}
	info, err := tmpFile.Stat()
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(tmpFile, info.Size())
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	var match *zip.File
	for _, f := range zr.File {
		if f.Mode().IsDir() {
			continue
		}
		base := filepath.Base(f.Name)
		if base == binaryName || base == binaryName+".exe" || (len(base) > 0 && strings.TrimSuffix(base, ".exe") == binaryName) {
			match = f
			break
		}
	}
	if match == nil {
		if len(zr.File) == 1 && !zr.File[0].Mode().IsDir() {
			match = zr.File[0]
		} else {
			return fmt.Errorf("no binary matching %q in zip", binaryName)
		}
	}
	rc, err := match.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", destPath, err)
	}
	defer out.Close()
	if _, err := io.Copy(out, rc); err != nil {
		_ = os.Remove(destPath)
		return err
	}
	_ = out.Chmod(0755)
	return nil
}

func extractTarGz(r io.Reader, destPath, binaryName string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	var found bool
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar: %w", err)
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		base := path.Base(h.Name)
		if base == binaryName || base == binaryName+".exe" || (len(base) > 0 && strings.TrimSuffix(base, ".exe") == binaryName) {
			out, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("create %s: %w", destPath, err)
			}
			defer out.Close()
			if _, err := io.Copy(out, tr); err != nil {
				_ = os.Remove(destPath)
				return err
			}
			_ = out.Chmod(0755)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("no binary matching %q in tar.gz", binaryName)
	}
	return nil
}
