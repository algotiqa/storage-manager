//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package inventory

import (
	"encoding/json"
	"log/slog"

	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/storage-manager/pkg/backend"
)

//=============================================================================

func InitMessageListener() {
	slog.Info("Starting inventory message listener...")
	go msg.ReceiveMessages(msg.QuInventoryToStorage, handleMessage)
}

//=============================================================================

func handleMessage(m *msg.Message) bool {
	slog.Info("New message received", "source", m.Source, "type", m.Type)

	if m.Source == msg.SourceTradingSystem {
		tsm := TradingSystemMessage{}
		err := json.Unmarshal(m.Entity, &tsm)
		if err != nil {
			slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
			return true
		}

		if m.Type == msg.TypeCreate {
			return addTradingSystem(&tsm)
		}
		if m.Type == msg.TypeUpdate {
			return updateTradingSystem(&tsm)
		}
		if m.Type == msg.TypeDelete {
			return deleteTradingSystem(&tsm)
		}
	}

	slog.Error("Dropping message with unknown source/type!", "source", m.Source, "type", m.Type)
	return true
}

//=============================================================================

func addTradingSystem(tsm *TradingSystemMessage) bool {
	slog.Info("addTradingSystem: New trading system received", "id", tsm.TradingSystem.Id, "name", tsm.TradingSystem.Name)

	ts := &backend.TradingSystem{
		Id      : tsm.TradingSystem.Id,
		Username: tsm.TradingSystem.Username,
		Name    : tsm.TradingSystem.Name,
	}

	var err error

	if len(tsm.StoragePack) == 0 {
		err = backend.AddTradingSystem(ts)
	} else {
		err = backend.RestoreBackup(ts.Username, ts.Id, tsm.StoragePack)
	}

	if err != nil {
		slog.Error("addTradingSystem: Cannot add the trading system", "id", tsm.TradingSystem.Id, "error", err.Error())
	} else {
		slog.Info("addTradingSystem: Operation complete", "id", tsm.TradingSystem.Id)
	}

	return err == nil
}

//=============================================================================

func updateTradingSystem(tsm *TradingSystemMessage) bool {
	slog.Info("updateTradingSystem: Trading system change received", "id", tsm.TradingSystem.Id, "name", tsm.TradingSystem.Name)

	ts := &backend.TradingSystem{
		Id:       tsm.TradingSystem.Id,
		Username: tsm.TradingSystem.Username,
		Name:     tsm.TradingSystem.Name,
	}
	err := backend.UpdateTradingSystem(ts)

	if err != nil {
		slog.Error("updateTradingSystem: Cannot update the trading system", "id", tsm.TradingSystem.Id, "error", err.Error())
	} else {
		slog.Info("updateTradingSystem: Operation complete", "id", tsm.TradingSystem.Id)
	}

	return err == nil
}

//=============================================================================

func deleteTradingSystem(tsm *TradingSystemMessage) bool {
	slog.Info("deleteTradingSystem: Trading system deletion received", "id", tsm.TradingSystem.Id)

	err := backend.DeleteTradingSystem(tsm.TradingSystem.Id, tsm.TradingSystem.Username)

	if err != nil {
		slog.Error("deleteTradingSystem: Raised error while deleting trading system", "error", err.Error())
	} else {
		slog.Info("deleteTradingSystem: Operation complete", "id", tsm.TradingSystem.Id)
	}

	return err == nil
}

//=============================================================================
