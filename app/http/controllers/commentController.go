package controllers

import (
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
)

type CommentController struct{}

func (ctrl CommentController) FetchChild(id int32) []map[string]any {
	var values []map[string]any

	if err := utils.DB.Table(utils.SingularName+"_comments").Where("parent_id = ?", id).Find(&values).Error; err != nil {
		return values
	}

	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + utils.SingularName + "_comments/find.json")
	customResponses := utils.MultiMapValuesShifter(values, transformer)

	for _, value := range customResponses {
		value["childs"] = ctrl.FetchChild(value["id"].(int32))
	}

	return customResponses
}

func (ctrl CommentController) Find(ctx *gin.Context) {
	value := map[string]any{}
	columns := []string{utils.SingularName + "_comments.*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + utils.SingularName + "_comments/find.json")
	query := utils.DB.Table(utils.SingularName + "_comments")

	utils.SetJoin(query, transformer, &columns)

	if err := query.Select(columns).Where(utils.SingularName+"_comments."+"id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", utils.SingularName+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachJoin(transformer, value)

	if transformer["id"] != nil {
		transformer["childs"] = ctrl.FetchChild(transformer["id"].(int32))
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+utils.SingularName+" success", transformer))
}

func (ctrl CommentController) FindAll(ctx *gin.Context) {
	values := []map[string]any{}
	columns := []string{utils.SingularName + "_comments.*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + utils.SingularName + "_comments/find.json")
	query := utils.DB.Table(utils.SingularName + "_comments")

	utils.SetJoin(query, transformer, &columns)

	filterable, _ := utils.JsonFileParser("setting/filter/" + utils.SingularName + "_comments/find.json")
	filter := utils.FilterByQueries(query, filterable, ctx)
	pagination := utils.SetPagination(query, ctx)

	if err := query.Select(columns).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	customResponses := utils.MultiMapValuesShifter(values, transformer)

	if ctx.Query("include_childs") != "" {
		for _, value := range customResponses {
			if value["id"] != nil {
				value["childs"] = ctrl.FetchChild(value["id"].(int32))
			}
		}
	}

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+utils.PluralName+" success", customResponses, pagination, filter))
}

func (ctrl CommentController) Create(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + utils.SingularName + "_comments/create.json")
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

	if err := utils.DB.Table(utils.SingularName + "_comments").Create(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+utils.SingularName+" success", transformer))
}

func (ctrl CommentController) Update(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + utils.SingularName + "_comments/update.json")
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

	if err := utils.DB.Table(utils.SingularName+"_comments").Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "update "+utils.SingularName+" success", transformer))
}

// todo : need to check constraint error
func (ctrl CommentController) Delete(ctx *gin.Context) {
	if err := utils.DB.Table(utils.SingularName+"_comments").Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+utils.SingularName+" success", nil))
}
