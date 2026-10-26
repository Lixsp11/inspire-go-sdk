// Package config manages named CLI profiles without storing session credentials.
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Profile struct {
	Name           string `json:"name"`
	Endpoint       string `json:"endpoint,omitempty"`
	CookieFile     string `json:"cookie_file,omitempty"`
	Socks5Proxy    string `json:"socks5_proxy,omitempty"`
	HttpProxy      string `json:"http_proxy,omitempty"`
	HttpsProxy     string `json:"https_proxy,omitempty"`
	NoProxy        string `json:"no_proxy,omitempty"`
	ReadTimeout    int    `json:"read_timeout,omitempty"`
	ConnectTimeout int    `json:"connect_timeout,omitempty"`
	RetryCount     int    `json:"retry_count,omitempty"`
}

func (p Profile) String() string {
	return "<Inspire profile: credential paths and proxy values redacted>"
}
func (p Profile) GoString() string { return p.String() }

type Configuration struct {
	CurrentProfile string    `json:"current"`
	Profiles       []Profile `json:"profiles"`
}

func NewConfiguration() *Configuration {
	return &Configuration{CurrentProfile: "default", Profiles: []Profile{}}
}
func LoadConfigurationFromFile(path string) (*Configuration, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return NewConfiguration(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot open CLI configuration")
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, 1<<20+1))
	if err != nil || len(data) > 1<<20 {
		return nil, fmt.Errorf("cannot read CLI configuration")
	}
	c := NewConfiguration()
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err = d.Decode(c); err != nil {
		return nil, fmt.Errorf("invalid CLI configuration JSON")
	}
	if err = d.Decode(new(interface{})); err != io.EOF {
		return nil, fmt.Errorf("invalid CLI configuration JSON")
	}
	seen := map[string]bool{}
	for _, p := range c.Profiles {
		if p.Name == "" || seen[p.Name] || p.ReadTimeout < 0 || p.ConnectTimeout < 0 || p.RetryCount < 0 {
			return nil, fmt.Errorf("invalid CLI profile")
		}
		seen[p.Name] = true
	}
	return c, nil
}
func (c *Configuration) GetProfile(name string) (Profile, bool) {
	for _, p := range c.Profiles {
		if p.Name == name {
			return p, true
		}
	}
	return Profile{Name: name}, false
}
func (c *Configuration) PutProfile(profile Profile) {
	for i, p := range c.Profiles {
		if p.Name == profile.Name {
			c.Profiles[i] = profile
			return
		}
	}
	c.Profiles = append(c.Profiles, profile)
}
func (c *Configuration) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot encode CLI configuration")
	}
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create configuration directory")
	}
	f, err := os.CreateTemp(dir, ".config-*")
	if err != nil {
		return fmt.Errorf("cannot create configuration file")
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if err = f.Chmod(0600); err == nil {
		_, err = f.Write(append(data, '\n'))
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err == nil {
		err = os.Rename(tmp, path)
	}
	if err != nil {
		return fmt.Errorf("cannot save CLI configuration")
	}
	return nil
}
