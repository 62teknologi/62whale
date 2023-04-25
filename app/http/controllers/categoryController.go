package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"whale/62teknologi-golang-utility/utils"

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

	fmt.Println("table")
	fmt.Println(ctrl.Table)

	utils.SetJoin(query, transformer, &columns)

	if err := query.Select(columns).Where(ctrl.Table+".id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.SingularName+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachJoin(transformer, value)

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
	order := "id desc"
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	query := utils.DB.Table(ctrl.Table)
	filter := utils.SetFilterByQuery(query, transformer, ctx)
	pagination := utils.SetPagination(query, ctx)
	utils.SetJoin(query, transformer, &columns)

	if err := query.Select(columns).Order(order).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.PluralLabel+" not found", nil))
		return
	}

	customResponses := utils.MultiMapValuesShifter(values, transformer)

	if ctx.Query("include_childs") != "" {
		var total int32 = 1
		for _, value := range customResponses {
			if value["id"] != nil {
				value["childs"] = ctrl.FetchChild(value["id"].(int32), []string{}, &total)
			}
		}
		fmt.Printf("total queries for "+ctrl.Table+" is %d\n", total)
	}

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+ctrl.PluralLabel+" success", customResponses, pagination, filter))
}

func (ctrl *CategoryController) Create(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/create.json")
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

	var name string

	if transformer["name"] != nil {
		name, _ = transformer["name"].(string)
		transformer["slug"] = slug.Make(name)
	} else {
		transformer["slug"] = uuid.New()
	}

	if err := utils.DB.Table(ctrl.Table).Create(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *CategoryController) Update(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/update.json")
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

	var name string

	if transformer["name"] != nil {
		name, _ = transformer["name"].(string)
		// not sure is it needed or not, may confusing if slug changes
		transformer["slug"] = slug.Make(name)
	}

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
	customResponses := utils.MultiMapValuesShifter(values, transformer)

	for _, value := range customResponses {
		value["childs"] = ctrl.FetchChild(value["id"].(int32), sequence, total)
		delete(value, "filterable")
	}

	return customResponses
}
