package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
    "log"

	"github.com/mattn/go-runewidth"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
}

func AskForConfirmation() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	switch strings.ToLower(response) {
	case "y\n", "yes\n", "\n":
		return true, nil
	}
	return false, nil
}

func PrintCatMessage(msg string) {
	width := 9 // min width
	msgWidth := runewidth.StringWidth(msg)
	if msgWidth > width {
		width = msgWidth
	}

	msg += strings.Repeat(" ", width-runewidth.StringWidth(msg))
	message := fmt.Sprintf(" %s\n", strings.Repeat("-", width+2))
	message += fmt.Sprintf("< %s >\n", msg)
	message += fmt.Sprintf(" %s\n", strings.Repeat("-", width+2))

	fmt.Print(message)
	fmt.Println(`  (\__/)||`)
	fmt.Println(`_ (•ㅅ•)||`)
	fmt.Println(` \/o   \っ`)
}
