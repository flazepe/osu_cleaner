package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := strings.Join(os.Args[1:], " ")

	if dir == "" {
		fmt.Print("Please provide your osu! beatmaps directory.")
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
					importantFiles[strings.ToLower(line[15:])] = true
				} else if strings.HasPrefix(line, "0,0,") {
					filename := strings.Split(line, ",")[2]

					if strings.HasPrefix(filename, `"`) {
						filename = filename[1 : len(filename)-1]
					}

					importantFiles[strings.ToLower(filename)] = true
				}
			}
		}

		// Remove all unimportant files and directories
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".osu") || importantFiles[strings.ToLower(file.Name())] {
				continue
			}

			fmt.Println(filepath.Join(beatmap.Name(), file.Name()))

			if err := os.RemoveAll(filepath.Join(beatmapDir, file.Name())); err != nil {
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
