package ansiterm

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
)

var parser *AnsiParser

type AnsiParser struct {
	state        State
	eventHandler AnsiEventHandler
	context      *AnsiContext
}

func CreateParser(initialState State, evtHandler AnsiEventHandler) *AnsiParser {
	logrus.Infof("CreateParser")

	parser = &AnsiParser{
		state:        initialState,
		eventHandler: evtHandler,
		context:      &AnsiContext{},
	}

	return parser
}

func (ap *AnsiParser) Parse(bytes []byte) (int, error) {
	for i, b := range bytes {
		if err := ap.handle(b); err != nil {
			return i, err
		}
	}

	return len(bytes), nil
}

func (ap *AnsiParser) handle(b byte) error {
	newState, err := ap.state.Handle(b)
	if err != nil {
		return err
	}

	if newState == nil {
		logrus.Warning("newState is nil")
		return errors.New(fmt.Sprintf("New state of 'nil' is invalid."))
	}

	if newState != ap.state {
		if err := ap.changeState(newState); err != nil {
			return err
		}
	}

	return nil
}

func (ap *AnsiParser) changeState(newState State) error {
	logrus.Infof("ChangeState %s --> %s", ap.state.Name(), newState.Name())

	// Exit old state
	if err := ap.state.Exit(); err != nil {
		logrus.Infof("Exit state '%s' failed with : '%v'", ap.state.Name(), err)
		return err
	}

	// Perform transition action
	if err := ap.state.Transition(newState); err != nil {
		logrus.Infof("Transition from '%s' to '%s' failed with: '%v'", ap.state.Name(), newState.Name, err)
		return err
	}

	// Enter new state
	if err := newState.Enter(); err != nil {
		logrus.Infof("Enter state '%s' failed with: '%v'", newState.Name(), err)
		return err
	}

	ap.state = newState
	return nil
}
