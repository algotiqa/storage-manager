//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package app

import "github.com/algotiqa/core"

//=============================================================================

type Storage struct {
	Folder string
}

//=============================================================================

type Config struct {
	core.Application
	core.Authentication
	core.Messaging
	Storage
}

//=============================================================================
