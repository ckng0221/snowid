package snowid

import (
	"testing"
	"time"
)

const (
	binaryStringId = "0110010001000010101011010010000010000001000110111110000000000110"
	stringId       = "7224527107372343302"
	machineId      = 30
	dataCenterId   = 13
)

// Date
func TestCustomEpochDateTimeDefaultEpoch(t *testing.T) {
	s, err := NewSnowIdGenerator(machineId, dataCenterId, DefaultEpoch)
	if err != nil {
		t.Error(err.Error())
	}
	id, _ := s.GenerateId()
	id_date := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	if id_date != expectedDate {
		t.Errorf("Expected %s, got %s", expectedDate, id_date)
		return
	}
}

func TestCustomEpochDateTimeEarlierEpoch(t *testing.T) {
	earlierEpoch := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
	s, err := NewSnowIdGenerator(machineId, dataCenterId, earlierEpoch)
	if err != nil {
		t.Error(err.Error())
	}
	id, _ := s.GenerateId()
	id_date := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	if id_date != expectedDate {
		t.Errorf("Expected %s, got %s", expectedDate, id_date)
		return
	}
}

func TestCustomEpochDateTimeUnixEpoch(t *testing.T) {
	unixEpoch := time.Unix(0, 0).UTC()
	s, err := NewSnowIdGenerator(machineId, dataCenterId, unixEpoch)
	if err != nil {
		t.Error(err.Error())
	}
	id, _ := s.GenerateId()
	id_date := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	if id_date != expectedDate {
		t.Errorf("Expected %s, got %s", expectedDate, id_date)
		return
	}
}

// Parsing
func TestParseStringBinaryDataCenterId(t *testing.T) {
	id, err := ParseIdBinary(binaryStringId, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	if id.DataCenterId != dataCenterId {
		t.Errorf("Expected DataCenterId 13, got %d", id.DataCenterId)
		return
	}
}

func TestParseStringBinaryMachineId(t *testing.T) {
	id, err := ParseIdBinary(binaryStringId, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := machineId
	if id.MachineId != int8(expectedValue) {
		t.Errorf("Expected MachineId %d, got %d", expectedValue, id.MachineId)
		return
	}
}

func TestParseStringBinaryDatetime(t *testing.T) {
	id, err := ParseIdBinary(binaryStringId, time.Unix(0, 0).UTC())
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := time.Date(2024, 7, 31, 21, 31, 27, 0, time.UTC)
	if id.Datetime().UTC().String()[:19] != expectedValue.String()[:19] {
		t.Errorf("Expected Datetime %s, got %s", expectedValue.String()[:19], id.Datetime().UTC().String()[:19])
		return
	}
}

func TestParseStringDatetime(t *testing.T) {
	id, err := ParseId(stringId, time.Unix(0, 0).UTC())
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := time.Date(2024, 7, 31, 21, 31, 27, 0, time.UTC)
	if id.Datetime().UTC().String()[:19] != expectedValue.String()[:19] {
		t.Errorf("Expected Datetime %s, got %s", expectedValue.String()[:19], id.Datetime().UTC().String()[:19])
		return
	}
}

func TestParseStringBinarySequenceNumber(t *testing.T) {
	id, err := ParseIdBinary(binaryStringId, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := 6
	if id.SequenceNumber != int16(expectedValue) {
		t.Errorf("Expected SequenceNumber %d, got %d", expectedValue, id.SequenceNumber)
		return
	}
}

func TestParseBinaryStringToDecimal(t *testing.T) {
	id, err := ParseIdBinary(binaryStringId, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := stringId
	if id.String() != expectedValue {
		t.Errorf("Expected ID %s, got %s", expectedValue, id.String())
		return
	}
}

func TestParseStringToBinaryString(t *testing.T) {
	id, err := ParseId(stringId, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := binaryStringId
	if id.StringBinary() != expectedValue {
		t.Errorf("Expected ID %s, got %s", expectedValue, id.String())
		return
	}
}

func TestGenerateIdEqualParseId(t *testing.T) {
	s, err := NewSnowIdGenerator(dataCenterId, machineId, DefaultEpoch)
	if err != nil {
		t.Error(err.Error())
	}
	id, err := s.GenerateId()
	if err != nil {
		t.Errorf("failed to generate ID")
		return
	}
	id_copy := id.String()

	id2, err := ParseId(id_copy, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	if id.String() != id2.String() {
		t.Errorf("Expected ID %s, got %s", id.String(), id2.String())
		return
	}
}

// ID Generation
func TestGenerateMany(t *testing.T) {
	s, err := NewSnowIdGenerator(dataCenterId, machineId, DefaultEpoch)
	if err != nil {
		t.Error(err.Error())
	}
	qty := 1000
	ids := []*SnowID{}
	for range qty {
		id, _ := s.GenerateId()
		ids = append(ids, id)
	}
	if len(ids) != qty {
		t.Errorf("Expected quantity %d, got %d", qty, len(ids))
	}
}

func TestGenerateExceedSequenceLimitIdCount(t *testing.T) {
	s, err := NewSnowIdGenerator(dataCenterId, machineId, DefaultEpoch)
	if err != nil {
		t.Error(err.Error())
	}
	ids := []*SnowID{}
	timestamp := time.Now().UnixMilli()
	for range SEQUENCE_CAP {
		id, _ := s.generateId(timestamp)
		ids = append(ids, id)
	}
	if len(ids) != SEQUENCE_CAP {
		t.Errorf("Expected quantity %d, got %d", SEQUENCE_CAP, len(ids))
	}

	_, err = s.generateId(timestamp)
	if err == nil {
		t.Errorf("Expected to throw error when exceeding sequence max")
	}
}

func TestGenerateExceedSequenceLimitIdCountWithSleep(t *testing.T) {
	s, err := NewSnowIdGenerator(dataCenterId, machineId, DefaultEpoch)
	if err != nil {
		t.Error(err.Error())
	}
	ids := []*SnowID{}
	for range SEQUENCE_CAP {
		id, err := s.GenerateId()
		if err != nil {
			t.Errorf("Failed to generate ID: %s", err.Error())
		}
		ids = append(ids, id)
	}
	time.Sleep(1 * time.Millisecond)
	for range SEQUENCE_CAP {
		id, err := s.GenerateId()
		if err != nil {
			t.Errorf("Failed to generate ID: %s", err.Error())
		}
		ids = append(ids, id)
	}
	// check total
	if len(ids) != 2*SEQUENCE_CAP {
		t.Errorf("Expected quantity %d, got %d", 2*SEQUENCE_CAP, len(ids))
	}

	// check length
	last := ids[len(ids)-1]

	if len(last.StringBinary()) > TOTAL_BIT {
		t.Errorf("Expected length %d, got %d", TOTAL_BIT, len(last.StringBinary()))
	}

}
