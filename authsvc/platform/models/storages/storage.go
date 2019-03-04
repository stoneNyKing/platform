package storages

import (
	"errors"
	"strings"
)

var (
	ErrStroageAlreadyExist = errors.New("Stroage already exist")
	ErrStorageNotExist     = errors.New("Stroage does not exist")
)

type Storage struct {
	Key   string  `xorm:"unique not null"`
	Value []uint8 `xorm:MEDIUMBLOB`
}

type ErrorResp struct {
	Ret int    `json:"Ret"`
	Msg string `json:"Msg"`
}

func IsKeyExist(k string) (bool, error) {
	if len(k) == 0 {
		return false, nil
	}
	return orm.Get(&Storage{Key: strings.ToLower(k)})
}

func GetStorage(k string) (*Storage, error) {
	if len(k) == 0 {
		return nil, ErrStorageNotExist
	}
	storage := &Storage{Key: strings.ToLower(k)}
	has, err := orm.Get(storage)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrStorageNotExist
	}
	return storage, nil
}

func SetStorage(storage *Storage) (*Storage, error) {
	isExist, err := IsKeyExist(strings.ToLower(storage.Key))
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrStroageAlreadyExist
	}

	if _, err = orm.Insert(storage); err != nil {
		return nil, err
	}
	return storage, err
}

type StorageKeyResp struct {
	Ret int
	Key string
	W   uint
	H   uint
}
