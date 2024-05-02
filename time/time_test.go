package my_time

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeLocalFormat(t *testing.T) {
	format := "2006-01-02 15:04:05"
	nowtime := time.Now().UTC()

	// UTC
	fmt.Printf("UTC时间：\t\t%s\n", nowtime.Format(format))

	// Asia/Chongqing
	fmt.Printf("Chongqing时间：\t%s\n", TimeLocalFormat())

	// America/New_York
	timelocal, _ := time.LoadLocation("America/New_York")
	fmt.Printf("New_York时间：\t%s\n", nowtime.In(timelocal).Format(format))
	fmt.Println("")

}
