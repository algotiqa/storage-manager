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

package inventory

import (
	"encoding/json"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/storage-manager/pkg/backend"
	"log/slog"
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
		Id:       tsm.TradingSystem.Id,
		Username: tsm.TradingSystem.Username,
		Name:     tsm.TradingSystem.Name,
	}
	err := backend.AddTradingSystem(ts)

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
