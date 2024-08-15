package main

import (
	"errors"

	"github.com/bp-chat/bp-tui/commands"
)

type CommandQueue struct {
	onAir   [MaxNumberOfCommands]bool
	waiting []commands.IOut
}

func CreateCommandQueue() CommandQueue {
	return CommandQueue{
		onAir:   [MaxNumberOfCommands]bool{},
		waiting: make([]commands.IOut, 0),
	}
}

func (queue *CommandQueue) TakeSlot() (int, error) {
	i := 0
	for i < MaxNumberOfCommands {

		if queue.onAir[i] == false {
			queue.onAir[i] = true
			return i, nil
		}
	}
	return 0, errors.New("No command slot available")
}

func (queue *CommandQueue) Free(slot int) {
	queue.onAir[slot] = false
}

func (queue *CommandQueue) Enqueue(out commands.IOut) {
	queue.waiting = append(queue.waiting, out)
}

func (queue *CommandQueue) Pop() (*commands.IOut, error) {
	if len(queue.waiting) == 0 {
		return nil, errors.New("The list is empty")
	}
	out := queue.waiting[0]
	queue.waiting = queue.waiting[1:]
	return &out, nil
}
