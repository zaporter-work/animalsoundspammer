package templategomodule

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/edaniels/golog"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"

	"go.viam.com/utils"
)

var (
	Model = resource.NewModel("zack", "random", "animalsoundspammer")
)

func init() {
	registration := resource.Registration[resource.Resource, *Config]{
		Constructor: func(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger golog.Logger) (resource.Resource, error) {
			return createComponent(ctx, deps, conf, logger)
		},
	}
	resource.RegisterComponent(generic.API, Model, registration)
}

type component struct {
	resource.Named
	resource.AlwaysRebuild
	cfg *Config

	cancelCtx     context.Context
	cancelFunc    func()
	moduleDataDir string
	vs     vision.Service

	logger golog.Logger
}

func createComponent(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger golog.Logger) (resource.Resource, error) {
    
	logger.Warnln("create Component")
	ctx, span := trace.StartSpan(ctx, "zaporter::New")
	defer span.End()
	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return nil, errors.Wrap(err, "create component failed due to config parsing")
	}
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	instance := &component{
		Named:         conf.ResourceName().AsNamed(),
		cfg:           newConf,
		cancelCtx:     cancelCtx,
		cancelFunc:    cancelFunc,
		moduleDataDir: os.Getenv("VIAM_MODULE_DATA"),
		logger:        logger,
	}
	if instance.moduleDataDir != "" {
		instance.logger.Infoln("data_dir", instance.moduleDataDir)
		content := "hello - zack"
		os.WriteFile(filepath.Join(instance.moduleDataDir, "info.txt"), []byte(content), os.ModePerm)
	}

	instance.logger.Infoln("message", newConf.Message)
	if len(instance.cfg.Animals) == 0 {
		instance.cfg.Animals = []string{"cow", "pig", "goat", "elephant"}
	}
	instance.startBgProcess()
	return instance, nil
}
func (c *component) startBgProcess() {
	utils.PanicCapturingGo(func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if os.Getenv("HI_ZACK") == "hi" {
					c.logger.Info(fmt.Sprintf("zack says: %s", c.cfg.Message))
				}
				c.logger.Info(fmt.Sprintf("the %s says: %s", c.cfg.Animals[rand.Intn(len(c.cfg.Animals))], c.cfg.Message))
			case <-c.cancelCtx.Done():
				c.logger.Info("shutdown")
				return
			}

		}

	})

}

// // Reconfigure must reconfigure the resource atomically and in place. If this
// // cannot be guaranteed, then usage of AlwaysRebuild or TriviallyReconfigurable
// // is permissible.
// func (c *component) Reconfigure(ctx context.Context, deps Dependencies, conf Config) error {
// return nil;
// }

// DoCommand sends/receives arbitrary data
func (c *component) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	c.logger.Info("docommand")
	return make(map[string]interface{}), nil
}

// Close must safely shut down the resource and prevent further use.
// Close must be idempotent.
// Later reconfiguration may allow a resource to be "open" again.
func (c *component) Close(ctx context.Context) error {
	c.logger.Info("close")
	c.cancelFunc()
	return nil
}
