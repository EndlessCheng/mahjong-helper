package main

type flagKV map[string]string

func parseArgs(args []string) (flags flagKV, restArgs []string) {
	flags = flagKV{}
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

func (f flagKV) Bool(flagNames ...string) bool {
	for _, name := range flagNames {
		if _, ok := f[name]; ok {
			return true
		}
	}
	return false
}

func (f flagKV) String(flagNames ...string) string {
	for _, name := range flagNames {
		if val, ok := f[name]; ok {
			return val
		}
	}
	return ""
}
