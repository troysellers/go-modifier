package main

import (
	"encoding/csv"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"
)

func TestMod(t *testing.T) {

	c := 20000
	c1 := c / 5000
	m := int(math.Mod(float64(c), 5000))

	for i := 1; i <= c1; i++ {
		log.Printf("get %d", 5000*i)
	}
	log.Printf("get the remaining %d", m)

}

func TestNothing(t *testing.T) {
	names := getData("/tmp/mockaroo-data/names.csv")
	log.Printf("%d names\n", len(names))

	leads := getData("/tmp/mockaroo-data/Lead-query-modified.csv")
	log.Printf("%d leads\n", len(leads))

	var newD [][]string
	newD = append(newD, leads[0])
	for _, row := range leads[1:] {
		row[1] = getRandom(0, names[1:])
		row[2] = getRandom(1, names[1:])
		row[3] = getRandom(2, names[1:])
		newD = append(newD, row)
	}
	write("/tmp/mockaroo-data/lead-update.csv", newD)
}

func write(f string, d [][]string) {

	file, err := os.Create(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.WriteAll(d)
	if err != nil {
		panic(err)
	}
}
func getRandom(col int, d [][]string) string {

	return d[rand.Intn(len(d))][col]
}

func getData(n string) [][]string {
	var leadData [][]string

	f, err := os.Open(n)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	leadData, _ = csv.NewReader(f).ReadAll()

	return leadData

}
