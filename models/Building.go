package models

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/KevinMcIntyre/tu-nav-server/utils"
)

type Building struct {
	ID          int     `json:"id"`
	UID         *string `json:"uid,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ImageRef    *string `json:"imageRef,omitempty"`
	Latitude    float64 `json:"lat,omitempty"`
	Longitude   float64 `json:"long,omitempty"`
	Address     *string `json:"address,omitempty"`
	IndoorID    *string `json:"indoorID,omitempty"`
}

type SeedBuilding struct {
	UID       string
	Name      string
	Latitude  string
	Longitude string
}

func GetBuildings(db *sql.DB) (*[]Building, error) {
	rows, err := db.Query(`
		SELECT
		id,
		uid,
		name,
		description,
		image_ref,
		longitude,
		latitude,
		address,
		indoorid
		FROM buildings ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	var buildings []Building
	defer rows.Close()
	for rows.Next() {
		building := new(Building)
		if err := rows.Scan(
			&building.ID,
			&building.UID,
			&building.Name,
			&building.Description,
			&building.ImageRef,
			&building.Latitude,
			&building.Longitude,
			&building.Address,
			&building.IndoorID,
		); err != nil {
			return nil, err
		}
		buildings = append(buildings, *building)
	}
	return &buildings, nil
}

func parseBuildingsFile(file *os.File) (*[]SeedBuilding, error) {
	r := csv.NewReader(file)
	// Skip the header record
	_, err := r.Read()
	if err != nil {
		return nil, err
	}
	var buildings []SeedBuilding
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		buildings = append(buildings, SeedBuilding{UID: record[0], Name: record[1], Latitude: record[2], Longitude: record[3]})
	}

	return &buildings, nil
}

func SeedBuildingsFile(db *sql.DB, filePath string) error {
	seeded, err := checkIfBuildingsAreSeeded(db)
	if seeded {
		return nil
	}

	buildingsFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	buildings, err := parseBuildingsFile(buildingsFile)
	if err != nil {
		return err
	}

	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	statement, err := transaction.Prepare(`
		INSERT INTO buildings (uid, name, longitude, latitude)
		VALUES ($1, $2, $3, $4);
	`)

	for _, building := range *buildings {
		lat, err := strconv.ParseFloat(building.Latitude, 64)
		if err != nil {
			fmt.Println("Error parsing latitude for: " + building.Name)
			return err
		}

		long, err := strconv.ParseFloat(building.Longitude, 64)
		if err != nil {
			fmt.Println("Error parsing longitude for: " + building.Name)
			return err
		}
		buildingUID := sql.NullString{String: building.UID, Valid: !utils.IsEmptyString(building.UID)}
		buildingName := sql.NullString{String: building.Name, Valid: !utils.IsEmptyString(building.Name)}

		_, err = statement.Exec(buildingUID, buildingName, long, lat)
		if err != nil {
			return err
		}
	}

	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		return err
	}
	return nil
}

func checkIfBuildingsAreSeeded(db *sql.DB) (bool, error) {
	var buildingCount int

	err := db.QueryRow(`
        SELECT COUNT(id) FROM buildings;
    `).Scan(&buildingCount)
	if err != nil {
		return false, err
	}

	if buildingCount > 0 {
		return true, nil
	}

	return false, nil
}
