package timetracker

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type DataPoint struct {
	Filename  string
	Project   string
	Timestamp float64
}

func _find_git_config_file(file_path string, point *DataPoint) error {
	file_path = filepath.Clean(file_path)
	stat, err := os.Stat(file_path)
	if err != nil {
		return err
	}
	if stat.Mode().IsRegular() {
		file_path = filepath.Dir(file_path)
	}
	stat, err = os.Stat(filepath.Join(file_path, ".git", "config"))
	if err == nil && stat.Mode().IsRegular() {
		point.Project = filepath.Base(file_path)
		return nil
	}

	dirname, base_name := filepath.Split(file_path)
	if base_name == "" {
		return errors.New("No git folder found")
	}

	return _find_git_config_file(dirname, point)
}

func WriteDataPoint(filename string, timestamp float64) error {
	log_file := GetOutputPath()
	fd, err := os.OpenFile(log_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	point := &DataPoint{Filename: filename, Timestamp: timestamp}
	_find_git_config_file(filename, point)

	serialized, err := json.Marshal(point)
	if err != nil {
		return err
	}

	serialized = append(serialized, []byte("\n")...)
	_, err = fd.Write(serialized)
	return err
}

func GetOutputPath() string {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == "VELOTRACKER_LOG" {
			return pair[1]
		}
	}

	return os.ExpandEnv("$HOME/.timetracker.log")
}
