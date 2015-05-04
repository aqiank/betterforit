package net

import (
        "fmt"
        "io"
        "log"
        "net/http"
        "os"
        "path"
)

func Download(destDir, src string) error {
	var file *os.File
	var resp *http.Response
	var err error
	var n int64

	// Setup a connection to the file
	if resp, err = http.Get(src); err != nil {
		return fmt.Errorf("Download: %v", err)
	}

	// Prepare directory
	if err = os.MkdirAll(destDir, os.ModeDir | 0700); err != nil {
		return fmt.Errorf("Download: %v", err)
	}

	// Check if file already exists
	name := path.Base(src)
	if _, err = os.Stat(destDir + "/" + name); err == nil {
		log.Println(name, "already exists")
		return nil
	}

	// Create destination file to download to
	if file, err = os.Create(destDir + "/" + name); err != nil {
		return fmt.Errorf("Download: %v", err)
	}
	defer file.Close()

	// Start download
	if n, err = io.Copy(file, resp.Body); err != nil {
                return fmt.Errorf("Download: %v", err)
	}

	log.Printf("Downloaded %s (%dKB)\n", name, n / 1024)
        return nil
}

func DownloadAs(destDir, src, as string) error {
	var file *os.File
	var resp *http.Response
	var err error
	var n int64

	// Setup a connection to the file
	if resp, err = http.Get(src); err != nil {
		return fmt.Errorf("Download: %v", err)
	}

	// Prepare directory
	if err = os.MkdirAll(destDir, os.ModeDir | 0700); err != nil {
		return fmt.Errorf("Download: %v", err)
	}

	// Check if file already exists
	if _, err = os.Stat(destDir + "/" + as); err == nil {
		log.Println(as, "already exists")
		return nil
	}

	// Create destination file to download to
	if file, err = os.Create(destDir + "/" + as); err != nil {
		return fmt.Errorf("Download: %v", err)
	}
	defer file.Close()

	// Start download
	if n, err = io.Copy(file, resp.Body); err != nil {
                return fmt.Errorf("Download: %v", err)
	}

	log.Printf("Downloaded %s (%dKB)\n", as, n / 1024)
        return nil
}
