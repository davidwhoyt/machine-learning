package main

import (
	"fmt"
	"github.com/buskersguidetotheuniverse.org/types"
	"flag"
	"log"
	"os"
	"sync"
	"github.com/buskersguidetotheuniverse.org/noaa"
	"github.com/buskersguidetotheuniverse.org/hbase"
)

var wg sync.WaitGroup

// Fetch the weather from a series of NOAA stations and save the results to a local hbase instance.
func main() {

	printWeather := flag.Bool("report", false, "print conditions for each stations to console")
	flag.Parse()

	stations := flag.Args()

	numStations := len(stations)
	if numStations == 0 {
		fmt.Printf("No stations passed in.")
		os.Exit(0)
	}

	log.Printf("%v", os.Args)
	log.Printf("-p: %v\n", *printWeather)
	log.Printf("tail: %v\n", stations)

	for _, station := range stations {
		s := station
		log.Printf("station: %v", s)
		wg.Add(1)
		go func(station string, pw bool) {
			defer wg.Done()
			processStation(station, pw)
			log.Println("processing returned.")
		}(s, *printWeather)
	}

	log.Println("waiting for all threads to return.")
	wg.Wait()
}

func processStation(station string, printWeather bool) {
	log.Println("processing..." + station)
	currentConditions, err := noaa.CurrentConditions(station)
	if err != nil {
		log.Printf("WARN: couldn't look up weather for %v:  %v", station, err)

		// not a fatal error, but obviously we don't want to try to persist the result
		return
	}

	//TODO: Use ErrorGroup or whatever instead.
	err = hbase.SaveObservation(&currentConditions)
	if err != nil {
		log.Fatalf("Error saving weather to HBase: %v", err)
		os.Exit(1)
	}

	if printWeather {
		printCurrentConditions(&currentConditions)
	}

	log.Println("done processing")

}

func printCurrentConditions(currentConditions *types.CurrentConditionsResponse) {
	fmt.Printf("Current Conditions:\n")
	fmt.Println(currentConditions.Props.Station)
	fmt.Println(currentConditions.Props.Timestamp)
	//fmt.Println(currentConditions.Props.BarometricPressure)
	fmt.Println(currentConditions.Props.Temperature)
	fmt.Println(currentConditions.Props.WindSpeed)
	fmt.Println(currentConditions.Props.WindDirection)
	fmt.Println(currentConditions.Props.TextDescription)

}
