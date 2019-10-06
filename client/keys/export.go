package keys

import (
	"bufio"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
)

func exportKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export private keys",
		Long:  `Export a private key from the local keybase in ASCII-armored encrypted format.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runExportCmd,
	}
	cmd.Flags().Bool(flags.FlagLegacyKeybase, false, "Use legacy on-disk keybase")
	return cmd
}

func runExportCmd(cmd *cobra.Command, args []string) error {
	var err error
	var kb keys.Keybase
	var encryptPassword string

	inBuf := bufio.NewReader(cmd.InOrStdin())
	decryptPassword := DefaultKeyPass
	if !viper.GetBool(flags.FlagLegacyKeybase) {
		kb = NewKeyring(bufio.NewReader(inBuf))
	} else {
		cmd.PrintErrln(deprecatedKeybaseWarning)
		kb, err = NewKeyBaseFromHomeFlag()
		if err != nil {
			return err
		}

		decryptPassword, err = input.GetPassword("Enter passphrase to decrypt your key:", inBuf)
		if err != nil {
			return err
		}
	}

	encryptPassword, err = input.GetPassword("Enter passphrase to encrypt the exported key:", inBuf)
	if err != nil {
		return err
	}

	armored, err := kb.ExportPrivKey(args[0], decryptPassword, encryptPassword)
	if err != nil {
		return err
	}

	cmd.Println(armored)
	return nil
}
