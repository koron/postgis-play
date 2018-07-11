package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/paulmach/osm"
)

func pstr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func run(r io.Reader, dst string) error {
	d := xml.NewDecoder(r)
	m := osm.OSM{}
	err := d.Decode(&m)
	if err != nil {
		return err
	}
	fmt.Println(len(m.Nodes))

	db, err := sql.Open("postgres", dst)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS postgis`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS osm_nodes (
		id    BIGINT NOT NULL PRIMARY KEY,
		lat   DOUBLE PRECISION NOT NULL,
		lon   DOUBLE PRECISION NOT NULL,
		shop  TEXT,
		name  TEXT,
		brand TEXT
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS osm_nodes_brand_idx ON osm_nodes (brand)`)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, n := range m.Nodes {
		shop := n.Tags.Find(`shop`)
		name := n.Tags.Find(`name`)
		brand := n.Tags.Find(`brand`)
		_, err := tx.Exec(
			`INSERT INTO osm_nodes (id, lat, lon, shop, name, brand)
				VALUES($1, $2, $3, $4, $5, $6)
				ON CONFLICT ON CONSTRAINT osm_nodes_pkey
				DO UPDATE SET lat=$2, lon=$3, shop=$4, name=$5, brand=$6`,
			n.ID, n.Lat, n.Lon, pstr(shop), pstr(name), pstr(brand))
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func main() {
	err := run(os.Stdin, `postgres://postgres@127.0.0.1:5432/postgres?sslmode=disable`)
	if err != nil {
		log.Fatal(err)
	}
}
