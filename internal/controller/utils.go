package controller

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
)

func createFile(filepath string, data []byte, perm int, owner string) error {
	dir := path.Dir(filepath)
	if err := os.MkdirAll(dir, fs.FileMode(perm)); err != nil {
		return err
	}

	err := os.WriteFile(filepath, data, fs.FileMode(perm))
	if err != nil {
		return err
	}

	ownerList := strings.Split(owner, ":")
	if len(ownerList) != 2 {
		err := fmt.Errorf("the owner <%s> is not correct, it must satisfy the mask '$userName:$groupName'", owner)
		return err
	}

	group, err := user.LookupGroup(ownerList[1])
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	user, err := user.Lookup(ownerList[0])
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	userUid, err := strconv.Atoi(user.Uid)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	groupUid, err := strconv.Atoi(group.Gid)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	err = os.Chown(filepath, userUid, groupUid)
	return err
}

func getMap(i any) (map[string]any, error) {
	rArgs := make(map[string]any)
	err := fmt.Errorf("args converting is not available")
	args, ok := i.(map[any]any)
	if !ok {
		return rArgs, err
	}
	for k, v := range args {
		key := strings.ToLower(fmt.Sprint(k))
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
