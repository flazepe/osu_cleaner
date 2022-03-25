package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pborman/getopt/v2"
)

func main() {
	getopt.StringLong("dir", 'd', "", "the osu! directory")
	getopt.Parse()

	dir := getopt.GetValue("dir")

	if dir == "" {
		getopt.Usage()
		os.Exit(1)
	}

	beatmaps, err := os.ReadDir(dir)

	if err != nil {
		fmt.Printf(`Could not open the directory "%s".`, dir)
		os.Exit(1)
	}

	var savedSpace int64

	for _, beatmap := range beatmaps {
		if !beatmap.IsDir() {
			continue
		}

		beatmapDir := filepath.Join(dir, beatmap.Name())
		files, err := os.ReadDir(beatmapDir)

		if err != nil {
			panic(err)
		}

		importantFiles := map[string]bool{}

		// Find important files
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".osu") {
				continue
			}

			osuContent, err := os.ReadFile(filepath.Join(beatmapDir, file.Name()))

			if err != nil {
				panic(err)
			}

			for _, line := range strings.Split(string(osuContent), "\n") {
				line = strings.TrimSpace(line)

				if strings.HasPrefix(line, "AudioFilename: ") {
					importantFiles[line[15:]] = true
				} else if strings.HasPrefix(line, `0,0,"`) && strings.HasSuffix(line, `",0,0`) {
					importantFiles[line[5:len(line)-5]] = true
				}
			}
		}

		// Remove all unimportant files and directories
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".osu") || importantFiles[file.Name()] {
				continue
			}

			fmt.Println(filepath.Join(beatmap.Name(), file.Name()))

			err := os.RemoveAll(filepath.Join(beatmapDir, file.Name()))

			if err != nil {
				panic(err)
			}

			info, err := file.Info()

			if err != nil {
				panic(err)
			}

			savedSpace += info.Size()
		}
	}

	fmt.Printf("Saved %d MB!", savedSpace/1024/1024)
}
