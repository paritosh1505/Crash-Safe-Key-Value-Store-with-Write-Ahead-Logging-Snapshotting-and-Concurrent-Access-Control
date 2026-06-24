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

type SnapShotFile struct {
	file *os.File
}
type ManifestFile struct {
	file *os.File
}
type SnapShotHeader struct {
	MagicNumber uint32
	Version     uint32
	EntryCount  uint32
}
type SnapShotEntry struct {
	Keylen uint32
	Vallen uint32
	Key    string
	Val    string
}
type ManifestData struct {
	Manifest_version   int    `json:"version"`
	Last_compact_index int    `json:"last_index"`
	Snapshot_file      string `json:"file"`
	Snapshot_checksum  uint32 `json:"checksum"`
	Created_at         int64  `json:"CreatedAt"`
	Entry_count        int    `json:"entryCount"`
}

type SnapShotFooter struct {
	Crcval uint32
}

func (m *ManifestData) AddingEntryToManifest() error {
	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(manifest, data, 0644)
}
func CreateManifest() (*ManifestFile, error) {
	file, err := os.OpenFile(manifest, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	manifest := ManifestData{
		Manifest_version:   1,
		Last_compact_index: 0,
		Snapshot_file:      "",
		Snapshot_checksum:  0,
		Created_at:         0,
		Entry_count:        0,
	}
	manifest.AddingEntryToManifest()
	return &ManifestFile{file}, nil
}
func CreateSnapShot() (*SnapShotFile, error) {
	file, err := os.OpenFile(snap_path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Error while creating snapshot file ", err)
	}
	return &SnapShotFile{file}, nil
}
func (e *SnapShotEntry) WriteSnapShotEntry(mw io.Writer) error {

	for key, val := range mapval {
		e.Keylen = uint32(len(key))
		e.Vallen = uint32(len(val))
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
func DeleteSnapBinary() {
	fmt.Println("Deleting the Binary since we have encountered the error")
	err := os.Remove(snap_path)
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
func CheckManifest() {
	_, errStat := os.Stat(manifest)
	if err := os.IsNotExist(errStat); err {
		if _, err := CreateManifest(); err != nil {
			fmt.Println("File error", err)
		}
	}
}

func FetchManifestIndex() int {
	CheckManifest()
	var manifestStruct ManifestData
	data, _ := os.ReadFile(manifest)
	err := json.Unmarshal(data, &manifestStruct)
	if err != nil {
		panic(err)
	}
	index := manifestStruct.Last_compact_index
	//filename := fmt.Sprintf("WAL_LOG/wal-%06d.log", index)
	return index
}

func WriteToSnap() {
	var latest_snapShot int
	var name string
	hasher := crc32.NewIEEE() //used for incremental hashing
	for _, p := range sealedFile {
		filepath := dirPath + "/" + p
		//name = strings.Split(p, ".")[0]
		snapIndex := strings.Split(strings.Split(p, ".")[0], "-")[1]
		latest_snapShot, _ = strconv.Atoi(snapIndex)
		os.Remove(filepath)
	}
	headerval := SnapShotHeader{
		MagicNumber: magicNumber,
		Version:     uint32(versionNum),
		EntryCount:  uint32(len(mapval)),
	}
	entry := SnapShotEntry{}
	file, err := CreateSnapShot()
	if err != nil {
		log.Fatal("Error while creating in file ", err)
	}

	mw := io.MultiWriter(file.file, hasher)
	if err := headerval.WriteSnapShotHeader(mw); err != nil {
		DeleteSnapBinary()
	}
	if err := entry.WriteSnapShotEntry(mw); err != nil {
		log.Fatal("Error while writng the entry to snap=>", err)
	}
	footer := SnapShotFooter{}
	var hashVal uint32
	if hashVal, err = footer.WriteSnapShotFooter(file.file, hasher); err != nil {
		log.Fatal("Error in footer-->", err)
	}

	manifest := ManifestData{
		Manifest_version:   1,
		Last_compact_index: latest_snapShot,
		Snapshot_file:      name,
		Snapshot_checksum:  hashVal,
		Created_at:         time.Now().Unix(),
		Entry_count:        len(mapval),
	}
	manifest.AddingEntryToManifest()
}

func ScanFileUsingManifest() string {
	data, err := os.ReadFile(manifest)
	if err != nil {
		log.Fatal("Manifest file not found")
	}
	var manifestData ManifestData
	err = json.Unmarshal(data, &manifestData)
	if err != nil {
		log.Fatal("Error while fetching the data from manifest json")
	}
	return CreateWALSegment(manifestData.Last_compact_index)
}

func LoadSnapDataToMemory() error {
	file, err := os.Open(snap_path)
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
		mapval[string(keyByte)] = string(valByte)
	}
	return nil
}
