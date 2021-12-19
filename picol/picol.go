package picol

import (
	"errors"
	"fmt"
	"strings"
)

var (
	PICOL_RETURN   = errors.New("RETURN")
	PICOL_BREAK    = errors.New("BREAK")
	PICOL_CONTINUE = errors.New("CONTINUE")
)

type Var string
type CmdFunc func(i *Interp, argv []string, privdata interface{}) (string, error)
type Cmd struct {
	fn       CmdFunc
	privdata interface{}
}
type CallFrame struct {
	vars   map[string]Var
	parent *CallFrame
}
type Interp struct {
	level     int
	callframe *CallFrame
	commands  map[string]Cmd
}

func InitInterp() *Interp {
	return &Interp{
		level:     0,
		callframe: &CallFrame{vars: make(map[string]Var)},
		commands:  make(map[string]Cmd),
	}
}

func (i *Interp) Var(name string) (Var, bool) {
	for frame := i.callframe; frame != nil; frame = frame.parent {
		v, ok := frame.vars[name]
		if ok {
			return v, ok
		}
	}
	return "", false
}
func (i *Interp) SetVar(name, val string) {
	i.callframe.vars[name] = Var(val)
}

func (i *Interp) UnsetVar(name string) {
	delete(i.callframe.vars, name)
}

func (i *Interp) Command(name string) *Cmd {
	v, ok := i.commands[name]
	if !ok {
		return nil
	}
	return &v
}

func (i *Interp) RegisterCommand(name string, fn CmdFunc, privdata interface{}) error {
	c := i.Command(name)
	if c != nil {
		return fmt.Errorf("Command '%s' already defined", name)
	}

	i.commands[name] = Cmd{fn, privdata}
	return nil
}

/* EVAL! */
func (i *Interp) Eval(t string) (string, error) {
	p := InitParser(t)
	var result string
	var err error

	argv := []string{}

	for {
		prevtype := p.Type
		// XXX
		t = p.GetToken()
		if p.Type == PT_EOF {
			break
		}

		switch p.Type {
		case PT_VAR:
			v, ok := i.Var(t)
			if !ok {
				return "", fmt.Errorf("No such variable '%s'", t)
			}
			t = string(v)
		case PT_CMD:
			result, err = i.Eval(t)
			if err != nil {
				return result, err
			} else {
				t = result
			}
		case PT_ESC:
			// XXX: escape handling missing!
		case PT_SEP:
			prevtype = p.Type
			continue
		}

		// We have a complete command + args. Call it!
		if p.Type == PT_EOL {
			prevtype = p.Type
			if len(argv) != 0 {
				c := i.Command(argv[0])
				if c == nil {
					return "", fmt.Errorf("No such command '%s'", argv[0])
				}
				result, err = c.fn(i, argv, c.privdata)
				if err != nil {
					return result, err
				}
			}
			// Prepare for the next command
			argv = []string{}
			continue
		}

		// We have a new token, append to the previous or as new arg?
		if prevtype == PT_SEP || prevtype == PT_EOL {
			argv = append(argv, t)
		} else { // Interpolation
			argv[len(argv)-1] = strings.Join([]string{argv[len(argv)-1], t}, "")
		}
		prevtype = p.Type
	}
	return result, nil
}
