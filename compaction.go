package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func (c *CentralStorage) AddingEntryToManifest() error {
	data, err := json.MarshalIndent(c.manifestJson, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.manifest, data, 0644)
}
func (c *CentralStorage) CreateManifest() (*ManifestFile, error) {
	file, err := os.OpenFile(c.manifest, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	c.manifestJson = ManifestJson{
		Manifest_version:   1,
		Last_compact_index: 0,
		Snapshot_file:      "",
		Snapshot_checksum:  0,
		Created_at:         0,
		Entry_count:        0,
	}
	c.AddingEntryToManifest()
	return &ManifestFile{file}, nil
}
func (c *CentralStorage) CreateSnapShot() (*SnapShotFile, error) {
	file, err := os.OpenFile(c.snap_path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Error while creating snapshot file ", err)
	}
	return &SnapShotFile{file}, nil
}
func (e *SnapShotEntry) WriteSnapShotEntry(mw io.Writer, mapval map[string]string) error {

	for key, val := range mapval {
		e.Keylen = int32(len(key))
		e.Vallen = int32(len(val))
		e.Key = key
		e.Val = val
		if err := binary.Write(mw, binary.BigEndian, e.Keylen); err != nil {
			return err
		}
		if err := binary.Write(mw, binary.BigEndian, e.Vallen); err != nil {
			return err
		}
		mw.Write([]byte(e.Key))
		mw.Write([]byte(e.Val))

	}
	return nil
}
func (c *CentralStorage) DeleteSnapBinary() {
	fmt.Println("Deleting the Binary since we have encountered the error")
	err := os.Remove(c.snap_path)
	if err != nil {
		log.Fatal("Error while deleting the file ", err)
	}
}
func (h *SnapShotHeader) WriteSnapShotHeader(mw io.Writer) error {
	if err := binary.Write(mw, binary.BigEndian, h.MagicNumber); err != nil {
		return err

	}
	if err := binary.Write(mw, binary.BigEndian, h.Version); err != nil {
		return err
	}
	if err := binary.Write(mw, binary.BigEndian, h.EntryCount); err != nil {
		return err
	}
	return nil
}

func (h *SnapShotFooter) WriteSnapShotFooter(file *os.File, hasher hash.Hash32) (uint32, error) {
	crcCal := hasher.Sum32()
	if err := binary.Write(file, binary.BigEndian, crcCal); err != nil {
		return 0, err
	}
	return crcCal, nil
}
func (c *CentralStorage) CheckManifest() {
	_, errStat := os.Stat(c.manifest)
	if err := os.IsNotExist(errStat); err {
		if _, err := c.CreateManifest(); err != nil {
			fmt.Println("File error", err)
		}
	}
}

func (c *CentralStorage) FetchManifestIndex() int {
	c.CheckManifest()
	var manifestStruct ManifestJson
	data, _ := os.ReadFile(c.manifest)
	err := json.Unmarshal(data, &manifestStruct)
	if err != nil {
		panic(err)
	}
	index := manifestStruct.Last_compact_index
	return index
}

func (c *CentralStorage) TriggerCompaction() {
	if c.isCompactionRunning {
		return
	} else {
		c.isCompactionRunning = true
	}
	err := c.RotateWal()
	if err != nil {
		c.isCompactionRunning = false
		return
	}
	go c.RunCompaction()
}
func (c *CentralStorage) RotateWal() error {
	getCurrentFile := c.file.Name()
	c.sealedEntry = append(c.sealedEntry, getCurrentFile)
	if err := c.file.Close(); err != nil {
		return fmt.Errorf("Error while closing the file %s", err)
	}
	FileCount := strings.Split(strings.Split(getCurrentFile, "-")[1], ".")[0]
	fileIndex, err := strconv.Atoi(FileCount)
	if err != nil {
		return fmt.Errorf("error while converting ascii to Int %w", err)
	}

	path := fmt.Sprintf("WAL_LOG/wal-%06d.log", fileIndex+1)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Error In RotoateWAL function")
	}
	c.file = file
	c.currentOperationFile = file.Name()
	return nil
}
func (c *CentralStorage) RunCompaction() {

}
func (c *CentralStorage) CopyLiveMapData() map[string]string {
	privateSnap := make(map[string]string, len(c.mapval))
	for key, val := range c.mapval {
		privateSnap[key] = val
	}
	return privateSnap
}
func (c *CentralStorage) WriteToSnap() {
	var latest_snapShot int
	var name string
	hasher := crc32.NewIEEE() //used for incremental hashing

	headerval := SnapShotHeader{
		MagicNumber: magicNumber,
		Version:     uint32(versionNum),
		EntryCount:  uint32(len(c.mapval)),
	}
	entry := SnapShotEntry{}

	file, err := c.CreateSnapShot()
	if err != nil {
		log.Fatal("Error while creating in file ", err)
	}

	mw := io.MultiWriter(file.file, hasher)
	if err := headerval.WriteSnapShotHeader(mw); err != nil {
		c.DeleteSnapBinary()
	}

	if err := entry.WriteSnapShotEntry(mw, c.mapval); err != nil {
		log.Fatal("Error while writng the entry to snap=>", err)
	}
	footer := SnapShotFooter{}
	var hashVal uint32
	if hashVal, err = footer.WriteSnapShotFooter(file.file, hasher); err != nil {
		log.Fatal("Error in footer-->", err)
	}
	file.file.Sync()

	/*for _, p := range c.sealedFile {
		filepath := dirPath + "/" + p
		name = strings.Split(p, ".")[0]
		snapIndex := strings.Split(strings.Split(p, ".")[0], "-")[1]
		latest_snapShot, _ = strconv.Atoi(snapIndex)
		sealedFileList = append(sealedFileList, filepath)
	}*/
	c.manifestJson = ManifestJson{
		Manifest_version:   1,
		Last_compact_index: latest_snapShot,
		Snapshot_file:      name,
		Snapshot_checksum:  hashVal,
		Created_at:         time.Now().Unix(),
		Entry_count:        len(mapval),
	}
	c.AddingEntryToManifest()

	for file := range c.sealedFileSet {
		os.Remove(file)
	}
}

