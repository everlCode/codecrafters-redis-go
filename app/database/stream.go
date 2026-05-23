package database

import "github.com/codecrafters-io/redis-starter-go/app/helpers"

type Stream struct {
	entries []StreamEntry
	lastId  string
}

type StreamEntry struct {
	id   string
	data map[string]string
}

func (s Stream) GetLastId() string {
	return s.lastId
}

func (s Stream) GetEntries(params ...string) []StreamEntry {
	if len(params) == 0 {
		return s.entries
	}

	start := params[0]
	var startMilliseconds, startSeqId int
	if start == "-" {
		startMilliseconds, startSeqId = 0, 0
	} else {
		startMilliseconds, startSeqId = helpers.GetStreamIdParts(start)
	}

	var end string
	if len(params) > 1 {
		end = params[1]
	} else {
		end = "+"
	}
	var endMilliseconds, endSeqId int
	isEndExist := end != "+"
	if isEndExist {
		endMilliseconds, endSeqId = helpers.GetStreamIdParts(end)
	}

	var response []StreamEntry
	for _, streamEntry := range s.entries {
		entryId := streamEntry.GetId()
		miliseconds, seqId := helpers.GetStreamIdParts(entryId)

		if miliseconds >= startMilliseconds && seqId >= startSeqId &&
			(!isEndExist || (miliseconds <= endMilliseconds && seqId <= endSeqId)) {
			response = append(response, streamEntry)
		}
	}

	return response
}

func (entry StreamEntry) GetId() string {
	return entry.id
}

func (entry StreamEntry) GeData() map[string]string {
	return entry.data
}

func CreateStream() *Stream {
	stream := NewStream()

	return &stream
}

func NewStream() Stream {
	return Stream{
		entries: []StreamEntry{},
	}
}

func (stream *Stream) Add(id string, data map[string]string) {
	entry := StreamEntry{}
	entry.id = id
	entry.data = data
	stream.entries = append(stream.entries, entry)
	stream.lastId = id
}