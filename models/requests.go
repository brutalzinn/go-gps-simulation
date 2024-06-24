package models

type RequestPayload struct {
	PointA Coordinates `json:"pointA"`
	PointB Coordinates `json:"pointB"`
	Device string      `json:"device"`
}
