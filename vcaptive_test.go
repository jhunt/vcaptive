package vcaptive_test

import (
	"testing"

	"github.com/jhunt/vcaptive"
)

func TestApplication(t *testing.T) {
	var (
		topic string
		err error
		app vcaptive.Application
	)

	topic = "docs.cloudfoundry.org"
	app, err = vcaptive.ParseApplication(`
{
  "instance_id": "fe98dc76ba549876543210abcd1234",
  "instance_index": 0,
  "host": "0.0.0.0",
  "port": 61857,
  "started_at": "2013-08-12 00:05:29 +0000",
  "started_at_timestamp": 1376265929,
  "start": "2013-08-12 00:05:29 +0000",
  "state_timestamp": 1376265929,
  "limits": {
    "mem": 512,
    "disk": 1024,
    "fds": 16384
  },
  "application_version": "ab12cd34-5678-abcd-0123-abcdef987654",
  "application_name": "styx-james",
  "application_uris": [
    "my-app.example.com"
  ],
  "version": "ab12cd34-5678-abcd-0123-abcdef987654",
  "name": "my-app",
  "uris": [
    "my-app.example.com"
  ],
  "users": null
}
`)
	if err != nil {
		t.Fatalf("[%s] failed to parse VCAP_APPLICATION: %s", topic, err)
	}
	if app.Name != "styx-james" {
		t.Errorf("[%s] expected name to be 'styx-james', but got '%s'", topic, app.Name)
	}
	if app.Version != "ab12cd34-5678-abcd-0123-abcdef987654" {
		t.Errorf("[%s] unexpected version '%s'", topic, app.Version)
	}
	if len(app.URIs) != 1 {
		t.Fatalf("[%s] unexpected number of URIs; expected %d but got %d", topic, 1, len(app.URIs))
	}
	if app.URIs[0] != "my-app.example.com" {
		t.Errorf("[%s] unexpected uri[0] '%s'", topic, app.URIs[0])
	}

	topic = "live example"
	app, err = vcaptive.ParseApplication(`
{
  "application_id": "e233016d-3bce-4e1e-9269-b1ad1555cf99",
  "application_name": "my-test-app",
  "application_uris": [
   "my-test-app.cfapps.io"
  ],
  "application_version": "35c179da-ae9a-4cb6-b787-98261b3bb183",
  "cf_api": "https://api.cfapps.io",
  "limits": {
   "disk": 1024,
   "fds": 16384,
   "mem": 1024
  },
  "name": "my-test-app",
  "space_id": "1afffefc-6318-4b72-8383-7bac3fdc6ec6",
  "space_name": "stark-and-wayne",
  "uris": [
   "my-test-app.cfapps.io"
  ],
  "users": null,
  "version": "35c179da-ae9a-4cb6-b787-98261b3bb183"
}
`)
	if err != nil {
		t.Fatalf("[%s] failed to parse VCAP_APPLICATION: %s", topic, err)
	}
	if app.Name != "my-test-app" {
		t.Errorf("[%s] expected name to be 'styx-james', but got '%s'", topic, app.Name)
	}
	if app.Version != "35c179da-ae9a-4cb6-b787-98261b3bb183" {
		t.Errorf("[%s] unexpected version '%s'", topic, app.Version)
	}
	if len(app.URIs) != 1 {
		t.Fatalf("[%s] unexpected number of URIs; expected %d but got %d", topic, 1, len(app.URIs))
	}
	if app.URIs[0] != "my-test-app.cfapps.io" {
		t.Errorf("[%s] unexpected uri[0] '%s'", topic, app.URIs[0])
	}
}

