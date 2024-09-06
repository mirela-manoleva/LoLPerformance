package limiter

/*
	File description:
	Defines the functionality record requests made in order to check if any limits are broken.
	Also defines the functions that save and load requests from a file to ensure that new requests
	are not breaking any limits because of previous executions of the program.

	Note:
	The user of the package shouldn't know about these records.
	Don't export anything around records!
*/

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"time"
)

/*
	Name of the file where records are written into and read from.
*/
var saveFilePath = filepath.Join("temp", "records.tmp")

/*
	Slice to hold all requests that have happened.
*/
var requestRecords []time.Time

/*
	Adds a records to requestRecords.
	If only addRecord is used to add records the slice requestRecords is ensured to be sorted.
*/
func addRecord(record time.Time) {
	if len(requestRecords) != 0 && record.Before(requestRecords[len(requestRecords)-1])  {
		i := 0
		for ; i < len(requestRecords) ; i++ {
			if record.Before(requestRecords[i])  {
				requestRecords = slices.Insert(requestRecords, i, record)
				return
			}
		}
	}
	requestRecords = append(requestRecords, record)
}

/*
	Deletes all previously stored records. Used to free space.
*/
func clearAllRecords() {
	// Setting the slice to nil will release the underlying memory to the garbage collector.
	requestRecords = nil
}

/*
	Saves the recent requests made into a temporary file to ensure that the next start of the program doesn't break any defined limits.
	Freely use as it overwrittes the temporary file.
*/
func SaveRequestsMade() error {
	if len(requestRecords) == 0 {
		return nil
	}

	err := os.MkdirAll(filepath.Dir(saveFilePath), 0777)
	if err != nil {
		return err
	}

	fileWriter, err := os.Create(saveFilePath)
	if err != nil {
		return err
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

/*
	Loads the requests made by a previous iteration of the program.
	Should be used to ensure that the first few API calls don't break any set limits.
	This can happen if the previous iteration of the program made alot of API calls before it closed.
	Use it at the start of the program.
*/
func LoadRequestsMade() error {
	fileReader, err := os.Open(saveFilePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
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