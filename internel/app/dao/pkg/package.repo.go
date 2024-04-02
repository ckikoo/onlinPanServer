package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/pkg/contextx"

	"gorm.io/gorm"
)

type PkgRepo struct {
	DB *gorm.DB
	RD *redisx.Redisx
}

var rdKey string = "%v:packageInfo"

func (repo *PkgRepo) GetPackInfo(ctx context.Context) ([]Pkg, error) {
	c := config.C

	packInfos, err := repo.RD.HGetAll(ctx, fmt.Sprintf(rdKey, c.AppName))
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}
	fmt.Printf("packInfos: %v\n", packInfos)
	var items []Pkg

	if packInfos == nil || len(packInfos) == 0 {
		err := repo.DB.Find(&items).Error
		if err != nil {
			return nil, err
		}

		if items == nil || len(items) == 0 {
			return nil, nil
		}
		packInfo := make(map[string]interface{}, 0)
		for _, v := range items {
			temp, _ := json.Marshal(v)
			packInfo[fmt.Sprint(v.Id)] = temp
		}
		effect, err := repo.RD.HMset(ctx, fmt.Sprintf(rdKey, c.AppName), packInfo)
		if err != nil {
			fmt.Printf("effect: %v err:%v\n", effect, err.Error())
			return nil, errors.New("internal fail:")
		}

		return items, nil
	} else {

		var item Pkg
		for _, v := range packInfos {
			if err := json.Unmarshal([]byte(v), &item); err != nil {
				return nil, fmt.Errorf("json format error:%v", err)
			}
			items = append(items, item)
		}
		return items, nil
	}

}

func (repo *PkgRepo) GetPackInfoByID(ctx context.Context, packId string) (any, error) {
	c := config.C
	res, err := repo.RD.HMGet(ctx, fmt.Sprintf(rdKey, c.AppName), packId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("res: %v\n", res)
	item, ok := res[0].(string)
	if !ok {
		return nil, errors.New("type error")
	}
	var v Pkg
	json.Unmarshal([]byte(item), &v)

	return v, nil
}

func (repo *PkgRepo) BuySpace(ctx context.Context, uid string, spaceinfo BuySpace) (bool, error) {
	key := "onlinePan:space:%v"

	db := GetSpaceDB(ctx, repo.DB)

	err := db.Create(spaceinfo).Error
	if err != nil {
		return false, err
	}

	item, _ := json.Marshal(spaceinfo)
	flag, err := repo.RD.HSet(ctx, fmt.Sprintf(key, uid), spaceinfo.SpaceId, item)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return flag, err
	}

	repo.RD.Delete(ctx, fmt.Sprintf("user:space:%v", contextx.FromUserEmail(ctx)))

	return flag, nil

}
