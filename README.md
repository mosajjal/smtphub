# smtphub

SMTPhub is an SMTP server designed to function as an email hook.

## Installation

```bash
$ go install github.com/mosajjal/smtphub@latest 
```

## Configuration

A sample configuration file is provided in `conf.yaml`.

sample config

```yaml
server:
  listen: "0.0.0.0:1025"
  useTLS: false
  tlsCert: ""
  tlsKey: ""
  appName: "MailHog"
  hostname: "localhost"
  auth:
    allowAnon: false
    users:
      - username: "user1"
        password: "pass1"
      - username: "user2"
        password: "pass2"
   
hooks:
  - name: "hook1"
    conditions:
      - subject: ".*?test1.*?"
      # - body: ".*? test2 .*?"
      # - from: ".*? test3 .*?"
      # - to: ".*? test4 .*?"
      # - remoteAddr: ".*? test5 .*?"
    actions:
      - type: exec
        command: "sleep 1; echo '{{.Subject}}' '{{.Body}}' '{{.From}}' '{{.To}}' '{{.RemoteAddr}}'"
        timeout: 2s
        env:
          - name: "TEST"
            value: "test"
```

After installation, a sample configuration file named conf.yaml is available. The configuration file includes server settings such as the listening address and TLS certificate, as well as authentication options for users.
A section of the configuration file called "hooks" enables the user to specify conditions for incoming emails, such as a particular subject line, and to define an action, such as executing a command.
The conditions specified for each hook are connected with an AND logic operator. The example configuration provided will launch an SMTP server and execute an echo command for every email received that contains test1 in the subject.

