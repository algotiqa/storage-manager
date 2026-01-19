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

package service

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/storage-manager/pkg/business"
)

//=============================================================================

func getDocumentation(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		var res *business.DocumentationResponse
		res, err = business.GetDocumentation(c, tsId)
		if err == nil {
			_ = c.ReturnObject(res)
			return
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func setDocumentation(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		docReq := business.DocumentationRequest{}
		err = c.BindParamsFromBody(&docReq)

		if err == nil {
			err = business.SetDocumentation(c, tsId, &docReq)
			if err == nil {
				_ = c.ReturnObject("")
				return
			}
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func getEquityChart(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()
	if err == nil {
		chartType := c.GetParamAsString("type", "unknown")
		var data []byte
		data, err = business.GetEquityChart(c, tsId, chartType)
		if err == nil {
			_ = c.ReturnData("image/png", data)
			return
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func setEquityCharts(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		equReq := business.NewEquityRequest()
		err = c.BindParamsFromBody(equReq)

		if err == nil {
			err = business.SetEquityCharts(c, tsId, equReq)
			if err == nil {
				_ = c.ReturnObject("")
				return
			}
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func deleteEquityCharts(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		equReq := business.NewEquityRequest()
		err = c.BindParamsFromBody(equReq)

		if err == nil {
			err = business.DeleteEquityCharts(c, tsId, equReq)
			if err == nil {
				_ = c.ReturnObject("")
				return
			}
		}
	}

	c.ReturnError(err)
}

//=============================================================================
