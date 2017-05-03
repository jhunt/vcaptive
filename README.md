vcaptive
========

`vcaptive` is a small Go library for consuming Cloud Foundry
`$VCAP_SERVICES`-style service and credential definitions.

It takes this:

```
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
```

And lets you do this:

```
package main

import (
  "fmt"
  "os"

  "github.com/jhunt/vcaptive"
)

func main() {
  services, err := vcaptive.Parse(os.Getenv("VCAP_SERVICES"))
  if err != nil {
    fmt.Fprintf(os.Stderr, "VCAP_SERVICES: %s\n", err)
    os.Exit(1)
  }

  instance, found := services.Tagged("postgres", "postgresql")
  if !found {
    fmt.Fprintf(os.Stderr, "VCAP_SERVICES: no 'postgres' service found\n")
    os.Exit(2)
  }

  uri, ok := instance.Get("uri")
  if !ok {
    fmt.Fprintf(os.Stderr, "VCAP_SERVICES: '%s' service has no 'uri' credential\n", instance.Label)
    os.Exit(3)
  }

  fmt.Printf("Connecting to %s...\n", uri)
  // ...
}
```

If you don't have tags, you can also retrieve the first service
that has a given set of credentials:

```
inst, found := services.WithCredentials("smtp_host", "smtp_username")
```

Resources
---------

- [Cloud Foundry (OSS) Documentation on `VCAP_SERVICES`][1]

Contributing
------------

I wrote this tool because I needed it, and no one else had written
it.  My hope is that you find it useful as well.  If it's close,
but not 100% there, why not fork it, fix what needs fixing, and
submit a pull request?

Happy Hacking!


[1]: https://docs.cloudfoundry.org/devguide/deploy-apps/environment-variable.html#VCAP-SERVICES
