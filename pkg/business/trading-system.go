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

package business

import (
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
