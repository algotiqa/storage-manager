//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package backend

import (
	"encoding/json"
	"github.com/tradalia/core"
	"github.com/tradalia/storage-manager/pkg/app"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//=============================================================================

const (
	Code   = "code"
	Image  = "image"
	Report = "report"

	InfoFile    = "info.json"
	DocFile     = "documentation.txt"
	EquityChart = "equity-chart.png"
)

var Dirs = []string{ Code, Image, Report }

//=============================================================================

var folder         string
var defEquityChart []byte

//=============================================================================
//===
//=== Init functions
//===
//=============================================================================

func InitStorage(cfg *app.Config) {
	folder = cfg.Storage.Folder

	err := os.MkdirAll(folder, 0700)
	core.ExitIfError(err)

	defEquityChart, err = os.ReadFile("default/"+ EquityChart)
	core.ExitIfError(err)
}

//=============================================================================
//===
//=== Public functionsk
//===
//=============================================================================

func AddTradingSystem(ts *TradingSystem) error {
	sId := strconv.Itoa(int(ts.Id))
	path:= folder +"/"+ ts.Username +"/"+ sId +"/"

	for _, dir := range Dirs {
		err := os.MkdirAll(path + dir, 0700)
		if err != nil {
			return err
		}
	}

	err := SetTradingSystemInfo(ts)
	if err != nil {
		return err
	}

	return SetTradingSystemDoc(ts.Username, ts.Id, "")
}

//=============================================================================

func UpdateTradingSystem(ts *TradingSystem) error {
	return SetTradingSystemInfo(ts)
}

//=============================================================================

func DeleteTradingSystem(id uint, username string) error {
	sId := strconv.Itoa(int(id))
	return os.RemoveAll(folder +"/"+ username +"/"+ sId)
}

//=============================================================================
//=== Equity chart
//=============================================================================

func GetEquityChartTypes(username string, id uint) ([]string, error) {
	path := []string{
		folder,
		username,
		strconv.Itoa(int(id)),
	}

	files,err := getFiles(path...)

	if err != nil {
		return nil, err
	}

	var types []string

	for _, file := range files {
		if isEquityChartName(file.Name()) {
			types = append(types, getChartType(file.Name()))
		}
	}

	return types, nil
}

//=============================================================================

func ReadEquityChart(username string, id uint, chartType string) ([]byte,error) {
	path := []string{
		folder,
		username,
		strconv.Itoa(int(id)),
		buildEquityChartName(chartType),
	}

	return readFile(path...)
}

//=============================================================================

func WriteEquityChart(username string, id uint, data []byte, chartType string) error {
	path := []string{
		folder,
		username,
		strconv.Itoa(int(id)),
		buildEquityChartName(chartType),
	}

	return writeFile(data, path...)
}

//=============================================================================

func DeleteEquityChart(username string, id uint, chartType string) error {
	path := []string{
		folder,
		username,
		strconv.Itoa(int(id)),
		buildEquityChartName(chartType),
	}

	return deleteFile(path...)
}

//=============================================================================

func GetDefaultEquityChart() []byte {
	return defEquityChart
}

//=============================================================================
//=== Documentation
//=============================================================================

func GetTradingSystemDoc(username string, id uint) (string, error) {
	path := []string{
		folder,
		username,
		strconv.Itoa(int(id)),
		DocFile,
	}

	data, err := readFile(path...)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

//=============================================================================

func SetTradingSystemDoc(username string, id uint, doc string) error {
	path := []string{
		folder,
		username,
		strconv.Itoa(int(id)),
		DocFile,
	}

	return writeFile([]byte(doc), path...)
}

//=============================================================================
//=== Information
//=============================================================================

func GetTradingSystemInfo(username string, id uint) (*TradingSystem, error) {
	path := []string{
		folder,
		username,
		strconv.Itoa(int(id)),
		InfoFile,
	}

	data, err := readFile(path...)
	if err != nil {
		return nil, err
	}

	ts := TradingSystem{}
	err = json.Unmarshal(data, &ts)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}

//=============================================================================

func SetTradingSystemInfo(ts *TradingSystem) error {
	path := []string{
		folder,
		ts.Username,
		strconv.Itoa(int(ts.Id)),
		InfoFile,
	}

	data, err := json.Marshal(ts)
	if err != nil {
		return err
	}

	return writeFile(data, path...)
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getFiles(path ...string) ([]os.DirEntry, error) {
	dir := filepath.Join(path...)
	return os.ReadDir(dir)
}

//=============================================================================

func readFile(path ...string) ([]byte, error) {
	file := filepath.Join(path...)
	return os.ReadFile(file)
}

//=============================================================================

func writeFile(data []byte, path ...string) error {
	file := filepath.Join(path...)
	err  := os.WriteFile(file +".temp", data, 0600)

	if err != nil {
		return err
	}

	_,err = os.Stat(file)

	if err == nil {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}

	return os.Rename(file +".temp", file)
}

//=============================================================================

func deleteFile(path ...string) error {
	file := filepath.Join(path...)
	return os.Remove(file)
}

//=============================================================================

func buildEquityChartName(chartType string) string {
	return chartType +"-"+ EquityChart
}

//=============================================================================

func isEquityChartName(fileName string) bool {
	return strings.HasSuffix(fileName, "-"+ EquityChart)
}

//=============================================================================

func getChartType(fileName string) string {
	index := strings.Index(fileName, "-"+ EquityChart)
	return fileName[0:index]
}

//=============================================================================
