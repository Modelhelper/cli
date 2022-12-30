package model

import "context"

// THESE MODELS ARE NOW OBSOLETE.
// The models are moved to the modelhelper.models.go file.
type ModelConverter interface {
	ToModel(ctx context.Context) interface{}
}
