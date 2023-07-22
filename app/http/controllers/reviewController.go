package controllers

import (
	"net/http"

	"github.com/62teknologi/62whale/62golib/utils"
	"github.com/62teknologi/62whale/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type ReviewController struct {
	SingularName  string
	PluralName    string
	SingularLabel string
	PluralLabel   string
	Table         string
}

func (ctrl *ReviewController) Init(ctx *gin.Context) {
	ctrl.SingularName = utils.Pluralize.Singular(ctx.Param("table"))
	ctrl.PluralName = utils.Pluralize.Plural(ctx.Param("table"))
	ctrl.SingularLabel = ctrl.SingularName + " review"
	ctrl.PluralLabel = ctrl.SingularName + " reviews"
	ctrl.Table = ctrl.SingularName + "_reviews"
}

func (ctrl ReviewController) Find(ctx *gin.Context) {
	ctrl.Init(ctx)

	value := map[string]any{}
	columns := []string{ctrl.Table + ".*"}
	order := "id desc"
	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/response/" + ctrl.Table + "/find.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	query := utils.DB.Table(ctrl.Table)

	utils.SetBelongsTo(query, transformer, &columns, ctx)
	delete(transformer, "filterable")

	if err := query.Select(columns).Order(order).Where(ctrl.Table+"."+"id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.SingularLabel+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachBelongsTo(transformer, value)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl ReviewController) FindAll(ctx *gin.Context) {
	ctrl.Init(ctx)

	values := []map[string]any{}
	columns := []string{ctrl.Table + ".*"}
	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/response/" + ctrl.Table + "/find.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	query := utils.DB.Table(ctrl.Table)
	filter := utils.SetFilterByQuery(query, transformer, ctx)
	search := utils.SetGlobalSearch(query, transformer, ctx)

	utils.SetOrderByQuery(query, ctx)
	utils.SetBelongsTo(query, transformer, &columns, ctx)

	delete(transformer, "filterable")
	delete(transformer, "searchable")

	pagination := utils.SetPagination(query, ctx)

	if err := query.Select(columns).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.PluralLabel+" not found", nil))
		return
	}

	customResponses := utils.MultiMapValuesShifter(transformer, values)
	summary := utils.GetSummary(transformer, values)

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+ctrl.PluralLabel+" success", customResponses, pagination, filter, search, summary))
}

func (ctrl ReviewController) Create(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/request/" + ctrl.Table + "/create.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	input := utils.ParseForm(ctx)

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "validation", validation.Errors))
		return
	}

	if input["name"] != nil && transformer["slug"] == "" {
		name, _ := input["name"].(string)
		transformer["slug"] = slug.Make(name)
	} else if transformer["slug"] == "" {
		transformer["slug"] = uuid.New()
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	if err := utils.DB.Table(ctrl.Table).Create(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl ReviewController) Update(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/request/" + ctrl.Table + "/update.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	input := utils.ParseForm(ctx)

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "validation", validation.Errors))
		return
	}

	// not sure is it needed or not, may confusing if slug changes
	if input["name"] != nil && transformer["slug"] == "" {
		name, _ := input["name"].(string)
		transformer["slug"] = slug.Make(name)
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	if err := utils.DB.Table(ctrl.Table).Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "update "+ctrl.SingularLabel+" success", transformer))
}

// todo : need to check constraint error
func (ctrl ReviewController) Delete(ctx *gin.Context) {
	ctrl.Init(ctx)

	if err := utils.DB.Table(ctrl.Table).Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+ctrl.SingularLabel+" success", nil))
}

func (ctrl ReviewController) DeleteByQuery(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/request/" + ctrl.PluralName + "/delete.json")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	query := utils.DB.Table(ctrl.PluralName)
	utils.SetFilterByQuery(query, transformer, ctx)

	if err := query.Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+ctrl.SingularLabel+" success", nil))
}
