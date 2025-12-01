package parse

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/stpotter16/biodata/internal/types"
)

func ParseEntryPost(r *http.Request) (types.Entry, error) {
	body := struct {
		Date   string `json:"date"`
		Weight string `json:"weight"`
		Waist  string `json:"waist"`
		BP     string `json:"bp"`
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Printf("Invalid new entry request: %v", err)
		return types.Entry{}, err
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		log.Printf("Could not parse payload date field: %v", err)
		return types.Entry{}, err
	}

	weight, err := strconv.ParseFloat(body.Weight, 64)
	if err != nil {
		log.Printf("Could not parse payload weight field: %v", err)
		return types.Entry{}, err
	}

	waist, err := strconv.ParseFloat(body.Waist, 64)
	if err != nil {
		log.Printf("Could not parse payload waist field: %v", err)
		return types.Entry{}, err
	}

	entry := types.Entry{
		Date:   date,
		Weight: weight,
		Waist:  waist,
		BP:     body.BP,
	}

	return entry, nil
}
