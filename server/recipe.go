package server

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/shopspring/decimal"

	"gitlab.kenda.com.tw/kenda/mcom"
	mcomErr "gitlab.kenda.com.tw/kenda/mcom/errors"
	"gitlab.kenda.com.tw/kenda/mcom/impl/orm/models"
	"gitlab.kenda.com.tw/kenda/mcom/utils/types"

	pb "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
	pbTypes "gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/types"
)

func parseRecipeOptionFlow(optFlows []*pb.OptionalFlow) []*mcom.RecipeOptionalFlow {
	flows := make([]*mcom.RecipeOptionalFlow, len(optFlows))
	for i, optFlow := range optFlows {
		flows[i] = &mcom.RecipeOptionalFlow{
			Name:           optFlow.GetName(),
			OIDs:           optFlow.GetProcessOids(),
			MaxRepetitions: optFlow.GetMaxRepetitions(),
		}
	}
	return flows
}

func parseProcesses(processes []*pb.Process) []*mcom.Process {
	procs := make([]*mcom.Process, len(processes))
	for i, process := range processes {
		procs[i] = &mcom.Process{
			OID:           process.GetReferenceOid(),
			OptionalFlows: parseRecipeOptionFlow(process.OptionalFlows),
		}
	}
	return procs
}

func parseRecipes(reqRecipes []*pb.Recipe) ([]mcom.Recipe, error) {
	recipes := make([]mcom.Recipe, len(reqRecipes))
	for i, recipe := range reqRecipes {
		processDefinitions, err := parseProcessDefinition(recipe.GetProcessDefs())
		if err != nil {
			return nil, fmt.Errorf("recipe %s, error=%s", recipe.GetId(), err.Error())
		}

		t := time.Unix(0, 0) // compatible usage
		if releasedAt := recipe.GetVersion().GetReleasedAt(); releasedAt != nil {
			t, err = ptypes.Timestamp(recipe.GetVersion().GetReleasedAt())
			if err != nil {
				return []mcom.Recipe{}, err
			}
		}

		recipes[i] = mcom.Recipe{
			ID: recipe.GetId(),
			Product: mcom.Product{
				ID:   recipe.GetProductId(),
				Type: recipe.GetProductType(),
			},
			Version: mcom.RecipeVersion{
				Major:      recipe.GetVersion().GetMajor(),
				Minor:      recipe.GetVersion().GetMinor(),
				Stage:      recipe.GetVersion().GetStage(),
				ReleasedAt: types.ToTimeNano(t),
			},
			Processes:          parseProcesses(recipe.Processes),
			ProcessDefinitions: processDefinitions,
		}
	}
	return recipes, nil
}

func parseRecipeTool(pbTools []*pb.RecipeTool) []*mcom.RecipeTool {
	tools := make([]*mcom.RecipeTool, len(pbTools))
	for i, t := range pbTools {
		tools[i] = &mcom.RecipeTool{
			Type:     t.GetType(),
			ID:       t.GetId(),
			Required: t.GetRequired(),
		}
	}
	return tools
}

type typeParameter struct {
	high *pbTypes.Decimal
	mid  *pbTypes.Decimal
	low  *pbTypes.Decimal
}

type parameter struct {
	high *decimal.Decimal
	mid  *decimal.Decimal
	low  *decimal.Decimal
}

func parseParameter(param typeParameter) (*parameter, error) {
	high, err := toPointerDecimal(param.high)
	if err != nil {
		return nil, fmt.Errorf("high value error=%s", err.Error())
	}

	mid, err := toPointerDecimal(param.mid)
	if err != nil {
		return nil, fmt.Errorf("mid value error=%s", err.Error())
	}

	low, err := toPointerDecimal(param.low)
	if err != nil {
		return nil, fmt.Errorf("low value error=%s", err.Error())
	}

	return &parameter{
		high: high,
		mid:  mid,
		low:  low,
	}, nil
}

func parseRecipeMaterial(recipeMaterials []*pb.RecipeMaterial) ([]*mcom.RecipeMaterial, error) {
	materials := make([]*mcom.RecipeMaterial, len(recipeMaterials))
	for i, mtrl := range recipeMaterials {
		param, err := parseParameter(typeParameter{
			high: mtrl.GetValue().GetHigh(),
			mid:  mtrl.GetValue().GetMid(),
			low:  mtrl.GetValue().GetLow(),
		})
		if err != nil {
			return nil, fmt.Errorf("recipe material %d:%s, error=%s", i, mtrl.GetName(), err.Error())
		}

		materials[i] = &mcom.RecipeMaterial{
			Name:  mtrl.GetName(),
			Grade: mtrl.GetGrade(),
			Value: mcom.RecipeMaterialParameter{
				High: param.high,
				Mid:  param.mid,
				Low:  param.low,
				Unit: mtrl.GetValue().GetUnit(),
			},
			Site:             mtrl.GetSite(),
			RequiredRecipeID: mtrl.GetRequiredRecipeId(),
		}
	}
	return materials, nil
}

