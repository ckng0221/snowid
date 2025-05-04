package snowid

import (
	"testing"
	"time"
)

const (
	binaryStringID = "0110010001000010101011010010000010000001000110111110000000000110"
	stringID       = "7224527107372343302"
	machineID      = 30
	dataCenterID   = 13
)

// Date
func TestCustomEpochDateTimeDefaultEpoch(t *testing.T) {
	s := NewSnowIDGenerator(machineID, dataCenterID, DefaultEpoch)
	id, err := s.GenerateID()
	if err != nil {
		t.Errorf("Failed to generate ID: %s", err.Error())
		return
	}
	id_date := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	if id_date != expectedDate {
		t.Errorf("Expected %s, got %s", expectedDate, id_date)
		return
	}
}

func TestCustomEpochDateTimeEarlierEpoch(t *testing.T) {
	earlierEpoch := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
	s := NewSnowIDGenerator(machineID, dataCenterID, earlierEpoch)
	id, err := s.GenerateID()
	if err != nil {
		t.Errorf("Failed to generate ID: %s", err.Error())
		return
	}
	id_date := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	if id_date != expectedDate {
		t.Errorf("Expected %s, got %s", expectedDate, id_date)
		return
	}
}

func TestCustomEpochDateTimeUnixEpoch(t *testing.T) {
	unixEpoch := time.Unix(0, 0).UTC()
	s := NewSnowIDGenerator(machineID, dataCenterID, unixEpoch)
	id, err := s.GenerateID()
	if err != nil {
		t.Errorf("Failed to generate ID: %s", err.Error())
		return
	}
	id_date := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	if id_date != expectedDate {
		t.Errorf("Expected %s, got %s", expectedDate, id_date)
		return
	}
}

// Parsing
func TestParseStringBinaryDataCenterID(t *testing.T) {
	id, err := ParseIDBinary(binaryStringID, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	if id.DataCenterID != dataCenterID {
		t.Errorf("Expected DataCenterID 13, got %d", id.DataCenterID)
		return
	}
}

func TestParseStringBinaryMachineID(t *testing.T) {
	id, err := ParseIDBinary(binaryStringID, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := machineID
	if id.MachineID != int8(expectedValue) {
		t.Errorf("Expected MachineID %d, got %d", expectedValue, id.MachineID)
		return
	}
}

func TestParseStringBinaryDatetime(t *testing.T) {
	id, err := ParseIDBinary(binaryStringID, time.Unix(0, 0).UTC())
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
	id, err := ParseID(stringID, time.Unix(0, 0).UTC())
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
	id, err := ParseIDBinary(binaryStringID, DefaultEpoch)
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
	id, err := ParseIDBinary(binaryStringID, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := stringID
	if id.String() != expectedValue {
		t.Errorf("Expected ID %s, got %s", expectedValue, id.String())
		return
	}
}

func TestParseStringToBinaryString(t *testing.T) {
	id, err := ParseID(stringID, DefaultEpoch)
	if err != nil {
		t.Errorf("failed to parse ID")
		return
	}
	expectedValue := binaryStringID
	if id.StringBinary() != expectedValue {
		t.Errorf("Expected ID %s, got %s", expectedValue, id.String())
		return
	}
}

func TestGenerateIDEqualParseID(t *testing.T) {
	s := NewSnowIDGenerator(dataCenterID, machineID, DefaultEpoch)
	id, err := s.GenerateID()
	if err != nil {
		t.Errorf("failed to generate ID")
		return
	}
	id_copy := id.String()

	id2, err := ParseID(id_copy, DefaultEpoch)
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
	s := NewSnowIDGenerator(dataCenterID, machineID, DefaultEpoch)
	qty := 1000
	ids := []*SnowID{}
	for range qty {
		id, err := s.GenerateID()
		if err != nil {
			t.Errorf("Failed to generate ID: %s", err.Error())
			return
		}
		ids = append(ids, id)
	}
	if len(ids) != qty {
		t.Errorf("Expected quantity %d, got %d", qty, len(ids))
	}
}

func TestGenerateExceedSequenceLimitIDCount(t *testing.T) {
	s := NewSnowIDGenerator(dataCenterID, machineID, DefaultEpoch)
	ids := []*SnowID{}
	timestamp := time.Now().UnixMilli()
	for range SEQUENCE_CAP {
		id, err := s.generateID(timestamp)
		if err != nil {
			t.Errorf("Failed to generate ID: %s", err.Error())
			return
		}

		ids = append(ids, id)
	}
	if len(ids) != SEQUENCE_CAP {
		t.Errorf("Expected quantity %d, got %d", SEQUENCE_CAP, len(ids))
	}

	id, err := s.generateID(timestamp)
	if err == nil {
		t.Errorf("Expected to throw error when exceeding sequence max")
	}
	if id != nil {
		t.Errorf("Expect ID to be nil")
	}
}

func TestGenerateExceedSequenceLimitIDCountWithSleep(t *testing.T) {
	s := NewSnowIDGenerator(dataCenterID, machineID, DefaultEpoch)
	ids := []*SnowID{}
	for range SEQUENCE_CAP {
		id, err := s.GenerateID()
		if err != nil {
			t.Errorf("Failed to generate ID: %s", err.Error())
			return
		}
		ids = append(ids, id)
	}
	time.Sleep(2 * time.Millisecond)
	for range SEQUENCE_CAP {
		id, err := s.GenerateID()
		if err != nil {
			t.Errorf("Failed to generate ID: %s", err.Error())
			return
		}
		ids = append(ids, id)
	}
	// check total
	if len(ids) != 2*SEQUENCE_CAP {
		t.Errorf("Expected quantity %d, got %d", 2*SEQUENCE_CAP, len(ids))
		return
	}

	// check length
	last := ids[len(ids)-1]

	if len(last.StringBinary()) > TOTAL_BIT {
		t.Errorf("Expected length %d, got %d", TOTAL_BIT, len(last.StringBinary()))
		return
	}

}
