package main

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"sort"
	"strings"
	"time"
)

type Label string

// RecordInput - Parse the request input into this object (no createdAt)
type RecordInput struct {
	ID     string  `json:"id"`
	Labels []Label `json:"labels"`
	Object struct {
		Tag string `json:"tag"`
		Url string `json:"url"`
	} `json:"object"`
}

// Record - The struct we keep in the storage
type Record struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Labels    []Label   `json:"labels"`
	Object    struct {
		Tag string `json:"tag"`
		Url string `json:"url"`
	} `json:"object"`
}

type SortByID struct {
	Records []Record
	SortAsc bool
}

func (s SortByID) Len() int      { return len(s.Records) }
func (s SortByID) Swap(i, j int) { s.Records[i], s.Records[j] = s.Records[j], s.Records[i] }
func (s SortByID) Less(i, j int) bool {
	if s.SortAsc {
		return s.Records[i].ID < s.Records[j].ID
	}
	return s.Records[i].ID > s.Records[j].ID
}

type SortByCreated struct {
	Records []Record
	SortAsc bool
}

func (s SortByCreated) Len() int      { return len(s.Records) }
func (s SortByCreated) Swap(i, j int) { s.Records[i], s.Records[j] = s.Records[j], s.Records[i] }
func (s SortByCreated) Less(i, j int) bool {

	slog.Info("comparing",
		slog.Any("a", s.Records[i].CreatedAt.UnixMicro()),
		slog.Any("b", s.Records[j].CreatedAt.UnixMicro()),
	)

	if s.SortAsc {
		return s.Records[i].CreatedAt.UnixMicro() < s.Records[j].CreatedAt.UnixMicro()
	}
	return s.Records[i].CreatedAt.UnixMicro() > s.Records[j].CreatedAt.UnixMicro()
}

type ListOrdering struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type ListInput struct {
	Ordering  ListOrdering
	Filtering []Label
}

// The in-memory storage
var storage = make(map[string]Record)

func storeRecord(in *Record) {
	if in == nil {
		slog.Warn("unable to save nil record")
	}

	// TODO: check for presence and return a specific error
	// TODO: use a mutex
	storage[in.ID] = *in
	slog.Info("count updated", slog.Any("count", len(storage)))
}

func findRecordByID(id string) (*Record, error) {
	r, found := storage[id]

	if !found {
		return nil, fmt.Errorf("record not found with id: %s", id)
	}
	return &r, nil
}

func deleteRecord(id string) error {
	record, err := findRecordByID(id)

	if err != nil {
		return fmt.Errorf("delete record error: %w", err)
	}

	delete(storage, id)
	slog.Info("deleted record", slog.Any("record", record))
	return nil
}

func createRecord(in RecordInput) *Record {
	if in.ID == "" {
		slog.Warn("using a generated ID")
		in.ID = uuid.NewString()
	}

	//timestamp, err := time.Parse(time.RFC3339, in.CreatedAt)

	return &Record{
		ID:        in.ID,
		CreatedAt: time.Now(),
		Labels:    in.Labels,
		Object:    in.Object,
	}
}

func getFilterMap(in []Label) map[Label]struct{} {
	res := map[Label]struct{}{}
	for _, f := range in {
		res[f] = struct{}{}
	}

	return res
}

func sortRecords(records []Record, listOrdering ListOrdering) []Record {
	if strings.ToLower(listOrdering.ID) == "asc" {
		sort.Sort(SortByID{Records: records, SortAsc: true})
		return records
	}

	if strings.ToLower(listOrdering.ID) == "desc" {
		sort.Sort(SortByID{Records: records, SortAsc: false})
		return records
	}

	if strings.ToLower(listOrdering.CreatedAt) == "desc" {
		sort.Sort(SortByCreated{Records: records, SortAsc: false})
		return records
	}

	// default sorting
	sort.Sort(SortByCreated{Records: records, SortAsc: true})

	return records
}

func listRecords(listInput ListInput) ([]Record, error) {
	res := []Record{}

	// optimize the filtering by pre-computing a map
	filterMap := getFilterMap(listInput.Filtering)

	slog.Info("prepare filtering", slog.Any("filterMap", filterMap))

	for _, record := range storage {
		slog.Info("checking record", slog.Any("record", record))

		if len(filterMap) > 0 {

			for _, recordLbl := range record.Labels {
				slog.Info("WTF check", slog.Any("recordLbl", recordLbl))

				_, found := filterMap[recordLbl]

				if found {
					res = append(res, record)
				}
			}
		} else {
			res = append(res, record)
		}
	}

	res = sortRecords(res, listInput.Ordering)

	if len(res) == 0 {
		slog.Info("no records found")
	}

	slog.Info("sorted records", slog.Any("records", res))
	return res, nil
}
