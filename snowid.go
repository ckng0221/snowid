package snowid

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Unique 64-bit ID generator based on Twitter Snowflake
// ID generator
// |0|timestamp|datacenterID|machineID|sequenceNumber
// |1|41|5|5|12| bits

// 41 bit timestamp in miliseconds, can have max around 69years
// use 2025-01-01 UTC as the starting epoch, instead of unix epoch
// max until year 2095

const (
	PLACEHOLDER_BIT = 1
	TIMESTAMP_BIT   = 41
	DATACENTER_BIT  = 5
	MACHINE_BIT     = 5
	SEQUENCE_BIT    = 12
	TOTAL_BIT       = PLACEHOLDER_BIT + TIMESTAMP_BIT + DATACENTER_BIT + MACHINE_BIT + SEQUENCE_BIT // 64

	DEFAULT_PLACEHOLDER_BIT = "0"

	DATACENTER_CAP = (1 << DATACENTER_BIT) // 32
	MACHINE_CAP    = (1 << MACHINE_BIT)    // 32
	SEQUENCE_CAP   = (1 << SEQUENCE_BIT)   // 4096

	MAX_DATACENTER_ID = DATACENTER_CAP - 1 // 31
	MAX_MACHINE_ID    = MACHINE_CAP - 1    // 31
	MAX_SEQUENCE_ID   = SEQUENCE_CAP - 1   // 4095
)

var (
	DefaultEpoch = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
)

type SnowIDGenerator struct {
	l            sync.RWMutex
	Records      map[int64]int
	DataCenterID int8
	MachineID    int8
	Epoch        time.Time
}

type SnowID struct {
	Timestamp      int64     `json:"timestamp"`
	DataCenterID   int8      `json:"datacenter_id"`
	MachineID      int8      `json:"machine_id"`
	SequenceNumber int16     `json:"sequence_number"`
	Epoch          time.Time `json:"epoch"`
}

// Initialize the ID Generator
//
// dataCenterID. min 0, max 31
//
// machineID. min 0, max 31
//
// epoch: The epoch time to start generating IDs. Could use the DefaultEpoch.
func NewSnowIDGenerator(dataCenterID, machineID int, epoch time.Time) *SnowIDGenerator {
	// validation
	if dataCenterID < 0 || dataCenterID > MAX_DATACENTER_ID {
		panic("dataCenterID must be between 0 and 31")
	}
	if machineID < 0 || machineID > MAX_MACHINE_ID {
		panic("machineID must be between 0 and 31")
	}

	s := &SnowIDGenerator{
		Records:      make(map[int64]int),
		DataCenterID: int8(dataCenterID),
		MachineID:    int8(machineID),
		Epoch:        epoch,
	}
	return s
}

// Generate Snowflake ID object
//
// Return error if the sequence number is equal to or greater than the max
// sequence in 1 milisecond
func (s *SnowIDGenerator) GenerateID() (*SnowID, error) {
	currentTimestamp := time.Since(s.Epoch).Milliseconds()

	return s.generateID(currentTimestamp)
}

// Generate ID with timestamp input
func (s *SnowIDGenerator) generateID(timestamp int64) (*SnowID, error) {
	id := &SnowID{
		Timestamp:    timestamp,
		DataCenterID: s.DataCenterID,
		MachineID:    s.MachineID,
		Epoch:        s.Epoch,
	}
	s.l.Lock()
	defer s.l.Unlock()
	count, ok := s.Records[timestamp]
	if ok {
		count++
		if count > MAX_SEQUENCE_ID {
			return nil, errors.New("sequence number greater than max limit")
		}

		s.Records[timestamp] = count
		id.SequenceNumber = int16(count)
		return id, nil
	} else {
		// sequence starts from 0
		s.Records[timestamp] = 0
		return id, nil
	}
}

