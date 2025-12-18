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
		log.Printf("Invalid new entry request %+v: %v", body, err)
		return types.Entry{}, err
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		log.Printf("Could not parse payload date field %s: %v", body.Date, err)
		return types.Entry{}, err
	}

	weight, err := parseWeightField(body.Weight)
	if err != nil {
		return types.Entry{}, err
	}

	waist, err := parseWaistField(body.Waist)
	if err != nil {
		return types.Entry{}, err
	}

	bp, err := parseBPField(body.BP)
	if err != nil {
		return types.Entry{}, err
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
	dateStr := r.PathValue("date")
	entryDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Could not parse entry date %s: %v", dateStr, err)
		return types.Entry{}, nil
	}

	body := struct {
		Weight string `json:"weight"`
		Waist  string `json:"waist"`
		BP     string `json:"bp"`
	}{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Printf("Invalid new entry request %+v: %v", body, err)
		return types.Entry{}, err
	}

	weight, err := parseWeightField(body.Weight)
	if err != nil {
		return types.Entry{}, err
	}

	waist, err := parseWaistField(body.Waist)
	if err != nil {
		return types.Entry{}, err
	}

	bp, err := parseBPField(body.BP)
	if err != nil {
		return types.Entry{}, err
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
		log.Printf("Invalid BP string received. Got %s", bps)
		return types.BP{}, fmt.Errorf("expected string with systolic and diastolic. Received %s", bps)
	}
	systolicFloat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		log.Printf("Invalid systolic portion of BP string %s: %v", parts[0], err)
		return types.BP{}, err
	}
	diastolicFloat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		log.Printf("Invalid diastolic portion of BP string %s: %v", parts[1], err)
		return types.BP{}, err
	}
	return types.BP{
		Systolic:  &systolicFloat,
		Diastolic: &diastolicFloat,
	}, nil
}

func parseWeightField(weightStr string) (types.Weight, error) {
	var weight types.Weight
	if weightStr == "" {
		log.Printf("No weight field in payload")
		weight.Value = nil
		return weight, nil
	}

	weightValue, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		log.Printf("Could not parse payload weight field %s: %v", weightStr, err)
		return types.Weight{}, err
	}
	return types.NewWeight(weightValue), nil
}

func parseWaistField(waistStr string) (types.Waist, error) {
	var waist types.Waist
	if waistStr == "" {
		log.Printf("No waist field in payload")
		waist.Value = nil
		return waist, nil
	}

	waistValue, err := strconv.ParseFloat(waistStr, 64)
	if err != nil {
		log.Printf("Could not parse payload waist field %s: %v", waistStr, err)
		return types.Waist{}, err
	}
	return types.NewWaist(waistValue), nil
}

func parseBPField(bpStr string) (types.BP, error) {
	if bpStr == "" {
		log.Printf("No BP field in payload")
		return types.BP{Systolic: nil, Diastolic: nil}, nil
	}
	return parseBPString(bpStr)
}
