package nibbler_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	mock "github.com/blorticus/go-test-mocks"
	"github.com/blorticus/nibbler"
)

type nibblerExpectedResult struct {
	expectedByte  byte
	expectAnError bool
	expectEOF     bool
}

type nibblerTestCase struct {
	operation      string // "ReadByte", "UnreadByte", "PeekAtNextByte"
	expectedResult *nibblerExpectedResult
}

func (testCase *nibblerTestCase) runTestCaseAgainst(nibber nibbler.ByteNibbler) (testCaseError error) {
	var err error
	var b byte

	switch testCase.operation {
	case "ReadByte":
		b, err = nibber.ReadByte()

	case "PeekAtNextByte":
		b, err = nibber.PeekAtNextByte()

	case "UnreadByte":
		err = nibber.UnreadByte()
	}

	if testCase.expectedResult.expectEOF {
		if err == nil {
			return fmt.Errorf("expected EOF, got no error returned")
		} else if err != io.EOF {
			return fmt.Errorf("expected EOF, got different error returned (%s)", err.Error())
		}
	} else {
		if err == io.EOF {
			return fmt.Errorf("expected no EOF, got EOF")
		}

		if testCase.expectedResult.expectAnError {
			if err == nil {
				return fmt.Errorf("expected an error, no error returned")
			}
		} else if err != nil {
			return fmt.Errorf("expected no error, got an error (%s)", err.Error())
		}
	}

	if !testCase.expectedResult.expectAnError && !testCase.expectedResult.expectEOF {
		if testCase.expectedResult.expectedByte != b {
			return fmt.Errorf("expected byte value (%d), got (%d)", testCase.expectedResult.expectedByte, b)
		}
	}

	return nil
}

func TestByteSliceNibbler(t *testing.T) {
	nib := nibbler.NewByteSliceNibbler([]byte{})

	for testIndex, testCase := range []*nibblerTestCase{
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: true}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: true}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: true, expectEOF: false}},
	} {
		if err := testCase.runTestCaseAgainst(nib); err != nil {
			t.Errorf("(ByteSliceNibbler with Empty slice) (test %d) %s", testIndex+1, err.Error())
		}

	}

	nib = nibbler.NewByteSliceNibbler([]byte{0, 1, 2, 3, 4, 5})
	for testIndex, testCase := range []*nibblerTestCase{
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: true, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: false}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: false, expectEOF: false}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: true, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 1, expectAnError: false, expectEOF: false}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 2, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 2, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 3, expectAnError: false, expectEOF: false}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 4, expectAnError: false, expectEOF: false}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: false, expectEOF: false}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 3, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 3, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 4, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 5, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: true}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: true}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: false, expectEOF: false}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 5, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 5, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: true}},
	} {
		if err := testCase.runTestCaseAgainst(nib); err != nil {
			t.Errorf("(ByteSliceNibbler with 6 values in slice) (test %d) %s", testIndex+1, err.Error())
		}

	}
}

func TestByteReaderNibbler(t *testing.T) {
	reader := mock.NewReader().AddGoodRead([]byte{0, 1, 2, 3}).AddGoodRead([]byte{4, 5}).AddEOF()
	nibbler := nibbler.NewByteReaderNibbler(reader)
	for testIndex, testCase := range []*nibblerTestCase{
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: true, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: false}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: false, expectEOF: false}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: true, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 1, expectAnError: false, expectEOF: false}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 2, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 2, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 3, expectAnError: false, expectEOF: false}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 4, expectAnError: false, expectEOF: false}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: false, expectEOF: false}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 3, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 3, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 4, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 5, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: true}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: true}},
		{operation: "UnreadByte", expectedResult: &nibblerExpectedResult{expectAnError: false, expectEOF: false}},
		{operation: "PeekAtNextByte", expectedResult: &nibblerExpectedResult{expectedByte: 5, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 5, expectAnError: false, expectEOF: false}},
		{operation: "ReadByte", expectedResult: &nibblerExpectedResult{expectedByte: 0, expectAnError: false, expectEOF: true}},
	} {
		if err := testCase.runTestCaseAgainst(nibbler); err != nil {
			t.Errorf("(ByteReaderNibbler with 6 values then EOF) (test %d) %s", testIndex+1, err.Error())
		}

	}
}

