package storage
import "os"
import "path/filepath"

func savedata1(path string, data []byte) error {
	fp, error := os.OpenFile( path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0664,)

	if err !=nil{
		return err
	}

	defer fp.Close()

	_, err = fp.Write(data)

	if err !=nil{
		return err
	}

	return fp.Sync()
}

func savedata2(path string, data []byte) error {
	dir := filepath.Dir(path)

	fp, err := os.CreateTemp(dir, "database-*.tmp")
	if err !=nil{
		return err
	}

	tmp := fp.Name()
	renamed := false

	defer func() {
		fp.Close()
		if !renamed {
			os.Remove(tmp)
		}
	}()

	_. err = fp.Write(data)
	if err != nil{
		return err
	}

	err = fp.Close()
	if err !=nil{
		return err
	}

	err = os.Rename(tmp, path)
	if err !=nil{
		return err
	}

	renamed = true
	return nil
}
