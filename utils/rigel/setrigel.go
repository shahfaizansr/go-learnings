package rigel

import (
	"log"

	"github.com/remiges-tech/rigel"
	"github.com/remiges-tech/rigel/etcd"
	"github.com/shahfaizansr/models"
)

func SetRigel(appConfig models.AppConfig) *rigel.Rigel {

	// Create a new EtcdStorage instance
	//	logger.Log("Creating a new instance of EtcdStorage")
	etcdStorage, err := etcd.NewEtcdStorage([]string{"http://localhost:2380"})
	if err != nil {
		//	logger.Err().Error(err).Log("Error while Creating new instance of EtcdStorage")
		log.Fatalf("Failed to create EtcdStorage: %v", err)
	}
	//	logger.LogActivity("Created a new instance of EtcdStorage with endpoints", appConfig.EtcdEndpoint)

	rigel := rigel.New(etcdStorage, appConfig.Rigel.AppName, appConfig.Rigel.ModuleName, appConfig.Rigel.VersionNumber, appConfig.Rigel.ConfigName)

	return rigel
}
