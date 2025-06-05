package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/shahfaizansr/initilizer"
	sqldb "github.com/shahfaizansr/migrate/db"
	"github.com/shahfaizansr/models"
	"github.com/shahfaizansr/services/mybulkcalcservice"
	"github.com/shahfaizansr/services/mycalcservice"
	"github.com/shahfaizansr/services/myhelloservice"
	"github.com/shahfaizansr/utils"
	"github.com/shahfaizansr/utils/cvlconstant"
	"github.com/shahfaizansr/utils/rigel"
)

func init() {
	initilizer.LoadEnvFile()
	initilizer.LoadDataBase()
}

func main() {

	var (
		l    *logharbour.Logger
		lctx *logharbour.LoggerContext
	)

	ctx := context.Background()

	environment := flag.String("env", "dev_env", "setup the environment using CLI")
	flag.Parse()

	appConfig := utils.SetConfigEnvironment(models.Environment(*environment))

	// Open a file for logging.
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new Rigel instance
	rigelClient := rigel.New(appConfig)

	loggerPriority, err := rigelClient.GetString(ctx, cvlconstant.LOGGER_PRIORITY_KEY)
	if err != nil {
		log.Fatalf("Failed to get log priority from rigel: %v\n", err)
	}
	currentLogPriority := utils.GetCurrentLoggerPriority(loggerPriority)

	if appConfig.KafkaConfig.IsKafkaOn {
		// Create a Kafka writer
		// Define your Kafka configuration
		kafkaConfig := logharbour.KafkaConfig{
			Brokers: appConfig.KafkaConfig.KafkaBrokers,
			Topic:   appConfig.KafkaConfig.KafkaTopic,
		}

		// Define the maximum number of connections in the pool
		poolSize := 10

		kafkaWriter, err := logharbour.NewKafkaWriter(kafkaConfig, logharbour.WithPoolSize(poolSize))
		if err != nil {
			log.Fatalf("Failed to create Kafka writer: %v", err)
		}
		// Create a fallback writer that uses stdout as the fallback.
		fallbackWriter := logharbour.NewFallbackWriter(kafkaWriter, os.Stdout)

		// Create a logger context with the default priority.
		lctx := logharbour.NewLoggerContext(currentLogPriority)

		// Initialize the logger with the context, validator, default priority, and fallback writer.
		l = logharbour.NewLoggerWithFallback(lctx, "CVL-KRA", fallbackWriter)
		//	writer = kafkaWriter
	} else {
		//	writer = os.Stdout
		lctx = logharbour.NewLoggerContext(currentLogPriority)
		l = logharbour.NewLogger(lctx, cvlconstant.APP_NAME, os.Stdout)
		lctx.ChangeMinLogPriority(currentLogPriority)

	}

	// get data from rigel
	connURL, err := rigel.GetConnURLFromRigel(rigelClient, ctx, l)
	if err != nil {
		l.Err().Error(err).Log("Error while getting data from rigel")
		log.Fatalf("Failed to get data from rigel: %v", err)
		return
	}

	// Create a LogLevel instance to control pgx logging level at runtime.
	logLevel := &models.LogLevel{}
	// Set initial log level (e.g., LogLevelInfo).
	logLevel.Set(tracelog.LogLevelDebug)

	dbconnn, err := sqldb.NewSQLServerHandler(l, rigelClient)

	if err != nil {
		l.Err().Log("Error while establishes a connection with database")
		log.Fatalln("Failed to establishes a connection with database")
		return
	}
	l.LogActivity("Connection with database established", connURL)

	// redis setup
	l.Log("Establishing the connection with Redis")
	// redisClient := redis_utils.NewRedisClient(ctx, l, appConfig)
	// l.Log("Connection with redis established")

	// // MinIO Setup
	// l.Log("Establishing the connection with Minio")

	// minioClientModel, err := minio_utils.BuildMinIOSetup(ctx, appConfig.MinIOInfo)
	// if err != nil {
	// 	l.Err().Error(err).Log("Failed to initialise MinIO Setup")
	// 	log.Fatalln("Failed to initialise MinIO Setup", err)
	// 	return
	// }
	// l.Log("Connection with minio established")

	// jm := jobs.NewJobManager(dbconnn, redis_utils.GetRedisClient(), minio_utils.GetMinioClient(), l, &jobs.JobManagerConfig{
	// 	BatchStatusCacheDurSec: 0,
	// 	BatchOutputBucket:      "cvlkra",
	// })

	// // Register the custom initializer for batch type "mybulkcalc".
	// err = jm.RegisterInitializer(cvlconstant.MyBulkCalc.String(), &mybatch.MyCalcInitializer{})
	// if err != nil {
	// 	l.Err().Error(err).Log("Failed to register initializer")
	// 	log.Fatal("Failed to register initializer:", err)
	// 	return
	// }

	// // Register the processor that will be used to handle each input line for batch type "mybulkcalc" and process name "bulk_operation".
	// err = jm.RegisterProcessorBatch(cvlconstant.MyBulkCalc.String(), cvlconstant.BULK_OPERATION, &mybatch.MyBulkCalcBathProcessor{})
	// if err != nil {
	// 	l.Err().Error(err).Log("Failed to register process batch")
	// 	log.Fatal("Failed to register process batch: ", err)
	// }

	// go jm.Run()

	router := gin.Default()

	myHelloService := service.NewService(router).
		WithLogHarbour(l)
	myHelloService.RegisterRoute(cvlconstant.POST_METHOD, cvlconstant.FORWARD_SLASH+cvlconstant.MyHello.String(), myhelloservice.HelloHandler)

	myCalcService := service.NewService(router).
		WithLogHarbour(l).
		WithDatabase(dbconnn)
	myCalcService.RegisterRoute(cvlconstant.POST_METHOD, cvlconstant.FORWARD_SLASH+cvlconstant.MyCalc.String(), mycalcservice.CalculatorHandler)

	myBulkCalcService := service.NewService(router).
		WithLogHarbour(l).
		WithDatabase(dbconnn)
	myBulkCalcService.RegisterRoute(cvlconstant.POST_METHOD, cvlconstant.FORWARD_SLASH+cvlconstant.MyBulkCalc.String(), mybulkcalcservice.HandleMyBulkCalc)

	// myBatchService := service.NewService(router).
	// 	WithLogHarbour(l).
	// 	WithDatabase(dbconnn).
	// 	WithDependency("redisClient", redisClient).
	// 	WithDependency("minio", minioClientModel).
	// 	WithDependency("logger", l).
	// 	WithDependency("appConfig", appConfig).
	// 	WithRigelConfig(rigelClient)
	// 	// WithDependency("jobmanager", jm)
	// myBatchService.RegisterRoute(cvlconstant.POST_METHOD, cvlconstant.FORWARD_SLASH+cvlconstant.MyBulkCalcBatch.String(), mybatchservice.HandleMyBulkCalcBatch)

	wscutils.SetMsgIDInvalidJSON(1001)
	wscutils.SetErrCodeInvalidJSON(cvlconstant.INVALID_JSON_FORMAT)

	router.Run()
}
