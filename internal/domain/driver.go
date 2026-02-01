package domain

type Driver struct {
	ID        int64  `json:"id"`
	Name      string `json:"name" binding:"required"`
	VehicleID int64  `json:"vehicle_id" binding:"required"`
	Score     int    `json:"score" binding:"required"`
}

type FullDriver struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Vehicle Vehicle `json:"vehicle"`
	Score   int     `json:"score"`
}

type Vehicle struct {
	ID     int64  `json:"id"`
	Type   string `json:"type"`
	Vendor string `json:"vendor"`
	Model  string `json:"model"`
}

type DriverIdRequest struct {
	ID int64 `json:"id" uri:"id" binding:"required"`
}

type AddDriverResponse struct {
	ID int64 `json:"id"`
}
