package psmdbackup

import (
	"context"
	"io"

	pbapi "github.com/percona/percona-backup-mongodb/proto/api"
	"github.com/percona/percona-server-mongodb-operator/pkg/psmdb/backup"
	"google.golang.org/grpc"
)

// newBackup creates new Backup
func newBackup(crName, backupName, storageName string) (Backup, error) {
	b := Backup{}
	conn, err := grpc.Dial(crName+backup.GetCoordinatorSuffix(), nil)
	if err != nil {
		return b, err
	}
	client := pbapi.NewApiClient(conn)
	b = Backup{
		Client:      client,
		StorageName: storageName,
		Name:        backupName,
	}
	return b, nil
}

// Backup implements BC
type Backup struct {
	Client      pbapi.ApiClient
	Status      string
	Start       int64
	End         int64
	Type        string
	Name        string
	StorageName string
}

// BackupExist is check if backup exist and update it status if true
func (b *Backup) BackupExist() (bool, error) {
	stream, err := b.Client.BackupsMetadata(context.TODO(), &pbapi.BackupsMetadataParams{})
	if err != nil {
		return false, err
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return false, err
		}
		if msg.Metadata.Description == b.Name {
			if msg.Metadata.EndTs > 0 {
				b.Status = "ready"
			}
			return true, nil

		}
	}
	return false, nil
}

// StartBackup is for starting new backup
func (b *Backup) StartBackup() error {
	msg := &pbapi.RunBackupParams{
		CompressionType: pbapi.CompressionType_COMPRESSION_TYPE_NO_COMPRESSION,
		Cypher:          pbapi.Cypher_CYPHER_NO_CYPHER,
		Description:     b.Name,
		StorageName:     b.StorageName,
	}
	b.Client.RunBackup(context.Background(), msg)
	return nil
}

// CheckBackupStatus is for checking backup status
func (b *Backup) CheckBackupStatus() error {
	return nil
}
