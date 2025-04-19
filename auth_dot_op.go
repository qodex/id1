package main

import (
	"fmt"
	"strings"
)

func authDotOp(id string, cmd Command) bool {
	if len(id) == 0 {
		return false
	}
	roles := getRoles(id, cmd.Key)
	parentKey := K(cmd.Key.Parent)
	for len(parentKey.String()) != 0 {
		dotOpKey := KK(parentKey, fmt.Sprintf(".%s", cmd.Op))
		if data, err := CmdGet(dotOpKey).Exec(); err == nil {
			lines := strings.Split(string(data), "\n")
			if ContainsAny(lines, roles) {
				return true
			}
		}
		parentKey = K(parentKey.Parent)
	}

	return false
}

func getRoles(id string, key Id1Key) []string {
	result := []string{"*", id}
	parentKey := K(key.Parent)
	for {
		rolesKey := KK(parentKey, ".roles", id)
		if data, err := CmdGet(rolesKey).Exec(); err == nil {
			lines := strings.Split(string(data), "\n")
			result = append(result, lines...)
		}
		if len(parentKey.String()) > 0 {
			parentKey = K(parentKey.Parent)
		} else {
			break
		}
	}
	return result
}

func ContainsAny(slice1, slice2 []string) bool {
	itemSet := make(map[string]bool)
	for _, item := range slice2 {
		itemSet[item] = true
	}
	for _, item := range slice1 {
		if itemSet[item] {
			return true
		}
	}
	return false
}
