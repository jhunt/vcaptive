package vcaptive

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type Application struct {
	ID      string   `json:"application_id"`
	Name    string   `json:"application_name"`
	Version string   `json:"application_version"`
	URIs    []string `json:"application_uris"`
}

type Services map[string][]Instance

type Instance struct {
	Name           string      `json:"name"`
	Label          string      `json:"label"`
	Tags           []string    `json:"tags"`
	Plan           string      `json:"plan"`
	Credentials    Credentials `json:"credentials"`
	Provider       interface{} `json:"provider"`
	SyslogDrainURL interface{} `json:"syslog_drain_url"`
}

type Credentials map[string]interface{}

func ParseServices(interf interface{}) (Services, error) {
	var s string
	switch interf.(type) {
	case string:
		s = interf.(string)
	default:
		tem, err := json.Marshal(interf)
		if err != nil {
			os.Exit(1)
		}
		s = string(tem)
	}
	var ss Services
	return ss, json.Unmarshal([]byte(s), &ss)
}

func ParseApplication(s string) (Application, error) {
	var a Application
	return a, json.Unmarshal([]byte(s), &a)
}

func (ss Services) Tagged(tags ...string) (Instance, bool) {
	for _, list := range ss {
		for _, svc := range list {
			for _, have := range svc.Tags {
				for _, want := range tags {
					if have == want {
						return svc, true
					}
				}
			}
		}
	}
	return Instance{}, false
}

func (ss Services) WithCredentials(keys ...string) (Instance, bool) {
	for _, list := range ss {
		for _, svc := range list {
			found := true
			for _, want := range keys {
				if _, ok := svc.Get(want); !ok {
					found = false
					break
				}
			}
			if found {
				return svc, true
			}
		}
	}
	return Instance{}, false
}

func (inst Instance) Get(key string) (interface{}, bool) {
	var o interface{}

	o = inst.Credentials
	for _, p := range strings.Split(key, ".") {
		switch o.(type) {
		case Credentials:
			v, ok := o.(Credentials)[p]
			if !ok {
				return nil, false
			}
			o = v

		case map[string]interface{}:
			v, ok := o.(map[string]interface{})[p]
			if !ok {
				return nil, false
			}
			o = v

		case []interface{}:
			u, err := strconv.ParseUint(p, 10, 0)
			if err != nil {
				return nil, false
			}
			i := int(u)
			if i >= len(o.([]interface{})) {
				return nil, false
			}
			o = o.([]interface{})[i]

		default:
			return nil, false
		}
	}

	return o, true
}