type namedCharacterSetMatchTestCase struct {
	matchType             string // "matching" or "nonMatching"
	setName               string
	expectError           bool
	expectEOF             bool
	expectedReturnedBytes []byte
}

func (testCase *namedCharacterSetMatchTestCase) runTestCaseAgainstNibbler(nib nibbler.ByteNibbler) error {
	var returnedBytes []byte
	var returnedError error

	switch testCase.matchType {
	case "matching":
		returnedBytes, returnedError = nib.ReadNextBytesMatchingSet(testCase.setName)

	case "notMatching":
		returnedBytes, returnedError = nib.ReadNextBytesNotMatchingSet(testCase.setName)

	default:
		return fmt.Errorf("internal test error")
	}

	if returnedError != nil {
		if returnedError == io.EOF {
			if !testCase.expectEOF {
				return fmt.Errorf("expected error, got EOF")
			}

			return nil
		} else {
			if testCase.expectEOF {
				return fmt.Errorf("expected EOF, got error = (%s)", returnedError.Error())
			}

			return nil
		}
	} else if testCase.expectEOF {
		return fmt.Errorf("expected EOF, got no EOF")
	} else if testCase.expectError {
		return fmt.Errorf("expected error, got no error")
	}

	if bytes.Compare(testCase.expectedReturnedBytes, returnedBytes) != 0 {
		var returnedBytesInError string
		var expectedBytesInError string

		if len(returnedBytes) > 10 {
			returnedBytesInError = byteSliceToSanitizedString(returnedBytes[:10]) + "..."
		} else {
			returnedBytesInError = byteSliceToSanitizedString(returnedBytes)
		}

		if len(testCase.expectedReturnedBytes) > 10 {
			expectedBytesInError = byteSliceToSanitizedString(testCase.expectedReturnedBytes[:10]) + "..."
		} else {
			expectedBytesInError = byteSliceToSanitizedString(testCase.expectedReturnedBytes)
		}

		return fmt.Errorf("returned bytes (%s) do not match expected bytes (%s)", returnedBytesInError, expectedBytesInError)
	}

	return nil
}

func byteSliceToSanitizedString(preStream []byte) string {
	sanitized := ""
	for _, b := range preStream {
		switch b {
		case '\t':
			sanitized += `\t`

		case '\n':
			sanitized += `\n`

		case '\r':
			sanitized += `\r`

		default:
			if b < 32 || b > 127 {
				sanitized += `\x` + fmt.Sprintf("%02x", b)
			} else {
				sanitized += string(b)
			}
		}
	}

	return sanitized
}

// test against stream: "abc \tD12\r21D "
func testAnyByteNibblerForNamedSetMatchers(byteNibbler nibbler.ByteNibbler, baseTestName string, t *testing.T) {
	byteSetMap := nibbler.NewNamedCharacterSetsMap().
		AddNamedCharacterSetFromString("set-abcdefg", "abcdefg").
		AddNamedCharacterSetFromByteArray("set-12", []byte{'1', '2'}).
		AddNamedCharacterSetFromString("set-whitespace", " \t\r\n")

	testCases := []*namedCharacterSetMatchTestCase{
		{matchType: "matching", setName: "foo", expectError: true},
		{matchType: "matching", setName: "set-12", expectedReturnedBytes: []byte{}},
		{matchType: "notMatching", setName: "foo", expectError: true},
		{matchType: "notMatching", setName: "set-abcdefg", expectedReturnedBytes: []byte{}},
		{matchType: "matching", setName: "set-abcdefg", expectedReturnedBytes: []byte{'a', 'b', 'c'}},
		{matchType: "matching", setName: "set-abcdefg", expectedReturnedBytes: []byte{}},
		{matchType: "notMatching", setName: "set-whitespace", expectedReturnedBytes: []byte{}},
		{matchType: "notMatching", setName: "set-12", expectedReturnedBytes: []byte{' ', '\t', 'D'}},
		{matchType: "matching", setName: "set-12", expectedReturnedBytes: []byte{'1', '2'}},
		{matchType: "matching", setName: "set-123", expectError: true},
		{matchType: "matching", setName: "set-whitespace", expectedReturnedBytes: []byte{'\r'}},
		{matchType: "notMatching", setName: "set-12", expectedReturnedBytes: []byte{}},
		{matchType: "matching", setName: "set-12", expectedReturnedBytes: []byte{'2', '1'}},
		{matchType: "notMatching", setName: "set-whitespace", expectedReturnedBytes: []byte{'D'}},
		{matchType: "matching", setName: "set-whitespace", expectedReturnedBytes: []byte{' '}, expectEOF: true},
		{matchType: "matching", setName: "set-whitespace", expectedReturnedBytes: []byte{}, expectEOF: true},
		{matchType: "matching", setName: "set-12", expectedReturnedBytes: []byte{}, expectEOF: true},
		{matchType: "matching", setName: "set-abcdefg", expectedReturnedBytes: []byte{}, expectEOF: true},
	}

	if _, err := byteNibbler.ReadNextBytesMatchingSet("set-12"); err == nil {
		t.Errorf("(%s) (before adding CharacterSetMap) expected error on ReadNextBytesMatchingSet, got no error", baseTestName)
	}

	if _, err := byteNibbler.ReadNextBytesNotMatchingSet("set-12"); err == nil {
		t.Errorf("(%s) (before adding CharacterSetMap) expected error on ReadNextBytesNotMatchingSet, got no error", baseTestName)
	}

	byteNibbler.AddNamedCharacterSetsMap(byteSetMap)

	for testCaseIndex, testCase := range testCases {
		if err := testCase.runTestCaseAgainstNibbler(byteNibbler); err != nil {
			t.Errorf("(%s) (testSet index %d) %s", baseTestName, testCaseIndex, err.Error())
		}
	}
}

