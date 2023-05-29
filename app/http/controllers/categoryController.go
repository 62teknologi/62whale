package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/62teknologi/62whale/62golib/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type CategoryController struct {
	SingularName  string
	PluralName    string
	SingularLabel string
	PluralLabel   string
	Table         string
}

func (ctrl *CategoryController) Init(ctx *gin.Context) {
	ctrl.SingularName = utils.Pluralize.Singular(ctx.Param("table"))
	ctrl.PluralName = utils.Pluralize.Plural(ctx.Param("table"))
	ctrl.SingularLabel = ctrl.SingularName + " category"
	ctrl.PluralLabel = ctrl.SingularName + " categories"
	ctrl.Table = ctrl.SingularName + "_categories"
}

func (ctrl *CategoryController) Find(ctx *gin.Context) {
	ctrl.Init(ctx)

	value := map[string]any{}
	columns := []string{ctrl.Table + ".*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	query := utils.DB.Table(ctrl.Table)

	utils.SetOrderByQuery(query, ctx)
	utils.SetBelongsTo(query, transformer, &columns)

	delete(transformer, "filterable")

	if err := query.Select(columns).Where(ctrl.Table+".id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.SingularName+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachBelongsTo(transformer, value)

	if transformer["id"] != nil {
		total := int32(1)
		transformer["childs"] = ctrl.FetchChild(transformer["id"].(int32), []string{}, &total)
		fmt.Printf("total queries for "+ctrl.Table+" where parent id %d is %d\n", transformer["id"], total)
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *CategoryController) FindAll(ctx *gin.Context) {
	ctrl.Init(ctx)

	values := []map[string]any{}
	columns := []string{ctrl.Table + ".*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	query := utils.DB.Table(ctrl.Table)
	filter := utils.SetFilterByQuery(query, transformer, ctx)
	search := utils.SetGlobalSearch(query, transformer, ctx)

	utils.SetOrderByQuery(query, ctx)
	utils.SetBelongsTo(query, transformer, &columns)

	delete(transformer, "filterable")
	delete(transformer, "searchable")

	pagination := utils.SetPagination(query, ctx)

	if err := query.Select(columns).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.PluralLabel+" not found", nil))
		return
	}

	customResponses := utils.MultiMapValuesShifter(transformer, values)
	summary := utils.GetSummary(transformer, values)

	if ctx.Query("include_childs") != "" {
		var total int32 = 1
		for _, value := range customResponses {
			if value["id"] != nil {
				value["childs"] = ctrl.FetchChild(value["id"].(int32), []string{}, &total)
			}
		}
		fmt.Printf("total queries for "+ctrl.Table+" is %d\n", total)
	}

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+ctrl.PluralLabel+" success", customResponses, pagination, filter, search, summary))
}

func (ctrl *CategoryController) Create(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/create.json")
	input := utils.ParseForm(ctx)

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
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

func (ctrl *CategoryController) Update(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/update.json")
	input := utils.ParseForm(ctx)

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
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
func (ctrl *CategoryController) Delete(ctx *gin.Context) {
	ctrl.Init(ctx)

	if err := utils.DB.Table(ctrl.Table).Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+ctrl.SingularLabel+" success", nil))
}

// todo : this will generate N queries of N row. need cached mechanism to prevent that.
func (ctrl *CategoryController) FetchChild(id int32, sequence []string, total *int32) []map[string]any {
	*total = *total + 1
	var values []map[string]any

	sequence = append(sequence, strconv.Itoa(int(int32(id))))

	if err := utils.DB.Table(ctrl.Table).Where("parent_id = ?", id).Where("id NOT IN ?", sequence).Find(&values).Error; err != nil {
		return values
	}

	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	customResponses := utils.MultiMapValuesShifter(transformer, values)

	for _, value := range customResponses {
		value["childs"] = ctrl.FetchChild(value["id"].(int32), sequence, total)
		delete(value, "filterable")
	}

	return customResponses
}
