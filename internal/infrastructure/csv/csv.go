package csv

import (
	"encoding/csv"
	"github.com/Lexv0lk/expense-tracker-tui/internal/infrastructure/files"
	"io"
	"os"
	"path/filepath"
)

const (
	saveFileName    = "expenses.csv"
	defaultSavePath = "%AppData%"
	appName         = "expense-tracker"
)

func SaveToCSV(data [][]string) error {
	file, err := files.OpenOrCreateFile(getSaveDir(), saveFileName)
	if err != nil {
		return err
	}

	return saveToCSV(file, data)
}

func GetSaveFilePath() string {
	return filepath.Join(getSaveDir(), saveFileName)
}

func saveToCSV(file io.WriteCloser, data [][]string) error {
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = ';'

	for _, record := range data {
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

func getSaveDir() string {
	path, err := os.UserConfigDir()

	if err != nil {
		path = defaultSavePath
	}

	return filepath.Join(path, appName)
}
