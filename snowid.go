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

type SnowIdGenerator struct {
	l            sync.RWMutex
	Records      map[int64]int
	DataCenterId int8
	MachineId    int8
	Epoch        time.Time
}

type SnowID struct {
	Timestamp      int64     `json:"timestamp"`
	DataCenterId   int8      `json:"datacenter_id"`
	MachineId      int8      `json:"machine_id"`
	SequenceNumber int16     `json:"sequence_number"`
	Epoch          time.Time `json:"epoch"`
}

// Initialize the ID Generator
//
// dataCenterId. min 0, max 31
//
// machineId. min 0, max 31
//
// epoch: The epoch time to start generating IDs. Could use the DefaultEpoch
func NewSnowIdGenerator(dataCenterId, machineId int, epoch time.Time) (*SnowIdGenerator, error) {
	// validation
	if dataCenterId < 0 || dataCenterId > MAX_DATACENTER_ID {
		return nil, errors.New("datacenterId must be between 0 and 31")
	}
	if machineId < 0 || machineId > MAX_MACHINE_ID {
		return nil, errors.New("machineId must be between 0 and 31")
	}

	s := &SnowIdGenerator{
		Records:      make(map[int64]int),
		DataCenterId: int8(dataCenterId),
		MachineId:    int8(machineId),
		Epoch:        epoch,
	}
	return s, nil
}

// Generate Snowflake ID object
//
// Return error if the sequence number is equal to or greater than the max
// sequence in 1 milisecond
func (s *SnowIdGenerator) GenerateId() (*SnowID, error) {
	currentTimestamp := time.Since(s.Epoch).Milliseconds()

	return s.generateId(currentTimestamp)
}

// Generate ID with timestamp input
func (s *SnowIdGenerator) generateId(timestamp int64) (*SnowID, error) {
	id := &SnowID{
		Timestamp:    timestamp,
		DataCenterId: s.DataCenterId,
		MachineId:    s.MachineId,
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
func (s *SnowIdGenerator) AutoResetRecords(duration time.Duration) {
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
func (s *SnowIdGenerator) ResetRecords() {
	s.l.Lock()
	s.Records = make(map[int64]int)
	s.l.Unlock()
}

// Return binary string
func (id *SnowID) StringBinary() string {
	timestamp_bin := fmt.Sprintf("%041b", id.Timestamp)
	dataCenterId_bin := fmt.Sprintf("%05b", id.DataCenterId)
	machineId_bin := fmt.Sprintf("%05b", id.MachineId)
	sequenceNumber_bin := fmt.Sprintf("%012b", id.SequenceNumber)

	return fmt.Sprintf("%s%s%s%s%s", DEFAULT_PLACEHOLDER_BIT, timestamp_bin, dataCenterId_bin, machineId_bin, sequenceNumber_bin)
}

func (id *SnowID) String() string {
	id_int, _ := strconv.ParseInt(id.StringBinary(), 2, 64)
	return fmt.Sprint(id_int)
}

func (id *SnowID) Datetime() time.Time {
	unixEpoch := time.Unix(0, 0).UTC()

	return time.UnixMilli(id.Timestamp + id.Epoch.UnixMilli() - unixEpoch.UnixMilli()).UTC()
}

// Parse ID in decimal string
//
// eg. ParseId("0000001001011100001100001111001101011110100000100001000000000000")
func ParseIdBinary(id string, customEpoch time.Time) (*SnowID, error) {
	if len(id) != TOTAL_BIT {
		return nil, errors.New("invalid ID length. The ID should be 64-bit binary string")
	}
	// get binary string
	timestamp_start := 0 + PLACEHOLDER_BIT
	datacenter_start := timestamp_start + TIMESTAMP_BIT
	machine_start := datacenter_start + DATACENTER_BIT
	sequence_start := machine_start + MACHINE_BIT

	timestamp := id[timestamp_start:datacenter_start]
	datacenterId := id[datacenter_start:machine_start]
	machineId := id[machine_start:sequence_start]
	sequenceNumber := id[sequence_start:]

	// convert to integer
	timestamp_int, err := strconv.ParseInt(timestamp, 2, 64)
	if err != nil {
		return nil, errors.New("invalid timestamp. The timestamp should be a number")
	}
	datacenterId_int, err := strconv.ParseInt(datacenterId, 2, 8)
	if err != nil {
		return nil, errors.New("datacenter id should be a number")
	}
	machineId_int, err := strconv.ParseInt(machineId, 2, 8)
	if err != nil {
		return nil, errors.New("machine id should be a number")
	}
	sequenceNumber_int, err := strconv.ParseInt(sequenceNumber, 2, 16)
	if err != nil {
		return nil, errors.New("sequence number should be a number")
	}

	idObj := &SnowID{
		Timestamp:      timestamp_int,
		DataCenterId:   int8(datacenterId_int),
		MachineId:      int8(machineId_int),
		SequenceNumber: int16(sequenceNumber_int),
		Epoch:          customEpoch,
	}

	return idObj, nil
}

// Parse ID in decimal string
//
// eg. ParseId("170064707754004481")
func ParseId(id string, customEpoch time.Time) (*SnowID, error) {
	id_int, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("the id is not a number")
	}

	id_bin := fmt.Sprintf("%064b", id_int)
	if len(id_bin) != TOTAL_BIT {
		return nil, errors.New("invalid id length")
	}
	idObj, err := ParseIdBinary(id_bin, customEpoch)
	if err != nil {
		return nil, err
	}

	return idObj, nil
}
