package converter

import (
	"github.com/Modelhelper/cli/v3/modelhelper"
	"github.com/Modelhelper/cli/v3/modelhelper/models"
)

type projectModelConverter struct{}

// ToProjectModel implements modelhelper.ProjectModelConverter
func (p *projectModelConverter) ToProjectModel(cfg *models.Config, options *models.ProjectTemplateCreateOptions) *models.ProjectTemplateModel {
	panic("unimplemented")
}

func NewProjectModelConverter() modelhelper.ProjectModelConverter {
	return &projectModelConverter{}
}

// type ProjectModelConverter interface {
// 	ToProjectModel(cfg *models.Config, options *models.ProjectTemplateCreateOptions) *models.ProjectTemplateModel
// }
