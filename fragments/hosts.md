## Example: Generated SSH Configuration and a Ping Check

This is a longer example, where a "source of truth" file is used

- to generate an SSH configuration, and 
- to generate `ping` commands to see if the hosts are up and reachable.

### Source of Truth

Here is a file that defines a bunch of (imaginary) hosts. It is the substrate for the generation of the SSH configuration `~/.ssh/config` and for a ping test.

Each entry in the list is a map that defines for a host:

- A description (for comment purposes)
- A hostname or IP address (required)
- A short hame. When given, this name is taken as an alias that `ssh` understands.
- A user name. When given, `ssh` connections are set up as this user, when not, the logged-in user will be taken.
- A port. When given, the `ssh` connections run over the stated port.
- Whether the host supports X11. When yes, X11 connections are forwarded over the `ssh` connection.
- An identity file. When given, that specific identity file is used for authentication.

You can also find this as `examples/hosts/hosts`.

```
{{{/* Configuration of ssh-able hosts */}}

{{
    $hosts := list
        (map "description"   "My computer at work"
             "hostname"      "ws12345.example.com"
             "shortname"     "ws"
             "user"          "user12345"
             "idfile"        "/home/user/.ssh/specific_id_rsa_file"
             "hasX11"        true)

        (map "description"   "Bastion, emergency access from outside"
             "hostname"      "123.456.789.012"
             "shortname"     "bastion"
             "user"          "emergency"
             "port"          2222)

        (map "description"   "Office DHCP/Router"
             "hostname"      "192.168.1.1"
             "user"          "pi")

        (map "description"   "Raspberry Pi DNS/Blackhole"
             "hostname"      "192.168.1.10"
             "shortname"     "pi"
             "user"          "pi")
}}
```

### Generation of an SSH Configuration

In order to generate a configuration that `ssh` understands, the following template is used. It iterates over the above list, and for each map it picks up keys and values. Additionally, a few defaults are injected into each SSH host configuration: to use compression, to multiplex connections once logged in, and so on.

You can also find this as `examples/hosts/ssh-config`.

```
{{/* Expansion of hosts */}}

{{/* Use as follows:
    gtpl -re hosts ssh-config > ~/.ssh/config
*/}}

{{ range $h := $hosts }}
#
# {{ getval $h "description" }}
# ------------------------------------------------------------------------
Host {{ if contains $h "shortname"}}{{ getval $h "shortname"}}{{ else }}{{ getval $h "hostname" }}{{ end }}
    Hostname {{ getval $h "hostname" }}
    {{ if contains $h "user" }}User {{ getval $h "user" }}{{ end }}
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
{{ if contains $h "idfile" }}
    IdentityFile {{ getval $h "idfile" }}
    IdentitiesOnly yes
{{ end }}
{{ if contains $h "hasX11" }}
    ForwardX11Trusted yes
{{ end }}
{{ if contains $h "port" }}
    Port {{ getval $h "port" }}
{{ end }}
{{ end }}
```

Sample output:

```
#
# My computer at work
# ------------------------------------------------------------------------
Host ws
    Hostname ws12345.example.com
    User user12345
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
    IdentityFile /home/user/.ssh/specific_id_rsa_file
    IdentitiesOnly yes
    ForwardX11Trusted yes
#
# Bastion, emergency access from outside
# ------------------------------------------------------------------------
Host bastion
    Hostname 123.456.789.012
    User emergency
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
    Port 2222
#
# Office DHCP/Router
# ------------------------------------------------------------------------
Host 192.168.1.1
    Hostname 192.168.1.1
    User pi
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
#
# Raspberry Pi DNS/Blackhole
# ------------------------------------------------------------------------
Host pi
    Hostname 192.168.1.10
    User pi
    Compression yes
    ControlMaster auto
    ControlPath ~/.ssh/ctrl-%C
    ControlPersist yes
    ServerAliveInterval 300
```

### Pinging the Configured Hosts

The following template picks up hosts information and generates `ping` comands to see if the hosts are accessible.

You can also find this as `examples/hosts/ping-test`.

```
{{/* Ping test for hosts. */}}

{{/* Use as follows:
    gtpl -re hosts ping-test | sh
*/}}

{{ range $h := $hosts }}
ping -c3 -t3 {{ getval $h "hostname" }} >/dev/null 2>&1 || echo {{ getval $h "hostname" }} is unreachable!
{{ end }}
```

Example output (which can be fed to the `/bin/sh`):

```
ping -c3 -t3 ws12345.example.com >/dev/null 2>&1 || echo ws12345.example.com is unreachable!
ping -c3 -t3 123.456.789.012 >/dev/null 2>&1 || echo 123.456.789.012 is unreachable!
ping -c3 -t3 192.168.1.1 >/dev/null 2>&1 || echo 192.168.1.1 is unreachable!
ping -c3 -t3 192.168.1.10 >/dev/null 2>&1 || echo 192.168.1.10 is unreachable!
```
