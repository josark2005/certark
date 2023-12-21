package certark

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/josark2005/certark/ark"
)

var Confdir string
var Tasks = map[string]TaskProfile{}
var TaskNotValid = map[string]TaskProfile{}
var AcmeUsers = map[string]AcmeUserProfile{}
var AcmeUsersNotValid = map[string]AcmeUserProfile{}
var States = map[string]StateProfile{}

// standalone mode entry
func Standalone(confDir string) {
	ark.Debug().Str("dir", confDir).Msg("Running in standalone mode")
	Confdir = confDir
	Load(Confdir)
	checkCertsValidation()
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
	stateDir := Confdir + "/state"

	// load tasks
	Tasks = loadTasks(taskDir)

	// load acme users
	AcmeUsers = loadAcmeUsers(acmeDir)

	// load states
	States = loadStates(stateDir)

	// check validity
	checkAcmeUserValidity()
	checkTaskValidity()
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

// load states
func loadStates(stateDir string) map[string]StateProfile {
	stateFiles := []string{}
	if !IsDir(stateDir) {
		err := errors.New("state directory not found")
		ark.Error().Err(err).Str("dir", stateDir).Msg("Failed to load states")
	} else {
		err := filepath.Walk(stateDir, func(path string, info os.FileInfo, err error) error {
			if path == stateDir || info.IsDir() {
				return nil
			}
			stateFiles = append(stateFiles, path)
			return nil
		})
		if err != nil {
			ark.Error().Err(err).Msg("Failed to load states")
		}
	}

	if len(stateFiles) == 0 {
		return map[string]StateProfile{}
	}

	s := map[string]StateProfile{}
	// read task Files
	for _, path := range stateFiles {
		state := StateProfile{}
		err := ReadFileAndParseJson(path, &state)
		if err != nil {
			ark.Warn().Err(err).Str("task", path[len(stateDir)-1:]).Msg("State skipped")
		} else {
			s[state.TaskName] = state
			ark.Debug().Str("task", state.TaskName).Msg("State loaded")

		}
	}
	return s
}

// check certs validation
func checkCertsValidation() {

}
