package relay

import (
	"fmt"
	"net/http"
	"one-api/common"
	"one-api/common/utils"
	"one-api/model"
	"one-api/types"
	"sort"
  "strings"

	"github.com/gin-gonic/gin"
)

// https://platform.openai.com/docs/api-reference/models/list
type OpenAIModels struct {
	Id      string  `json:"id"`
	Object  string  `json:"object"`
	Created int     `json:"created"`
	OwnedBy *string `json:"owned_by"`
}
// https://ai.google.dev/api/models#rest-resource:-models
type GeminiModels struct {
  Name      string  `json:"name"`
  BaseModelId string `json:"baseModelId"`
  Version   string  `json:"version"`
  DisplayName string `json:"displayName"`
  Description string `json:"description"`
  InputTokenLimit int `json:"inputTokenLimit"`
  OutputTokenLimit int `json:"outputTokenLimit"`
  SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
  Temperature float64 `json:"temperature"`
  MaxTemperature float64 `json:"maxTemperature"`
  TopP float64 `json:"topP"`
  TopK int `json:"topK"`
}

// func ListModelsByToken(c *gin.Context) {
// 	groupName := c.GetString("token_group")
// 	if groupName == "" {
// 		groupName = c.GetString("group")
// 	}

// 	if groupName == "" {
// 		common.AbortWithMessage(c, http.StatusServiceUnavailable, "分组不存在")
// 		return
// 	}

// 	models, err := model.ChannelGroup.GetGroupModels(groupName)
// 	if err != nil {
// 		c.JSON(200, gin.H{
// 			"object": "list",
// 			"data":   []string{},
// 		})
// 		return
// 	}
// 	sort.Strings(models)

// 	var groupOpenAIModels []*OpenAIModels
// 	for _, modelName := range models {
// 		groupOpenAIModels = append(groupOpenAIModels, getOpenAIModelWithName(modelName))
// 	}

// 	// 根据 OwnedBy 排序
// 	sort.Slice(groupOpenAIModels, func(i, j int) bool {
// 		if groupOpenAIModels[i].OwnedBy == nil {
// 			return true // 假设 nil 值小于任何非 nil 值
// 		}
// 		if groupOpenAIModels[j].OwnedBy == nil {
// 			return false // 假设任何非 nil 值大于 nil 值
// 		}
// 		return *groupOpenAIModels[i].OwnedBy < *groupOpenAIModels[j].OwnedBy
// 	})

// 	c.JSON(200, gin.H{
// 		"object": "list",
// 		"data":   groupOpenAIModels,
// 	})
// }

