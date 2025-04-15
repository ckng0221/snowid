package snowid

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Unique ID generator based on Twitter Snowflake
// ID generator
// |0|timestamp|datacenterID|machineID|sequenceNumber
// |1|41|5|5|12| bits
// maximum 4096 sequence

// 41 bit timestamp in miliseconds, can have max around 69years
// use 2024-01-01 UTC as the starting epoch, instead of unix epoch
// max until year 2093
var (
	DefaultEpoch = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
)

type SnowIdGenerator struct {
	l            sync.RWMutex
	Records      map[int]int
	DataCenterId int8
	MachineId    int8
	Epoch        time.Time
}

// Initialize the ID Generator
//
// dataCenterId. min 0, max 31
//
// machineId. min 0, max 31
//
// echo: The epoch time to start generating IDs. Could use idg.DefaultEpoch
func NewSnowIdGenerator(dataCenterId, machineId int, epoch time.Time) (*SnowIdGenerator, error) {
	// validation
	if dataCenterId < 0 || dataCenterId > 31 {
		return nil, errors.New("datacenterId must be between 0 and 31")
	}
	if machineId < 0 || machineId > 31 {
		return nil, errors.New("machineId must be between 0 and 31")
	}

	idg := &SnowIdGenerator{
		Records:      make(map[int]int, 2),
		DataCenterId: int8(dataCenterId),
		MachineId:    int8(machineId),
		Epoch:        epoch,
	}
	return idg, nil
}

func (idg *SnowIdGenerator) GenerateId() *ID {
	// base
	currentTimestamp := int(time.Since(idg.Epoch).Milliseconds())

	id := &ID{
		Timestamp:    int64(currentTimestamp),
		DataCenterId: idg.DataCenterId,
		MachineId:    idg.MachineId,
		Epoch:        idg.Epoch,
	}
	idg.l.RLock()
	count, ok := idg.Records[currentTimestamp]
	idg.l.RUnlock()
	if ok {
		idg.l.Lock()
		defer idg.l.Unlock()
		count++

		idg.Records[currentTimestamp] = count
		id.SequenceNumber = int16(count)
		return id
	} else {
		idg.l.Lock()
		defer idg.l.Unlock()
		// sequence starts from 0
		idg.Records[currentTimestamp] = 0
		return id
	}
}

// AutoResetRecords will reset the records every n seconds
//
// This is useful on server
func (idg *SnowIdGenerator) AutoResetRecords(duration time.Duration) {
	resetRecordsOnSchedule := func() {
		// clean up every n second
		ticker := time.NewTicker(duration)
		for range ticker.C {
			idg.l.Lock()
			idg.Records = make(map[int]int)
			idg.l.Unlock()

			// log.Println("records cleaned up.")
		}
	}
	go resetRecordsOnSchedule()
}

// Reset all hashtable records
func (idg *SnowIdGenerator) ResetRecords() {
	idg.l.Lock()
	idg.Records = make(map[int]int)
	idg.l.Unlock()
}

type ID struct {
	Timestamp      int64     `json:"timestamp"`
	DataCenterId   int8      `json:"datacenter_id"`
	MachineId      int8      `json:"machine_id"`
	SequenceNumber int16     `json:"sequence_number"`
	Epoch          time.Time `json:"epoch"`
}

// Return binary string
func (id *ID) StringBinary() string {
	initialBit := "0"

	//uniqueId
	timestamp_bin := fmt.Sprintf("%041b", id.Timestamp)
	dataCenterId_bin := fmt.Sprintf("%05b", id.DataCenterId)
	machineId_bin := fmt.Sprintf("%05b", id.MachineId)
	sequenceNumber_bin := fmt.Sprintf("%012b", id.SequenceNumber)

	return fmt.Sprintf("%s%s%s%s%s", initialBit, timestamp_bin, dataCenterId_bin, machineId_bin, sequenceNumber_bin)
}

func (id *ID) String() string {
	id_int, _ := strconv.ParseInt(id.StringBinary(), 2, 64)
	return fmt.Sprint(id_int)
}

func (id *ID) Datetime() time.Time {
	unixEpoch := time.Unix(0, 0).UTC()

	return time.UnixMilli(id.Timestamp + id.Epoch.UnixMilli() - unixEpoch.UnixMilli()).UTC()
}

// Parse ID in decimal string
//
// eg. ParseId("0000001001011100001100001111001101011110100000100001000000000000")
func ParseIdBinary(id string, customEpoch time.Time) (*ID, error) {
	if len(id) != 64 {
		return nil, errors.New("invalid id length")
	}
	// get binary string
	timestamp := id[1 : 1+41]
	datacenterId := id[42 : 42+5]
	machineId := id[47 : 47+5]
	sequenceNumber := id[52 : 52+12]

	// convert to integer
	timestamp_int, _ := strconv.ParseInt(timestamp, 2, 64)
	datacenterId_int, _ := strconv.ParseInt(datacenterId, 2, 8)
	fmt.Println(machineId)
	machineId_int, _ := strconv.ParseInt(machineId, 2, 8)
	sequenceNumber_int, _ := strconv.ParseInt(sequenceNumber, 2, 16)
	idObj := &ID{
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
func ParseId(id string, customEpoch time.Time) (*ID, error) {
	id_int, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("the id is not a number")
	}

	id_bin := fmt.Sprintf("%064b", id_int)
	if len(id_bin) != 64 {
		return nil, errors.New("invalid id length")
	}
	idObj, err := ParseIdBinary(id_bin, customEpoch)
	if err != nil {
		return nil, err
	}

	return idObj, nil
}
