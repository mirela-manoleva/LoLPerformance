package main

/*
	File description:
	Helper functions used in excel.go.
	Mainly used to reduce code size and improve readability.

	Note that the functions don't save the file.
	That is intentional. We don't want any changes unless all changes are possible.
	The functions that use the helpers are the ones that should save.

	They also take excelize.File as argument without checking if it's valid.
	That is also intentional, the helper isn't responsible for it, the function that uses the helper is.

	Another note is that the errors are not returned "as is" because the excelize package has very vague descriptions.
*/

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	rowFirst    = "1"
	columnFirst = columnDate
	columnLast  = columnReview
)

/*
Riot API returns a number as a queue type.
This function returns the queue type in a format we want for the excel table.
Reference on them can be found here:
https://static.developer.riotgames.com/docs/lol/queues.json
*/
func queueToString(queueType int) string {
	switch queueType {
	case 400:
		return "Draft Pick"
	case 420:
		return "Ranked Solo"
	case 430:
		return "Blind Pick"
	case 440:
		return "Ranked Flex"
	case 450:
		return "ARAM"
	case 490:
		return "Quickplay"
	case 700:
		return "Clash"
	case 720:
		return "ARAM Clash"
	default:
		return "Other"
	}
}

/*
Riot API returns bool type.
This function returns the outcome in a format we want for the excel table.
*/
func outcomeToString(isWin bool) string {
	if isWin {
		return "Win"
	} else {
		return "Loss"
	}
}

/*
Used when creating new sheet.
Sets the column width.
*/
func setColumnFormat(file *excelize.File, sheetName string) error {
	for i := 0; i < len(columnNames); i++ {
		if i >= len(columnSizes) { // Rest of the columns are left default
			return nil
		}

		currentColumn, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return fmt.Errorf("error finding column name at index %d - %s", i+1, err.Error())
		}
		err = file.SetColWidth(sheetName, currentColumn, currentColumn, columnSizes[i])
		if err != nil {
			return fmt.Errorf("error setting column %s width to %f - %s", currentColumn, columnSizes[i], err.Error())
		}
	}

	return nil
}

/*
Used when creating new sheet.
Sets the header of the sheet with all the column names.
*/
func setFirstRow(file *excelize.File, sheetName string) error {
	firstCell := columnFirst + rowFirst
	lastCell := columnLast + rowFirst

	err := file.SetSheetRow(sheetName, firstCell, &columnNames)
	if err != nil {
		return fmt.Errorf("error setting the values of the column names in sheet \"%s\" - %s", sheetName, err.Error())
	}

	columnNamesStyle, err := file.NewStyle(&firstRowStyle)
	if err != nil {
		return fmt.Errorf("error creating style for column names - %s", err.Error())
	}

	err = file.SetCellStyle(sheetName, firstCell, lastCell, columnNamesStyle)
	if err != nil {
		return fmt.Errorf("error setting style for the column names - %s", err.Error())
	}

	return nil
}

/*
Used when adding a new row.
*/
func setValuesNewRow(file *excelize.File, sheetName string, row string, game GameData, summoner SummonerData) error {
	err := file.SetCellValue(sheetName, columnDate+row, time.Unix(0, game.unixStartTimestamp*int64(time.Millisecond)))
	if err != nil {
		return fmt.Errorf("error when setting date on row %s - %s", row, err.Error())
	}

	err = file.SetCellStr(sheetName, columnRank+row, summoner.rank)
	if err != nil {
		return fmt.Errorf("error when setting rank on row %s - %s", row, err.Error())
	}

	err = file.SetCellStr(sheetName, columnQueueType+row, queueToString(game.queueType))
	if err != nil {
		return fmt.Errorf("error when setting queue type on row %s - %s", row, err.Error())
	}

	err = file.SetCellStr(sheetName, columnOutcome+row, outcomeToString(game.isWin))
	if err != nil {
		return fmt.Errorf("error when setting outcome on row %s - %s", row, err.Error())
	}

	err = file.SetCellStr(sheetName, columnRole+row, game.role)
	if err != nil {
		return fmt.Errorf("error when setting role on row %s - %s", row, err.Error())
	}

	err = file.SetCellStr(sheetName, columnChampion+row, game.champion)
	if err != nil {
		return fmt.Errorf("error when setting champion on row %s - %s", row, err.Error())
	}

	err = file.SetCellInt(sheetName, columnKills+row, game.kills)
	if err != nil {
		return fmt.Errorf("error when setting kills on row %s - %s", row, err.Error())
	}

	err = file.SetCellInt(sheetName, columnDeaths+row, game.deaths)
	if err != nil {
		return fmt.Errorf("error when setting deaths on row %s - %s", row, err.Error())
	}

	err = file.SetCellInt(sheetName, columnAssists+row, game.assists)
	if err != nil {
		return fmt.Errorf("error when setting assists on row %s - %s", row, err.Error())
	}

	err = file.SetCellFloat(sheetName, columnKP+row, game.killParticipation, decimalPlacesKP, 64)
	if err != nil {
		return fmt.Errorf("error when setting KP on row %s - %s", row, err.Error())
	}

	err = file.SetCellFloat(sheetName, columnKDA+row, game.kda, decimalPlaces, 64)
	if err != nil {
		return fmt.Errorf("error when setting KDA on row %s - %s", row, err.Error())
	}

	// Note the 0.5 way of rounding. It work only cos we don't care much about the accuracy of +/-0.000001 seconds
	err = file.SetCellValue(sheetName, columnGameLength+row, time.Duration(int64(game.gameLengthSeconds+0.5)*int64(time.Second)))
	if err != nil {
		return fmt.Errorf("error when setting game length on row %s - %s", row, err.Error())
	}

	err = file.SetCellFloat(sheetName, columnDPM+row, game.damagePerMinute, decimalPlaces, 64)
	if err != nil {
		return fmt.Errorf("error when setting DPM on row %s - %s", row, err.Error())
	}

	err = file.SetCellFloat(sheetName, columnGPM+row, game.goldPerMinute, decimalPlaces, 64)
	if err != nil {
		return fmt.Errorf("error when setting GPM on row %s - %s", row, err.Error())
	}

	err = file.SetCellFloat(sheetName, columnCSPM+row, game.csPerMinute, decimalPlaces, 64)
	if err != nil {
		return fmt.Errorf("error when setting CSPM on row %s - %s", row, err.Error())
	}

	return nil
}

