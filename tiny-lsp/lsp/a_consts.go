// Copyright 2024 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

const (
	k_initialize  = "initialize"
	k_initialized = "initialized"
	k_shutdown    = "shutdown"
	k_exit        = "exit"

	k_workspace_executeCommand            = "workspace/executeCommand"
	k_workspace_didChangeConfiguration    = "workspace/didChangeConfiguration"
	k_workspace_workspaceFolders          = "workspace/workspaceFolders"
	k_workspace_didChangeWorkspaceFolders = "workspace/didChangeWorkspaceFolders"

	k_textDocument_didOpen        = "textDocument/didOpen"
	k_textDocument_didChange      = "textDocument/didChange"
	k_textDocument_didSave        = "textDocument/didSave"
	k_textDocument_didClose       = "textDocument/didClose"
	k_textDocument_formatting     = "textDocument/formatting"
	k_textDocument_documentSymbol = "textDocument/documentSymbol"
	k_textDocument_completion     = "textDocument/completion"
	k_textDocument_definition     = "textDocument/definition"
	k_textDocument_references     = "textDocument/references"
	k_textDocument_hover          = "textDocument/hover"
	k_textDocument_codeAction     = "textDocument/codeAction"

	k_window_logMessage = "window/logMessage"
)
