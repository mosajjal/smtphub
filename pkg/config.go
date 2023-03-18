package smtphub

import "time"

type Config struct {
	Server struct {
		Listen   string `yaml:"listen"`
		UseTLS   bool   `yaml:"useTLS"`
		TLSCert  string `yaml:"tlsCert"`
		TLSKey   string `yaml:"tlsKey"`
		AppName  string `yaml:"appName"`
		Hostname string `yaml:"hostname"`
		Auth     struct {
			AllowAnon bool `yaml:"allowAnon"`
			Users     []struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"users"`
		} `yaml:"auth"`
	} `yaml:"server"`
	Hooks []Hook `yaml:"hooks"`
}

type HookCondition struct {
	Subject    string `yaml:"subject"`
	Body       string `yaml:"body"`
	From       string `yaml:"from"`
	To         string `yaml:"to"`
	RemoteAddr string `yaml:"remoteAddr"`
}

type HookAction struct {
	Type    string            `yaml:"type"`
	Command string            `yaml:"command"`
	Timeout time.Duration     `yaml:"timeout"`
	Env     map[string]string `yaml:"env"`
}

type Hook struct {
	Name       string          `yaml:"name"`
	Conditions []HookCondition `yaml:"conditions"`
	Actions    []HookAction    `yaml:"actions"`
}
