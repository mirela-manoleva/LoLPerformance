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

	if err := setColumnFormat(file, sheetName); err != nil {
		return fmt.Errorf("error creating new file \"%s\" - %s", fileName, err.Error())
	}

	if err := setFirstRow(file, sheetName); err != nil {
		return fmt.Errorf("error creating new file \"%s\" - %s", fileName, err.Error())
	}

	if err := file.SaveAs(fileName); err != nil {
		return err
	}

	return nil
}

/*
Used if the user wants to create a new sheet in an existing file.
*/
func AddSheet(fileName string, sheetName string) (Err error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			Err = errors.Join(Err, err)
		}
	}()

	index, err := file.NewSheet(sheetName)
	if err != nil {
		return err
	}
	file.SetActiveSheet(index)

	if err := setColumnFormat(file, sheetName); err != nil {
		return fmt.Errorf("error adding a new sheet \"%s\" - %s", sheetName, err.Error())
	}

	if err := setFirstRow(file, sheetName); err != nil {
		return fmt.Errorf("error adding a new sheet \"%s\" - %s", sheetName, err.Error())
	}

	if err := file.Save(); err != nil {
		return err
	}

	return nil
}

/*
Adds a new row with all the game information and formats the data.
*/
func AddGameRecord(fileName string, sheetName string, game GameData, summoner SummonerData) (Err error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			Err = errors.Join(Err, err)
		}
	}()

	/*
		To get the number of the row where we want to write we check which the first empty cell in the columnDate is.
		That is because columnImprove can be written ahead of time. Checking the last free row won't work.
		It's assuming that columnImprove != columnDate.
	*/
	columnDateIndex, err := excelize.ColumnNameToNumber(columnDate)
	if err != nil {
		return fmt.Errorf("error adding a new game record - %s", err.Error())
	}
	columnDateIndex-- // Will be traversing an array with it.

	cols, err := file.GetCols(sheetName)
	if err != nil {
		return fmt.Errorf("error adding a new game record - %s", err.Error())
	}

	var row string
	for i := 0; i < len(cols[columnDateIndex]); i++ {
		if len(cols[columnDateIndex][i]) == 0 {
			row = strconv.Itoa(i + 1)
			break
		}
	}
	if row == "" { // all columns so far are full
		row = strconv.Itoa(len(cols[columnDateIndex]) + 1)
	}

	err = setValuesNewRow(file, sheetName, row, game, summoner)
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
