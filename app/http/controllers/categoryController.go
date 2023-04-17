package controllers

import (
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type CategoryController struct{}

func (ctrl CategoryController) FetchChild(id int32) []map[string]any {
	var values []map[string]any

	if err := utils.DB.Table(utils.SingularName+"_categories").Where("parent_id = ?", id).Find(&values).Error; err != nil {
		return values
	}

	transformer, _ := utils.JsonFileParser("transformers/response/" + utils.SingularName + "_categories/find.json")
	customResponses := utils.MultiMapValuesShifter(values, transformer)

	for _, value := range customResponses {
		value["childs"] = ctrl.FetchChild(value["id"].(int32))
	}

	return customResponses
}

func (ctrl CategoryController) Find(ctx *gin.Context) {
	var value map[string]any

	if err := utils.DB.Table(utils.SingularName+"_categories").Where("id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if value["id"] == nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", utils.SingularName+" not found", nil))
		return
	}

	transformer, _ := utils.JsonFileParser("transformers/response/" + utils.SingularName + "_categories/find.json")
	utils.MapValuesShifter(transformer, value)

	transformer["childs"] = ctrl.FetchChild(transformer["id"].(int32))

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+utils.SingularName+" success", transformer))
}

func (ctrl CategoryController) FindAll(ctx *gin.Context) {
	var values []map[string]any

	if err := utils.DB.Table(utils.SingularName + "_categories").Find(&values).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	transformer, _ := utils.JsonFileParser("transformers/response/" + utils.SingularName + "_categories/find.json")
	customResponses := utils.MultiMapValuesShifter(values, transformer)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+utils.PluralName+" success", customResponses))
}

func (ctrl CategoryController) Create(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("transformers/request/" + utils.SingularName + "_categories/create.json")
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

	name, _ := transformer["name"].(string)
	transformer["slug"] = slug.Make(name)

	if err := utils.DB.Table(utils.SingularName + "_categories").Create(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+utils.SingularName+" success", transformer))
}

func (ctrl CategoryController) Update(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("transformers/request/" + utils.SingularName + "_categories/update.json")
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

	name, _ := transformer["name"].(string)
	transformer["slug"] = slug.Make(name)

	if err := utils.DB.Table(utils.SingularName+"_categories").Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	// todo : make a better response!
	ctrl.Find(ctx)
}

// todo : need to check constraint error
func (ctrl CategoryController) Delete(ctx *gin.Context) {
	if err := utils.DB.Table(utils.SingularName+"_groups").Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+utils.SingularName+" success", nil))
}
