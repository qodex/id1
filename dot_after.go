package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*
scans a folder for .after.<time> files, if time is after now, reads end executes command inside, deletes the file.
*
*/
func dotAfter(dir string) {
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		isDotAfterFile := strings.HasPrefix(info.Name(), ".after.")
		if !isDotAfterFile {
			return nil
		}
		timestampMS, _ := strconv.Atoi(strings.Split(info.Name(), ".")[2])
		timestampIsPast := timestampMS > 0 && time.Now().UnixMilli() > int64(timestampMS)
		if !timestampIsPast {
			return nil
		}
		dotAfterContent, _ := os.ReadFile(path)
		dotAfterCommand, parseError := ParseCommand(dotAfterContent)
		if parseError == nil && auth(dotAfterCommand.Args["x-id"], dotAfterCommand) {
			dotAfterCommand.Exec()
		} else {
			log.Printf("unauthorised .after command by '%s': %s %s", dotAfterCommand.Args["x-id"], dotAfterCommand.Op, dotAfterCommand.Key)
		}
		os.Remove(path)
		time.Sleep(time.Millisecond)
		return nil
	})
}
