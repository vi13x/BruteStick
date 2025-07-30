package state

import (
	"encoding/json"
	"io/ioutil"
)

type BruteState struct {
	CurrentLength   int
	CurrentPassword string
	Stopped         bool `json:"-"`
}

func NewBruteState() *BruteState {
	return &BruteState{
		CurrentLength:   1,
		CurrentPassword: "",
		Stopped:         false,
	}
}

func LoadState(path string, key []byte) (*BruteState, error) {
	encData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data, err := Decrypt(encData, key)
	if err != nil {
		return nil, err
	}

	var s BruteState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func SaveState(s *BruteState, path string, key []byte) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	enc, err := Encrypt(data, key)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, enc, 0600)
}

func (s *BruteState) Stop() {
	s.Stopped = true
}
