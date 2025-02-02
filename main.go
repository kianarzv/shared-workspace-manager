package main

import (
	"fmt"
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
	ID           int
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
func NewDesk(id int, features string) *Desk {
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
		result += strconv.Itoa(desk.ID) + " "
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
		GetMessage(time.UserName + " got desk " + strconv.Itoa(desk.ID) + " for " + strconv.Itoa(price))
	} else if price != 0{
		GetMessage(time.UserName + " got desk " + strconv.Itoa(desk.ID) + " for " + strconv.Itoa(price))
	} else {
		GetMessage(time.UserName + " got desk " + strconv.Itoa(desk.ID))
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
}
