//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"archive/zip"
	"bytes"
	"io"
	"strconv"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/storage-manager/pkg/backend"
)

//=============================================================================

func GetDocumentation(c *auth.Context, id uint) (*DocumentationResponse, error) {
	c.Log.Info("GetDocumentation: Getting documentation for trading system", "id", id)

	doc, err := backend.GetTradingSystemDoc(c.Session.Username, id)
	if err != nil {
		c.Log.Error("GetDocumentation: Cannot retrieve documentation for trading system", "id", id, "error", err)
		return nil, err
	}

	var info *backend.TradingSystem
	info, err = backend.GetTradingSystemInfo(c.Session.Username, id)
	if err != nil {
		c.Log.Error("GetDocumentation: Cannot retrieve info for trading system", "id", id, "error", err)
		return nil, err
	}

	c.Log.Info("GetDocumentation: Operation complete", "id", id)

	return &DocumentationResponse{
		Id:            id,
		Name:          info.Name,
		Documentation: doc,
	}, nil
}

//=============================================================================

func SetDocumentation(c *auth.Context, id uint, r *DocumentationRequest) error {
	c.Log.Info("SetDocumentation: Setting documentation for trading system", "id", id)
	err := backend.SetTradingSystemDoc(c.Session.Username, id, r.Documentation)

	if err != nil {
		c.Log.Info("SetDocumentation: Cannot store documentation for trading system", "id", id, "error", err)
	} else {
		c.Log.Info("SetDocumentation: Operation complete", "id", id)
	}

	return err
}

//=============================================================================

func GetEquityChart(c *auth.Context, id uint, chartType string) ([]byte, error) {
	data, err := backend.ReadEquityChart(c.Session.Username, id, chartType)

	if err != nil {
		return backend.GetDefaultEquityChart(), nil
	}

	return data, err
}

//=============================================================================
// Called by Portfolio trader

func SetEquityCharts(c *auth.Context, id uint, r *EquityRequest) error {
	c.Log.Info("SetEquityCharts: Setting equity charts for trading system", "id", id)

	for chartType, data := range r.Images {
		err := backend.WriteEquityChart(r.Username, id, data, chartType)
		if err != nil {
			c.Log.Info("SetEquityCharts: Can't write equity chart", "id", id, "error", err, "type", chartType)
			return err
		}
	}

	c.Log.Info("SetEquityCharts: Equity charts set", "id", id)
	return nil
}

//=============================================================================
// Called by Portfolio trader

func DeleteEquityCharts(c *auth.Context, id uint, r *EquityRequest) error {
	c.Log.Info("DeleteEquityCharts: Delete equity chart for trading system", "id", id, "username", r.Username)

	types, err := backend.GetEquityChartTypes(r.Username, id)
	if err == nil {
		for _, ct := range types {
			err = backend.DeleteEquityChart(r.Username, id, ct)

			if err != nil {
				c.Log.Error("DeleteEquityCharts: Cannot delete equity chart", "id", id, "username", r.Username, "error", err, "type", ct)
				return err
			}
		}

		c.Log.Error("DeleteEquityCharts: Equity charts deleted", "id", id, "username", r.Username)
	}

	return err
}

//=============================================================================

func ExportTradingSystems(c *auth.Context, ids []uint) ([]byte, error){
	c.Log.Info("ExportTradingSystems: Exporting trading systems", "count", len(ids))

	buf       := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, id := range ids {
		data,err := backend.CreateBackup(c.Session.Username, id)
		if err != nil {
			_=zipWriter.Close()
			return nil, err
		}

		filename := strconv.FormatUint(uint64(id), 10) +".zip"
		err = writeData(zipWriter, filename, data)
		if err != nil {
			_=zipWriter.Close()
			return nil, err
		}
	}

	_=zipWriter.Close()
	return buf.Bytes(), nil
}

//=============================================================================

func writeData(zipWriter *zip.Writer, filename string, data []byte) error {
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)
	_, err = io.Copy(writer, buffer)
	return err
}

//=============================================================================
