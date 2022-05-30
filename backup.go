package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/plus3it/gorecurcopy"
	_ "github.com/plus3it/gorecurcopy"
)

type Backup struct {
	Path   string
	Object string
	Name   string
}
type BackupInfo struct {
	File      string `json:"file"`
	NewFile   string `json:"newFile,omitempty"`
	Operation string `json:"operation"`
}

func NewBackup(path, name string) Backup {
	return Backup{
		Path:   filepath.Dir(path),
		Object: filepath.Base(path),
		Name:   name,
	}
}

func (b Backup) Init() error {
	backupPath := filepath.Join(GetConfig().BackupDirectory, b.Name)
	isExist, err := IsFolderExist(backupPath)
	if err != nil {
		return err
	}
	if !isExist {
		err = CreateFolder(backupPath)
		if err != nil {
			return err
		}
		err = b.Copy()
		if err != nil {
			return err
		}
	}
	return nil
}
func (b Backup) Copy() error {
	srcPath := filepath.Join(b.Path, b.Object)
	backupPath := filepath.Join(GetConfig().BackupDirectory, b.Name, GetTimestamp())
	err := CreateFolder(backupPath)
	if err != nil {
		return err
	}
	err = gorecurcopy.CopyDirectory(srcPath, backupPath)
	if err != nil {
		return err
	}
	return nil
}

func (b Backup) Changes(file string) error {
	srcPath := filepath.Join(b.Path, b.Object)
	changedFile, err := filepath.Rel(srcPath, file)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(changedFile)
	fileName := filepath.Base(changedFile)

	backupPath := filepath.Join(GetConfig().BackupDirectory, b.Name, GetTimestamp())
	if dirPath != "." {
		backupPath = filepath.Join(backupPath, dirPath)
	}

	err = CreateFolder(backupPath)
	if err != nil {
		return err
	}

	backupFile := filepath.Join(backupPath, fileName)
	err = gorecurcopy.Copy(file, backupFile)
	if err != nil {
		return err
	}
	return nil
}
func (b Backup) Delete(file string) error {
	srcPath := filepath.Join(b.Path, b.Object)
	changedFile, err := filepath.Rel(srcPath, file)
	if err != nil {
		return err
	}

	backupPath := filepath.Join(GetConfig().BackupDirectory, b.Name, GetTimestamp()+".json")
	f, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer f.Close()

	info := BackupInfo{
		File:      changedFile,
		Operation: "Deleted",
	}

	err = json.NewEncoder(f).Encode(info)
	if err != nil {
		return err
	}

	return nil
}
func (b Backup) Rename(fileOld, fileNew string) error {
	srcPath := filepath.Join(b.Path, b.Object)

	changedOldFile, err := filepath.Rel(srcPath, fileOld)
	if err != nil {
		return err
	}
	changedNewFile, err := filepath.Rel(srcPath, fileNew)
	if err != nil {
		return err
	}

	backupPath := filepath.Join(GetConfig().BackupDirectory, b.Name, GetTimestamp()+".json")
	f, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer f.Close()

	info := BackupInfo{
		File:      changedOldFile,
		NewFile:   changedNewFile,
		Operation: "Renamed",
	}

	err = json.NewEncoder(f).Encode(info)
	if err != nil {
		return err
	}

	return nil
}
func (b Backup) Create(file string) error {
	srcPath := filepath.Join(b.Path, b.Object)
	changedFile, err := filepath.Rel(srcPath, file)
	if err != nil {
		return err
	}

	backupPath := filepath.Join(GetConfig().BackupDirectory, b.Name, GetTimestamp()+".json")
	f, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer f.Close()

	info := BackupInfo{
		File:      changedFile,
		Operation: "Created",
	}

	err = json.NewEncoder(f).Encode(info)
	if err != nil {
		return err
	}

	return nil
}