func TestByteSliceNibblerNamedSet(t *testing.T) {
	byteSliceNibbler := nibbler.NewByteSliceNibbler([]byte("abc \tD12\r21D "))
	testAnyByteNibblerForNamedSetMatchers(byteSliceNibbler, "TestByteSliceNibblerNamedSet", t)
}

func TestByteReaderNibblerNamedSet(t *testing.T) {
	completeStream := "abc \tD12\r21D "
	reader := mock.NewReader().
		AddGoodRead([]byte(completeStream[0:4])).
		AddGoodRead([]byte(completeStream[4:11])).
		AddGoodRead([]byte(completeStream[11:12])).
		AddGoodRead([]byte(completeStream[12:13])).
		AddEOF()

	byteReaderNibbler := nibbler.NewByteReaderNibbler(reader)
	testAnyByteNibblerForNamedSetMatchers(byteReaderNibbler, "TestByteSliceNibblerNamedSet", t)
}

func TestByteReaderNibblerNamedSetWithEmptyRead(t *testing.T) {
	// When a ByteReaderNibbler attempts an underlying Read(), and that returns no error, no EOF but
	// also no data, an error should be thrown.  Ensure that is is reflected here.
	completeStream := "abc \tD12\r21D "

	reader := mock.NewReader().
		AddGoodRead([]byte(completeStream[0:4])).
		AddEmptyRead().
		AddGoodRead([]byte(completeStream[4:11])).
		AddGoodRead([]byte(completeStream[11:12])).
		AddGoodRead([]byte(completeStream[12:13])).
		AddEOF()

	byteNibbler := nibbler.NewByteReaderNibbler(reader)

	byteSetMap := nibbler.NewNamedCharacterSetsMap().
		AddNamedCharacterSetFromString("set-abcdefg", "abcdefg").
		AddNamedCharacterSetFromByteArray("set-12", []byte{'1', '2'}).
		AddNamedCharacterSetFromString("set-whitespace", " \t\r\n")

	byteNibbler.AddNamedCharacterSetsMap(byteSetMap)

	testCases := []*namedCharacterSetMatchTestCase{
		{matchType: "matching", setName: "set-12", expectedReturnedBytes: []byte{}},
		{matchType: "matching", setName: "set-abcdefg", expectedReturnedBytes: []byte{'a', 'b', 'c'}},
		{matchType: "notMatching", setName: "set-12", expectError: true},
	}

	for testCaseIndex, testCase := range testCases {
		if err := testCase.runTestCaseAgainstNibbler(byteNibbler); err != nil {
			t.Errorf("(TestByteReaderNibblerNamedSetWithEmptyRead) (testSet index %d) %s", testCaseIndex, err.Error())
		}
	}
}
