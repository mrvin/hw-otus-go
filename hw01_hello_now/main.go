// A program that prints the current local time
// and the exact time obtained using the NTP library
// (github.com/beevik/ntp) in the format:
// 		current time: <time>
// 		exact time: <time>
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	var srvNTP string // ntp server name.

	flag.StringVar(&srvNTP, "n", "2.ru.pool.ntp.org", "ntp server name")
	flag.Parse()
	log.Printf("server name: %s\n", srvNTP)

	currentTime := time.Now().Round(0)
	fmt.Printf("current time: %v\n", currentTime)

	exactTime, err := ntp.Time(srvNTP)
	if err != nil {
		log.Fatalf("ntp.Time: %s", err)
	}
	exactTime = exactTime.Round(0)
	fmt.Printf("exact time: %s\n", exactTime)
}
