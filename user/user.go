package user

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	userPath = "users/"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func New(name, password string) *User {
	sha := sha256.New()
	sha.Write([]byte(password))
	data := sha.Sum(nil)
	encodedPW := base64.StdEncoding.EncodeToString(data)
	return &User{
		Name:     name,
		Password: encodedPW,
	}
}

func Load(name string) (*User, error) {
	file, err := os.Open(userPath + name + ".json")
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	user := New("", "")
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (user *User) Equals(other *User) bool {
	if other == nil {
		return false
	}
	return user.Name == other.Name && user.Password == other.Password
}

func (user *User) Remove() error {
	return os.Remove(userPath + user.Name + ".json")
}

func (user *User) Store() error {
	tempFile, err := os.Open(userPath + user.Name + ".json")
	if err == nil {
		_ = tempFile.Close()
		return fmt.Errorf("user %s already exists. Try to call user.Remove() first", user.Name)
	}

	file, err := os.Create(userPath + user.Name + ".json")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
