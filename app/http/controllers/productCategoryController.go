package controllers

import (
	"errors"
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func init() {
	label = "product category"
	table = "product_categories"
}

func FindProductCategory(ctx *gin.Context) {
	var value map[string]interface{}
	err := utils.DB.Table(table).Where("id = ?", ctx.Param("id")).Take(&value).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if value["id"] == nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", label+" not found", nil))
		return
	}

	transformer, _ := utils.JsonFileParser("transformers/response/" + label + "/find.json")
	customResponse := transformer["product"]

	utils.MapValuesShifter(transformer, value)

	if customResponse != nil {
		utils.MapValuesShifter(customResponse.(map[string]any), value)
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+label+" success", transformer))
}

func FindProductCategories(ctx *gin.Context) {
	var values []map[string]interface{}
	err := utils.DB.Table(table).Find(&values).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	var customResponses []map[string]any
	for _, value := range values {
		transformer, _ := utils.JsonFileParser("transformers/response/product/find.json")
		customResponse := transformer["product"]

		utils.MapValuesShifter(transformer, value)
		if customResponse != nil {
			utils.MapValuesShifter(customResponse.(map[string]any), value)
		}
		customResponses = append(customResponses, transformer)
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+label+"s success", customResponses))
}

func UpdateProductCategory(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("transformers/request/" + label + "/update.json")
	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	queryResult := utils.DB.Table(table).Where("id = ?", ctx.Param("id")).Updates(transformer)

	if queryResult.Error != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(queryResult.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", queryResult.Error.Error(), nil))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", queryResult.Error.Error(), nil))
		return
	}

	// todo : make a better response!
	FindProductCategory(ctx)
}

func CreateProductCategory(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("transformers/request/" + label + "/create.json")
	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	queryResult := utils.DB.Table(table).Create(&transformer)

	if queryResult.Error != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(queryResult.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", queryResult.Error.Error(), nil))
			return
		}

		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", queryResult.Error.Error(), nil))
		return
	}

	/*
		todo :
		- make a better response!
		- find hout how to return last ID without model
	*/
	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+label+" success", transformer))
}
