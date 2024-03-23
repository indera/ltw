package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

func handleCreate(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON payload into a struct
	var newRec RecordInput
	err := json.NewDecoder(r.Body).Decode(&newRec)

	if err != nil {
		slog.Error("bad input", slog.Any("error", err))
		http.Error(w, "Failed to decode JSON input", http.StatusBadRequest)
		return
	}

	storageRec := createRecord(newRec)
	storeRecord(storageRec)
	slog.Info("record saved")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("record saved"))
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	err := deleteRecord(idString)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("record deleted"))
}

// getLabelFilters - transform a string like "a,b,c" => ["a", "b", "c"]
func getLabelFilters(rawFilter string) []Label {
	filters := []Label{}
	rawFilter = strings.ReplaceAll(rawFilter, " ", "")

	if rawFilter == "" {
		return filters
	}

	if strings.Contains(rawFilter, ",") {
		for _, f := range strings.Split(rawFilter, ",") {
			filters = append(filters, Label(f))
		}
	} else {
		filters = append(filters, Label(rawFilter))
	}

	return filters
}

// getOrdering - helper for computing ordering details
// Examples of input for valid output:
//   - ordering=createdAt:asc
//   - ordering=id:desc
func getOrdering(rawOrdering string) (*ListOrdering, error) {
	regexValidSorting := regexp.MustCompile(`(?i).*:(asc|desc)$`)

	if rawOrdering != "" && !regexValidSorting.MatchString(rawOrdering) {
		return nil, errors.New("invalid `ordering` suffix - allowed values are `:asc` and `:desc`")
	}

	parts := strings.Split(rawOrdering, ":")
	if len(parts) != 2 {
		// no ordering specified
		return &ListOrdering{}, nil
	}

	sortIdDirection := ""
	sortCreatedAtDirection := ""

	// correct sorting specified
	if strings.ToLower(parts[0]) == "id" {
		sortIdDirection = parts[1]
	} else if strings.ToLower(parts[0]) == "createdat" {
		sortCreatedAtDirection = parts[1]
	} else {
		slog.Warn("ignore specified ordering", slog.Any("rawOrdering", rawOrdering))
	}

	return &ListOrdering{
		ID:        sortIdDirection,
		CreatedAt: sortCreatedAtDirection,
	}, nil
}

func handleList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	ordering := queryParams.Get("ordering")
	filtering := queryParams.Get("filtering")

	listOrdering, err := getOrdering(ordering)
	if err != nil {
		slog.Error("invalid list request", slog.Any("error", err))
		http.Error(w, "invalid list request", http.StatusBadRequest)
		return
	}

	filters := getLabelFilters(filtering)
	listInput := ListInput{
		Ordering:  *listOrdering,
		Filtering: filters,
	}

	if err != nil {
		slog.Error("list input issues", slog.Any("error", err))
		http.Error(w, "Failed to decode JSON input", http.StatusBadRequest)
		return
	}

	searchResults, err := listRecords(listInput)
	if err != nil {
		http.Error(w, "no results found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(searchResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
