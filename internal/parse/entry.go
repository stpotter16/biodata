package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	var weight types.Weight
	if body.Weight == "" {
		log.Printf("No weight field in payload")
		weight.Value = nil
	} else {
		weightValue, err := strconv.ParseFloat(body.Weight, 64)
		if err != nil {
			log.Printf("Could not parse payload weight field: %v", err)
			return types.Entry{}, err
		}
		weight = types.NewWeight(weightValue)
	}

	var waist types.Waist
	if body.Waist == "" {
		log.Printf("No waist field in payload")
		waist.Value = nil
	} else {
		waistValue, err := strconv.ParseFloat(body.Waist, 64)
		if err != nil {
			log.Printf("Could not parse payload waist field: %v", err)
			return types.Entry{}, err
		}
		waist = types.NewWaist(waistValue)
	}

	var bp types.BP
	if body.BP == "" {
		log.Printf("No BP field in payload")
		bp.Systolic = nil
		bp.Diastolic = nil
	} else {
		bp, err = parseBPString(body.BP)
		if err != nil {
			return types.Entry{}, err
		}
	}

	entry := types.Entry{
		Date:   date,
		Weight: weight,
		Waist:  waist,
		BP:     bp,
	}

	return entry, nil
}

func ParseEntryPut(r *http.Request) (types.Entry, error) {
	// TODO - get date from url
	dateStr := r.PathValue("date")
	entryDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Could not parse entry date %s", dateStr)
		return types.Entry{}, nil
	}

	body := struct {
		Weight string `json:"weight"`
		Waist  string `json:"waist"`
		BP     string `json:"bp"`
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Printf("Invalid new entry request: %v", err)
		return types.Entry{}, err
	}

	var weight types.Weight
	if body.Weight == "" {
		log.Printf("No weight field in payload")
		weight.Value = nil
	} else {
		weightValue, err := strconv.ParseFloat(body.Weight, 64)
		if err != nil {
			log.Printf("Could not parse payload weight field: %v", err)
			return types.Entry{}, err
		}
		weight = types.NewWeight(weightValue)
	}

	var waist types.Waist
	if body.Waist == "" {
		log.Printf("No waist field in payload")
		waist.Value = nil
	} else {
		waistValue, err := strconv.ParseFloat(body.Waist, 64)
		if err != nil {
			log.Printf("Could not parse payload waist field: %v", err)
			return types.Entry{}, err
		}
		waist = types.NewWaist(waistValue)
	}

	var bp types.BP
	if body.BP == "" {
		log.Printf("No BP field in payload")
		bp.Systolic = nil
		bp.Diastolic = nil
	} else {
		var err error
		bp, err = parseBPString(body.BP)
		if err != nil {
			return types.Entry{}, err
		}
	}

	entry := types.Entry{
		Date:   entryDate,
		Weight: weight,
		Waist:  waist,
		BP:     bp,
	}

	return entry, nil
}

func ParseEntryDTO(entryDTO types.EntryDTO) (types.Entry, error) {
	var weight types.Weight
	if entryDTO.Weight.Valid {
		weight.Value = &entryDTO.Weight.Float64
	} else {
		weight.Value = nil
	}

	var waist types.Waist
	if entryDTO.Waist.Valid {
		waist.Value = &entryDTO.Waist.Float64
	} else {
		waist.Value = nil
	}

	var bp types.BP
	if entryDTO.Bp.Valid {
		var err error
		bp, err = parseBPString(entryDTO.Bp.String)
		if err != nil {
			bp.Systolic = nil
			bp.Diastolic = nil
		}
	} else {
		bp.Systolic = nil
		bp.Diastolic = nil
	}
	entry := types.Entry{
		Date:   entryDTO.Date,
		Weight: weight,
		Waist:  waist,
		BP:     bp,
	}
	return entry, nil
}

func parseBPString(bps string) (types.BP, error) {
	parts := strings.Split(bps, "/")
	if len(parts) != 2 {
		log.Printf("Invalid BP string received")
		return types.BP{}, fmt.Errorf("Expected string with systolic and diastolic. Received %s", bps)
	}
	systolicFloat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		log.Printf("Invalid systolic portion of BP string")
		return types.BP{}, err
	}
	diastolicFloat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		log.Printf("Invalid diastolic portion of BP string")
		return types.BP{}, err
	}
	return types.BP{
		Systolic:  &systolicFloat,
		Diastolic: &diastolicFloat,
	}, nil
}
