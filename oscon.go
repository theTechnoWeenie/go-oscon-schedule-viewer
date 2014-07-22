package main
import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"sort"
)

type Event struct {
	Serial int `json:"serial"`
	Name string `json:"name"`
	EventType string `json:"event_type"`
	Start string `json:"time_start"`
	Stop string `json:"time_stop"`
	Venue int `json:"venue_serial"`
	Description string `json:"description"`
	WebSite string `json:"website_url"`
	Speakers []int `json:"speakers"`
	Categories []string `json:"categories"`
}

type Venue struct {
	Serial int `json:"serial"`
	Name string `json:"name"`
	Category string `json:"category"`
}

type SortableEvents []Event
func (a SortableEvents) Len() int { return len(a) }
func (a SortableEvents) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortableEvents) Less(i, j int) bool {
	one, _ := time.Parse("2006-01-02 15:04:05", a[i].Start)
	two, _ := time.Parse("2006-01-02 15:04:05", a[j].Start)
	return one.Before(two)
}

func main(){
	fmt.Printf("Retrieving schedule from OSCON Site...\n")
	resp, err := http.Get("http://www.oreilly.com/pub/sc/osconfeed")
	defer resp.Body.Close()
	if(err == nil){
		body, parseErr := ioutil.ReadAll(resp.Body)
		if(parseErr == nil){
			var rootNode map[string]*json.RawMessage
			json.Unmarshal(body, &rootNode)

			var schedule map[string]*json.RawMessage
			json.Unmarshal(*rootNode["Schedule"], &schedule)

			var venues []Venue
			json.Unmarshal(*schedule["venues"], &venues)

			var events []Event
			json.Unmarshal(*schedule["events"], &events)

			filtered := filterAndSortEvents(events)
			now := time.Now()
			loc, _ := time.LoadLocation("Local")
			for _, event := range filtered {
				start, _ := time.ParseInLocation("2006-01-02 15:04:05", event.Start, loc)
				end , _ := time.ParseInLocation("2006-01-02 15:04:05", event.Stop, loc)

				fmt.Println("======")
				if(start.Before(now) && end.After(now)){
					fmt.Println("CURRENTLY RUNNING")
				}
				fmt.Println("Name:", event.Name)
				fmt.Println("Starting:", event.Start)
				fmt.Println("Ends:", event.Stop)
				fmt.Println("Happening in:", lookUpVenue(venues, event.Venue))
			}

		}else{
			fmt.Printf("There was an error parsing the body")
		}
	}else{
		fmt.Printf("There was an error retrieving the req")
	}
}

func lookUpVenue(venues []Venue, desiredVenue int) string {
	for _, venue := range venues {
		if(venue.Serial == desiredVenue){
			return venue.Name
		}
	}
	return "Unknown"
}

func filterAndSortEvents(events []Event) []Event {
	mySessionSerials := []int {33947,33627,34252,34187,34332,34137,34371,34094,34921,34588,33590,34535,34757,34422,34575,35038,34434}
	var filtered []Event
	for _, event := range events {
		for _, myEventId := range mySessionSerials {
			if(myEventId == event.Serial){
				filtered = append(filtered, event)
			}
		}
	}
	sort.Sort(SortableEvents(filtered))
	return filtered
}
