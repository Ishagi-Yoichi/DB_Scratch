
package storage

import (
	"os"
	"path/filepath"
)

func SaveData1(path string, data []byte) error {
	fp, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0664,
	)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(data)
	if err != nil {
		return err
	}

	return fp.Sync()
}

// SaveData2 writes new data to a temporary file
// and then replaces the original file.
func SaveData2(path string, data []byte) error {
	// 1. Get the directory containing the target file.
	dir := filepath.Dir(path)

	// 2. Create a temporary file in the same directory.
	fp, err := os.CreateTemp(dir, "database-*.tmp")
	if err != nil {
		return err
	}

	tmp := fp.Name()
	renamed := false

	// 3. Clean up the temporary file if replacement fails.
	defer func() {
		fp.Close()
		if !renamed {
			os.Remove(tmp)
		}
	}()

	// 4. Write new data into the temporary file.
	_, err = fp.Write(data)
	if err != nil {
		return err
	}

	// 5. Synchronize temporary file contents.
	err = fp.Sync()
	if err != nil {
		return err
	}

	// 6. Close before renaming.
	// This is important for Windows compatibility.
	err = fp.Close()
	if err != nil {
		return err
	}

	// 7. Replace the original file.
	err = os.Rename(tmp, path)
	if err != nil {
		return err
	}

	renamed = true
	return nil
}
