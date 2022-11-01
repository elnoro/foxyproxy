# FoxyProxy

```
fxpr is a CLI tool to quickly spin up and destroy DigitalOcean servers

Usage:
  fxpr [command]

Available Commands:
  proxy         Start a droplet and an SSH tunnel on localhost. Hit Ctrl-C to destroy the droplet
  test          Start a droplet you can SSH into. Hit Ctrl-C to destroy the droplet
  list          Shows the list of existing droplets
```


## Configuration

### Config file
Put the config file in `"$HOME/.config/fxpr/config.json"`
You need to create a Digital Ocean token and register an ssh key.

Config file example:
```
{
  "do_token": "generated-digital-ocean-token",
  "fingerprint": "registered-ssh-key-fingerprint",
  "port": 1337
}
```

### Firefox

[Use this instruction from Digital Ocean](https://www.digitalocean.com/community/tutorials/how-to-route-web-traffic-securely-without-a-vpn-using-a-socks-tunnel#step-2-mdash-configuring-firefox-to-use-the-tunnel)