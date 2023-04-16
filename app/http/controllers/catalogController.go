package controllers

import (
	"fmt"
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"

	"gorm.io/gorm"
)

func FindCatalog(ctx *gin.Context) {
	var value map[string]interface{}
	err := utils.DB.Table(utils.PluralName).Where("id = ?", ctx.Param("id")).Take(&value).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if value["id"] == nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", utils.SingularName+" not found", nil))
		return
	}

	transformer, _ := utils.JsonFileParser("transformers/response/" + utils.SingularName + "/find.json")
	customResponse := transformer["catalog"]

	utils.MapValuesShifter(transformer, value)

	if customResponse != nil {
		utils.MapValuesShifter(customResponse.(map[string]any), value)
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+utils.SingularName+" success", transformer))
}

func FindCatalogues(ctx *gin.Context) {
	var values []map[string]interface{}
	err := utils.DB.Table(utils.PluralName).Find(&values).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	var customResponses []map[string]any
	for _, value := range values {
		transformer, _ := utils.JsonFileParser("transformers/response/" + utils.SingularName + "/find.json")
		customResponse := transformer["catalog"]

		utils.MapValuesShifter(transformer, value)
		if customResponse != nil {
			utils.MapValuesShifter(customResponse.(map[string]any), value)
		}
		customResponses = append(customResponses, transformer)
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+utils.PluralName+" success", customResponses))
}

func CreateCatalog(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("transformers/request/" + utils.SingularName + "/create.json")
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

func UpdateCatalog(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("transformers/request/" + utils.SingularName + "/update.json")
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

	fmt.Println(transformer)

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
	FindCatalog(ctx)
}

func DeleteCatalog(ctx *gin.Context) {
	if err := utils.DB.Table(utils.PluralName).Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+utils.SingularName+" success", nil))
}
