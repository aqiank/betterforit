package net

import (
        "fmt"
        "io"
        "log"
        "net/http"
        "os"
        "path"
)

func Download(destDir, src, name string) error {
        var file *os.File
	var resp *http.Response
	var err error
	var n int64

	// setup a connection to the file
	if resp, err = http.Get(src); err != nil {
		return fmt.Errorf("DownloadAs: %v", err)
	}

	// prepare directory
	if err = os.MkdirAll(destDir, os.ModeDir | 0700); err != nil {
		return fmt.Errorf("DownloadAs: %v", err)
	}

        // use different name if available
        if name == "" {
                name = path.Base(src)
        }

        // check if file already exists
	if _, err = os.Stat(destDir + "/" + name); err == nil {
		log.Println(name, "already exists")
		return nil
	}

	// create destination file to download to
	if file, err = os.Create(destDir + "/" + name); err != nil {
		return fmt.Errorf("DownloadAs: %v", err)
	}
	defer file.Close()

	// start download
	if n, err = io.Copy(file, resp.Body); err != nil {
                return fmt.Errorf("DownloadAs: %v", err)
	}

	log.Printf("downloaded %s (%dKB)\n", name, n / 1024)
        return nil
}

func SaveText(destDir, text, name string) error {
        var file *os.File
        var err error

	// prepare directory
	if err = os.MkdirAll(destDir, os.ModeDir | 0700); err != nil {
		return fmt.Errorf("SaveText: %v", err)
	}

        // use different name if available
        name = path.Base(name)

        // check if file already exists
	if _, err = os.Stat(destDir + "/" + name); err == nil {
		log.Println(name, "already exists")
		return nil
	}

	// create destination file to download to
	if file, err = os.Create(destDir + "/" + name); err != nil {
		return fmt.Errorf("SaveText: %v", err)
	}
	defer file.Close()

        if _, err = io.WriteString(file, text); err != nil {
                return fmt.Errorf("SaveText: %v", err)
        }

        return nil
}
