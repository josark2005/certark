package certark

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/josark2005/certark/acme"
)

// get acme user filepath
func GetAcmeUserFilepath(name string) string {
	return AcmeUserDir + "/" + name
}

// check if acme user exists
func CheckAcmeUserExists(name string) bool {
	res := FileOrDirExists(GetAcmeUserFilepath(name))
	return res
}

// get acme user
func GetAcmeUser(name string) (AcmeUserProfile, error) {
	profilePath := GetAcmeUserFilepath(name)

	if !CheckAcmeUserExists(name) {
		err := errors.New("acme user profile does not exist")
		return AcmeUserProfile{}, err
	}

	profile := AcmeUserProfile{}
	err := ReadFileAndParseJson(profilePath, &profile)
	if err != nil {
		return AcmeUserProfile{}, err
	}

	return profile, nil
}

// get acme user json
func GetAcmeUserJson(name string) ([]byte, error) {
	profilePath := GetAcmeUserFilepath(name)

	content, err := os.ReadFile(profilePath)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

// get acme user json pretty
func GetAcmeUserJsonPretty(name string) (string, error) {
	profileContent, err := GetAcmeUserJson(name)
	if err != nil {
		return "", err
	}

	var jsonBuff bytes.Buffer
	if err = json.Indent(&jsonBuff, profileContent, "", ""); err != nil {
		return "", err
	}

	return jsonBuff.String(), nil
}

// list acme users
func ListAcmeUsers() ([]string, error) {
	users := []string{}
	err := filepath.Walk(AcmeUserDir, func(path string, info os.FileInfo, err error) error {
		// skip dir itself
		if path == AcmeUserDir {
			return nil
		}
		// skip dirs
		if info.IsDir() {
			return nil
		}
		users = append(users, path[len(AcmeUserDir)+1:])
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return users, nil
}

// add acme user
func AddAcmeUser(name string, email string) error {
	if CheckAcmeUserExists(name) {
		err := errors.New("user existed")
		return err
	}

	profileFilepath := GetAcmeUserFilepath(name)

	profile := DefaultAcmeUserProfile
	profile.Email = email

	err := WriteStructToFile(profile, profileFilepath)
	if err != nil {
		return err
	}
	return nil
}

// remove acme user
func RemoveAcmeUser(name string) error {
	if !CheckAcmeUserExists(name) {
		err := errors.New("user does not exist")
		return err
	}

	// remove profile
	err := os.Remove(GetAcmeUserFilepath(name))
	if err != nil {
		return err
	}
	return nil
}

// set acme user private key in file
func SetAcmeUserPrivateKeyInFile(name string, privateKeyFilepath string) error {
	if !CheckAcmeUserExists(name) {
		err := errors.New("user does not exist")
		return err
	}

	// read private key
	privatekey, err := os.ReadFile(privateKeyFilepath)
	if err != nil {
		return err
	}

	profile, err := GetAcmeUser(name)
	if err != nil {
		return err
	}

	// update private key
	profile.PrivateKey = string(privatekey)
	err = WriteStructToFile(profile, GetAcmeUserFilepath(name))
	return err
}

// register acme user
func RegisterAcmeUser(name string) error {
	if !CheckAcmeUserExists(name) {
		err := errors.New("user does not exist")
		return err
	}

	profile, err := GetAcmeUser(name)
	if err != nil {
		return err
	}
	privateKey := ""
	if CurrentConfig.Mode == MODE_PROD {
		privateKey = acme.RegisterAcmeUser(profile.Email, acme.MODE_PRODUCTION)
	} else {
		privateKey = acme.RegisterAcmeUser(profile.Email, acme.MODE_STAGING)
	}

	profile.PrivateKey = privateKey

	err = WriteStructToFile(profile, GetAcmeUserFilepath(name))
	return err
}
