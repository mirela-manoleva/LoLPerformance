package main

/*
	File description:
	Defines functionality for maintaining a game records excel sheet.
	The user can create a file, sheet or add a game.
*/

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

/*
Creates excel file and adds a formated sheet.
*/
func CreateGameRecordFile(fileName string, sheetName string) (Err error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			Err = errors.Join(Err, err)
		}
	}()

	if err := file.SetSheetName("Sheet1", sheetName); err != nil {
		return err
	}

	if err := setColumnWidth(file, sheetName); err != nil {
		return fmt.Errorf("error creating new file \"%s\" - %s", fileName, err.Error())
	}

	if err := setColumnNames(file, sheetName); err != nil {
		return fmt.Errorf("error creating new file \"%s\" - %s", fileName, err.Error())
	}

	if err := file.SaveAs(fileName); err != nil {
		return err
	}

	return nil
}

/*
Adds a new row with all the game information and formats the data.
*/
func AddGameRecord(fileName string, sheetName string, game PlayerGameData, rank Rank, summonerName string) (Err error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			Err = errors.Join(Err, err)
		}
	}()

	// Find the row where we want to add the game
	var row string

	if sheetIndex, err := file.GetSheetIndex(sheetName); err != nil {
		return err
	} else if sheetIndex == -1 {
		err := AddSheet(file, sheetName)
		if err != nil {
			return err
		}
		row = "2"
	} else {
		/*
			If the sheet exists we check which the first empty cell in the columnDate is.
			That is because some columns can be written ahead of time. The columnDate cannot be written ahead of time.
			Unless you're a seer.
		*/
		columnDateIndex, err := excelize.ColumnNameToNumber(columnDate)
		if err != nil {
			return fmt.Errorf("error adding a new game record - %s", err.Error())
		}
		columnDateIndex-- // Will be traversing an array with it.

		columns, err := file.GetCols(sheetName)
		if err != nil {
			return fmt.Errorf("error adding a new game record - %s", err.Error())
		}

		for i := 0; i < len(columns[columnDateIndex]); i++ {
			if len(columns[columnDateIndex][i]) == 0 {
				row = strconv.Itoa(i + 1)
				break
			}
		}
		if row == "" { // all columns so far are full
			row = strconv.Itoa(len(columns[columnDateIndex]) + 1)
		}
	}

	err = setValuesNewRow(file, sheetName, row, game, rank, summonerName)
	if err != nil {
		return err
	}

	err = setStyleNewRow(file, sheetName, row)
	if err != nil {
		return err
	}

	if err := file.Save(); err != nil {
		return err
	}

	return nil
}
