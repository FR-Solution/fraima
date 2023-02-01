package generator

import (
	"fmt"
	"os/exec"
)

func getMap(i any) (map[string]any, error) {
	rArgs := make(map[string]any)
	err := fmt.Errorf("args converting is not available")
	args, ok := i.(map[any]any)
	if !ok {
		return rArgs, err
	}
	for k, v := range args {
		key := fmt.Sprint(k)
		if nArgs, ok := v.(map[any]any); ok {
			rArgs[key], err = getMap(nArgs)
			if err != nil {
				return rArgs, err
			}
			continue
		}

		rArgs[key] = v
	}
	return rArgs, nil
}

func startService(name string) error {
	err := exec.Command("systemctl", "enable", fmt.Sprintf("%s.service", name)).Run()
	if err != nil {
		return err
	}
	err = exec.Command("systemctl", "start", fmt.Sprintf("%s.service", name)).Run()
	if err != nil {
		return err
	}
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		return err
	}
	return nil
}
