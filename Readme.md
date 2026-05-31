## Building A WAL-Based Key-Value Store From Scratch — The Concepts Behind Bitcask

--WAL will be in binary format and this is format i am following 

[1 byte]  operation type  (0x01 = SET, 0x02 = DELETE)
[8 bytes] timestamp       (uint64, Unix nanoseconds, big-endian)
[4 bytes] key length      (uint32, big-endian)
[4 bytes] value length    (uint32, big-endian; 0 for DELETE)
[N bytes] key             (raw UTF-8 bytes)
[M bytes] value           (raw bytes; absent for DELETE)
[4 bytes] CRC32 checksum  (computed over all preceding bytes in this entry)

### Hex Dump Examples

#### 1. SET Entry (Key: "key1", Value: "hello")
Assuming timestamp: `1717140000000000000` (0x17D784EE6B280000)

```text
01                          // Operation: SET (1 byte)
17 D7 84 EE 6B 28 00 00     // Timestamp: 1717140000000000000 (8 bytes)
00 00 00 04                 // Key Length: 4 (4 bytes)
00 00 00 05                 // Value Length: 5 (4 bytes)
6B 65 79 31                 // Key: "key1" (4 bytes)
68 65 6C 6C 6F              // Value: "hello" (5 bytes)
A4 B2 C1 D0                 // CRC32 Checksum (4 bytes)
```

#### 2. DELETE Entry (Key: "key2")
Assuming timestamp: `1717140000000000001` (0x17D784EE6B280001)

```text
02                          // Operation: DELETE (1 byte)
17 D7 84 EE 6B 28 00 01     // Timestamp: 1717140000000000001 (8 bytes)
00 00 00 04                 // Key Length: 4 (4 bytes)
00 00 00 00                 // Value Length: 0 (4 bytes)
6B 65 79 32                 // Key: "key2" (4 bytes)
F1 A2 B3 C4                 // CRC32 Checksum (4 bytes)
```
