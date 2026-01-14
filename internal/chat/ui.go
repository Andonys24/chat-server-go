package chat

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CleanConsole() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// /c ejecuta el comando y luego termina
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func GenerateTitle(title string, clean bool) {
	if clean {
		if err := CleanConsole(); err != nil {
			fmt.Println("Error al limpiar la consola")
		}
	}

	asterick := strings.Repeat("*", len(title)*3)
	spaces := strings.Repeat(" ", len(title)-1)

	fmt.Printf("\n%s\n", asterick)
	fmt.Printf("*%s%s%s*\n", spaces, title, spaces)
	fmt.Printf("%s\n\n", asterick)

}