func ListModelsByToken(c *gin.Context) {
  // Get Token for aquiring models using Gemini endpoint method
  // token := c.GetString("key")
  // if token == "" {
  //   if c.Contains("/gemini/") {
  //     common.AbortWithMessage(c, http.StatusServiceUnavailable, "令牌不存在")
  //     return
  //   }
  //   return
  // }
  if c.Contains("/gemini/") {
    token := c.Param("key")
    if token == "" {
      common.AbortWithMessage(c, http.StatusServiceUnavailable, "令牌不存在")
      return
    }

    // Get group name from token
    groupName := c.GetString("token_group")
    if groupName == "" {
      groupName = c.GetString("group")
    }

    if groupName == "" {
      common.AbortWithMessage(c, http.StatusServiceUnavailable, "分组不存在")
      return
    }

    models, err := model.ChannelGroup.GetGroupModels(groupName)
    if err != nil {
      c.JSON(200, gin.H{
        "object": "list",
        "data":   []string{},
      })
      return
    }
    sort.Strings(models)

    var groupGeminiModels []*GeminiModels
    for _, modelName := range models {
      groupGeminiModels = append(groupGeminiModels, getGeminiModelWithName(modelName))
    }
    c.JSON(200, gin.H{
      "object": "list",
      "data":   groupGeminiModels,
    })


  // Get group name from token
  groupName := c.GetString("token_group")
  if groupName == "" {
    groupName = c.GetString("group")
  }

  if groupName == "" {
    common.AbortWithMessage(c, http.StatusServiceUnavailable, "分组不存在")
    return
  }

  models, err := model.ChannelGroup.GetGroupModels(groupName)
  if err != nil {
    c.JSON(200, gin.H{
      "object": "list",
      "data":   []string{},
    })
    return
  }
  sort.Strings(models)

  var groupOpenAIModels []*OpenAIModels
  for _, modelName := range models {
    groupOpenAIModels = append(groupOpenAIModels, getOpenAIModelWithName(modelName))
  }

  c.JSON(200, gin.H{
    "object": "list",
    "data":   groupOpenAIModels,
  })
}

func ListModelsForAdmin(c *gin.Context) {
	prices := model.PricingInstance.GetAllPrices()
	var openAIModels []OpenAIModels
	for modelId, price := range prices {
		openAIModels = append(openAIModels, OpenAIModels{
			Id:      modelId,
			Object:  "model",
			Created: 1677649963,
			OwnedBy: getModelOwnedBy(price.ChannelType),
		})
	}
	// 根据 OwnedBy 排序
	sort.Slice(openAIModels, func(i, j int) bool {
		if openAIModels[i].OwnedBy == nil {
			return true // 假设 nil 值小于任何非 nil 值
		}
		if openAIModels[j].OwnedBy == nil {
			return false // 假设任何非 nil 值大于 nil 值
		}
		return *openAIModels[i].OwnedBy < *openAIModels[j].OwnedBy
	})

	c.JSON(200, gin.H{
		"object": "list",
		"data":   openAIModels,
	})
}

// func RetrieveModel(c *gin.Context) {
//   // Get Token for aquiring models using Gemini endpoint method
//   if c.Contains("/gemini/") {
//     token := c.Param("key")
//     if token == "" {
//       common.AbortWithMessage(c, http.StatusServiceUnavailable, "令牌不存在")
//       return
//     }

//     // curl https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash?key=$GEMINI_API_KEY

//     modelName := c.
//     groupName := c.GetString("token_group")
//     if groupName == "" {
//       groupName = c.GetString("group")
//     }

//   }
// 	modelName := c.Param("model")
// 	openaiModel := getOpenAIModelWithName(modelName)
// 	if *openaiModel.OwnedBy != model.UnknownOwnedBy {
// 		c.JSON(200, openaiModel)
// 	} else {
// 		openAIError := types.OpenAIError{
// 			Message: fmt.Sprintf("The model '%s' does not exist", modelName),
// 			Type:    "invalid_request_error",
// 			Param:   "model",
// 			Code:    "model_not_found",
// 		}
// 		c.JSON(200, gin.H{
// 			"error": openAIError,
// 		})
// 	}
// }
func RetrieveModel(c *gin.Context) {
  // Check if this is a Gemini model request
  if strings.Contains(c.Request.URL.Path, "/gemini/") {
    apiKey := c.Query("key")
    if apiKey == "" {
      common.AbortWithMessage(c, http.StatusUnauthorized, "API key is required")
      return
    }

    modelName := c.Param("model")
    if modelName == "" {
      common.AbortWithMessage(c, http.StatusBadRequest, "Model name is required")
      return
    }

    // Get model details
    geminiModel := getGeminiModelWithName(modelName)
    c.JSON(200, geminiModel)
    return
  }

  // Standard OpenAI model handling
  modelName := c.Param("model")
  openaiModel := getOpenAIModelWithName(modelName)
  if *openaiModel.OwnedBy != model.UnknownOwnedBy {
    c.JSON(200, openaiModel)
  } else {
    openAIError := types.OpenAIError{
      Message: fmt.Sprintf("The model '%s' does not exist", modelName),
      Type:    "invalid_request_error",
      Param:   "model",
      Code:    "model_not_found",
    }
    c.JSON(200, gin.H{
      "error": openAIError,
    })
  }
}
func getModelOwnedBy(channelType int) (ownedBy *string) {
	ownedByName := model.ModelOwnedBysInstance.GetName(channelType)
	if ownedByName != "" {
		return &ownedByName
	}

	return &model.UnknownOwnedBy
}

func getOpenAIModelWithName(modelName string) *OpenAIModels {
	price := model.PricingInstance.GetPrice(modelName)

	return &OpenAIModels{
		Id:      modelName,
		Object:  "model",
		Created: 1677649963,
		OwnedBy: getModelOwnedBy(price.ChannelType),
	}
}

func getGeminiModelWithName(modelName string) *GeminiModels {
  price := model.PricingInstance.GetPrice(modelName)

  return &GeminiModels{
    Name:         modelName,
    BaseModelId:  BaseModelId || ,
    Version:      Version || "none",
    DisplayName:  DisplayName || "none",
    Description:  Description || "none",
    InputTokenLimit:  InputTokenLimit || ,

func GetModelOwnedBy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    model.ModelOwnedBysInstance.GetAll(),
	})
}

type ModelPrice struct {
	Type   string  `json:"type"`
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

type AvailableModelResponse struct {
	Groups  []string     `json:"groups"`
	OwnedBy string       `json:"owned_by"`
	Price   *model.Price `json:"price"`
}

func AvailableModel(c *gin.Context) {
	groupName := c.GetString("group")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    getAvailableModels(groupName),
	})
}

func GetAvailableModels(groupName string) map[string]*AvailableModelResponse {
	return getAvailableModels(groupName)
}

func getAvailableModels(groupName string) map[string]*AvailableModelResponse {
	publicModels := model.ChannelGroup.GetModelsGroups()
	publicGroups := model.GlobalUserGroupRatio.GetPublicGroupList()
	if groupName != "" && !utils.Contains(groupName, publicGroups) {
		publicGroups = append(publicGroups, groupName)
	}

	availableModels := make(map[string]*AvailableModelResponse, len(publicModels))

	for modelName, group := range publicModels {
		groups := []string{}
		for _, publicGroup := range publicGroups {
			if group[publicGroup] {
				groups = append(groups, publicGroup)
			}
		}

		if len(groups) == 0 {
			continue
		}

		if _, ok := availableModels[modelName]; !ok {
			price := model.PricingInstance.GetPrice(modelName)
			availableModels[modelName] = &AvailableModelResponse{
				Groups:  groups,
				OwnedBy: *getModelOwnedBy(price.ChannelType),
				Price:   price,
			}
		}
	}

	return availableModels
}
