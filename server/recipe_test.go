package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"gitlab.kenda.com.tw/kenda/mcom"
	mcomErr "gitlab.kenda.com.tw/kenda/mcom/errors"
	"gitlab.kenda.com.tw/kenda/mcom/impl/orm/models"
	"gitlab.kenda.com.tw/kenda/mcom/mock"
	"gitlab.kenda.com.tw/kenda/mcom/utils/types"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
	pbTypes "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/types"
)

const (
	testRecipeID           = "recipe_a"
	testProductID          = "product_a"
	testProductType        = "type_c"
	testProcessAOid        = "100fc317-1dd2-11b2-8001-080027b246c3"
	testProcessAName       = "process_a"
	testProcessBOid        = "010fc317-1dd2-11b2-8001-080027b246c3"
	testProcessBName       = "process_b"
	testProcessCOid        = "001fc317-1dd2-11b2-8001-080027b246c3"
	testProcessCName       = "process_c"
	testStation            = "GG"
	testTool               = "tool_a"
	testToolType           = "tool_type"
	testMaterial           = "mtrl_a"
	testMaterialType       = "mtrl_type"
	testSite               = "site_a"
	testControlName        = "ctrl_a"
	testMeasurementsName   = "meas_a"
	testCommonsControlName = "com_ctrl_a"
	testSubstitutionName   = "000001"
	testSubstitutionType   = "B"
	testProportion         = "123.2"
	testBatchSize          = "123.456"
)

var (
	proportion = decimal.RequireFromString(testProportion)
	batchSize  = decimal.RequireFromString(testBatchSize)
)

