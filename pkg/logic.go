package smtphub

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/mhale/smtpd"
)

var config Config

var handler = smtpd.Handler(func(remoteAddr net.Addr, from string, to []string, data []byte) error {
	fmt.Printf("Received mail from %s to %s", from, to)
	// run the email through all the hook conditions
	// if any of the conditions match, run the hook actions
	// if none of the conditions match, do nothing
	for _, hook := range config.Hooks {
		fmt.Printf("Checking hook %s\n", hook.Name)
		for _, condition := range hook.Conditions {
			if condition.Subject != "" {
				// do a regex match on the subject
				regexp, err := regexp.Compile(condition.Subject)
				if err != nil {
					log.Fatal(err)
				}
				if !regexp.MatchString(string(data)) {
					// subject doesn't match
					continue
				}
			}
			if condition.Body != "" {
				// do a regex match on the body
				regexp, err := regexp.Compile(condition.Body)
				if err != nil {
					log.Fatal(err)
				}
				if !regexp.MatchString(string(data)) {
					// body doesn't match
					continue
				}
			}
			if condition.From != "" {
				// do a regex match on the from
				regexp, err := regexp.Compile(condition.From)
				if err != nil {
					log.Fatal(err)
				}
				if !regexp.MatchString(from) {
					// from doesn't match
					continue
				}
			}
			if condition.To != "" {
				// do a regex match on the to
				regexp, err := regexp.Compile(condition.To)
				if err != nil {
					log.Fatal(err)
				}
				if !regexp.MatchString(strings.Join(to, ",")) {
					// to doesn't match
					continue
				}
			}
			if condition.RemoteAddr != "" {
				// do a regex match on the remoteAddr
				regexp, err := regexp.Compile(condition.RemoteAddr)
				if err != nil {
					log.Fatal(err)
				}
				if !regexp.MatchString(remoteAddr.String()) {
					// remoteAddr doesn't match
					continue
				}
			}
			// all conditions match, run the actions
			for _, action := range hook.Actions {
				fmt.Printf("Running action %s\n", action.Type)
				switch action.Type {
				case "exec":
					// run the command
					fmt.Printf("Running command %s\n", action.Command)
					env := os.Environ()
					for k, v := range action.Env {
						env = append(env, fmt.Sprintf("%s=%s", k, v))
					}
					// command has a template that needs to be filled in
					// example: echo {{.Subject}} {{.Body}} {{.From}} {{.To}} {{.RemoteAddr}}

					template, err := template.New("command").Parse(action.Command)
					if err != nil {
						log.Println(err)
					}
					var command strings.Builder
					if err := template.Execute(&command, struct {
						Subject    string
						Body       string
						From       string
						To         string
						RemoteAddr string
					}{
						Subject:    string(data),
						Body:       string(data),
						From:       from,
						To:         strings.Join(to, ","),
						RemoteAddr: remoteAddr.String(),
					}); err != nil {
						log.Println(err)
					}
					ctx, cancel := context.WithTimeout(context.Background(), action.Timeout)
					defer cancel()
					cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command.String())
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Env = env
					err = cmd.Run()
					if err != nil {
						log.Println(err)
					}
				default:
					fmt.Printf("Unknown action type %s\n", action.Type)
				}

			}
		}
	}
	return nil
})

var authHandler = smtpd.AuthHandler(func(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	if config.Auth.AllowAnon {
		return true, nil
	}
	for _, user := range config.Auth.Users {
		if user.Username == string(username) && user.Password == string(password) {
			return true, nil
		}
	}
	return false, nil
})

func Run(c Config) {
	// start a SMTP server
	fmt.Println("Starting SMTP server...")
	config = c

	srv := smtpd.Server{
		Addr:        config.Server.Listen,
		Appname:     config.Server.AppName,
		AuthHandler: authHandler,
		Handler:     handler,
		Hostname:    config.Server.Hostname,
	}
	if config.Server.UseTLS {
		err := srv.ConfigureTLS(config.Server.TLSCert, config.Server.TLSKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Fatal(srv.ListenAndServe())

}
