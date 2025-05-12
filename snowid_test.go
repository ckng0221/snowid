package snowid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TestSuit struct {
	suite.Suite
	S              *SnowIDGenerator
	SE             *SnowIDGenerator
	SU             *SnowIDGenerator
	BinaryStringID string
	StringID       string
	DataCenterID   int
	MachineID      int
	SequenceNumber int16
}

func (suite *TestSuit) SetupTest() {
	earlierEpoch := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
	unixEpoch := time.Unix(0, 0).UTC()

	// Sample values used in test cases
	suite.BinaryStringID = "0110010001000010101011010010000010000001000110111110000000000110"
	suite.StringID = "7224527107372343302"
	suite.MachineID = 30
	suite.DataCenterID = 13
	suite.SequenceNumber = 6

	suite.S = NewSnowIDGenerator(suite.MachineID, suite.DataCenterID, DefaultEpoch)
	suite.SE = NewSnowIDGenerator(suite.MachineID, suite.DataCenterID, earlierEpoch)
	suite.SU = NewSnowIDGenerator(suite.MachineID, suite.DataCenterID, unixEpoch)

}

// Date
func (suite *TestSuit) TestCustomEpochDateTimeDefaultEpoch() {
	id, err := suite.S.GenerateID()
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}

	id_date := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	suite.Equal(expectedDate, id_date)
}

func (suite *TestSuit) TestCustomEpochDateTimeEarlierEpoch() {
	id, err := suite.SE.GenerateID()
	suite.Equal(nil, err, "Failed to generate ID")

	idDate := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	suite.Equal(expectedDate, idDate)
}

func (suite *TestSuit) TestCustomEpochDateTimeUnixEpoch() {
	id, err := suite.SU.GenerateID()
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}

	idDate := id.Datetime().UTC().String()[:10]
	expectedDate := time.Now().UTC().String()[:10]

	suite.Equal(expectedDate, idDate)
}

// Parsing
func (suite *TestSuit) TestParseStringBinaryDataCenterID() {
	id, err := ParseIDBinary(suite.BinaryStringID, DefaultEpoch)
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}

	suite.Equal(suite.DataCenterID, int(id.DataCenterID))
}

func (suite *TestSuit) TestParseStringBinaryMachineID() {
	id, err := ParseIDBinary(suite.BinaryStringID, DefaultEpoch)
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}

	suite.Equal(suite.MachineID, int(id.MachineID))
}

func (suite *TestSuit) TestParseStringBinaryDatetime() {
	id, err := ParseIDBinary(suite.BinaryStringID, time.Unix(0, 0).UTC())
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}
	expectedValue := time.Date(2024, 7, 31, 21, 31, 27, 0, time.UTC)
	suite.Equal(expectedValue.String()[:19], id.Datetime().UTC().String()[:19])
}

func (suite *TestSuit) TestParseStringDatetime() {
	id, err := ParseID(suite.StringID, time.Unix(0, 0).UTC())
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}
	expectedValue := time.Date(2024, 7, 31, 21, 31, 27, 0, time.UTC)
	suite.Equal(expectedValue.String()[:19], id.Datetime().UTC().String()[:19])
}

func (suite *TestSuit) TestParseStringBinarySequenceNumber() {
	id, err := ParseIDBinary(suite.BinaryStringID, DefaultEpoch)
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}
	expectedValue := suite.SequenceNumber
	suite.Equal(int16(expectedValue), id.SequenceNumber)
}

func (suite *TestSuit) TestParseBinaryStringToDecimal() {
	id, err := ParseIDBinary(suite.BinaryStringID, DefaultEpoch)
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}
	expectedValue := suite.StringID
	suite.Equal(expectedValue, id.String())
}

func (suite *TestSuit) TestParseStringToBinaryString() {
	id, err := ParseID(suite.StringID, DefaultEpoch)
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}
	expectedValue := suite.BinaryStringID
	suite.Equal(expectedValue, id.StringBinary())
}

func (suite *TestSuit) TestGenerateIDEqualParseID() {
	id, err := suite.S.GenerateID()
	suite.Equal(nil, err, "Failed to generate ID")
	if err != nil {
		return
	}
	idCopy := id.String()

	id2, err := ParseID(idCopy, DefaultEpoch)
	suite.Equal(nil, err, "Failed to parse ID")
	if err != nil {
		return
	}
	suite.Equal(id.String(), id2.String())
}

// ID Generation
func (suite *TestSuit) TestGenerateMany() {
	qty := 1000
	ids := []*SnowID{}
	for range qty {
		id, err := suite.S.GenerateID()
		suite.Equal(nil, err, "Failed to generate ID")
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	suite.Equal(qty, len(ids))
}

func (suite *TestSuit) TestGenerateExceedSequenceLimitIDCount() {
	ids := []*SnowID{}
	timestamp := time.Now().UnixMilli()
	for range SEQUENCE_CAP {
		id, err := suite.S.generateID(timestamp)
		suite.Equal(nil, err, "Failed to generate ID")
		if err != nil {
			return
		}

		ids = append(ids, id)
	}
	suite.Equal(SEQUENCE_CAP, len(ids))

	id, err := suite.S.generateID(timestamp)
	suite.NotEqual(nil, err, "Expected to throw error when exceeding sequence max")
	suite.NotEqual(nil, id, "Expect ID to be nil")
}

func (suite *TestSuit) TestGenerateExceedSequenceLimitIDCountWithSleep() {
	ids := []*SnowID{}
	for range SEQUENCE_CAP {
		id, err := suite.S.GenerateID()
		suite.Equal(nil, err, "Failed to generate ID")
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	time.Sleep(2 * time.Millisecond)
	for range SEQUENCE_CAP {
		id, err := suite.S.GenerateID()
		suite.Equal(nil, err, "Failed to generate ID")
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	// check total
	suite.Equal(2*SEQUENCE_CAP, len(ids))

	// check length
	last := ids[len(ids)-1]

	suite.Equal(TOTAL_BIT, len(last.StringBinary()))
}

func TestDefaultSuite(t *testing.T) {
	suite.Run(t, new(TestSuit))
}
