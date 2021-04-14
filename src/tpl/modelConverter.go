package tpl

import (
	"modelhelper/cli/types"
)

type ModelConverter interface {
	ToDataModel() (types.EntityImportModel, error)
}
