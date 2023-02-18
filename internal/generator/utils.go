package generator

import (
	"fmt"
)

func getMap(i any) (map[string]any, error) {
	if i == nil {
		return nil, nil
	}
	args, ok := i.(map[any]any)
	if !ok {
		return nil, errArgsUnavailable
	}

	rArgs := make(map[string]any)
	for k, v := range args {
		key := fmt.Sprint(k)
		if nArgs, ok := v.(map[any]any); ok {
			var err error
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
