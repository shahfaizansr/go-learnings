package rigel

import (
	"context"
	"fmt"
	"log"

	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/remiges-tech/rigel"
	"github.com/remiges-tech/rigel/etcd"
	"github.com/shahfaizansr/models"
)

func New(appConfig models.AppConfig) *rigel.Rigel {

	// Create a new EtcdStorage instance
	//	logger.Log("Creating a new instance of EtcdStorage")
	etcdStorage, err := etcd.NewEtcdStorage(appConfig.EtcdEndpoint)
	if err != nil {
		//	logger.Err().Error(err).Log("Error while Creating new instance of EtcdStorage")
		log.Fatalf("Failed to create EtcdStorage: %v", err)
	}
	//	logger.LogActivity("Created a new in	stance of EtcdStorage with endpoints", appConfig.EtcdEndpoint)

	rigel := rigel.New(etcdStorage, appConfig.Rigel.AppName, appConfig.Rigel.ModuleName, appConfig.Rigel.VersionNumber, appConfig.Rigel.ConfigName)

	return rigel
}

// GetConnURLFromRigel retrieves the connection URL from Rigel.
//
// Parameters:
// - rigel: a pointer to a Rigel instance.
// - ctx: the context.Context object.
// - l: a pointer to a logharbour.Logger instance.
//
// Returns:
// - string: the connection URL.
// - error: an error if any occurred.
func GetConnURLFromRigel(rigel *rigel.Rigel, ctx context.Context, l *logharbour.Logger) (string, error) {

	dbHost, err := rigel.GetString(ctx, "db_host")
	if err != nil {
		return "", err
	}

	dbPort, err := rigel.GetInt(ctx, "db_port")
	if err != nil {
		return "", err
	}

	dbUser, err := rigel.GetString(ctx, "db_user")
	if err != nil {
		return "", err
	}

	dbPassword, err := rigel.GetString(ctx, "db_password")
	if err != nil {
		return "", err
	}

	dbName, err := rigel.GetString(ctx, "db_name")
	if err != nil {
		return "", err
	}
	l.Log("Retrieves the configuration data from rigel")

	connURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	return connURL, nil
}

func GetSeversageConfig(ctx context.Context, rigelClient *rigel.Rigel, serversage *models.Serversage) error {
	var err error
	serversage.IsInstrumentApplication, err = rigelClient.GetBool(ctx, IS_INSTRUMENT_APPLICATION)
	if err != nil {
		return err
	}
	serversage.OtelEndpoint, err = rigelClient.GetString(ctx, OTEL_ENDPOINT)
	if err != nil {
		return err
	}

	serversage.ServiceName, err = rigelClient.GetString(ctx, SERVICE_NAME)
	if err != nil {
		return err
	}
	return nil
}
