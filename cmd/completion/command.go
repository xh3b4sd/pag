package completion

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	name             = "completion"
	shortDescription = "Generate shell completions."
	longDescription  = `Generating zsh completion for Oh My Zsh can be done by writing the
generated completion to the custom plugin folder.

    mkdir -p ~/.oh-my-zsh/custom/plugins/workflow && workflow completion zsh > ~/.oh-my-zsh/custom/plugins/workflow/_workflow

	`
)

type Config struct {
	Logger logger.Interface
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var c *cobra.Command
	{
		r := &runner{
			logger: config.Logger,
		}

		c = &cobra.Command{
			Use:                   name,
			Short:                 shortDescription,
			Long:                  longDescription,
			DisableFlagsInUseLine: true,
			ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
			Args:                  cobra.ExactValidArgs(1),
			RunE:                  r.Run,
		}
	}

	return c, nil
}
