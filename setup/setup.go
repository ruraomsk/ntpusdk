package setup

var (
	Set *Setup
)

type Setup struct {
	LogPath string `toml:"logpath"`
	NtpPort int    `toml:"ntp-port"`
}

func init() {
}