func Test_CreateRecipes(t *testing.T) {
	assert := assert.New(t)
	ctx := contextWithFactoryIDs(context.Background(), factoryIDs...)

	{ // CreateRecipes: insufficient request
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateRecipes,
			Input: mock.Input{
				Request: mcom.CreateRecipesRequest{
					Recipes: []mcom.Recipe{},
				},
			},
			Output: mock.Output{
				Response: nil,
				Error: mcomErr.Error{
					Code: mcomErr.Code_INSUFFICIENT_REQUEST,
				},
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.CreateRecipes(ctx, &mesync.CreateRecipesRequest{})
		assert.Error(err)
		assert.Equal(parseError(testFactoryID, mcomErr.Error{
			Code: mcomErr.Code_INSUFFICIENT_REQUEST,
		}), err)
	}
	{ // CreateRecipes: good case
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateRecipes,
			Input: mock.Input{
				Request: mcom.CreateRecipesRequest{
					Recipes: []mcom.Recipe{{
						ID: testRecipeID,
						Product: mcom.Product{
							ID:   testProductID,
							Type: testProductType,
						},
						Version: mcom.RecipeVersion{
							Major:      "1",
							Minor:      "2",
							Stage:      "3",
							ReleasedAt: types.TimeNano(123),
						},
						Processes: []*mcom.Process{{
							OID: testProcessAOid,
							OptionalFlows: []*mcom.RecipeOptionalFlow{{
								Name:           testProcessAName,
								OIDs:           []string{testProcessBOid, testProcessCOid},
								MaxRepetitions: 0,
							}},
						}},
						ProcessDefinitions: []mcom.ProcessDefinition{{
							OID:  testProcessAOid,
							Name: testProcessAName,
							Type: testProductType,
							Configs: []*mcom.RecipeProcessConfig{{
								Stations:  []string{testStation},
								BatchSize: &batchSize,
								Unit:      "kg",
								Tools: []*mcom.RecipeTool{{
									Type:     testToolType,
									ID:       testTool,
									Required: false,
								}},
								Steps: []*mcom.RecipeProcessStep{{
									Materials: []*mcom.RecipeMaterial{{
										Name:  testMaterial,
										Grade: "A",
										Value: mcom.RecipeMaterialParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
										Site:             testSite,
										RequiredRecipeID: "",
									}},
									Controls: []*mcom.RecipeProperty{{
										Name: testControlName,
										Param: &mcom.RecipePropertyParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
									}},
									Measurements: []*mcom.RecipeProperty{{
										Name: testMeasurementsName,
										Param: &mcom.RecipePropertyParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
									}},
								}},
								CommonControls: []*mcom.RecipeProperty{{
									Name: testCommonsControlName,
									Param: &mcom.RecipePropertyParameter{
										High: types.Decimal.RequireFromString("100"),
										Mid:  types.Decimal.RequireFromString("50"),
										Low:  types.Decimal.RequireFromString("0"),
										Unit: "kg",
									},
								}},
								CommonProperties: []*mcom.RecipeProperty{{
									Name: testCommonsControlName,
									Param: &mcom.RecipePropertyParameter{
										High: types.Decimal.RequireFromString("100"),
										Mid:  types.Decimal.RequireFromString("50"),
										Low:  types.Decimal.RequireFromString("0"),
										Unit: "kg",
									},
								}},
							}},
							Output: mcom.OutputProduct{
								ID:   testProductID,
								Type: testProductType,
							},
							ProductValidPeriod: mcom.ProductValidPeriodConfig{
								Standing: 0,
								Expiry:   168,
							},
						}, {
							OID:  testProcessBOid,
							Name: testProcessBName,
							Type: testProductType,
							Configs: []*mcom.RecipeProcessConfig{{
								Stations:  []string{testStation},
								BatchSize: &batchSize,
								Unit:      "kg",
								Tools: []*mcom.RecipeTool{{
									Type:     testToolType,
									ID:       testTool,
									Required: false,
								}},
								Steps: []*mcom.RecipeProcessStep{{
									Materials: []*mcom.RecipeMaterial{{
										Name:  testMaterial,
										Grade: "A",
										Value: mcom.RecipeMaterialParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
										Site:             testSite,
										RequiredRecipeID: "",
									}},
									Controls: []*mcom.RecipeProperty{{
										Name: testControlName,
										Param: &mcom.RecipePropertyParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
									}},
									Measurements: []*mcom.RecipeProperty{{
										Name: testMeasurementsName,
										Param: &mcom.RecipePropertyParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
									}},
								}},
								CommonControls: []*mcom.RecipeProperty{{
									Name: testCommonsControlName,
									Param: &mcom.RecipePropertyParameter{
										High: types.Decimal.RequireFromString("100"),
										Mid:  types.Decimal.RequireFromString("50"),
										Low:  types.Decimal.RequireFromString("0"),
										Unit: "kg",
									},
								}},
								CommonProperties: []*mcom.RecipeProperty{{
									Name: testCommonsControlName,
									Param: &mcom.RecipePropertyParameter{
										High: types.Decimal.RequireFromString("100"),
										Mid:  types.Decimal.RequireFromString("50"),
										Low:  types.Decimal.RequireFromString("0"),
										Unit: "kg",
									},
								}},
							}},
							Output: mcom.OutputProduct{
								ID:   testProductID,
								Type: testProductType,
							},
							ProductValidPeriod: mcom.ProductValidPeriodConfig{
								Standing: 0,
								Expiry:   168,
							},
						}, {
							OID:  testProcessCOid,
							Name: testProcessCName,
							Type: testProductType,
							Configs: []*mcom.RecipeProcessConfig{{
								Stations:  []string{testStation},
								BatchSize: &batchSize,
								Unit:      "kg",
								Tools: []*mcom.RecipeTool{{
									Type:     testToolType,
									ID:       testTool,
									Required: false,
								}},
								Steps: []*mcom.RecipeProcessStep{{
									Materials: []*mcom.RecipeMaterial{{
										Name:  testMaterial,
										Grade: "A",
										Value: mcom.RecipeMaterialParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
										Site:             testSite,
										RequiredRecipeID: "",
									}},
									Controls: []*mcom.RecipeProperty{{
										Name: testControlName,
										Param: &mcom.RecipePropertyParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
									}},
									Measurements: []*mcom.RecipeProperty{{
										Name: testMeasurementsName,
										Param: &mcom.RecipePropertyParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
									}},
								}},
								CommonControls: []*mcom.RecipeProperty{{
									Name: testCommonsControlName,
									Param: &mcom.RecipePropertyParameter{
										High: types.Decimal.RequireFromString("100"),
										Mid:  types.Decimal.RequireFromString("50"),
										Low:  types.Decimal.RequireFromString("0"),
										Unit: "kg",
									},
								}},
								CommonProperties: []*mcom.RecipeProperty{{
									Name: testCommonsControlName,
									Param: &mcom.RecipePropertyParameter{
										High: types.Decimal.RequireFromString("100"),
										Mid:  types.Decimal.RequireFromString("50"),
										Low:  types.Decimal.RequireFromString("0"),
										Unit: "kg",
									},
								}},
							}},
							Output: mcom.OutputProduct{
								ID:   testProductID,
								Type: testProductType,
							},
							ProductValidPeriod: mcom.ProductValidPeriodConfig{
								Standing: 0,
								Expiry:   168,
							},
						}},
					}},
				},
			},
			Output: mock.Output{
				Response: nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		t, err := ptypes.TimestampProto(time.Unix(0, 123))
		assert.NoError(err)

		_, err = mockServer.CreateRecipes(ctx, &mesync.CreateRecipesRequest{
			Recipes: []*mesync.Recipe{{
				Id:          testRecipeID,
				ProductType: testProductType,
				ProductId:   testProductID,
				Version: &mesync.RecipeVersion{
					Major:      "1",
					Minor:      "2",
					Stage:      "3",
					ReleasedAt: t,
				},
				Processes: []*mesync.Process{{
					ReferenceOid: testProcessAOid,
					OptionalFlows: []*mesync.OptionalFlow{{
						Name:           testProcessAName,
						ProcessOids:    []string{testProcessBOid, testProcessCOid},
						MaxRepetitions: 0,
					}},
				}},
				ProcessDefs: []*mesync.ProcessDefinition{{
					Oid:  testProcessAOid,
					Name: testProcessAName,
					Type: testProductType,
					Configs: []*mesync.RecipeProcessConfig{{
						Stations:  []string{testStation},
						BatchSize: &pbTypes.Decimal{Value: testBatchSize},
						Unit:      "kg",
						Tools: []*mesync.RecipeTool{{
							Type:     testToolType,
							Id:       testTool,
							Required: false,
						}},
						Steps: []*mesync.RecipeProcessStep{{
							Materials: []*mesync.RecipeMaterial{{
								Name:  testMaterial,
								Type:  testMaterialType,
								Grade: "A",
								Value: &mesync.RecipeMaterialParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
								Site:             testSite,
								RequiredRecipeId: "",
							}},
							Controls: []*mesync.RecipeProperty{{
								Name: testControlName,
								Param: &mesync.RecipePropertyParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
							}},
							Measurements: []*mesync.RecipeProperty{{
								Name: testMeasurementsName,
								Param: &mesync.RecipePropertyParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
							}},
						}},
						CommonsControls: []*mesync.RecipeProperty{{
							Name: testCommonsControlName,
							Param: &mesync.RecipePropertyParameter{
								High: &pbTypes.Decimal{Value: "100"},
								Mid:  &pbTypes.Decimal{Value: "50"},
								Low:  &pbTypes.Decimal{Value: "0"},
								Unit: "kg",
							},
						}},
						CommonsProperties: []*mesync.RecipeProperty{{
							Name: testCommonsControlName,
							Param: &mesync.RecipePropertyParameter{
								High: &pbTypes.Decimal{Value: "100"},
								Mid:  &pbTypes.Decimal{Value: "50"},
								Low:  &pbTypes.Decimal{Value: "0"},
								Unit: "kg",
							},
						}},
					}},
					Product: &mesync.OutputProduct{
						Id:   testProductID,
						Type: testProductType,
					},
					ProductValidPeriod: &mesync.ProductValidPeriodConfig{
						Standing: 0,
						Expiry:   168,
					},
				}, {
					Oid:  testProcessBOid,
					Name: testProcessBName,
					Type: testProductType,
					Configs: []*mesync.RecipeProcessConfig{{
						Stations:  []string{testStation},
						BatchSize: &pbTypes.Decimal{Value: testBatchSize},
						Unit:      "kg",
						Tools: []*mesync.RecipeTool{{
							Type:     testToolType,
							Id:       testTool,
							Required: false,
						}},
						Steps: []*mesync.RecipeProcessStep{{
							Materials: []*mesync.RecipeMaterial{{
								Name:  testMaterial,
								Type:  testMaterialType,
								Grade: "A",
								Value: &mesync.RecipeMaterialParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
								Site:             testSite,
								RequiredRecipeId: "",
							}},
							Controls: []*mesync.RecipeProperty{{
								Name: testControlName,
								Param: &mesync.RecipePropertyParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
							}},
							Measurements: []*mesync.RecipeProperty{{
								Name: testMeasurementsName,
								Param: &mesync.RecipePropertyParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
							}},
						}},
						CommonsControls: []*mesync.RecipeProperty{{
							Name: testCommonsControlName,
							Param: &mesync.RecipePropertyParameter{
								High: &pbTypes.Decimal{Value: "100"},
								Mid:  &pbTypes.Decimal{Value: "50"},
								Low:  &pbTypes.Decimal{Value: "0"},
								Unit: "kg",
							},
						}},
						CommonsProperties: []*mesync.RecipeProperty{{
							Name: testCommonsControlName,
							Param: &mesync.RecipePropertyParameter{
								High: &pbTypes.Decimal{Value: "100"},
								Mid:  &pbTypes.Decimal{Value: "50"},
								Low:  &pbTypes.Decimal{Value: "0"},
								Unit: "kg",
							},
						}},
					}},
					Product: &mesync.OutputProduct{
						Id:   testProductID,
						Type: testProductType,
					},
					ProductValidPeriod: &mesync.ProductValidPeriodConfig{
						Standing: 0,
						Expiry:   168,
					},
				}, {
					Oid:  testProcessCOid,
					Name: testProcessCName,
					Type: testProductType,
					Configs: []*mesync.RecipeProcessConfig{{
						Stations:  []string{testStation},
						BatchSize: &pbTypes.Decimal{Value: testBatchSize},
						Unit:      "kg",
						Tools: []*mesync.RecipeTool{{
							Type:     testToolType,
							Id:       testTool,
							Required: false,
						}},
						Steps: []*mesync.RecipeProcessStep{{
							Materials: []*mesync.RecipeMaterial{{
								Name:  testMaterial,
								Type:  testMaterialType,
								Grade: "A",
								Value: &mesync.RecipeMaterialParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
								Site:             testSite,
								RequiredRecipeId: "",
							}},
							Controls: []*mesync.RecipeProperty{{
								Name: testControlName,
								Param: &mesync.RecipePropertyParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
							}},
							Measurements: []*mesync.RecipeProperty{{
								Name: testMeasurementsName,
								Param: &mesync.RecipePropertyParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
							}},
						}},
						CommonsControls: []*mesync.RecipeProperty{{
							Name: testCommonsControlName,
							Param: &mesync.RecipePropertyParameter{
								High: &pbTypes.Decimal{Value: "100"},
								Mid:  &pbTypes.Decimal{Value: "50"},
								Low:  &pbTypes.Decimal{Value: "0"},
								Unit: "kg",
							},
						}},
						CommonsProperties: []*mesync.RecipeProperty{{
							Name: testCommonsControlName,
							Param: &mesync.RecipePropertyParameter{
								High: &pbTypes.Decimal{Value: "100"},
								Mid:  &pbTypes.Decimal{Value: "50"},
								Low:  &pbTypes.Decimal{Value: "0"},
								Unit: "kg",
							},
						}},
					}},
					Product: &mesync.OutputProduct{
						Id:   testProductID,
						Type: testProductType,
					},
					ProductValidPeriod: &mesync.ProductValidPeriodConfig{
						Standing: 0,
						Expiry:   168,
					},
				}},
			}},
		})
		assert.NoError(err)
	}
	{ // good case: receive a nil releasedAt
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncCreateRecipes,
			Input: mock.Input{
				Request: mcom.CreateRecipesRequest{
					Recipes: []mcom.Recipe{{
						ID: testRecipeID,
						Product: mcom.Product{
							ID:   testProductID,
							Type: testProductType,
						},
						Version: mcom.RecipeVersion{
							Major:      "1",
							Minor:      "",
							Stage:      "3",
							ReleasedAt: 0,
						},
						Processes: []*mcom.Process{{
							OID:           testProcessAOid,
							OptionalFlows: []*mcom.RecipeOptionalFlow{},
						}},
						ProcessDefinitions: []mcom.ProcessDefinition{{
							OID:  testProcessAOid,
							Name: testProcessAName,
							Type: testProductType,
							Configs: []*mcom.RecipeProcessConfig{{
								Stations:  []string{testStation},
								BatchSize: &batchSize,
								Unit:      "kg",
								Tools: []*mcom.RecipeTool{{
									Type:     testToolType,
									ID:       testTool,
									Required: false,
								}},
								Steps: []*mcom.RecipeProcessStep{{
									Materials: []*mcom.RecipeMaterial{{
										Name:  testMaterial,
										Grade: "A",
										Value: mcom.RecipeMaterialParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
										Site:             testSite,
										RequiredRecipeID: "",
									}},
									Controls: []*mcom.RecipeProperty{{
										Name: testControlName,
										Param: &mcom.RecipePropertyParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
									}},
									Measurements: []*mcom.RecipeProperty{{
										Name: testMeasurementsName,
										Param: &mcom.RecipePropertyParameter{
											High: types.Decimal.RequireFromString("100"),
											Mid:  types.Decimal.RequireFromString("50"),
											Low:  types.Decimal.RequireFromString("0"),
											Unit: "kg",
										},
									}},
								}},
								CommonControls: []*mcom.RecipeProperty{{
									Name: testCommonsControlName,
									Param: &mcom.RecipePropertyParameter{
										High: types.Decimal.RequireFromString("100"),
										Mid:  types.Decimal.RequireFromString("50"),
										Low:  types.Decimal.RequireFromString("0"),
										Unit: "kg",
									},
								}},
								CommonProperties: []*mcom.RecipeProperty{{
									Name: testCommonsControlName,
									Param: &mcom.RecipePropertyParameter{
										High: types.Decimal.RequireFromString("100"),
										Mid:  types.Decimal.RequireFromString("50"),
										Low:  types.Decimal.RequireFromString("0"),
										Unit: "kg",
									},
								}},
							}},
							Output: mcom.OutputProduct{
								ID:   testProductID,
								Type: testProductType,
							},
							ProductValidPeriod: mcom.ProductValidPeriodConfig{
								Standing: 0,
								Expiry:   168,
							},
						}},
					}},
				},
			},
			Output: mock.Output{
				Response: nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.CreateRecipes(ctx, &mesync.CreateRecipesRequest{
			Recipes: []*mesync.Recipe{{
				Id:          testRecipeID,
				ProductType: testProductType,
				ProductId:   testProductID,
				Version: &mesync.RecipeVersion{
					Major: "1",
					Minor: "",
					Stage: "3",
				},
				Processes: []*mesync.Process{{
					ReferenceOid:  testProcessAOid,
					OptionalFlows: []*mesync.OptionalFlow{},
				}},
				ProcessDefs: []*mesync.ProcessDefinition{{
					Oid:  testProcessAOid,
					Name: testProcessAName,
					Type: testProductType,
					Configs: []*mesync.RecipeProcessConfig{{
						Stations:  []string{testStation},
						BatchSize: &pbTypes.Decimal{Value: testBatchSize},
						Unit:      "kg",
						Tools: []*mesync.RecipeTool{{
							Type:     testToolType,
							Id:       testTool,
							Required: false,
						}},
						Steps: []*mesync.RecipeProcessStep{{
							Materials: []*mesync.RecipeMaterial{{
								Name:  testMaterial,
								Type:  testMaterialType,
								Grade: "A",
								Value: &mesync.RecipeMaterialParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
								Site:             testSite,
								RequiredRecipeId: "",
							}},
							Controls: []*mesync.RecipeProperty{{
								Name: testControlName,
								Param: &mesync.RecipePropertyParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
							}},
							Measurements: []*mesync.RecipeProperty{{
								Name: testMeasurementsName,
								Param: &mesync.RecipePropertyParameter{
									High: &pbTypes.Decimal{Value: "100"},
									Mid:  &pbTypes.Decimal{Value: "50"},
									Low:  &pbTypes.Decimal{Value: "0"},
									Unit: "kg",
								},
							}},
						}},
						CommonsControls: []*mesync.RecipeProperty{{
							Name: testCommonsControlName,
							Param: &mesync.RecipePropertyParameter{
								High: &pbTypes.Decimal{Value: "100"},
								Mid:  &pbTypes.Decimal{Value: "50"},
								Low:  &pbTypes.Decimal{Value: "0"},
								Unit: "kg",
							},
						}},
						CommonsProperties: []*mesync.RecipeProperty{{
							Name: testCommonsControlName,
							Param: &mesync.RecipePropertyParameter{
								High: &pbTypes.Decimal{Value: "100"},
								Mid:  &pbTypes.Decimal{Value: "50"},
								Low:  &pbTypes.Decimal{Value: "0"},
								Unit: "kg",
							},
						}},
					}},
					Product: &mesync.OutputProduct{
						Id:   testProductID,
						Type: testProductType,
					},
					ProductValidPeriod: &mesync.ProductValidPeriodConfig{
						Standing: 0,
						Expiry:   168,
					},
				}},
			}},
		})
		assert.NoError(err)
	}
}

func Test_DeleteRecipe(t *testing.T) {
	assert := assert.New(t)
	ctx := contextWithFactoryIDs(context.Background(), factoryIDs...)

	{ // insufficient request
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncDeleteRecipe,
			Input: mock.Input{
				Request: mcom.DeleteRecipeRequest{
					IDs: []string{""},
				},
			},
			Output: mock.Output{
				Response: nil,
				Error: mcomErr.Error{
					Code: mcomErr.Code_INSUFFICIENT_REQUEST,
				},
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.DeleteRecipes(ctx, &mesync.DeleteRecipesRequest{
			Ids: []string{""},
		})
		assert.Error(err)
		assert.Equal(parseError(testFactoryID, mcomErr.Error{
			Code: mcomErr.Code_INSUFFICIENT_REQUEST,
		}), err)
	}
	{ // good case
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncDeleteRecipe,
			Input: mock.Input{
				Request: mcom.DeleteRecipeRequest{
					IDs: []string{testRecipeID},
				},
			},
			Output: mock.Output{
				Response: nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.DeleteRecipes(ctx, &mesync.DeleteRecipesRequest{
			Ids: []string{testRecipeID},
		})
		assert.NoError(err)
	}
}

func Test_CreateSubstitution(t *testing.T) {
	assert := assert.New(t)
	ctx := contextWithFactoryIDs(context.Background(), factoryIDs...)
	{ // insufficient request case
		mockServer, err := newMockServer([]mock.Script{
			{
				Name: mock.FuncAddSubstitutions,
				Input: mock.Input{
					Request: mcom.BasicSubstitutionRequest{
						ProductID: models.ProductID{
							ID:    testMaterial,
							Grade: testMaterialType,
						},
						Contents: []models.Substitution{
							{
								ID: testSubstitutionName,
								//no Grade
								Proportion: proportion,
							},
						},
					},
				},
				Output: mock.Output{
					Response: nil,
					Error: mcomErr.Error{
						Code:    mcomErr.Code_INSUFFICIENT_REQUEST,
						Details: "Grade is not found.",
					},
				},
			},
		})
		assert.NoError(err)

		_, err = mockServer.CreateSubstitutes(ctx, &mesync.CreateSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType,
			},
			Substitutions: []*mesync.Substitutions{{
				Material: &mesync.Material{
					Id: testSubstitutionName,
				},
				Proportion: &pbTypes.Decimal{Value: testProportion},
			}},
		})
		if assert.Error(err) {
			assert.Equal(parseError(testFactoryID, mcomErr.Error{
				Code:    mcomErr.Code_INSUFFICIENT_REQUEST,
				Details: "Grade is not found.",
			}), err)
		}
	}
	{ // invalid number case
		mockServer, err := newMockServer([]mock.Script{})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.CreateSubstitutes(ctx, &mesync.CreateSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType},
			Substitutions: []*mesync.Substitutions{{
				Material: &mesync.Material{
					Id:    testSubstitutionName,
					Grade: testMaterialType},
				Proportion: &pbTypes.Decimal{Value: "aaa"},
			},
			}})
		if assert.Error(err) {
			assert.Equal(mcomErr.Error{
				Code:    mcomErr.Code_INVALID_NUMBER,
				Details: "Name: 000001, Grade: mtrl_type, Proportion: aaa, can't convert aaa to decimal",
			}, err)
		}
	}
	{ // good case
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncAddSubstitutions,
			Input: mock.Input{
				Request: mcom.BasicSubstitutionRequest{
					ProductID: models.ProductID{
						ID:    testMaterial,
						Grade: testMaterialType,
					},
					Contents: []models.Substitution{
						{
							ID:         testSubstitutionName,
							Grade:      testSubstitutionType,
							Proportion: proportion,
						},
					},
				},
			},
			Output: mock.Output{
				Response: nil,
			},
		}})
		assert.NoError(err)

		_, err = mockServer.CreateSubstitutes(ctx, &mesync.CreateSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType,
			},
			Substitutions: []*mesync.Substitutions{{
				Material: &mesync.Material{
					Id:    testSubstitutionName,
					Grade: testSubstitutionType,
				},
				Proportion: &pbTypes.Decimal{Value: testProportion},
			}},
		})
		assert.NoError(err)
	}
	{ // internal error
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncAddSubstitutions,
			Input: mock.Input{
				Request: mcom.BasicSubstitutionRequest{
					ProductID: models.ProductID{
						ID:    testMaterial,
						Grade: testMaterialType,
					},
					Contents: []models.Substitution{
						{
							ID:         testSubstitutionName,
							Grade:      testSubstitutionType,
							Proportion: proportion,
						},
					},
				},
			},
			Output: mock.Output{
				Error: errors.New("internal error"),
			},
		}})
		assert.NoError(err)

		_, err = mockServer.CreateSubstitutes(ctx, &mesync.CreateSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType,
			},
			Substitutions: []*mesync.Substitutions{{
				Material: &mesync.Material{
					Id:    testSubstitutionName,
					Grade: testSubstitutionType,
				},
				Proportion: &pbTypes.Decimal{Value: testProportion},
			}},
		})
		if assert.Error(err) {
			assert.Equal(parseError(testFactoryID, errors.New("internal error")), err)
		}
	}
}

