package config

import "strconv"

type FlagKV map[string]string

func ParseArgs(args []string) (flags FlagKV, restArgs []string) {
	flags = FlagKV{}
	for _, arg := range args {
		if len(arg) >= 2 && arg[0] == '-' {
			found := false
			for i, c := range arg {
				if c == '=' {
					flags[arg[1:i]] = arg[i+1:]
					found = true
					break
				}
			}
			if !found {
				flags[arg[1:]] = ""
			}
		} else {
			restArgs = append(restArgs, arg)
		}
	}
	return
}

func (f FlagKV) Bool(flagNames ...string) bool {
	for _, name := range flagNames {
		if _, ok := f[name]; ok {
			return true
		}
	}
	return false
}

func (f FlagKV) String(flagNames ...string) string {
	for _, name := range flagNames {
		if val, ok := f[name]; ok {
			return val
		}
	}
	return ""
}

func (f FlagKV) Int(flagNames ...string) (int, error) {
	for _, name := range flagNames {
		if val, ok := f[name]; ok {
			return strconv.Atoi(val)
		}
	}
	return -1, nil
}
