//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package inventory

//=============================================================================
//=== Entities
//=============================================================================

type TradingSystem struct {
	Id                uint    `json:"id"`
	Username          string  `json:"username"`
	Name              string  `json:"name"`
	Timeframe         int     `json:"timeframe"`
	DataProductId     uint    `json:"dataProductId"`
	BrokerProductId   uint    `json:"brokerProductId"`
	TradingSessionId  uint    `json:"tradingSessionId"`
	AgentProfileId    uint    `json:"agentProfileId" gorm:"default:null"`
	StrategyType      string  `json:"strategyType"`
	Overnight         bool    `json:"overnight"`
	Tags              string  `json:"tags"`
	ExternalRef       string  `json:"externalRef"`
	Finalized         bool    `json:"finalized"`
}

//=============================================================================
//=== Messages
//=============================================================================

type TradingSystemMessage struct {
	TradingSystem  TradingSystem  `json:"tradingSystem"`
	PortfolioPack  []byte         `json:"portfolioPack"`
	StoragePack    []byte         `json:"storagePack"`
}

//=============================================================================
