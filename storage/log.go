
package storage

import (
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"
)

const (
	headerSize    = 8
	maxRecordSize = 64 * 1024 * 1024 // 64 MB
)

type Log struct {
	file    *os.File
	records [][]byte
}

// OpenLog opens an existing log and recovers valid records.
func OpenLog(path string) (*Log, error) {
	fp, err := os.OpenFile(
		path,
		os.O_RDWR|os.O_CREATE,
		0664,
	)
	if err != nil {
		return nil, err
	}

	log := &Log{
		file: fp,
	}

	info, err := fp.Stat()
	if err != nil {
		fp.Close()
		return nil, err
	}

	fileSize := info.Size()
	offset := int64(0)
	validEnd := int64(0)

	// Read records one by one.
	for offset < fileSize {
		header := make([]byte, headerSize)

		n, err := fp.ReadAt(header, offset)
		if err != nil || n != headerSize {
			break
		}

		size := binary.LittleEndian.Uint32(header[0:4])
		checksum := binary.LittleEndian.Uint32(header[4:8])

		// Reject invalid record sizes.
		if size > maxRecordSize {
			break
		}

		recordEnd := offset + headerSize + int64(size)

		// Incomplete record.
		if recordEnd > fileSize {
			break
		}

		data := make([]byte, int(size))

		n, err = fp.ReadAt(data, offset+headerSize)
		if err != nil && err != io.EOF {
			fp.Close()
			return nil, err
		}

		if n != len(data) {
			break
		}

		// Verify checksum.
		if crc32.ChecksumIEEE(data) != checksum {
			break
		}

		log.records = append(log.records, data)

		offset = recordEnd
		validEnd = recordEnd
	}

	// Remove an incomplete or corrupted trailing record.
	if validEnd != fileSize {
		if err := fp.Truncate(validEnd); err != nil {
			fp.Close()
			return nil, err
		}

		if err := fp.Sync(); err != nil {
			fp.Close()
			return nil, err
		}
	}

	// Position the file cursor at the end.
	if _, err := fp.Seek(0, io.SeekEnd); err != nil {
		fp.Close()
		return nil, err
	}

	return log, nil
}

// Append adds a new record to the end of the log.
func (l *Log) Append(data []byte) error {
	if len(data) > maxRecordSize {
		return os.ErrInvalid
	}

	// Create the record header.
	header := make([]byte, headerSize)

	binary.LittleEndian.PutUint32(
		header[0:4],
		uint32(len(data)),
	)

	binary.LittleEndian.PutUint32(
		header[4:8],
		crc32.ChecksumIEEE(data),
	)

	// Combine header and payload.
	record := make([]byte, headerSize+len(data))

	copy(record[:headerSize], header)
	copy(record[headerSize:], data)

	// Write the complete record.
	for len(record) > 0 {
		n, err := l.file.Write(record)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		record = record[n:]
	}

	// Synchronize the appended record.
	if err := l.file.Sync(); err != nil {
		return err
	}

	l.records = append(l.records, append([]byte(nil), data...))

	return nil
}

// Replay returns all valid records in order.
func (l *Log) Replay() [][]byte {
	result := make([][]byte, len(l.records))

	for i, record := range l.records {
		result[i] = append([]byte(nil), record...)
	}

	return result
}

// Close closes the log file.
func (l *Log) Close() error {
	return l.file.Close()
}
