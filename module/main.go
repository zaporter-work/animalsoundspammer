package main

import (
	"context"

	"github.com/edaniels/golog"
	templategomodule "github.com/zaporter-work/animalsoundspammer"
	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/module"

	"go.viam.com/utils"
)

var (
	Version     = "development"
	GitRevision = ""
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("animalSoundSpammer"))
}
func mainWithArgs(ctx context.Context, args []string, logger golog.Logger) error {
	var versionFields []interface{}
	if Version != "" {
		versionFields = append(versionFields, "version", Version)
	}
	if GitRevision != "" {
		versionFields = append(versionFields, "git_rev", GitRevision)
	}
	if len(versionFields) != 0 {
		logger.Infow("animalSoundSpammer", versionFields...)
	} else {
		logger.Info("animalSoundSpammer" + " built from source; version unknown")
	}
	mod, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	mod.AddModelFromRegistry(ctx, generic.API, templategomodule.Model)

	mod.Start(ctx)
	defer mod.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
