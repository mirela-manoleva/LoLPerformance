package api

import (
	"encoding/gob"
	"errors"
	"os"
	"time"
)

/*
	Slice to hold all requests that have happened.
*/
var requestRecords []time.Time
var saveFile = "records.tmp"

/*
	Deletes all previously stored records. Used to free space.
*/
func clearAllRecords() {
	// Setting the slice to nil will release the underlying memory to the garbage collector.
	requestRecords = nil
}

func SaveRecords() error {
	if len(requestRecords) == 0 {
		return nil
	}

	fileWriter, err := os.Create(saveFile)
	if err != nil {
		return errors.New("error during encoding - " + err.Error())
	}

	encoder := gob.NewEncoder(fileWriter)
	err = encoder.Encode(requestRecords)
	if err != nil {
		errorMessage := "error during encoding - " + err.Error()
		if err := fileWriter.Close(); err != nil {
			return errors.New("Two errors occured! First is " + errorMessage + ". Second is error while closing file - " + err.Error())
		}
		return errors.New(errorMessage)
	}

	if err := fileWriter.Close(); err != nil {
		return errors.New("error while closing file - " + err.Error())
	}

	return nil
}

func LoadRecords() error {
	fileReader, err := os.Open(saveFile)
	if err != nil {
		return err
	}

	// No point in checking the error as we can't really do anything about it
	defer fileReader.Close()

	decoder := gob.NewDecoder(fileReader)
	err = decoder.Decode(&requestRecords)
	if err != nil {
		return errors.New("error during decoding - " + err.Error())
	}

	return nil
}