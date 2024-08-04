package colorful

import (
	"fmt"
	"testing"
)

func TestSetColor(t *testing.T) {
	fmt.Println(Blue("hello"))
	fmt.Println(Red("hello"))
	fmt.Println(Green("hello"))
	fmt.Println(Yellow("hello"))
	fmt.Println(Magenta("hello"))
	fmt.Println(Cyan("hello"))
	fmt.Println(White("hello"))
}
