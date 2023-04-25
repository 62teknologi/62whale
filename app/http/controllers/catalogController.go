package controllers

import (
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"

	"gorm.io/gorm"
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

func (ctrl *CatalogController) Find(ctx *gin.Context) {
	ctrl.Init(ctx)

	value := map[string]any{}
	columns := []string{ctrl.PluralName + ".*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.PluralName + "/find.json")
	query := utils.DB.Table(ctrl.PluralName)

	utils.SetJoin(query, transformer, &columns)

	if err := query.Select(columns).Where(ctrl.PluralName+".id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.SingularLabel+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachJoin(transformer, value)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *CatalogController) FindAll(ctx *gin.Context) {
	ctrl.Init(ctx)

	values := []map[string]any{}
	columns := []string{ctrl.PluralName + ".*"}
	order := "id desc"
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.PluralName + "/find.json")
	query := utils.DB.Table(ctrl.Table)
	filter := utils.SetFilterByQuery(query, transformer, ctx)
	pagination := utils.SetPagination(query, ctx)
	utils.SetJoin(query, transformer, &columns)

	if err := query.Select(columns).Order(order).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.PluralLabel+" not found", nil))
		return
	}

	customResponses := utils.MultiMapValuesShifter(values, transformer)

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+ctrl.PluralLabel+" success", customResponses, pagination, filter))
}

func (ctrl *CatalogController) Create(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + ctrl.PluralName + "/create.json")
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

	item, item_exist := input["items"]
	group, groups_exist := input["groups"]

	delete(transformer, "items")
	delete(transformer, "groups")

	var name string

	if transformer["name"] != nil {
		name, _ = transformer["name"].(string)
		transformer["slug"] = slug.Make(name)
	} else {
		transformer["slug"] = uuid.New()
	}

	if err := utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(ctrl.PluralName).Create(&transformer).Error; err != nil {
			return err
		}

		if item_exist || groups_exist {
			tx.Table(ctrl.PluralName).Where("slug = ?", transformer["slug"]).Take(&transformer)

			if item_exist {
				items := utils.Prepare1toM(ctrl.SingularName+"_id", transformer["id"], item)

				if err := tx.Table(ctrl.SingularName + "_items").Create(&items).Error; err != nil {
					return err
				}

				transformer["items"] = items
			}

			if groups_exist {
				groups := utils.PrepareMtoM(ctrl.SingularName+"_id", transformer["id"], ctrl.SingularName+"_group_id", group)

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

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *CatalogController) Update(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + ctrl.PluralName + "/update.json")
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

	item, item_exist := transformer["items"]
	group, groups_exist := transformer["groups"]

	delete(transformer, "items")
	delete(transformer, "groups")

	var name string

	if transformer["name"] != nil {
		name, _ = transformer["name"].(string)
		// not sure is it needed or not, may confusing if slug changes
		transformer["slug"] = slug.Make(name)
	}

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
func (ctrl *CatalogController) Delete(ctx *gin.Context) {
	if err := utils.DB.Table(ctrl.PluralName).Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+ctrl.SingularLabel+" success", nil))
}