func Test_UpdateSubstitution(t *testing.T) {
	assert := assert.New(t)
	ctx := contextWithFactoryIDs(context.Background(), factoryIDs...)
	{ //Code_INSUFFICIENT_REQUEST
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncUpdateSubstitutions,
			Input: mock.Input{
				Request: mcom.BasicSubstitutionRequest{
					ProductID: models.ProductID{
						ID:    testMaterial,
						Grade: testMaterialType,
					},
					Contents: []models.Substitution{
						{
							ID: testSubstitutionName,
							//no Grade
							Proportion: proportion,
						},
					},
				},
			},
			Output: mock.Output{
				Response: nil,
				Error: mcomErr.Error{
					Code:    mcomErr.Code_INSUFFICIENT_REQUEST,
					Details: "Grade is not found.",
				},
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.UpdateSubstitutes(ctx, &mesync.UpdateSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType,
			},
			Substitutions: []*mesync.Substitutions{{
				Material: &mesync.Material{
					Id: testSubstitutionName,
				},
				Proportion: &pbTypes.Decimal{Value: testProportion},
			}},
		})
		if assert.Error(err) {
			assert.Equal(parseError(testFactoryID, mcomErr.Error{
				Code:    mcomErr.Code_INSUFFICIENT_REQUEST,
				Details: "Grade is not found.",
			}), err)
		}
	}
	{ //Code_INVALID_NUMBER
		mockServer, err := newMockServer([]mock.Script{{}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.UpdateSubstitutes(ctx, &mesync.UpdateSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType,
			},
			Substitutions: []*mesync.Substitutions{{
				Material: &mesync.Material{
					Id:    testSubstitutionName,
					Grade: testSubstitutionType,
				},
				Proportion: &pbTypes.Decimal{Value: "aaa"},
			}},
		})
		if assert.Error(err) {
			assert.Equal(mcomErr.Error{
				Code:    mcomErr.Code_INVALID_NUMBER,
				Details: "Name: 000001, Grade: B, Proportion: aaa, can't convert aaa to decimal",
			}, err)
		}
	}
	{ // good case
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncUpdateSubstitutions,
			Input: mock.Input{
				Request: mcom.BasicSubstitutionRequest{
					ProductID: models.ProductID{
						ID:    testMaterial,
						Grade: testMaterialType,
					},
					Contents: []models.Substitution{
						{
							ID:         testSubstitutionName,
							Grade:      testSubstitutionType,
							Proportion: proportion,
						},
					},
				},
			},
			Output: mock.Output{
				Response: nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.UpdateSubstitutes(ctx, &mesync.UpdateSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType,
			},
			Substitutions: []*mesync.Substitutions{{
				Material: &mesync.Material{
					Id:    testSubstitutionName,
					Grade: testSubstitutionType,
				},
				Proportion: &pbTypes.Decimal{Value: testProportion},
			}},
		})
		assert.NoError(err)
	}
	{ // internal error
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncUpdateSubstitutions,
			Input: mock.Input{
				Request: mcom.BasicSubstitutionRequest{
					ProductID: models.ProductID{
						ID:    testMaterial,
						Grade: testMaterialType,
					},
					Contents: []models.Substitution{
						{
							ID:         testSubstitutionName,
							Grade:      testSubstitutionType,
							Proportion: proportion,
						},
					},
				},
			},
			Output: mock.Output{
				Error: errors.New("internal error"),
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.UpdateSubstitutes(ctx, &mesync.UpdateSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType,
			},
			Substitutions: []*mesync.Substitutions{{
				Material: &mesync.Material{
					Id:    testSubstitutionName,
					Grade: testSubstitutionType,
				},
				Proportion: &pbTypes.Decimal{Value: testProportion},
			}},
		})
		if assert.Error(err) {
			assert.Equal(parseError(testFactoryID, errors.New("internal error")), err)
		}
	}

}

