package certark

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/jokin1999/certark/ark"
)

var Confdir string
var Tasks = map[string]TaskProfile{}
var AcmeUsers = map[string]AcmeUserProfile{}

// standalone mode entry
func Standalone(confDir string) {
	ark.Debug().Str("dir", confDir).Msg("Running in standalone mode")
	Confdir = confDir
	Load(Confdir)
}

// server mode entry
func Server() {

}

// client mode entry
func Client() {

}

// reload configurations
func Reload() {
	Load(Confdir)
}

// load configurations
func Load(dir string) {
	taskDir := Confdir + "/task"
	acmeDir := Confdir + "/user"

	// load tasks
	Tasks = loadTasks(taskDir)

	// load acme users
	AcmeUsers = loadAcmeUsers(acmeDir)
}

// load tasks
func loadTasks(taskDir string) map[string]TaskProfile {
	taskFiles := []string{}
	if !IsDir(taskDir) {
		err := errors.New("task directory not found")
		ark.Error().Err(err).Str("dir", taskDir).Msg("Failed ot load tasks")
	} else {
		err := filepath.Walk(taskDir, func(path string, info os.FileInfo, err error) error {
			if path == taskDir || info.IsDir() {
				return nil
			}
			taskFiles = append(taskFiles, path)
			return nil
		})
		if err != nil {
			ark.Error().Err(err).Msg("Failed to load tasks")
		}
	}

	if len(taskFiles) == 0 {
		return map[string]TaskProfile{}
	}

	tasks := map[string]TaskProfile{}
	// read task Files
	for _, path := range taskFiles {
		task := TaskProfile{}
		err := ReadFileAndParseJson(path, &task)
		if err != nil {
			ark.Warn().Err(err).Str("task", path[len(taskDir)-1:]).Msg("Task skipped")
		} else {
			if task.Enabled {
				tasks[task.TaskName] = task
				ark.Debug().Str("task", task.TaskName).Msg("Task loaded")
			} else {
				ark.Debug().Str("task", task.TaskName).Str("reason", "disabled").Msg("Task skipped")
			}
		}
	}
	return tasks
}

// load acme users
func loadAcmeUsers(acmeDir string) map[string]AcmeUserProfile {
	acmeFiles := []string{}
	if !IsDir(acmeDir) {
		err := errors.New("acme user directory not found")
		ark.Error().Err(err).Str("dir", acmeDir).Msg("Failed ot load acme users")
	} else {
		err := filepath.Walk(acmeDir, func(path string, info os.FileInfo, err error) error {
			if path == acmeDir || info.IsDir() {
				return nil
			}
			acmeFiles = append(acmeFiles, path)
			return nil
		})
		if err != nil {
			ark.Error().Err(err).Msg("Failed to load acme users")
		}
	}

	if len(acmeFiles) == 0 {
		return map[string]AcmeUserProfile{}
	}

	aus := map[string]AcmeUserProfile{}
	// read task Files
	for _, path := range acmeFiles {
		acmeuser := AcmeUserProfile{}
		err := ReadFileAndParseJson(path, &acmeuser)
		if err != nil {
			ark.Warn().Err(err).Str("user", path[len(acmeDir)-1:]).Msg("Acme user skipped")
		} else {
			if acmeuser.Enabled {
				aus[acmeuser.Email] = acmeuser
				ark.Debug().Str("user", acmeuser.Email).Msg("Acme user loaded")
			} else {
				ark.Debug().Str("user", acmeuser.Email).Str("reason", "disabled").Msg("Acme user skipped")
			}
		}
	}
	return aus
}
