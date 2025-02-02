package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	//"golang.org/x/text/message"
	//"log"
	//"net/http"
	//"sync"
	//"encoding/json"
)

// Time struct represents a reservation or request time slot for a desk.
type Time struct {
	UserName   string
	StartTime  int
	FinishTime int
}

// Desk struct represents a desk with an ID, features, and reservation/request times.
type Desk struct {
	ID           string
	Features     string
	ReserveTimes []Time
	RequestTimes []Time
}

// Floor struct represents a floor containing desks.
type Floor struct {
	Type   string
	Desks  []*Desk
	NextID int
}

// Service struct represents the overall desk reservation system.
type Service struct {
	Floors          []*Floor
	Features        []int
	SpecialEntrance int
	NextID          int
}

// NewTime creates a new Time instance.
func NewTime(userName string, startTime, finishTime int) Time {
	return Time{
		UserName:   userName,
		StartTime:  startTime,
		FinishTime: finishTime,
	}
}

// NewDesk creates a new Desk instance.
func NewDesk(id string, features string) *Desk {
	return &Desk{
		ID:           id,
		Features:     features,
		ReserveTimes: []Time{},
		RequestTimes: []Time{},
	}
}

// NewFloor creates a new Floor instance.
func NewFloor(floorType string) *Floor {
	return &Floor{
		Type:   floorType,
		Desks:  []*Desk{},
		NextID: 1,
	}
}

// NewService creates a new Service instance with a special entrance cost.
func NewService(specialEntrance int) *Service {
	return &Service{
		Floors:          []*Floor{},
		Features:        []int{},
		SpecialEntrance: specialEntrance,
		NextID:          1,
	}
}

