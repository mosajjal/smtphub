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


## Sample Setup with a VPS

### Step 1: Get a VPS from Vultr

Firstly, you need to purchase a VPS (Virtual Private Server) from a provider like Vultr. You will be given an IP address for the VPS upon creation. Make a note of this IP address as you will need it for the next steps. Keep in mind that port 25 inbound and outbound is blocked by default on Vultr VPSes. You will need to open a support ticket to have this port unblocked.

### Step 2: Set up subdomain and MX records

Next, create a subdomain for your email server, such as mail.example.com. You can create this subdomain in the DNS management tool provided by your domain registrar.

If your subdomain is not already pointing to your VPS's IP address, you need to add an A record to your DNS settings to direct traffic to your VPS. To do this, you can create a new A record with your subdomain as the name and the IP address of your VPS as the value.

After creating the subdomain and pointing it to your VPS's IP address, add MX records for the subdomain to direct email traffic to your VPS. You can create new MX records with your subdomain as the name and your VPS's IP address as the value. You can also set the priority of the MX records, where the lowest value indicates the highest priority.

For example, you can add the following MX records for mail.example.com:

```
mail.example.com.   MX   10   your.vps.ip.address
mail.example.com.   MX   20   your.vps.ip.address
```

This will direct email traffic to your VPS with higher priority given to the MX record with a priority of 10.

### Step 3: Get a TLS certificate

To secure your email traffic, obtain a TLS certificate for your email domain. You can obtain a free certificate from Let's Encrypt, or purchase one from a certificate authority.

If you are using Let's Encrypt, you can install the Certbot tool on your VPS to automate the certificate acquisition process. You can run the following commands to install and run Certbot:

```bash
$ sudo apt update
$ sudo apt install certbot
$ sudo certbot certonly --standalone -d mail.example.com
```

This will install Certbot and obtain a TLS certificate for mail.example.com using the standalone plugin.

### Step 4: Install and configure SMTPhub

Finally, you can install SMTPhub on your VPS by executing the following command in a terminal:

```bash
$ go install github.com/mosajjal/smtphub@latest 
```

After installation, you need to configure SMTPhub to use the TLS certificate and listen on the subdomain created earlier. You can also configure user authentication, specify email conditions, and define actions in the configuration file.

For example, you can configure SMTPhub to use the TLS certificate obtained in Step 3 and listen on mail.example.com:25 by modifying the conf.yaml file as follows:

```yaml
server:
  listen: "mail.example.com:25"
  useTLS: true
  tlsCert: "/etc/letsencrypt/live/mail.example.com/fullchain.pem"
  tlsKey: "/etc/letsencrypt/live/mail.example.com/privkey.pem"
```

If you added user authentication in the configuration file, make sure to create those users and their passwords using the smtphub-users tool provided by SMTPhub.

Make sure to restart SMTPhub after any changes to the configuration file.