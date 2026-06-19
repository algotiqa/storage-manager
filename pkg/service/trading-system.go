//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package service

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
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

func exportTradingSystems(c *auth.Context) {
	ids,err := c.GetIdsFromUrl()
	if err == nil {
		if len(ids) == 0 {
			err = req.NewBadRequestError("Parameter 'id' is missing or empty")
		} else {
			var res []byte
			res,err = business.ExportTradingSystems(c, ids)
			if err == nil {
				_=c.ReturnData("application/zip", res)
				return
			}
		}
	}

	c.ReturnError(err)
}

//=============================================================================
