package completion

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type runner struct {
	logger logger.Interface
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.run(ctx, cmd, args)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error
	switch args[0] {
	case "bash":
		err = cmd.Root().GenBashCompletion(os.Stdout)
		if err != nil {
			return tracer.Mask(err)
		}
	case "zsh":
		err = cmd.Root().GenZshCompletion(os.Stdout)
		if err != nil {
			return tracer.Mask(err)
		}
	case "fish":
		err = cmd.Root().GenFishCompletion(os.Stdout, true)
		if err != nil {
			return tracer.Mask(err)
		}
	case "powershell":
		err = cmd.Root().GenPowerShellCompletion(os.Stdout)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
