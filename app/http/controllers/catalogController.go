package controllers

import (
	"github.com/62teknologi/62whale/62golib/utils"
	"github.com/62teknologi/62whale/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strings"
)

type CatalogController struct {
	SingularName  string
	PluralName    string
	SingularLabel string
	PluralLabel   string
	Table         string
}

func (ctrl *CatalogController) Init(ctx *gin.Context) {
	ctrl.SingularName = utils.Pluralize.Singular(ctx.Param("table"))
	ctrl.PluralName = utils.Pluralize.Plural(ctx.Param("table"))
	ctrl.SingularLabel = ctrl.SingularName
	ctrl.PluralLabel = ctrl.PluralName
	ctrl.Table = ctrl.PluralName
}

func (ctrl CatalogController) Find(ctx *gin.Context) {
	ctrl.Init(ctx)

	value := map[string]any{}
	columns := []string{ctrl.PluralName + ".*"}

	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/response/" + ctrl.PluralName + "/find.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	query := utils.DB.Table(ctrl.PluralName)

	utils.SetBelongsTo(query, transformer, &columns)
	delete(transformer, "filterable")
	field := "id"
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Param("slug")
		field = "slug"
	}

	if err := query.Select(columns).Where(ctrl.PluralName+"."+field+" = ?", id).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.SingularLabel+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachBelongsTo(transformer, value)
	utils.AttachHasMany(transformer)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl CatalogController) FindAll(ctx *gin.Context) {
	ctrl.Init(ctx)

	values := []map[string]any{}
	columns := []string{ctrl.PluralName + ".*"}

	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/response/" + ctrl.PluralName + "/find.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	query := utils.DB.Table(ctrl.Table)
	filter := utils.SetFilterByQuery(query, transformer, ctx)
	search := utils.SetGlobalSearch(query, transformer, ctx)

	utils.SetOrderByQuery(query, ctx)
	utils.SetBelongsTo(query, transformer, &columns)

	delete(transformer, "filterable")
	delete(transformer, "searchable")

	pagination := utils.SetPagination(query, ctx)

	if err := query.Select(columns).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	customResponses := utils.MultiMapValuesShifter(transformer, values)
	summary := utils.GetSummary(transformer, values)

	utils.MultiAttachHasMany(customResponses)

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+ctrl.PluralLabel+" success", customResponses, pagination, filter, search, summary))
}

func (ctrl CatalogController) Create(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/request/" + ctrl.PluralName + "/create.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	input := utils.ParseForm(ctx)

	for key := range transformer {
		_, ok := input[key]

		if reflect.TypeOf(transformer[key]).Kind() != reflect.Map && reflect.TypeOf(transformer[key]).Kind() != reflect.Slice {
			if !ok && !strings.Contains(transformer[key].(string), "required") {
				delete(transformer, key)
			}
		}
	}

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

	transformer["slug"] = ""

	hasMany := transformer["has_many"]
	delete(transformer, "has_many")

	manyToMany := transformer["many_to_many"]
	delete(transformer, "many_to_many")

	_, duplicateExist := transformer["duplicate"]

	if err = utils.DB.Transaction(func(tx *gorm.DB) error {
		if duplicateExist {
			for i := range hasMany.(map[string]any) {
				transformerValues := transformer[i]
				defaultItem := utils.FilterMap(transformerValues, func(item map[string]any) bool {
					_, itemDefaultExist := item["default"]
					if itemDefaultExist {
						isDefaultItem := item["default"].(bool)
						return isDefaultItem == true
					}
					return false
				})

				if len(defaultItem) != 0 {
					utils.SetDoubleRecord(transformer, defaultItem[0], i)
				} else {
					utils.SetDoubleRecord(transformer, transformerValues.([]any)[0].(map[string]any), i)
				}
				delete(transformer, "duplicate")
			}
		}

		hasManyItems := make(map[string]any)
		if hasMany != nil {
			for i := range hasMany.(map[string]any) {
				hasManyItems[i] = transformer[i]
				delete(transformer, i)
			}
		}

		hasManyToManyGroups := make(map[string]any)
		if hasManyToManyGroups != nil {
			for i := range manyToMany.(map[string]any) {
				hasManyToManyGroups[i] = transformer[i]
				delete(transformer, i)
			}
		}

		if err = tx.Table(ctrl.PluralName).Create(&transformer).Error; err != nil {
			return err
		}

		if hasMany != nil {
			for i, v := range hasMany.(map[string]any) {
				table := v.(map[string]any)["table"].(string)
				fk := v.(map[string]any)["fk"].(string)

				tx.Table(ctrl.PluralName).Where("slug = ?", transformer["slug"]).Take(&transformer)

				items := utils.Prepare1toM(fk, transformer["id"], hasManyItems[i])

				if err = tx.Table(table).Create(&items).Error; err != nil {
					return err
				}

				transformer[i] = items
			}
		}

		if manyToMany != nil {
			for i, v := range manyToMany.(map[string]any) {
				table := v.(map[string]any)["table"].(string)
				fk1 := v.(map[string]any)["fk_1"].(string)
				fk2 := v.(map[string]any)["fk_2"].(string)

				tx.Table(ctrl.PluralName).Where("slug = ?", transformer["slug"]).Take(&transformer)
				groups := utils.PrepareMtoM(fk1, transformer["id"], fk2, hasManyToManyGroups[i])

				if err = tx.Table(table).Create(&groups).Error; err != nil {
					return err
				}

				transformer[i] = groups
			}
		}

		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl CatalogController) Update(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, err := utils.JsonFileParser(config.Data.SettingPath + "/transformers/request/" + ctrl.PluralName + "/update.json")
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

	item, item_exist := transformer["items"]
	group, groups_exist := transformer["categories"]

	delete(transformer, "items")
	delete(transformer, "groups")

	if err := utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(ctrl.PluralName).Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
			return err
		}

		if item_exist || groups_exist {
			if item_exist {
				items := utils.Prepare1toM(ctrl.SingularName+"_id", ctx.Param("id"), item)

				if err := tx.Table(ctrl.SingularName + "_items").Create(&items).Error; err != nil {
					return err
				}

				transformer["items"] = items
			}

			if groups_exist {
				tx.Table(ctrl.PluralName+"_groups").Where(ctrl.SingularName+"_id = ?", ctx.Param("id")).Delete(map[string]any{})
				groups := utils.PrepareMtoM(ctrl.SingularName+"_id", ctx.Param("id"), ctrl.SingularName+"_group_id", group)

				if err := tx.Table(ctrl.PluralName + "_groups").Create(&groups).Error; err != nil {
					return err
				}

				transformer["groups"] = groups
			}
		}

		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "update "+ctrl.SingularLabel+" success", transformer))
}

// todo : need to check constraint error
func (ctrl CatalogController) Delete(ctx *gin.Context) {
	ctrl.Init(ctx)

	if err := utils.DB.Table(ctrl.PluralName).Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+ctrl.SingularLabel+" success", nil))
}
