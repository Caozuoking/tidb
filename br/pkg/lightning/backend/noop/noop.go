// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package noop

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pingcap/tidb/br/pkg/lightning/backend"
	"github.com/pingcap/tidb/br/pkg/lightning/backend/encode"
	"github.com/pingcap/tidb/br/pkg/lightning/config"
	"github.com/pingcap/tidb/br/pkg/lightning/verification"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/table"
	"github.com/pingcap/tidb/types"
)

// NewNoopBackend creates a new backend that does nothing.
func NewNoopBackend() backend.Backend {
	return backend.MakeBackend(noopBackend{})
}

type noopBackend struct{}

type noopRows struct{}

// SplitIntoChunks implements the Rows interface.
func (r noopRows) SplitIntoChunks(int) []encode.Rows {
	return []encode.Rows{r}
}

// Clear returns a new collection with empty content. It may share the
// capacity with the current instance. The typical usage is `x = x.Clear()`.
func (r noopRows) Clear() encode.Rows {
	return r
}

// Close the connection to the backend.
func (b noopBackend) Close() {}

// MakeEmptyRows creates an empty collection of encoded rows.
func (b noopBackend) MakeEmptyRows() encode.Rows {
	return noopRows{}
}

// RetryImportDelay returns the duration to sleep when retrying an import
func (b noopBackend) RetryImportDelay() time.Duration {
	return 0
}

// ShouldPostProcess returns whether KV-specific post-processing should be
// performed for this backend. Post-processing includes checksum and analyze.
func (b noopBackend) ShouldPostProcess() bool {
	return false
}

// NewEncoder creates an encoder of a TiDB table.
func (b noopBackend) NewEncoder(ctx context.Context, config *encode.EncodingConfig) (encode.Encoder, error) {
	return noopEncoder{}, nil
}

// OpenEngine creates a new engine file for the given table.
func (b noopBackend) OpenEngine(context.Context, *backend.EngineConfig, uuid.UUID) error {
	return nil
}

// CloseEngine closes the engine file, flushing any remaining data.
func (b noopBackend) CloseEngine(ctx context.Context, cfg *backend.EngineConfig, engineUUID uuid.UUID) error {
	return nil
}

// ImportEngine imports a closed engine file.
func (b noopBackend) ImportEngine(ctx context.Context, engineUUID uuid.UUID, regionSplitSize, regionSplitKeys int64) error {
	return nil
}

// CleanupEngine removes all data related to the engine.
func (b noopBackend) CleanupEngine(ctx context.Context, engineUUID uuid.UUID) error {
	return nil
}

// CheckRequirements performs the check whether the backend satisfies the
// version requirements
func (b noopBackend) CheckRequirements(context.Context, *backend.CheckCtx) error {
	return nil
}

// FetchRemoteTableModels obtains the models of all tables given the schema
// name. The returned table info does not need to be precise if the encoder,
// is not requiring them, but must at least fill in the following fields for
// TablesFromMeta to succeed:
//   - Name
//   - State (must be model.StatePublic)
//   - ID
//   - Columns
//   - Name
//   - State (must be model.StatePublic)
//   - Offset (must be 0, 1, 2, ...)
//   - PKIsHandle (true = do not generate _tidb_rowid)
func (b noopBackend) FetchRemoteTableModels(ctx context.Context, schemaName string) ([]*model.TableInfo, error) {
	return nil, nil
}

// FlushEngine ensures all KV pairs written to an open engine has been
// synchronized, such that kill-9'ing Lightning afterwards and resuming from
// checkpoint can recover the exact same content.
//
// This method is only relevant for local backend, and is no-op for all
// other backends.
func (b noopBackend) FlushEngine(ctx context.Context, engineUUID uuid.UUID) error {
	return nil
}

// FlushAllEngines performs FlushEngine on all opened engines. This is a
// very expensive operation and should only be used in some rare situation
// (e.g. preparing to resolve a disk quota violation).
func (b noopBackend) FlushAllEngines(ctx context.Context) error {
	return nil
}

// EngineFileSizes obtains the size occupied locally of all engines managed
// by this backend. This method is used to compute disk quota.
// It can return nil if the content are all stored remotely.
func (b noopBackend) EngineFileSizes() []backend.EngineFileSize {
	return nil
}

// ResetEngine clears all written KV pairs in this opened engine.
func (b noopBackend) ResetEngine(ctx context.Context, engineUUID uuid.UUID) error {
	return nil
}

// LocalWriter obtains a thread-local EngineWriter for writing rows into the given engine.
func (b noopBackend) LocalWriter(context.Context, *backend.LocalWriterConfig, uuid.UUID) (backend.EngineWriter, error) {
	return Writer{}, nil
}

// CollectLocalDuplicateRows collects duplicate rows from local backend.
func (b noopBackend) CollectLocalDuplicateRows(ctx context.Context, tbl table.Table, tableName string, opts *encode.SessionOptions) (bool, error) {
	panic("Unsupported Operation")
}

// CollectRemoteDuplicateRows collects duplicate rows from remote backend.
func (b noopBackend) CollectRemoteDuplicateRows(ctx context.Context, tbl table.Table, tableName string, opts *encode.SessionOptions) (bool, error) {
	panic("Unsupported Operation")
}

// ResolveDuplicateRows resolves duplicate rows.
func (b noopBackend) ResolveDuplicateRows(ctx context.Context, tbl table.Table, tableName string, algorithm config.DuplicateResolutionAlgorithm) error {
	return nil
}

// TotalMemoryConsume returns the total memory usage of the backend.
func (b noopBackend) TotalMemoryConsume() int64 {
	return 0
}

type noopEncoder struct{}

// Close the encoder.
func (e noopEncoder) Close() {}

// Encode encodes a row of SQL values into a backend-friendly format.
func (e noopEncoder) Encode([]types.Datum, int64, []int, int64) (encode.Row, error) {
	return noopRow{}, nil
}

type noopRow struct{}

// Size returns the size of the encoded row.
func (r noopRow) Size() uint64 {
	return 0
}

// ClassifyAndAppend classifies the row into the corresponding collection.
func (r noopRow) ClassifyAndAppend(*encode.Rows, *verification.KVChecksum, *encode.Rows, *verification.KVChecksum) {
}

// Writer define a local writer that do nothing.
type Writer struct{}

// AppendRows implements the EngineWriter interface.
func (w Writer) AppendRows(context.Context, string, []string, encode.Rows) error {
	return nil
}

// IsSynced implements the EngineWriter interface.
func (w Writer) IsSynced() bool {
	return true
}

// Close implements the EngineWriter interface.
func (w Writer) Close(context.Context) (backend.ChunkFlushStatus, error) {
	return trueStatus{}, nil
}

type trueStatus struct{}

// Flushed implements the ChunkFlushStatus interface.
func (s trueStatus) Flushed() bool {
	return true
}
