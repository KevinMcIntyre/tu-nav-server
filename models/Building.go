package models

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Building struct {
	UID       string `csv:"ID"`
	Name      string `csv:"NAME"`
	Latitude  string `csv:"LATITUDE"`
	Longitude string `csv:"LONGITUDE"`
}

func parseBuildingsFile(file *os.File) (*[]Building, error) {
	r := csv.NewReader(file)
	// Skip the header record
	_, err := r.Read()
	if err != nil {
		return nil, err
	}
	var buildings []Building
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		buildings = append(buildings, Building{UID: record[0], Name: record[1], Latitude: record[2], Longitude: record[3]})
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

	for _, building := range *buildings {
		transaction.Exec(`
            INSERT INTO buildings (uid, name, longitude, latitude)
            VALUES ($1, $2, $3, $4);
        `, building.UID, building.Name, building.Longitude, building.Latitude)
	}

	err = transaction.Commit()
	if err != nil {
		fmt.Println(err)
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

	return true, nil
}
