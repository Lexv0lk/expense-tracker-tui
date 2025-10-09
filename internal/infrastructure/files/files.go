package files

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const (
	saveFileName    = "expenses.json"
	defaultSavePath = "%AppData%"
	appName         = "expense-tracker"
)

func SaveToFile[T ~[]E, E any](data T) error {
	saveFile, err := OpenOrCreateFile(getSaveDir(), saveFileName)

	if err != nil {
		return err
	}

	return saveToFile(saveFile, data)
}

func GetFromFile[T ~[]E, E any]() (T, error) {
	saveFile, err := OpenOrCreateFile(getSaveDir(), saveFileName)

	if err != nil {
		return nil, err
	}

	return getFromFile[T](saveFile)
}

func saveToFile[T ~[]E, E any](file io.WriteCloser, data T) error {
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func getFromFile[T ~[]E, E any](file io.ReadCloser) (T, error) {
	defer file.Close()

	var data T
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&data)

	if err == io.EOF {
		err = nil

		if data == nil {
			data = make(T, 0)
		}
	}

	return data, err
}

func OpenOrCreateFile(dir string, filename string) (*os.File, error) {
	if err := ensureDirExists(dir); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(filepath.Join(dir, filename), os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func ensureDirExists(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}

	return err
}

func getSaveDir() string {
	path, err := os.UserConfigDir()

	if err != nil {
		path = defaultSavePath
	}

	return filepath.Join(path, appName)
}
