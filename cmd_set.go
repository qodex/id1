package id1

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (t *Command) set() error {
	if !preflightChecks(t) {
		return fmt.Errorf("failed preflight checks")
	}
	keyPath := filepath.Join(dbpath, t.Key.String())
	dir := filepath.Dir(keyPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(dir, 0770); mkdirErr != nil {
			return mkdirErr
		}
	}
	if err := os.WriteFile(keyPath, t.Data, 0644); err != nil {
		return err
	} else {
		pubsub.Publish(t)
		createDotTtl(t)
		return nil
	}
}

func preflightChecks(cmd *Command) bool {
	if strings.HasPrefix(cmd.Key.Name, ".after.") {
		if dotAfterCmd, err := ParseCommand(cmd.Data); err != nil {
			return false
		} else {
			dotAfterCmd.Args["x-id"] = cmd.Args["x-id"]
			cmd.Data = dotAfterCmd.Bytes()
		}
	}
	return true
}

// .filename.ttl contains .after.timestamp filename, which contains del:/filename command
func createDotTtl(cmd *Command) {
	ttlSec, _ := strconv.Atoi(cmd.Args["ttl"])
	if ttlSec == 0 {
		return
	}
	ttdMs := time.Now().UnixMilli() + (int64(ttlSec) * 1000) //time to die in Ms
	ttlKey := KK(cmd.Key.Parent, fmt.Sprintf(".ttl.%s", cmd.Key.Name))
	dotAfterKey := KK(cmd.Key.Parent, fmt.Sprintf(".after.%d", ttdMs))

	if oldDotAfter, err := CmdGet(ttlKey).Exec(); err == nil {
		CmdDel(K(string(oldDotAfter))).Exec()
	}

	dotAfterCommand := CmdDel(cmd.Key)
	dotAfterCommand.Args["x-id"] = cmd.Args["x-id"]
	CmdSet(ttlKey, map[string]string{"x-id": cmd.Args["x-id"]}, []byte(dotAfterKey.String())).Exec()
	CmdSet(dotAfterKey, map[string]string{"x-id": cmd.Args["x-id"]}, dotAfterCommand.Bytes()).Exec()
}
