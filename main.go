package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go"

	model "grant-canyon-entity-api/model"
	utils "grant-canyon-entity-api/utils"
)

var client *opensearch.Client

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	log.Println("Application running in environment: ", config.RuntimeSetup, " and on port: ", config.AppPort)

	client, _ = utils.CreateOpensearchPool()

	gin.SetMode(gin.ReleaseMode)

	// Logging to a file.
	f, _ := os.Create("entity-api.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		header := param.Request.Header
		errorMessageLen := len(param.ErrorMessage)
		if errorMessageLen == 0 {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" \nHeader: %s",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				header,
			)
		} else {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" Error: %s\"\nHeader: %s",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
				header,
			)
		}
	}))

	router.Use(cors.Default())

	router.GET("/IsAlive", isAlive)

	router.POST("/service/fund/v1/funds", postFunds)
	router.GET("/service/fund/v1/fund_id/:id", getFundById)
	router.POST("/service/fund/v1/search", searchFund)

	router.Run(config.ServerAddress + ":" + config.AppPort)
}

func isAlive(c *gin.Context) {
	c.JSON(http.StatusOK, "Is Alive")
}

func searchFund(c *gin.Context) {
	request := new(model.SearchFundRequest)

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(request); err != nil {
		return
	}

	funds := utils.SearchByTerm(client, "gc-fund-v1", request.Term, request.Fund)
	count := len(funds)
	programs := make(map[string]model.Program)
	for i := 0; i < count; i++ {
		fund := funds[i]
		if val, ok := fund["program"]; ok {
			program_name := val.(string)
			if prog, exist := programs[program_name]; exist {
				if val, ok := fund["external_id"]; ok {
					external_id := val.(string)
					prog.Funds[external_id] = fund
				}
			} else {
				prog.ProgramName = program_name
				prog.Funds = make(map[string]interface{})
				if val, ok := fund["external_id"]; ok {
					external_id := val.(string)
					prog.Funds[external_id] = fund
				}
				programs[program_name] = prog
			}
		}
	}
	programsLen := len(programs)
	res := new(model.SearchResults)
	res.Programs = make([]model.ProgramSearchResult, programsLen)
	i := 0
	for _, prog := range programs {
		resProg := new(model.ProgramSearchResult)
		resProg.ProgramName = prog.ProgramName
		fundsLen := len(prog.Funds)
		resProg.Funds = make([]interface{}, fundsLen)
		j := 0
		for _, fund := range prog.Funds {
			resProg.Funds[j] = fund
			j = j + 1
		}
		res.Programs[i] = *resProg
		i = i + 1
	}

	c.JSON(http.StatusOK, res)
}

func getFundById(c *gin.Context) {
	id := c.Param("id")

	sources := utils.SearchById(client, "gc-fund-v1", id)
	if sources != nil && len(sources) > 0 {
		c.JSON(http.StatusOK, sources)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
	}
}

func postFunds(c *gin.Context) {
	request := new(model.PutFundRequest)

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(request); err != nil {
		return
	}

	var docs []map[string]interface{}

	for _, fund := range request.Funds {
		recordBody, internalId := createFundDataOpensearch(fund)

		indexBody := map[string]interface{}{
			"_index": "gc-fund-v1",
			"_id":    internalId,
		}
		index := map[string]interface{}{
			"index": indexBody,
		}
		docs = append(docs, index)
		docs = append(docs, recordBody)

	}

	var buffer []byte
	newLineBytes := []byte("\n")
	for _, doc := range docs {
		marshelDoc, _ := json.Marshal(doc)
		buffer = append(buffer, marshelDoc...)
		buffer = append(buffer, newLineBytes...)
	}
	buffer = append(buffer, newLineBytes...)

	res, err := utils.PostBulk(client, "gc-fund-v1", string(buffer))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		c.IndentedJSON(http.StatusCreated, res)
	}
}

func createFundDataOpensearch(record model.Fund) (map[string]interface{}, string) {

	internalId := utils.GetMD5Hash(fmt.Sprintf("%s|%s", record.ExternalId, record.Fund))

	unmarsheledBody := map[string]interface{}{
		"internal_id":       internalId,
		"external_id":       record.ExternalId,
		"link":              record.Link,
		"name":              record.Name,
		"type":              record.Type,
		"fund":              record.Fund,
		"program":           record.Program,
		"call":              record.Call,
		"type_of_effort":    record.Call,
		"description":       record.Description,
		"total_budget":      record.TotalBudget,
		"grant_budget_low":  record.GrantBudgetLow,
		"grant_budget_high": record.GrantBudgetHigh,
		"currency":          record.Currency,
		"due_date":          record.DueDate,
		"submission_type":   record.SubmissionType,
		"keywords":          record.Keywords,
		"tags":              record.Tags,
	}

	return unmarsheledBody, internalId
}
