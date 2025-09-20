package resource

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidationResult contains the results of file validation
type ValidationResult struct {
	FolderPath string
	FileCount  int
	IsValid    bool
	Error      error
}

// ValidateFilesExist validates that all required game data folders and files exist
func ValidateFilesExist() error {
	// Validate core folders
	coreFolders := []string{BASE_FOLDER, COMMON_FOLDER, COMMON_BIN_FOLDER, COMMON_DOOR_FOLDER, PL_FOLDER}
	for _, folder := range coreFolders {
		if err := validateFolder(folder, "core"); err != nil {
			return err
		}
	}

	// Validate region-specific folders
	regionFolders := map[string]string{
		"RDT_FOLDER":         RDT_FOLDER,
		"COMMON_DATA_FOLDER": COMMON_DATA_FOLDER,
	}
	for key, folder := range regionFolders {
		if err := validateFolder(folder, "region-specific"); err != nil {
			return fmt.Errorf("region-specific folder validation failed for %s (%s): %w", key, folder, err)
		}
	}

	// Validate critical files
	criticalFiles := []string{
		ROOMCUT_FILE,
		ITEMDATA_FILE,
		ESPDATA1_FILE,
		ESPDATA2_FILE,
		LEON_MODEL_FILE,
		CORE_SPRITE_FILE,
		INVENTORY_FILE,
		MENU_IMAGE_FILE,
		MENU_TEXT_FILE,
		ITEMALL_FILE,
		SAVE_SCREEN_FILE,
	}

	for _, file := range criticalFiles {
		if err := validateFile(file); err != nil {
			return fmt.Errorf("critical file validation failed: %w", err)
		}
	}

	fmt.Println("[SUCCESS] All game data validation passed successfully")
	return nil
}

// validateFolder checks if a folder exists and is readable
func validateFolder(folderPath, folderType string) error {
	folderExists, err := PathExists(folderPath)
	if !folderExists {
		return fmt.Errorf("missing %s folder: %s: %w", folderType, folderPath, err)
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read %s folder %s: %w", folderType, folderPath, err)
	}

	fmt.Printf("[%s] Folder '%s' validated (%d files)\n", folderType, folderPath, len(files))
	return nil
}

// validateFile checks if a critical file exists
func validateFile(filePath string) error {
	fileExists, err := PathExists(filePath)
	if !fileExists {
		return fmt.Errorf("missing critical file: %s: %w", filePath, err)
	}

	// Check if file is readable by getting file info
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to access file %s: %w", filePath, err)
	}

	if info.Size() == 0 {
		return fmt.Errorf("critical file is empty: %s", filePath)
	}

	fmt.Printf("[FILE] Critical file '%s' validated (%d bytes)\n", filepath.Base(filePath), info.Size())
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
