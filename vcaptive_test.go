package vcaptive_test

import (
	"testing"

	"github.com/jhunt/vcaptive"
)

func TestVCAPtive(t *testing.T) {
	var (
		topic string
		err   error
		ok    bool
		v     interface{}

		ss   vcaptive.Services
		inst vcaptive.Instance
	)

	topic = "docs.cloudfoundry.org"
	ss, err = vcaptive.Parse(`
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
        "uri": "postgres://exampleuser:examplepass@babar.elephantsql.com:5432/exampleuser"
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
		t.Fatal("[%s] failed to parse VCAP_SERVICES: %s", topic, err)
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
		t.Fatal("[%s] did not find example 'postgres' service for instance testing", topic)
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
	v, ok = inst.Get("uri")
	if !ok {
		t.Fatalf("[%s] postgres service should have returned a value for the 'uri' cred, but did not", topic)
	}
	if v != "postgres://exampleuser:examplepass@babar.elephantsql.com:5432/exampleuser" {
		t.Errorf("[%s] postgres service returned the wrong value for the 'uri' cred: '%s'", topic, v)
	}


	topic = "multi-level credentials"
	ss, err = vcaptive.Parse(`
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
		t.Fatal("[%s] failed to parse VCAP_SERVICES: %s", topic, err)
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
	ss, err = vcaptive.Parse(`
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
		t.Fatal("[%s] failed to parse VCAP_SERVICES: %s", topic, err)
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