func (c *CentralStorage) ScanFileUsingManifest() string {
	data, err := os.ReadFile(c.manifest)
	if err != nil {
		log.Fatal("Manifest file not found")
	}
	var manifestData ManifestJson
	err = json.Unmarshal(data, &manifestData)
	if err != nil {
		log.Fatal("Error while fetching the data from manifest json")
	}
	return c.CreateWALSegment(manifestData.Last_compact_index)
}

func (c *CentralStorage) LoadSnapDataToMemory() error {
	file, err := os.Open(c.snap_path)
	magicbuff := make([]byte, 4)
	var snapEntry SnapShotEntry
	var header SnapShotHeader
	if err != nil {
		log.Fatal("Snap File has some issue==>", file)
	}
	err = binary.Read(file, binary.BigEndian, &header.MagicNumber)
	binary.BigEndian.PutUint32(magicbuff, header.MagicNumber)
	err = binary.Read(file, binary.BigEndian, &header.Version)
	err = binary.Read(file, binary.BigEndian, &header.EntryCount)

	if string(magicbuff) != "SNAP" {
		return fmt.Errorf("Invalid snapshot file since correct header is not present")
	}
	for i := 0; i < int(header.EntryCount); i++ {
		if err == io.EOF {
			break
		}
		binary.Read(file, binary.BigEndian, &snapEntry.Keylen)
		binary.Read(file, binary.BigEndian, &snapEntry.Vallen)
		keyByte := make([]byte, snapEntry.Keylen)
		valByte := make([]byte, snapEntry.Vallen)
		io.ReadFull(file, keyByte)
		io.ReadFull(file, valByte)

		fmt.Println(string(keyByte))
		fmt.Println(string(valByte))
		c.mapval[string(keyByte)] = string(valByte)
	}
	return nil
}
