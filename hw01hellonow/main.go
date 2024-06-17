// That program prints the current local time
// and the exact time obtained using the NTP library
// (github.com/beevik/ntp) in the format:
//
//	current time: <time>
//	exact time: <time>
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

var out io.Writer = os.Stdout // modified during testing.

var getTimeNTP = ntp.Time

var getTimeNow = time.Now

func main() {
	var srvNTP string // ntp server name.

	flag.StringVar(&srvNTP, "n", "2.ru.pool.ntp.org", "ntp server name")
	flag.Parse()
	log.Printf("server name: %s\n", srvNTP)

	currentTime := getTimeNow().Round(0)

	exactTime, err := getTimeNTP(srvNTP)
	if err != nil {
		log.Fatalf("getTimeNTP: %s", err)
	}
	exactTime = exactTime.Round(0)

	fmt.Fprintf(out, "current time: %v\n", currentTime)
	fmt.Fprintf(out, "exact time: %s\n", exactTime)
}