func TestServices(t *testing.T) {
	var (
		topic string
		err   error
		ok    bool
		v     interface{}

		ss   vcaptive.Services
		inst vcaptive.Instance
	)

	topic = "docs.cloudfoundry.org"
	ss, err = vcaptive.ParseServices(`
{
  "elephantsql": [
    {
      "name": "elephantsql-c6c60",
      "label": "elephantsql",
      "tags": [
        "postgres",
        "postgresql",
        "relational"
      ],
      "plan": "turtle",
      "credentials": {
        "uri": "postgres://exampleuser:examplepass@babar.elephantsql.com:5432/exampleuser",
        "port": 5432
      }
    }
  ],
  "sendgrid": [
    {
      "name": "mysendgrid",
      "label": "sendgrid",
      "tags": [
        "smtp"
      ],
      "plan": "free",
      "credentials": {
        "hostname": "smtp.sendgrid.net",
        "username": "QvsXMbJ3rK",
        "password": "HCHMOYluTv"
      }
    }
  ]
}
`)
	if err != nil {
		t.Fatalf("[%s] failed to parse VCAP_SERVICES: %s", topic, err)
	}

	_, ok = ss.Tagged("xyzzy")
	if ok {
		t.Errorf("[%s] found a service tagged 'xyzzy' when there shouldn't be one!", topic)
	}

	_, ok = ss.Tagged("smtp")
	if !ok {
		t.Errorf("[%s] did not find the SMTP service tagged 'smtp'", topic)
	}

	_, ok = ss.Tagged("smtp", "relational")
	if !ok {
		t.Errorf("[%s] did not find a service tagged EITHER 'smtp' OR 'relational'", topic)
	}

	inst, ok = ss.Tagged("postgres")
	if !ok {
		t.Fatalf("[%s] did not find example 'postgres' service for instance testing", topic)
	}
	if inst.Label != "elephantsql" {
		t.Errorf("[%s] postgres service should be labelled 'elephantsql', but was instead '%s'", topic, inst.Label)
	}
	if inst.Plan != "turtle" {
		t.Errorf("[%s] postgres service should be of plan 'turtle', but was instead '%s'", topic, inst.Plan)
	}

	v, ok = inst.Get("foo")
	if ok {
		t.Errorf("[%s] postgres service should not have returned anything for the 'foo' cred, but did: '%v'", topic, v)
	}
	v, ok = inst.GetString("foo")
	if ok {
		t.Errorf("[%s] postgres service should not have returned anything for the 'foo' cred, but did: '%v'", topic, v)
	}

	v, ok = inst.Get("uri")
	if !ok {
		t.Fatalf("[%s] postgres service should have returned a value for the 'uri' cred, but did not", topic)
	}
	if v != "postgres://exampleuser:examplepass@babar.elephantsql.com:5432/exampleuser" {
		t.Errorf("[%s] postgres service returned the wrong value for the 'uri' cred: '%s'", topic, v)
	}
	v, ok = inst.GetString("uri")
	if !ok {
		t.Fatalf("[%s] postgres service should have returned a value for the 'uri' cred, but did not", topic)
	}
	if v != "postgres://exampleuser:examplepass@babar.elephantsql.com:5432/exampleuser" {
		t.Errorf("[%s] postgres service returned the wrong value for the 'uri' cred: '%s'", topic, v)
	}

	if v, ok = inst.GetUint("uri"); ok {
		t.Errorf("[%s] postgres service should not have returned a number value for the 'uri' cred, but did: '%v'", topic, v)
	}

	v, ok = inst.GetUint("port")
	if !ok {
		t.Fatalf("[%s] postgres service should have returned a numeric value for the 'port' cred, but did not", topic)
	}
	if v != uint(5432) {
		t.Errorf("[%s] postgres service returned the wrong value for the 'port' cred: '%d'", topic, v)
	}


	topic = "multi-level credentials"
	ss, err = vcaptive.ParseServices(`
{
  "x": [
    {
      "name": "x",
      "label": "x",
      "tags": [ "x" ],
      "plan": "x",
      "credentials": {
        "ssl": {
          "uri":"https://user:pass@127.0.0.1",
          "ciphers":["foo", "bar"]
        },
        "plain": {
          "uri":"http://user:pass@127.0.0.1"
        }
      }
    }
  ]
}
`)
	if err != nil {
		t.Fatalf("[%s] failed to parse VCAP_SERVICES: %s", topic, err)
	}

	inst, ok = ss.Tagged("x")
	if !ok {
		t.Fatalf("[%s] did not find example 'x' service for instance credentials testing", topic)
	}
	v, ok = inst.Get("uri")
	if ok {
		t.Errorf("[%s] test service should not have returned anything for the 'uri' cred, but did: '%v'", topic, v)
	}

	v, ok = inst.Get("plain.uri")
	if !ok {
		t.Fatalf("[%s] test service should have returned a value for the 'plain.uri' cred, but did not", topic)
	}
    if v != "http://user:pass@127.0.0.1" {
		t.Errorf("[%s] test service returned the wrong value for the 'plain.uri' cred: '%s'", topic, v)
	}

	v, ok = inst.Get("ssl.uri")
	if !ok {
		t.Fatalf("[%s] test service should have returned a value for the 'ssl.uri' cred, but did not", topic)
	}
    if v != "https://user:pass@127.0.0.1" {
		t.Errorf("[%s] test service returned the wrong value for the 'ssl.uri' cred: '%s'", topic, v)
	}

    v, ok = inst.Get("plain.enoent")
	if ok {
		t.Errorf("[%s] test service should not have returned anything for the 'plain.enoent' cred, but did: '%v'", topic, v)
	}

    v, ok = inst.Get("plain.uri.enoent")
	if ok {
		t.Errorf("[%s] test service should not have returned anything for the 'plain.uri.enoent' cred, but did: '%v'", topic, v)
	}

	v, ok = inst.Get("ssl.ciphers.0")
	if !ok {
		t.Fatalf("[%s] test service should have returned a value for the 'ssl.ciphers.0' cred, but did not", topic)
	}
	if v != "foo" {
		t.Errorf("[%s test service returned the wrong value for the 'ssl.ciphers.0' cred: '%s'", topic, v)
	}
	v, ok = inst.Get("ssl.ciphers.1")
	if !ok {
		t.Fatalf("[%s] test service should have returned a value for the 'ssl.ciphers.1' cred, but did not", topic)
	}
	if v != "bar" {
		t.Errorf("[%s test service returned the wrong value for the 'ssl.ciphers.1' cred: '%s'", topic, v)
	}

	v, ok = inst.Get("ssl.ciphers.2")
	if ok {
		t.Fatalf("[%s] test service should not have returned anything for the 'ssl.ciphers.2' cred, but did: '%s'", topic, v)
	}

	v, ok = inst.Get("ssl.ciphers.ONE")
	if ok {
		t.Fatalf("[%s] test service should not have returned anything for the 'ssl.ciphers.ONE' cred, but did: '%s'", topic, v)
	}


	topic = "credentials-based"
	ss, err = vcaptive.ParseServices(`
{
  "x": [
    {
      "name": "vmail-120",
      "label": "vmail",
      "tags": [],
      "plan": "miniscule",
      "credentials": {
        "smtp_host"     : "127.0.0.1",
        "smtp_port"     : "587",
        "smtp_username" : "email",
        "smtp_password" : "secret"
      }
    }
  ]
}
`)
	if err != nil {
		t.Fatalf("[%s] failed to parse VCAP_SERVICES: %s", topic, err)
	}

	_, ok = ss.WithCredentials("www_foo")
	if ok {
		t.Errorf("[%s] found a service with appropriate credentials when there shouldn't be one!", topic)
	}

	inst, ok = ss.WithCredentials("smtp_host", "smtp_port", "smtp_username", "smtp_password")
	if !ok {
		t.Errorf("[%s] did not find the SMTP service", topic)
	}

	if inst.Label != "vmail" {
		t.Errorf("[%s] SMTP service should be labelled 'vmail', but was instead '%s'", topic, inst.Label)
	}
	if inst.Plan != "miniscule" {
		t.Errorf("[%s] SMTP service should be of plan 'miniscule', but was instead '%s'", topic, inst.Plan)
	}

	v, ok = inst.Get("web")
	if ok {
		t.Errorf("[%s] SMTP service should not have returned anything for the 'web' cred, but did: '%v'", topic, v)
	}
	v, ok = inst.Get("smtp_host")
	if !ok {
		t.Fatalf("[%s] SMTP service should have returned a value for the 'smtp_host' cred, but did not", topic)
	}
	if v != "127.0.0.1" {
		t.Errorf("[%s] postgres service returned the wrong value for the 'smtp_host' cred: '%s'", topic, v)
	}
}