func parseRecipeProperty(recipeProperty []*pb.RecipeProperty) ([]*mcom.RecipeProperty, error) {
	property := make([]*mcom.RecipeProperty, len(recipeProperty))
	for i, rp := range recipeProperty {
		param, err := parseParameter(typeParameter{
			high: rp.GetParam().GetHigh(),
			mid:  rp.GetParam().GetMid(),
			low:  rp.GetParam().GetLow(),
		})
		if err != nil {
			return nil, fmt.Errorf("recipe property %d:%s, error=%s", i, rp.GetName(), err.Error())
		}

		property[i] = &mcom.RecipeProperty{
			Name: rp.GetName(),
			Param: &mcom.RecipePropertyParameter{
				High: param.high,
				Mid:  param.mid,
				Low:  param.low,
				Unit: rp.GetParam().GetUnit(),
			},
		}
	}
	return property, nil
}

func parseRecipeProcessStep(pbSteps []*pb.RecipeProcessStep) ([]*mcom.RecipeProcessStep, error) {
	steps := make([]*mcom.RecipeProcessStep, len(pbSteps))
	for i, s := range pbSteps {
		materials, err := parseRecipeMaterial(s.GetMaterials())
		if err != nil {
			return nil, fmt.Errorf("recipe process step %d, error=%s", i, err.Error())
		}

		controls, err := parseRecipeProperty(s.GetControls())
		if err != nil {
			return nil, fmt.Errorf("recipe process step %d, controls error=%s", i, err.Error())
		}

		measurements, err := parseRecipeProperty(s.GetMeasurements())
		if err != nil {
			return nil, fmt.Errorf("recipe process step %d, mmeasurements error=%s", i, err.Error())
		}

		steps[i] = &mcom.RecipeProcessStep{
			Materials:    materials,
			Controls:     controls,
			Measurements: measurements,
		}
	}
	return steps, nil
}

func parseProcessConfigs(processConfigs []*pb.RecipeProcessConfig) ([]*mcom.RecipeProcessConfig, error) {
	configs := make([]*mcom.RecipeProcessConfig, len(processConfigs))
	for i, procConfig := range processConfigs {
		batchSize, err := decimal.NewFromString(procConfig.GetBatchSize().GetValue())
		if err != nil {
			return nil, fmt.Errorf("process config %d, batchSize error=%s", i, err.Error())
		}

		steps, err := parseRecipeProcessStep(procConfig.GetSteps())
		if err != nil {
			return nil, fmt.Errorf("process config %d, steps error=%s", i, err.Error())
		}

		commonControls, err := parseRecipeProperty(procConfig.GetCommonsControls())
		if err != nil {
			return nil, fmt.Errorf("process config %d, common controls error=%s", i, err.Error())
		}

		commonProperties, err := parseRecipeProperty(procConfig.GetCommonsProperties())
		if err != nil {
			return nil, fmt.Errorf("process config %d, common properties error=%s", i, err.Error())
		}

		configs[i] = &mcom.RecipeProcessConfig{
			Stations:         procConfig.GetStations(),
			BatchSize:        &batchSize,
			Unit:             procConfig.GetUnit(),
			Tools:            parseRecipeTool(procConfig.GetTools()),
			Steps:            steps,
			CommonControls:   commonControls,
			CommonProperties: commonProperties,
		}
	}
	return configs, nil
}

func parseProductValidPeriod(configs []*pb.ProductValidPeriodConfig) ([]*mcom.ProductValidPeriodConfig, error) {
	validPeriodConfigs := make([]*mcom.ProductValidPeriodConfig, len(configs))
	for i, config := range configs {
		validPeriodConfigs[i] = &mcom.ProductValidPeriodConfig{
			Standing: time.Duration(config.Standing),
			Expiry:   time.Duration(config.Expiry),
		}
	}
	return validPeriodConfigs, nil
}

