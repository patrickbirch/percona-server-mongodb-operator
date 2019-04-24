package backup

const (
	coordinatorAPIPort       = 10001
	coordinatorRPCPort       = 10000
	coordinatorDataMount     = "/data"
	coordinatorSuffix        = "-backup-coordinator"
	coordinatorContainerName = "-backup-coordinator"
	coordinatorDataVolume    = "backup-metadata"
	coordinatorRPCPortName   = "rpc"
	coordinatorAPIPortName   = "api"

	backupCtlContainerName = "backup-pmbctl"

	agentConfigDir              = "/etc/percona-backup-mongodb"
	agentStoragesConfigFile     = "storages.yml"
	agentContainerName          = "backup-agent"
	awsAccessKeySecretKey       = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKeySecretKey = "AWS_SECRET_ACCESS_KEY"
)

// GetCoordinatorSuffix is for getting ccordinator suffix
func GetCoordinatorSuffix() string {
	return coordinatorSuffix
}
