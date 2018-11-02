# `gh-rocketchat`
A service to translate GitHub webhooks to nice bot messages in rocket.chat.

## Why?
Not all of us are javascript experts or have admin privileges on the rocket.chat
instance in use. Hence, uploading/changing the webhook filter can be hard or
even impossible.

This program only requires an incoming webhook on the rocket.chat side to do
its work, easing deployment.

## How to run
Preferably, this program is run a service behind a TLS proxy like nginx or
Apache.

`gh-rocketchat` is configured with a yaml file, who's path can be specified with
the `-config` flag. This file contains the local listen address and all the
hook configurations.

```yaml
# This is the address that gh-rocketchat listens on locally
listenaddress: 127.0.0.1:8000

# This key specifies all the webhooks.
# Each webhook must have a unique name and unique endpoint, it must have a destination and optionally a secret to
# authenticate incoming webhooks

hooks:
  - name: myHook          # Must be unique between all hooks
    endpoint: /hook1      # Must also be unique. This is the URL endpoint where GitHub delivers the hook.
                          # In this case it would be machine.example.com/hook1
    secret: s3cr3t        # The shared secret for GitHub
    destination: https://rocketchat.example.com/hooks/SOMESUPERLONGSTRING
                          # Where to deliver the transformed webhooks
```

## Status
Very much Work-in-Progress. Right now `gh-rocketchat` can listen for webhooks
from GitHub and send Ping and Push webhooks to a downstream rocket.chat
instance.

### Future work
#### Before 1.0.0
- Handle Issue + issue comment events
- Handle Pullrequest events
- Handle Label events
- Handle Pullrequest Review + comment events
- Handle tag/branch creation and deletion

#### After 1.0.0
- Handle _all_ webhooks
- Be able to select repositories to process
- Be able to select events to process
- Combination of the two options above in a single hook
- Provide custom templates for events

## License
MIT