func parseProcessDefinition(processes []*pb.ProcessDefinition) ([]mcom.ProcessDefinition, error) {
	procs := make([]mcom.ProcessDefinition, len(processes))
	for i, process := range processes {
		configs, err := parseProcessConfigs(process.GetConfigs())
		if err != nil {
			return nil, fmt.Errorf("process %d:%s, error=%s", i, process.GetOid(), err.Error())
		}

		procs[i] = mcom.ProcessDefinition{
			OID:     process.GetOid(),
			Name:    process.GetName(),
			Type:    process.GetType(),
			Configs: configs,
			Output: mcom.OutputProduct{
				ID:   process.GetProduct().GetId(),
				Type: process.GetProduct().GetType(),
			},
			ProductValidPeriod: mcom.ProductValidPeriodConfig{
				Standing: time.Duration(process.ProductValidPeriod.GetStanding()),
				Expiry:   time.Duration(process.ProductValidPeriod.GetExpiry()),
			},
		}
	}
	return procs, nil
}

// CreateRecipes implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) CreateRecipes(ctx context.Context, req *pb.CreateRecipesRequest) (*empty.Empty, error) {
	recipes, err := parseRecipes(req.GetRecipes())
	if err != nil {
		return nil, mcomErr.Error{
			Code:    mcomErr.Code_INVALID_NUMBER,
			Details: err.Error(),
		}
	}

	return &empty.Empty{}, s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.CreateRecipes(ctx, mcom.CreateRecipesRequest{
			Recipes: recipes,
		})
	})
}

// DeleteRecipes implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) DeleteRecipes(ctx context.Context, req *pb.DeleteRecipesRequest) (*empty.Empty, error) {
	ids := req.GetIds()

	return &empty.Empty{}, s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.DeleteRecipe(ctx, mcom.DeleteRecipeRequest{
			IDs: ids,
		})
	})
}

// CreateSubstitutes implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) CreateSubstitutes(ctx context.Context, req *pb.CreateSubstitutionRequest) (*empty.Empty, error) {
	contents, err := parseSubstitutionContents(req.GetSubstitutions())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.AddSubstitutions(ctx, mcom.BasicSubstitutionRequest{
			ProductID: models.ProductID{ID: req.Material.GetId(), Grade: req.Material.GetGrade()},
			Contents:  contents,
		})
	})
}

// UpdateSubstitutes implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) UpdateSubstitutes(ctx context.Context, req *pb.UpdateSubstitutionRequest) (*empty.Empty, error) {
	contents, err := parseSubstitutionContents(req.GetSubstitutions())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.UpdateSubstitutions(ctx, mcom.BasicSubstitutionRequest{
			ProductID: models.ProductID{
				ID:    req.Material.GetId(),
				Grade: req.Material.GetGrade(),
			},
			Contents: contents,
		})
	})
}

// DeleteSubstitution implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s Server) DeleteSubstitution(ctx context.Context, req *pb.DeleteSubstitutionRequest) (*empty.Empty, error) {
	contents := make([]mcom.SubstitutionDeletionContent, len(req.GetMultiple().GetSubstitutions()))
	for i, substitution := range req.GetMultiple().GetSubstitutions() {
		contents[i] = mcom.SubstitutionDeletionContent{
			ProductID: models.ProductID{
				ID:    substitution.GetId(),
				Grade: substitution.GetGrade(),
			},
		}
	}
	return &empty.Empty{}, s.eachFactory(ctx, func(ctx context.Context, _ string, _ int, dm mcom.DataManager) error {
		return dm.DeleteSubstitutions(ctx, mcom.DeleteSubstitutionsRequest{
			ProductID: models.ProductID{
				ID:    req.Material.GetId(),
				Grade: req.Material.GetGrade(),
			},
			Contents:  contents,
			DeleteAll: req.GetAll(),
		})
	})
}

func parseSubstitutionContents(s []*pb.Substitutions) ([]models.Substitution, error) {
	contents := make([]models.Substitution, len(s))
	for i, substitution := range s {
		proportionValue := substitution.GetProportion().GetValue()
		proportion, err := decimal.NewFromString(proportionValue)
		if err != nil {
			return nil, mcomErr.Error{
				Code:    mcomErr.Code_INVALID_NUMBER,
				Details: fmt.Sprintf("Name: %s, Grade: %s, Proportion: %s, %v", substitution.Material.GetId(), substitution.Material.GetGrade(), proportionValue, err),
			}
		}

		contents[i] = models.Substitution{
			ID:         substitution.Material.GetId(),
			Grade:      substitution.Material.GetGrade(),
			Proportion: proportion,
		}
	}
	return contents, nil
}
