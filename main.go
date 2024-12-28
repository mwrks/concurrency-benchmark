package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
	"strings"
)

type Order struct {
	TicketID  int    `json:"ticket_id"`
	OrderedBy string `json:"ordered_by"`
}

func getRandomName() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	paragraph := "Lorem ipsum odor amet consectetuer adipiscing elit Tortor urna primis maximus habitasse vulputate nisi penatibus Lacus purus metus dapibus tempor dolor suspendisse Nostra pulvinar nostra integer ullamcorper faucibus bibendum Elementum fames nibh ipsum amet porttitor Ullamcorper eu in nostra in amet Aliquet mauris felis tristique rhoncus inceptos arcu Litora amet consequat mus aptent suspendisse metus donec Conubia leo hac eget nibh enim dapibus interdum Laoreet iaculis venenatis vehicula nunc elit aenean donec fusce Erat dui nullam vel elementum viverra nibh non Urna fringilla suspendisse scelerisque iaculis semper neque Consectetur feugiat ac fusce torquent diam senectus volutpat sociosqu In eros non ultricies bibendum nam curabitur vivamus nec Conubia platea ac ac turpis dolor ipsum Facilisi leo tellus purus urna ornare pharetra potenti risus Feugiat gravida ex faucibus nunc congue a ex consequat sed"
	names := strings.Fields(paragraph)
	return names[rnd.Intn(len(names))]
}

func makeOrderRequest(url string, ticketID int) error {
	order := Order{
		TicketID:  ticketID,
		OrderedBy: getRandomName(),
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Utilize all available CPU cores

	url := "http://localhost:3000/order" // Replace with your ticketing app's endpoint

	var wg sync.WaitGroup
	numRequests := 10000 // Increase the number of requests
	ticketID := 1        // Use a specific ticket ID for testing

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Introduce a small random delay to increase the likelihood of concurrent requests
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			err := makeOrderRequest(url, ticketID)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	wg.Wait()
	fmt.Println("All requests sent.")
}
