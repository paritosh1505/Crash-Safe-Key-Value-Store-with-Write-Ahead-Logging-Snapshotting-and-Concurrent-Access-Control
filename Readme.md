## Building A WAL-Based Key-Value Store From Scratch — The Concepts Behind Bitcask

---WAL will be in binary format and this is format i am following 
# Persistent KV Storage (WAL-based)

This project implements a high-performance, persistent Key-Value store inspired by the Bitcask design. It uses a segmented Write-Ahead Log (WAL) to ensure data durability and fast recovery.

## Core Features

- **Segmented WAL:** Instead of a single massive file, logs are split into manageable numbered segments.
- **Threshold-based Rotation:** Automatically creates a new log segment once the current file exceeds the configured size limit (109 bytes).
- **Binary Protocol:** Data is serialized into a compact binary format for efficiency.
- **Crash Recovery:** On startup, the system replays all WAL segments in chronological order to reconstruct the in-memory state.
- **Data Integrity:** Each entry includes a CRC32 checksum to detect data corruption.

## Segment Management

- **Naming Convention:** `wal-000001.log`, `wal-000002.log`, etc.
- **Storage Path:** All logs are stored in the `WAL_LOG/` directory.
- **Active Segment:** The latest segment is used for appending new writes. Once it hits the `threshold` (currently 109 bytes), it is closed and a new segment is initialized.

## WAL Binary Format

| Size | Field | Description |
| :--- | :--- | :--- |
| 1 byte | Operation Type | 0x01 = SET, 0x02 = DELETE |
| 8 bytes | Timestamp | uint64, Unix nanoseconds, big-endian |
| 4 bytes | Key Length | uint32, big-endian |
| 4 bytes | Value Length | uint32, big-endian (0 for DELETE) |
| N bytes | Key | Raw UTF-8 bytes |
| M bytes | Value | Raw bytes (absent for DELETE) |
| 4 bytes | CRC32 Checksum | Computed over all preceding bytes in the entry |

### Detailed Structure

1. **Operation Type:** Indicates if the record is a write or a deletion.
2. **Timestamp:** Used for resolution if the same key appears multiple times across segments.
3. **Lengths:** Explicit lengths for key and value allow for reading the raw bytes precisely.
4. **CRC32:** Ensures that the record has not been truncated or corrupted on disk.

## Operational Flow

### 1. Bootstrapping (Replay)
During startup, `ReplayAllWAL` scans the `WAL_LOG` directory. It reads every segment from the lowest index to the highest, decoding each `WALRecord` and applying the operations to the in-memory hash map.

### 2. Writing Data
When a write is requested:
1. The `ManageWALFile` logic determines if the current active segment has room.
2. If the file exceeds `threshold`, a new segment index is generated.
3. The record is encoded into binary, appended to the file, and `fsync` is called to ensure it hits the physical disk.
4. The in-memory map is updated.