// hasIntersection checks if two time intervals overlap.
func hasIntersection(a, b Time) bool {
	return max(a.StartTime, b.StartTime) < min(a.FinishTime, b.FinishTime)
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CheckHasIntersection checks if a given time conflicts with existing reservations or requests.
func (d *Desk) CheckHasIntersection(time Time) bool {
	for _, requestTime := range d.RequestTimes {
		if hasIntersection(time, requestTime) {
			return false
		}
	}
	for _, reserveTime := range d.ReserveTimes {
		if hasIntersection(time, reserveTime) {
			return false
		}
	}
	return true
}

// GetMessage prints a given message.
func GetMessage(msg string) {
	fmt.Println(msg)
}

// CheckHasFeature checks if a desk has a required feature.
func (d *Desk) CheckHasFeature(feature string) bool {
	if feature == "" {
		return true
	}
	return d.Features == feature
}

// SearchDesksForReserve finds available desks for reservation.
func (s *Service) SearchDesksForReserve(time Time, feature string, numDesks int) []*Desk {
	for _, floor := range s.Floors {
		if floor.Type == "special" {
			var availableDesks []*Desk
			for _, desk := range floor.Desks {
				if desk.CheckHasIntersection(time) && desk.CheckHasFeature(feature) {
					availableDesks = append(availableDesks, desk)
				}
			}
			if len(availableDesks) == numDesks {
				AddReserve(availableDesks, time, s)
				return availableDesks
			}
		}
	}
	if numDesks == 1 {
		GetMessage("No desk available")
	} else {
		GetMessage("Not enough desks available")
	}
	return nil
}

// SearchDesksForRequest finds an available desk for a request.
func (s *Service) SearchDesksForRequest(time Time, deskType string, feature string) *Desk {
	for _, floor := range s.Floors {
		if floor.Type == deskType {
			for _, desk := range floor.Desks {
				if desk.CheckHasIntersection(time) && desk.CheckHasFeature(feature) {
					AddRequest(desk, time, s, deskType)
					return desk
				}
			}
		}
	}
	GetMessage("No desk available")
	return nil
}

// GetFeaturePrice calculates the price based on desk features.
func (d *Desk) GetFeaturePrice(s *Service) int {
	featurePrice := 0
	for i := range d.Features {
		if d.Features[i] == '1' {
			featurePrice += s.Features[i]
		}
	}
	return featurePrice
}

// GetPrice calculates the total price for a set of reserved desks.
func GetPrice(desks []*Desk, s *Service, time Time) int {
	totalPrice := 0
	for _, desk := range desks {
		totalPrice += desk.GetFeaturePrice(s) * (time.FinishTime - time.StartTime)
	}
	return totalPrice
}

// AddReserve reserves desks and calculates the price.
func AddReserve(desks []*Desk, time Time, s *Service) {
	price := GetPrice(desks, s, time)
	price += s.SpecialEntrance * len(desks)
	result := time.UserName + " reserves desks "
	for _, desk := range desks {
		desk.ReserveTimes = append(desk.ReserveTimes, time)
		result += desk.ID + " "
	}	
	result += "for" + strconv.Itoa(price)
	GetMessage(result)
}

// AddRequest adds a request for a desk and calculates the price.
func AddRequest(desk *Desk, time Time, s *Service, deskType string) {
	desk.RequestTimes = append(desk.RequestTimes, time)
	var desks []*Desk
	desks = append(desks, desk)
	price := GetPrice(desks, s, time)
	if deskType == "special" {
		price += s.SpecialEntrance
		GetMessage(time.UserName + " got desk " + desk.ID + " for " + strconv.Itoa(price))
	} else if price != 0{
		GetMessage(time.UserName + " got desk " + desk.ID + " for " + strconv.Itoa(price))
	} else {
		GetMessage(time.UserName + " got desk " + desk.ID)
	}
}

// GetDeskByID returns desk by id
func (s *Service)GetDeskByID(id string) *Desk{
	parts := strings.Split(id, "-")
	floorNumber, err1 := strconv.Atoi(parts[0])
	deskNumber, err2  := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || s.NextID <= floorNumber || s.Floors[floorNumber].NextID >= deskNumber { 
		GetMessage("desk not found")
		return nil
	}
	return s.Floors[floorNumber].Desks[deskNumber]
}

// DeskStatus checks the status of a desk at a given time.
func (s *Service)DeskStatus(time int, id string) {
	desk := s.GetDeskByID(id)
	if desk == nil {
		return
	}
	Next := 10000
	for _, requestTime := range desk.RequestTimes {
		if time >= requestTime.StartTime && time < requestTime.FinishTime {
			GetMessage(requestTime.UserName + " got desk until " + strconv.Itoa(requestTime.FinishTime))
			return
		}
		if requestTime.StartTime > time { 
			Next = min(Next, requestTime.StartTime)
		}
	}
	for _, reserveTime := range desk.ReserveTimes {
		if time >= reserveTime.StartTime && time < reserveTime.FinishTime {
			GetMessage(reserveTime.UserName + " got desk until " + strconv.Itoa(reserveTime.FinishTime))
			return
		}
		if reserveTime.StartTime > time {
			Next = min(Next, reserveTime.StartTime)
		}
	}
	if Next == 10000 {
		GetMessage("desk is available")
		return;
	}
	GetMessage("desk available until " + strconv.Itoa(Next))
}



func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var service *Service
	var featurePrices []int

	// Read the number of features and their prices
	if scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		nFeatures, _ := strconv.Atoi(parts[0])
		if scanner.Scan() {
			featureParts := strings.Split(scanner.Text(), " ")
			for i := 0; i < nFeatures; i++ {
				price, _ := strconv.Atoi(featureParts[i]) 
				featurePrices[i] = price 
			}
			
		}
	}

	// Read the number of floors and special entrance fee
	if scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		nFloors, _ := strconv.Atoi(parts[0])
		specialEntrance, _ := strconv.Atoi(parts[1])
		service = NewService(specialEntrance)
		service.Features = featurePrices
		service.NextID = nFloors
		// Read floors and desks
		for i := 1; i <= nFloors; i++ {
			if scanner.Scan() {
				floorParts := strings.Split(scanner.Text(), " ")
				nDesks, _ := strconv.Atoi(floorParts[0])
				floorType := floorParts[1]
				floor := NewFloor(floorType)
				floor.NextID = nDesks
				if scanner.Scan() {
					featuresList := strings.Split(scanner.Text(), " ")
					for j := 1; j <= nDesks; j++ {
						floor.Desks = append(floor.Desks, NewDesk(strconv.Itoa(i) + "-" + strconv.Itoa(j), featuresList[j-1]))
					}
				}
				service.Floors = append(service.Floors, floor)
			}
		}
	}
	// Process commands
	for scanner.Scan() {
		line := scanner.Text()
		if line == "end" {
			break
		}

		parts := strings.Split(line, " ")
		cmd := parts[1]
		timestamp, _ := strconv.Atoi(parts[0])

		if cmd == "request_desk" {
			username := parts[2]
			deskType := parts[3]
			duration, _ := strconv.Atoi(parts[4])
			time := NewTime(username, timestamp, timestamp+duration)
			service.SearchDesksForRequest(time, deskType, "")

		} else if cmd == "reserve_desk" {
			username := parts[2]
			fromTime, _ := strconv.Atoi(parts[3])
			duration, _ := strconv.Atoi(parts[4])
			featureCode := parts[5]
			time := NewTime(username, fromTime, fromTime+duration)
			service.SearchDesksForReserve(time, featureCode, 1)

		} else if cmd == "reserve_multiple_desks" {
			username := parts[2]
			numDesks, _ := strconv.Atoi(parts[3])
			fromTime, _ := strconv.Atoi(parts[4])
			duration, _ := strconv.Atoi(parts[5])
			time := NewTime(username, fromTime, fromTime+duration)
			service.SearchDesksForReserve(time, "", numDesks)

		} else if cmd == "desk_status" {
			deskID := parts[2]
			service.DeskStatus(timestamp, deskID)
		}
	}
}
