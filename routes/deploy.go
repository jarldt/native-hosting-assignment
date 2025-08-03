package routes

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"static-site-hosting/services"
	"strings"
)

const (
	KB            = 1024
	MB            = 1024 * KB
	MaxUploadSize = 20 * MB // 20 MB
)

// DeployHandler handles the deployment of sites via ZIP file uploads
func DeployHandler(db services.DatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only allow POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Limit request body size
		r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

		// Parse multipart form
		err := r.ParseMultipartForm(MaxUploadSize)
		if err != nil {
			// Don't leak internal errors
			fmt.Printf("Error parsing multipart form: %v\n", err)
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		siteName := r.FormValue("siteName")
		if siteName == "" || strings.Contains(siteName, "..") || strings.Contains(siteName, "/") {
			http.Error(w, "Invalid site name", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("zipFile")
		if err != nil {
			fmt.Printf("Error retrieving zip file from body: %v\n", err)
			http.Error(w, "Missing zipFile", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save zip to a temp file
		tmpZip := filepath.Join(os.TempDir(), header.Filename)
		out, err := os.Create(tmpZip)
		if err != nil {
			fmt.Printf("Error creating temp file %s: %v\n", tmpZip, err)
			http.Error(w, "Could not save uploaded zip.", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpZip)
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Printf("Error writing to temp file %s: %v\n", tmpZip, err)
			http.Error(w, "Error writing file.", http.StatusInternalServerError)
			return
		}

		// Extract zip to deployments/siteName
		targetDir := filepath.Join("deployments", siteName)
		os.RemoveAll(targetDir) // Clean previous deployment
		err = Unzip(tmpZip, targetDir)
		if err != nil {
			fmt.Printf("Error extracting zip %s to %s: %v\n", tmpZip, targetDir, err)
			http.Error(w, "Failed to extract zip.", http.StatusInternalServerError)
			return
		}

		err = db.InsertSite(siteName)
		if err != nil {
			fmt.Printf("Error saving site %s to database: %v\n", siteName, err)
			http.Error(w, "Error occurred when saving site.", http.StatusInternalServerError)
			return
		}

		// Log the deployment
		ip := r.RemoteAddr
		if ipParts := strings.Split(ip, ":"); len(ipParts) > 0 {
			ip = ipParts[0] // Strip the port
		}

		ua := r.UserAgent()

		err = db.InsertDeploymentLog(siteName, ip, ua)
		if err != nil {
			fmt.Printf("Warning: Failed to log deployment for %s: %v\n", siteName, err)
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Site '%s' deployed successfully", siteName)
	}
}

// Unzip handles unzipping ZIP archives
func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		cleanName := filepath.Clean(f.Name)

		// Skip annoying junk files...
		if strings.HasPrefix(cleanName, "__MACOSX") || strings.Contains(cleanName, "/__MACOSX") {
			continue
		}

		base := filepath.Base(cleanName)
		if strings.HasPrefix(base, ".") {
			continue
		}

		// Basic zip slip attack prevention
		if strings.Contains(cleanName, "..") {
			continue
		}

		fp := filepath.Join(dest, cleanName)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fp, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		err := os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		if err != nil {
			return err
		}

		dstFile, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		_, err = io.Copy(dstFile, rc)

		dstFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	// After unzipping, check if the destination directory contains a single subdirectory.
	// If it does, move the contents of that subdirectory up to the destination directory.
	// Unfortunately, this is necessary due to some GUI archive utilities.
	entries, err := os.ReadDir(dest)
	if err != nil {
		return fmt.Errorf("error reading destination directory %s after unzip: %w", dest, err)
	}

	if len(entries) == 1 && entries[0].IsDir() {
		innerDirPath := filepath.Join(dest, entries[0].Name())

		innerEntries, err := os.ReadDir(innerDirPath)
		if err != nil {
			return fmt.Errorf("error reading inner directory %s after unzip: %w", innerDirPath, err)
		}

		// Move each item from the inner directory to dest
		for _, entry := range innerEntries {
			src := filepath.Join(innerDirPath, entry.Name())
			dst := filepath.Join(dest, entry.Name())

			err := os.Rename(src, dst)
			if err != nil {
				return fmt.Errorf("error moving %s to %s: %w", src, dst, err)
			}
		}

		// Remove the now-empty directory
		err = os.Remove(innerDirPath)
		if err != nil {
			return fmt.Errorf("error removing inner directory %s: %w", innerDirPath, err)
		}
	}

	return nil
}
