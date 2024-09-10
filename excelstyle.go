package main

/*
	File description:
	Defines all the stylistic details around the excel file.
*/

import (
	"github.com/xuri/excelize/v2"
)

const (
	columnDate         = "A"
	columnSummonerName = "B"
	columnRank         = "C"
	columnGoal         = "D"
	columnOutcome      = "E"
	columnRole         = "F"
	columnChampion     = "G"
	columnKills        = "H"
	columnDeaths       = "I"
	columnAssists      = "J"
	columnKP           = "K"
	columnKDA          = "L"
	columnGameLength   = "M"
	columnDPM          = "N"
	columnGPM          = "O"
	columnCSPM         = "P"
	columnFeeling      = "Q"
	columnReview       = "R"
)

const (
	decimalPlaces   = 2
	decimalPlacesKP = 4 // It is displayed as a percent with two decimal places. 0.2234 -> 22.34%
)

var columnNames = []string{
	"Date", "Summoner name", "Rank", "Learning Goal",
	"Outcome", "Role", "Champion",
	"Kills", "Deaths", "Assist", "KP", "KDA",
	"Game Length", "DPM", "GPM", "CSPM",
	"Feeling after game", "Review notes",
}
var columnSizes = []float64{
	13.0, 20.0, 13.0, 40.0,
	12.0, 12.0, 13.0,
	10.0, 10.0, 10.0, 10.0, 10.0,
	13.0, 10.0, 10.0, 10.0,
	75.0, 250.0,
}

const (
	orange     = "#FF9900"
	darkOrange = "#783F0F"
	red        = "#FF0000"
	yellow     = "#FFFF00"
	green      = "#00FF00"
)

const (
	minKP = "0"   // Percent
	midKP = "0.5" // Percent
	maxKP = "1"   // Percent

	minKDA = "0"
	midKDA = "1"
	maxKDA = "5"

	minDPM = "300"
	midDPM = "500"
	maxDPM = "1000"

	minGPM = "250"
	midGPM = "300"
	maxGPM = "360"

	minCSPM = "4"
	midCSPM = "6"
	maxCSPM = "8"
)

var firstRowStyle = excelize.Style{
	Alignment: &excelize.Alignment{Horizontal: "center"},
	Font:      &excelize.Font{Bold: true},
	Fill:      excelize.Fill{Type: "pattern", Color: []string{orange}, Pattern: 1},
	Border: []excelize.Border{
		{Type: "bottom", Color: darkOrange, Style: 5},
		{Type: "top", Color: darkOrange, Style: 5},
		{Type: "left", Color: darkOrange, Style: 5},
		{Type: "right", Color: darkOrange, Style: 5},
	},
}

var styleRegular = excelize.Style{
	Alignment: &excelize.Alignment{Horizontal: "center"},
	Border: []excelize.Border{
		{Type: "bottom", Color: darkOrange, Style: 5},
		{Type: "top", Color: darkOrange, Style: 5},
	},
}

var stylePercent = excelize.Style{
	Alignment: &excelize.Alignment{Horizontal: "center"},
	NumFmt:    10, // 0.00%
	Border: []excelize.Border{
		{Type: "bottom", Color: darkOrange, Style: 5},
		{Type: "top", Color: darkOrange, Style: 5},
	},
}

var styleDate = excelize.Style{
	Alignment: &excelize.Alignment{Horizontal: "center"},
	NumFmt:    15, // d-mmm-yy
	Border: []excelize.Border{
		{Type: "bottom", Color: darkOrange, Style: 5},
		{Type: "top", Color: darkOrange, Style: 5},
	},
}

var styleGameLength = excelize.Style{
	Alignment: &excelize.Alignment{Horizontal: "center"},
	NumFmt:    45, // mm:ss
	Border: []excelize.Border{
		{Type: "bottom", Color: darkOrange, Style: 5},
		{Type: "top", Color: darkOrange, Style: 5},
	},
}

var styleLast = excelize.Style{
	Alignment: &excelize.Alignment{Horizontal: "center"},
	Border: []excelize.Border{
		{Type: "bottom", Color: darkOrange, Style: 5},
		{Type: "top", Color: darkOrange, Style: 5},
		{Type: "right", Color: darkOrange, Style: 5},
	},
}

var styleWin = excelize.Style{
	Fill: excelize.Fill{Type: "pattern", Color: []string{green}, Pattern: 1},
}

var styleLoss = excelize.Style{
	Fill: excelize.Fill{Type: "pattern", Color: []string{red}, Pattern: 1},
}

var conditionalStyleKP = excelize.ConditionalFormatOptions{
	Type:     "3_color_scale",
	Criteria: "=",
	MinType:  "num",
	MidType:  "num",
	MaxType:  "num",
	MinValue: minKP,
	MidValue: midKP,
	MaxValue: maxKP,
	MinColor: red,
	MidColor: yellow,
	MaxColor: green,
}

var conditionalStyleKDA = excelize.ConditionalFormatOptions{
	Type:     "3_color_scale",
	Criteria: "=",
	MinType:  "num",
	MidType:  "num",
	MaxType:  "num",
	MinValue: minKDA,
	MidValue: midKDA,
	MaxValue: maxKDA,
	MinColor: red,
	MidColor: yellow,
	MaxColor: green,
}

var conditionalStyleDPM = excelize.ConditionalFormatOptions{
	Type:     "3_color_scale",
	Criteria: "=",
	MinType:  "num",
	MidType:  "num",
	MaxType:  "num",
	MinValue: minDPM,
	MidValue: midDPM,
	MaxValue: maxDPM,
	MinColor: red,
	MidColor: yellow,
	MaxColor: green,
}

var conditionalStyleGPM = excelize.ConditionalFormatOptions{
	Type:     "3_color_scale",
	Criteria: "=",
	MinType:  "num",
	MidType:  "num",
	MaxType:  "num",
	MinValue: minGPM,
	MidValue: midGPM,
	MaxValue: maxGPM,
	MinColor: red,
	MidColor: yellow,
	MaxColor: green,
}

var conditionalStyleCSPM = excelize.ConditionalFormatOptions{
	Type:     "3_color_scale",
	Criteria: "=",
	MinType:  "num",
	MidType:  "num",
	MaxType:  "num",
	MinValue: minCSPM,
	MidValue: midCSPM,
	MaxValue: maxCSPM,
	MinColor: red,
	MidColor: yellow,
	MaxColor: green,
}

// Must set Format when creating style
var conditionalStyleWin = excelize.ConditionalFormatOptions{
	Type:     "cell",
	Criteria: "equal to",
	Value:    "\"Win\"",
}

// Must set Format when creating style
var conditionalStyleLoss = excelize.ConditionalFormatOptions{
	Type:     "cell",
	Criteria: "equal to",
	Value:    "\"Loss\"",
}
