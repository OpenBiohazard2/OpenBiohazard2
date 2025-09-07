package game

import (
	"fmt"
	"os"
)

func ValidateFilesExist() error {
	folderList := []string{BASE_FOLDER, COMMON_FOLDER, COMMON_BIN_FOLDER, COMMON_DOOR_FOLDER, PL_FOLDER}
	for _, folderName := range folderList {
		folderExists, err := PathExists(folderName)
		if !folderExists {
			return fmt.Errorf("missing data error: unable to find folder %v: %w", folderName, err)
		}

		files, err := os.ReadDir(folderName)
		if err != nil {
			return fmt.Errorf("failed to read directory %v: %w", folderName, err)
		}
		fmt.Printf("Folder %v exists and has %v files\n", folderName, len(files))
	}

	regionSpecificFolders := map[string]string{
		"RDT_FOLDER":         RDT_FOLDER,
		"COMMON_DATA_FOLDER": COMMON_DATA_FOLDER,
	}
	for folderKey := range regionSpecificFolders {
		folderName := regionSpecificFolders[folderKey]
		folderExists, err := PathExists(folderName)
		if !folderExists {
			return fmt.Errorf(
				"missing data error: unable to find region specific folder %v. "+
					"You might need to change the value of this variable in resource.go. Current value of key %v = %v: %w",
				folderKey, folderKey, folderName, err,
			)
		}

		files, err := os.ReadDir(folderName)
		if err != nil {
			return fmt.Errorf("failed to read directory %v: %w", folderName, err)
		}
		fmt.Printf("Folder %v exists and has %v files\n", folderName, len(files))
	}

	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
