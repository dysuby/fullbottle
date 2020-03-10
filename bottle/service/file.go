package service

import (
	"fmt"
	"github.com/vegchic/fullbottle/bottle/dao"
	"github.com/vegchic/fullbottle/common/kv"
	"sort"
	"time"
)

// use for file name unique
const FileLockKey = "lock:folder_id=%d"

func CreateFile(file *dao.FileInfo, meta *dao.FileMeta) error {
	lock, err := kv.Obtain(fmt.Sprintf(FileLockKey, file.FolderId), 100*time.Millisecond)
	if err != nil {
		return err
	}
	defer lock.Release()
	files, err := dao.GetFilesByFolderId(file.OwnerId, file.FolderId, nil)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	name := file.Name
	repeat := 1
	for _, v := range files {
		if v.Name == name {
			name = fmt.Sprintf("%s (%d)", file.Name, repeat)
			repeat++
		}
	}
	file.Name = name
	file.FileId = meta.ID
	file.Size = meta.Size

	err = dao.CreateFile(file, meta)
	if err != nil {
		return err
	}

	return nil
}
