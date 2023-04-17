package controllers

import (
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"

	"gorm.io/gorm"
)

type CatalogController struct{}

func (ctrl CatalogController) Find(ctx *gin.Context) {
	var value map[string]interface{}

	if err := utils.DB.Table(utils.PluralName).Where("id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if value["id"] == nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", utils.SingularName+" not found", nil))
		return
	}

	transformer, _ := utils.JsonFileParser("transformers/response/" + utils.PluralName + "/find.json")
	utils.MapValuesShifter(transformer, value)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+utils.SingularName+" success", transformer))
}

func (ctrl CatalogController) FindAll(ctx *gin.Context) {
	var values []map[string]interface{}

	if err := utils.DB.Table(utils.PluralName).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	transformer, _ := utils.JsonFileParser("transformers/response/" + utils.PluralName + "/find.json")
	customResponses := utils.MultiMapValuesShifter(values, transformer)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+utils.PluralName+" success", customResponses))
}

func (ctrl CatalogController) Create(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("transformers/request/" + utils.PluralName + "/create.json")
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

	name, _ := transformer["name"].(string)
	transformer["slug"] = slug.Make(name)

	if err := utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.PluralName).Create(&transformer).Error; err != nil {
			return err
		}

		if item_exist || groups_exist {
			tx.Table(utils.PluralName).Where("slug = ?", transformer["slug"]).Take(&transformer)

			if item_exist {
				items := utils.Prepare1toM(utils.SingularName+"_id", transformer["id"], item)

				if err := tx.Table(utils.SingularName + "_items").Create(&items).Error; err != nil {
					return err
				}
			}

			if groups_exist {
				groups := utils.PrepareMtoM(utils.SingularName+"_id", transformer["id"], utils.SingularName+"_group_id", group)

				if err := tx.Table(utils.PluralName + "_groups").Create(&groups).Error; err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+utils.SingularName+" success", transformer))
}

func (ctrl CatalogController) Update(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("transformers/request/" + utils.PluralName + "/update.json")
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

	name, _ := transformer["name"].(string)
	transformer["slug"] = slug.Make(name)

	if err := utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.PluralName).Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
			return err
		}

		if item_exist || groups_exist {
			if item_exist {
				items := utils.Prepare1toM(utils.SingularName+"_id", ctx.Param("id"), item)

				if err := tx.Table(utils.SingularName + "_items").Create(&items).Error; err != nil {
					return err
				}
			}

			if groups_exist {
				tx.Table(utils.PluralName+"_groups").Where(utils.SingularName+"_id = ?", ctx.Param("id")).Delete(map[string]any{})
				groups := utils.PrepareMtoM(utils.SingularName+"_id", ctx.Param("id"), utils.SingularName+"_group_id", group)

				if err := tx.Table(utils.PluralName + "_groups").Create(&groups).Error; err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	// todo : make a better response!
	ctrl.Find(ctx)
}

// todo : need to check constraint error
func (ctrl CatalogController) Delete(ctx *gin.Context) {
	if err := utils.DB.Table(utils.PluralName).Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+utils.SingularName+" success", nil))
}