func Test_DeleteSubstitution(t *testing.T) {
	assert := assert.New(t)
	ctx := contextWithFactoryIDs(context.Background(), factoryIDs...)
	{ //Code_INSUFFICIENT_REQUEST
		{
			mockServer, err := newMockServer([]mock.Script{{
				Name: mock.FuncDeleteSubstitutions,
				Input: mock.Input{
					Request: mcom.DeleteSubstitutionsRequest{
						ProductID: models.ProductID{
							ID:    testMaterial,
							Grade: testMaterialType,
						},
						Contents: mcom.SubstitutionDeletionContents{
							mcom.SubstitutionDeletionContent{
								ProductID: models.ProductID{
									ID: testSubstitutionName,
								},
							},
						},
					},
				},
				Output: mock.Output{
					Response: nil,
					Error: mcomErr.Error{
						Code: mcomErr.Code_INSUFFICIENT_REQUEST,
					},
				},
			}})
			if !assert.NoError(err) {
				return
			}

			_, err = mockServer.DeleteSubstitution(ctx, &mesync.DeleteSubstitutionRequest{
				Material: &mesync.Material{
					Id:    testMaterial,
					Grade: testMaterialType},
				Type: &mesync.DeleteSubstitutionRequest_Multiple{Multiple: &mesync.DeleteSubstitutionRequest_Batch{
					Substitutions: []*mesync.Material{{
						Id: testSubstitutionName,
					},
					},
				}}},
			)
			if assert.Error(err) {
				assert.Equal(parseError(testFactoryID, mcomErr.Error{
					Code: mcomErr.Code_INSUFFICIENT_REQUEST,
				}), err)
			}
		}
	}
	{ // good case
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncDeleteSubstitutions,
			Input: mock.Input{
				Request: mcom.DeleteSubstitutionsRequest{
					ProductID: models.ProductID{
						ID:    testMaterial,
						Grade: testMaterialType,
					},
					Contents:  []mcom.SubstitutionDeletionContent{},
					DeleteAll: true,
				},
			},
			Output: mock.Output{
				Response: nil,
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.DeleteSubstitution(ctx, &mesync.DeleteSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType,
			},
			Type: &mesync.DeleteSubstitutionRequest_All{All: true},
		})
		assert.NoError(err)
	}
	{ // internal error
		mockServer, err := newMockServer([]mock.Script{{
			Name: mock.FuncDeleteSubstitutions,
			Input: mock.Input{
				Request: mcom.DeleteSubstitutionsRequest{
					ProductID: models.ProductID{
						ID:    testMaterial,
						Grade: testMaterialType,
					},
					Contents: mcom.SubstitutionDeletionContents{
						mcom.SubstitutionDeletionContent{
							ProductID: models.ProductID{
								ID:    testSubstitutionName,
								Grade: testSubstitutionType,
							},
						},
					},
					DeleteAll: false,
				},
			},
			Output: mock.Output{
				Error: errors.New("internal error"),
			},
		}})
		if !assert.NoError(err) {
			return
		}

		_, err = mockServer.DeleteSubstitution(ctx, &mesync.DeleteSubstitutionRequest{
			Material: &mesync.Material{
				Id:    testMaterial,
				Grade: testMaterialType},
			Type: &mesync.DeleteSubstitutionRequest_Multiple{Multiple: &mesync.DeleteSubstitutionRequest_Batch{
				Substitutions: []*mesync.Material{{
					Id:    testSubstitutionName,
					Grade: testSubstitutionType,
				},
				},
			}},
		})
		if assert.Error(err) {
			assert.Equal(parseError(testFactoryID, errors.New("internal error")), err)
		}
	}
}