/*
Used when adding a new row.
*/
func setStyleNewRow(file *excelize.File, sheetName string, row string) error {
	firstCell := columnFirst + row
	lastCell := columnLast + row

	styleRegularID, err := file.NewStyle(&styleRegular)
	if err != nil {
		return fmt.Errorf("error creating style for regular cells - %s", err.Error())
	}

	stylePercentID, err := file.NewStyle(&stylePercent)
	if err != nil {
		return fmt.Errorf("error creating style for percent cells - %s", err.Error())
	}

	styleDateID, err := file.NewStyle(&styleDate)
	if err != nil {
		return fmt.Errorf("error creating style for the date cell - %s", err.Error())
	}

	styleGameLengthID, err := file.NewStyle(&styleGameLength)
	if err != nil {
		return fmt.Errorf("error creating style for the game length cell - %s", err.Error())
	}

	styleLastID, err := file.NewStyle(&styleLast)
	if err != nil {
		return fmt.Errorf("error creating style for the last cells - %s", err.Error())
	}

	styleWinID, err := file.NewConditionalStyle(&styleWin)
	if err != nil {
		return fmt.Errorf("error creating conditional style for outcome cell (Win) - %s", err.Error())
	}

	styleLossID, err := file.NewConditionalStyle(&styleLoss)
	if err != nil {
		return fmt.Errorf("error creating conditional style for outcome cell (Loss) - %s", err.Error())
	}

	err = file.SetCellStyle(sheetName, firstCell, lastCell, styleRegularID)
	if err != nil {
		return fmt.Errorf("error setting regular style for row - %s", err.Error())
	}

	err = file.SetCellStyle(sheetName, columnDate+row, columnDate+row, styleDateID)
	if err != nil {
		return fmt.Errorf("error setting style for the date cell on row - %s", err.Error())
	}

	err = file.SetCellStyle(sheetName, columnGameLength+row, columnGameLength+row, styleGameLengthID)
	if err != nil {
		return fmt.Errorf("error setting style for the game length cell on row - %s", err.Error())
	}

	err = file.SetCellStyle(sheetName, columnKP+row, columnKP+row, stylePercentID)
	if err != nil {
		return fmt.Errorf("error setting style for the date cell on row - %s", err.Error())
	}

	err = file.SetCellStyle(sheetName, lastCell, lastCell, styleLastID)
	if err != nil {
		return fmt.Errorf("error setting style for the last cell on row - %s", err.Error())
	}

	err = file.SetConditionalFormat(sheetName, columnKP+row, []excelize.ConditionalFormatOptions{conditionalStyleKP})
	if err != nil {
		return fmt.Errorf("error setting conditional format for the KP cell on row - %s", err.Error())
	}

	err = file.SetConditionalFormat(sheetName, columnKDA+row, []excelize.ConditionalFormatOptions{conditionalStyleKDA})
	if err != nil {
		return fmt.Errorf("error setting conditional format for the KDA cell on row - %s", err.Error())
	}

	err = file.SetConditionalFormat(sheetName, columnDPM+row, []excelize.ConditionalFormatOptions{conditionalStyleDPM})
	if err != nil {
		return fmt.Errorf("error setting conditional format for the DPM cell on row - %s", err.Error())
	}

	err = file.SetConditionalFormat(sheetName, columnGPM+row, []excelize.ConditionalFormatOptions{conditionalStyleGPM})
	if err != nil {
		return fmt.Errorf("error setting conditional format for the GPM cell on row - %s", err.Error())
	}

	err = file.SetConditionalFormat(sheetName, columnCSPM+row, []excelize.ConditionalFormatOptions{conditionalStyleCSPM})
	if err != nil {
		return fmt.Errorf("error setting conditional format for the CSPM cell on row - %s", err.Error())
	}

	conditionalStyleWin.Format = styleWinID
	conditionalStyleLoss.Format = styleLossID

	err = file.SetConditionalFormat(sheetName, columnOutcome+row, []excelize.ConditionalFormatOptions{conditionalStyleWin, conditionalStyleLoss})
	if err != nil {
		return fmt.Errorf("error setting conditional formats for the Outcome cell on row - %s", err.Error())
	}

	return nil
}
