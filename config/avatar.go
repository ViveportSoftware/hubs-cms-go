package config

import (
	"fmt"
	"hubs-cms-go/utils"
)

func GetDirectusGetAvatarURI(avatarID string) string {
	return fmt.Sprintf("%s/items/avatar/%s", EnvVariable.DirectusBaseURI, avatarID)
}

func GetDirectusGetPublicAvatarURI(start int64, limit int64) string {
	return fmt.Sprintf(`%s/items/avatar?filter={"is_public":{"_eq":true}}&meta=*&%s`, EnvVariable.DirectusBaseURI, utils.GetPageParam(start, limit))
}

func GetDirectusGetMyAvatarURI(accountID string, start int64, limit int64) string {
	return fmt.Sprintf(`%s/items/avatar?filter={"owner":{"_eq":"%s"}}&meta=*&%s`, EnvVariable.DirectusBaseURI, accountID, utils.GetPageParam(start, limit))
}

func GetDirectusGetAssetURI(assetID string) string {
	return fmt.Sprintf("%s/assets/%s", EnvVariable.DirectusBaseURI, assetID)
}

func GetDirectusUploadAssetURI() string {
	return fmt.Sprintf("%s/files", EnvVariable.DirectusBaseURI)
}

func GetDirectusImportAssetURI() string {
	return fmt.Sprintf("%s/files/import", EnvVariable.DirectusBaseURI)
}

func GetDirectusCreateAvatarURI() string {
	return fmt.Sprintf("%s/items/avatar", EnvVariable.DirectusBaseURI)
}

func GetDirectusSingleAvatarURI(avatarID string) string {
	return fmt.Sprintf("%s/items/avatar/%s", EnvVariable.DirectusBaseURI, avatarID)
}