// AutoResetRecords will reset the records every n seconds
//
// Runs a goroutine that resets the records every n seconds.
// This is useful to avoid build up of records.
func (s *SnowIDGenerator) AutoResetRecords(duration time.Duration) {
	resetRecordsOnSchedule := func() {
		// clean up every n second
		ticker := time.NewTicker(duration)
		for range ticker.C {
			s.l.Lock()
			s.Records = make(map[int64]int)
			s.l.Unlock()
		}
	}
	go resetRecordsOnSchedule()
}

// Reset all hashtable records
func (s *SnowIDGenerator) ResetRecords() {
	s.l.Lock()
	s.Records = make(map[int64]int)
	s.l.Unlock()
}

// Return binary string of ID
func (id *SnowID) StringBinary() string {
	timestampBinStr := fmt.Sprintf("%041b", id.Timestamp)
	dataCenterIDBinStr := fmt.Sprintf("%05b", id.DataCenterID)
	machineIDBinStr := fmt.Sprintf("%05b", id.MachineID)
	sequenceNumberBinStr := fmt.Sprintf("%012b", id.SequenceNumber)

	return fmt.Sprintf("%s%s%s%s%s", DEFAULT_PLACEHOLDER_BIT, timestampBinStr, dataCenterIDBinStr, machineIDBinStr, sequenceNumberBinStr)
}

// return decimal integer of ID
func (id *SnowID) Int64() int64 {
	idInt, _ := strconv.ParseInt(id.StringBinary(), 2, 64)
	return idInt
}

// Return decimal string of ID
func (id *SnowID) String() string {
	return strconv.FormatInt(id.Int64(), 10)
}

func (id *SnowID) Datetime() time.Time {
	unixEpoch := time.Unix(0, 0).UTC()

	return time.UnixMilli(id.Timestamp + id.Epoch.UnixMilli() - unixEpoch.UnixMilli()).UTC()
}

// Parse ID in decimal string
//
// eg. ParseID("0000001001011100001100001111001101011110100000100001000000000000")
func ParseIDBinary(idStr string, customEpoch time.Time) (*SnowID, error) {
	if len(idStr) != TOTAL_BIT {
		return nil, errors.New("invalid ID length. The ID should be 64-bit binary string")
	}
	// get binary string
	timestampStart := 0 + PLACEHOLDER_BIT
	dataCenterStart := timestampStart + TIMESTAMP_BIT
	machineStart := dataCenterStart + DATACENTER_BIT
	sequenceStart := machineStart + MACHINE_BIT

	timestamp := idStr[timestampStart:dataCenterStart]
	dataCenterID := idStr[dataCenterStart:machineStart]
	machineID := idStr[machineStart:sequenceStart]
	sequenceNumber := idStr[sequenceStart:]

	// convert to integer
	timestampInt, err := strconv.ParseInt(timestamp, 2, 64)
	if err != nil {
		return nil, errors.New("invalid timestamp. The timestamp should be a number")
	}
	dataCenterIDInt, err := strconv.ParseInt(dataCenterID, 2, 8)
	if err != nil {
		return nil, errors.New("datacenter id should be a number")
	}
	machineIDInt, err := strconv.ParseInt(machineID, 2, 8)
	if err != nil {
		return nil, errors.New("machine id should be a number")
	}
	sequenceNumberInt, err := strconv.ParseInt(sequenceNumber, 2, 16)
	if err != nil {
		return nil, errors.New("sequence number should be a number")
	}

	id := &SnowID{
		Timestamp:      timestampInt,
		DataCenterID:   int8(dataCenterIDInt),
		MachineID:      int8(machineIDInt),
		SequenceNumber: int16(sequenceNumberInt),
		Epoch:          customEpoch,
	}

	return id, nil
}

// Parse ID in decimal string
//
// eg. ParseID("170064707754004481")
func ParseID(idStr string, customEpoch time.Time) (*SnowID, error) {
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("the ID is not a number")
	}

	idBinStr := fmt.Sprintf("%064b", idInt)
	if len(idBinStr) != TOTAL_BIT {
		return nil, errors.New("invalid ID length")
	}
	id, err := ParseIDBinary(idBinStr, customEpoch)
	if err != nil {
		return nil, err
	}

	return id, nil
}
