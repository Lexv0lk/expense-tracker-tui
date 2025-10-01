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
	saveFile, err := openOrCreateSaveFile()

	if err != nil {
		return err
	}

	return saveToFile(saveFile, data)
}

func GetFromFile[T ~[]E, E any]() (T, error) {
	saveFile, err := openOrCreateSaveFile()

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

func ensureSaveDirExists() error {
	_, err := os.Stat(getSaveDir())

	if os.IsNotExist(err) {
		return os.MkdirAll(getSaveDir(), 0755)
	}

	return err
}

func openOrCreateSaveFile() (*os.File, error) {
	if err := ensureSaveDirExists(); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(getSavePath(), os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func getSavePath() string {
	path, err := os.UserConfigDir()

	if err != nil {
		path = defaultSavePath
	}

	return filepath.Join(path, appName, saveFileName)
}

func getSaveDir() string {
	path, err := os.UserConfigDir()

	if err != nil {
		path = defaultSavePath
	}

	return filepath.Join(path, appName)
}
