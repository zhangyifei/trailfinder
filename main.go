package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"trailfinder/csvparser"
	"trailfinder/filter"

	"github.com/spf13/cobra"
)

// Embed the trailheads.csv file
//
//go:embed trailheads.csv
var embeddedCSVFile []byte

var (
	csvFile string
	port    string
	server  bool // Flag to indicate if server mode is enabled
)

func main() {
	var (
		restrooms, picnic, fishing, fee, bikerack, biketrail, dogtube, grills          string
		adatoilet, adafishing, adacamping, adapicnic, adatrail, adaparking, adafacilit string
		horsetrail, recyclebin, dogcompost, accessname, thleash                        string
	)

	filters := make(map[string]string)

	var rootCmd = &cobra.Command{
		Use:   "trailfinder",
		Short: "TrailFinder helps you find trails based on various criteria",
		Run: func(cmd *cobra.Command, args []string) {
			if server {
				startServer()
			} else {
				runCLI(filters, restrooms, picnic, fishing, fee, bikerack, biketrail, dogtube, grills,
					adatoilet, adafishing, adacamping, adapicnic, adatrail, adaparking, adafacilit,
					horsetrail, recyclebin, dogcompost, accessname, thleash)
			}
		},
	}

	// CLI flags
	rootCmd.Flags().BoolVarP(&server, "server", "s", false, "Run in server mode")
	rootCmd.Flags().StringVar(&restrooms, "restrooms", "", "Filter by restrooms availability")
	rootCmd.Flags().StringVar(&picnic, "picnic", "", "Filter by picnic areas")
	rootCmd.Flags().StringVar(&fishing, "fishing", "", "Filter by fishing availability")
	rootCmd.Flags().StringVar(&fee, "fee", "", "Filter by fee availability")
	rootCmd.Flags().StringVar(&bikerack, "bikerack", "", "Filter by bike rack availability")
	rootCmd.Flags().StringVar(&biketrail, "biketrail", "", "Filter by bike trail availability")
	rootCmd.Flags().StringVar(&dogtube, "dogtube", "", "Filter by dog tube availability")
	rootCmd.Flags().StringVar(&grills, "grills", "", "Filter by grill availability")
	rootCmd.Flags().StringVar(&adatoilet, "adatoilet", "", "Filter by ADA toilet availability")
	rootCmd.Flags().StringVar(&adafishing, "adafishing", "", "Filter by ADA fishing availability")
	rootCmd.Flags().StringVar(&adacamping, "adacamping", "", "Filter by ADA camping availability")
	rootCmd.Flags().StringVar(&adapicnic, "adapicnic", "", "Filter by ADA picnic availability")
	rootCmd.Flags().StringVar(&adatrail, "adatrail", "", "Filter by ADA trail availability")
	rootCmd.Flags().StringVar(&adaparking, "adaparking", "", "Filter by ADA parking availability")
	rootCmd.Flags().StringVar(&adafacilit, "adafacilit", "", "Filter by ADA facility availability")
	rootCmd.Flags().StringVar(&horsetrail, "horsetrail", "", "Filter by horse trail availability")
	rootCmd.Flags().StringVar(&recyclebin, "recyclebin", "", "Filter by recycle bin availability")
	rootCmd.Flags().StringVar(&dogcompost, "dogcompost", "", "Filter by dog compost availability")
	rootCmd.Flags().StringVar(&accessname, "accessname", "", "Filter by access name")
	rootCmd.Flags().StringVar(&thleash, "thleash", "", "Filter by TH leash availability")
	rootCmd.Flags().StringVarP(&csvFile, "csv", "c", "", "Path to the trailheads CSV file (default: embedded trailheads.csv)")
	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port for the HTTP server (default: 8080)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCLI(filters map[string]string, restrooms, picnic, fishing, fee, bikerack, biketrail, dogtube, grills,
	adatoilet, adafishing, adacamping, adapicnic, adatrail, adaparking, adafacilit, horsetrail, recyclebin,
	dogcompost, accessname, thleash string) {
	// Populate filters map
	filters["RESTROOMS"] = restrooms
	filters["PICNIC"] = picnic
	filters["FISHING"] = fishing
	filters["Fee"] = fee
	filters["BikeRack"] = bikerack
	filters["BikeTrail"] = biketrail
	filters["DogTube"] = dogtube
	filters["Grills"] = grills
	filters["ADAtoilet"] = adatoilet
	filters["ADAfishing"] = adafishing
	filters["ADAcamping"] = adacamping
	filters["ADApicnic"] = adapicnic
	filters["ADAtrail"] = adatrail
	filters["ADAparking"] = adaparking
	filters["ADAfacilit"] = adafacilit
	filters["HorseTrail"] = horsetrail
	filters["RecycleBin"] = recyclebin
	filters["DogCompost"] = dogcompost
	filters["AccessName"] = accessname
	filters["THLeash"] = thleash

	reader, err := getCSVReader()
	if err != nil {
		fmt.Println(err)
		return
	}

	trailChan := make(chan []filter.Trail)
	doneChan := make(chan error)

	go csvparser.ParseCSVConcurrently(reader, 100, trailChan, doneChan)

	var filteredTrails []filter.Trail
	for {
		select {
		case trails, ok := <-trailChan:
			if !ok {
				trailChan = nil
			} else {
				filteredChunk := filter.FilterTrailsParallel(trails, filters)
				filteredTrails = append(filteredTrails, filteredChunk...)
			}
		case err := <-doneChan:
			if err != nil {
				fmt.Println(err)
				return
			}
			doneChan = nil
		}

		if trailChan == nil && doneChan == nil {
			break
		}
	}

	if len(filteredTrails) == 0 {
		fmt.Println("No trails found matching the criteria.")
	} else {
		fmt.Println("Trails matching the criteria:")
		for _, trail := range filteredTrails {
			fmt.Printf("AccessName: %s\n", trail.AccessName)
			fmt.Printf("RESTROOMS: %s\n", trail.RESTROOMS)
			fmt.Printf("PICNIC: %s\n", trail.PICNIC)
			fmt.Printf("FISHING: %s\n", trail.FISHING)
			fmt.Printf("Fee: %s\n", trail.Fee)
			fmt.Printf("BikeRack: %s\n", trail.BikeRack)
			fmt.Printf("BikeTrail: %s\n", trail.BikeTrail)
			fmt.Printf("DogTube: %s\n", trail.DogTube)
			fmt.Printf("Grills: %s\n", trail.Grills)
			fmt.Printf("ADAtoilet: %s\n", trail.ADAtoilet)
			fmt.Printf("ADAfishing: %s\n", trail.ADAfishing)
			fmt.Printf("ADAcamping: %s\n", trail.ADAcamping)
			fmt.Printf("ADApicnic: %s\n", trail.ADApicnic)
			fmt.Printf("ADAtrail: %s\n", trail.ADAtrail)
			fmt.Printf("ADAparking: %s\n", trail.ADAparking)
			fmt.Printf("ADAfacilit: %s\n", trail.ADAfacilit)
			fmt.Printf("HorseTrail: %s\n", trail.HorseTrail)
			fmt.Printf("RecycleBin: %s\n", trail.RecycleBin)
			fmt.Printf("DogCompost: %s\n", trail.DogCompost)
			fmt.Printf("THLeash: %s\n", trail.THLeash)
			fmt.Println("----------------------------")
		}
	}
}

func startServer() {
	http.HandleFunc("/trails", trailsHandler)
	fmt.Println("Server is running on port", port)
	http.ListenAndServe(":"+port, nil)
}

func trailsHandler(w http.ResponseWriter, r *http.Request) {
	filters := make(map[string]string)

	// Populate filters from query parameters
	filters["RESTROOMS"] = r.URL.Query().Get("restrooms")
	filters["PICNIC"] = r.URL.Query().Get("picnic")
	filters["FISHING"] = r.URL.Query().Get("fishing")
	filters["Fee"] = r.URL.Query().Get("fee")
	filters["BikeRack"] = r.URL.Query().Get("bikerack")
	filters["BikeTrail"] = r.URL.Query().Get("biketrail")
	filters["DogTube"] = r.URL.Query().Get("dogtube")
	filters["Grills"] = r.URL.Query().Get("grills")
	filters["ADAtoilet"] = r.URL.Query().Get("adatoilet")
	filters["ADAfishing"] = r.URL.Query().Get("adafishing")
	filters["ADAcamping"] = r.URL.Query().Get("adacamping")
	filters["ADApicnic"] = r.URL.Query().Get("adapicnic")
	filters["ADAtrail"] = r.URL.Query().Get("adatrail")
	filters["ADAparking"] = r.URL.Query().Get("adaparking")
	filters["ADAfacilit"] = r.URL.Query().Get("adafacilit")
	filters["HorseTrail"] = r.URL.Query().Get("horsetrail")
	filters["RecycleBin"] = r.URL.Query().Get("recyclebin")
	filters["DogCompost"] = r.URL.Query().Get("dogcompost")
	filters["AccessName"] = r.URL.Query().Get("accessname")
	filters["THLeash"] = r.URL.Query().Get("thleash")

	reader, err := getCSVReader()
	if err != nil {
		http.Error(w, "Error reading CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}

	trailChan := make(chan []filter.Trail)
	doneChan := make(chan error)

	go csvparser.ParseCSVConcurrently(reader, 100, trailChan, doneChan)

	var filteredTrails []filter.Trail
	for {
		select {
		case trails, ok := <-trailChan:
			if !ok {
				trailChan = nil
			} else {
				fmt.Println(filters)
				filteredChunk := filter.FilterTrailsParallel(trails, filters)
				filteredTrails = append(filteredTrails, filteredChunk...)
			}
		case err := <-doneChan:
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			doneChan = nil
		}

		if trailChan == nil && doneChan == nil {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if len(filteredTrails) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("No trails found matching the criteria.")
	} else {
		json.NewEncoder(w).Encode(filteredTrails)
	}
}

func getCSVReader() (io.Reader, error) {
	if csvFile == "" {
		// Return the embedded CSV file as a reader
		return io.Reader(bytes.NewReader(embeddedCSVFile)), nil
	} else if isURL(csvFile) {
		resp, err := http.Get(csvFile)
		if err != nil {
			return nil, fmt.Errorf("error fetching CSV: %w", err)
		}
		return resp.Body, nil
	} else {
		file, err := os.Open(csvFile)
		if err != nil {
			return nil, fmt.Errorf("error opening CSV file: %w", err)
		}
		return file, nil
	}
}

func isURL(str string) bool {
	return len(str) > 4 && (str[:4] == "http" || str[:5] == "https")
}
