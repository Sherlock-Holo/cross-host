package main

import (
	"errors"
)

var (
	NoContainerId = errors.New("container id is empty")
	NoNetns       = errors.New("netns path is empty")
	NoIfaceName   = errors.New("iface name is empty")
)
