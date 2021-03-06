package gogo

import (
	"log"
	"os"
	"path"
)

var (
	// FindModeConfigFile returns config file for specified run mode.
	// You could custom your own run mode config file by overwriting.
	FindModeConfigFile = func(runMode, srcPath string) string {
		// adjust srcPath
		srcPath = path.Clean(srcPath)

		filename := "application.json"
		switch RunMode(runMode) {
		case Development:
			// try application.development.json
			filename = "application.development.json"

		case Test:
			// try application.test.json
			filename = "application.test.json"

		case Production:
			// skip

		}

		file := path.Join(srcPath, "config", filename)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			file = path.Join(srcPath, "config", "application.json")
		}

		return file
	}
)

// New creates application server with config resolved of run mode.
func New(runMode, srcPath string) *AppServer {
	// adjust app run mode
	mode := RunMode(runMode)
	if !mode.IsValid() {
		log.Fatalf("[GOGO] Invalid run mode, valid values are [%s|%s|%s]", Development, Test, Production)
	}

	// resolve config from application.json
	config, err := NewAppConfig(FindModeConfigFile(runMode, srcPath))
	if err != nil {
		log.Fatalf("[GOGO] NewAppConfig(%s): %v", FindModeConfigFile(runMode, srcPath), err)
	}
	config.SetMode(mode)

	// init default logger
	section := config.Section()
	logger := NewAppLogger(section.Logger.Output, runMode)
	logger.SetLevelByName(section.Logger.LevelName)
	logger.SetColor(!mode.IsProduction())

	logger.Printf("Initialized %s in %s mode", config.Name, config.Mode)

	return NewAppServer(mode, config, logger)
}

// NewWithLogger creates application server with provided Logger
func NewWithLogger(runMode, srcPath string, logger Logger) *AppServer {
	// adjust app run mode
	mode := RunMode(runMode)
	if !mode.IsValid() {
		log.Fatalf("[GOGO] Invalid run mode, valid values are [%s|%s|%s]", Development, Test, Production)
	}

	// resolve config from application.json
	config, err := NewAppConfig(FindModeConfigFile(runMode, srcPath))
	if err != nil {
		log.Fatalf("[GOGO] NewAppConfig(%s): %v", FindModeConfigFile(runMode, srcPath), err)
	}
	config.SetMode(mode)

	// overwrite logger level and colorful
	logger.SetLevelByName(config.Section().Logger.LevelName)
	logger.SetColor(!mode.IsProduction())

	logger.Printf("Initialized %s in %s mode", config.Name, config.Mode)

	return NewAppServer(mode, config, logger)
}
