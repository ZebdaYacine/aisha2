package container

import "log/slog"

func closeInfrastructure(infra Infrastructure) {
	if infra.Redis != nil {
		if err := infra.Redis.Close(); err != nil {
			slog.Error("close redis client", "error", err)
		}
	}
	if infra.Database != nil {
		infra.Database.Close()
	}
}
