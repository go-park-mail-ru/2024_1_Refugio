package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	dir := "cmd"

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)

			go func(file os.DirEntry) {
				defer wg.Done()

				cmdPath := filepath.Join(dir, file.Name(), "main.go")
				cmd := exec.Command("go", "run", cmdPath)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				fmt.Printf("Запуск микросервиса из каталога: %s\n", file.Name())
				if err := cmd.Run(); err != nil {
					fmt.Printf("Ошибка при запуске: %v\n", err)
				}
			}(file)
		}
	}

	wg.Wait()
}
