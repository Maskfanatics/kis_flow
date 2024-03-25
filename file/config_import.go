package file

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
	"kis-flow/kis"
	"os"
	"path"
	"path/filepath"
)

type allConfig struct {
	Flows map[string]*config.KisFlowConfig
	Funcs map[string]*config.KisFuncConfig
	Conns map[string]*config.KisConnConfig
}

func kisTypeFlowConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	flow := new(config.KisFlowConfig)
	if ok := yaml.Unmarshal(confData, flow); ok != nil {
		return errors.New(fmt.Sprintf("%s has wrong format kisType = %s", fileName, kisType))
	}

	// 如果FLow状态为关闭，则不做配置加载
	if common.KisOnOff(flow.Status) == common.FlowDisable {
		return nil
	}

	if _, ok := all.Flows[flow.FlowName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat flow_id:%s", fileName, flow.FlowName))
	}

	// 加入配置集合中
	all.Flows[flow.FlowName] = flow

	return nil
}

func kisTypeFuncConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	fn := new(config.KisFuncConfig)

	if ok := yaml.Unmarshal(confData, fn); ok != nil {
		return errors.New(fmt.Sprintf("%s has wrong format kisType = %s", fileName, kisType))
	}

	if _, ok := all.Flows[fn.FName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat function_id:%s", fileName, fn.FName))
	}

	all.Funcs[fn.FName] = fn

	return nil

}

func kisTypeConnConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	cc := new(config.KisConnConfig)

	if ok := yaml.Unmarshal(confData, cc); ok != nil {
		return errors.New(fmt.Sprintf("%s has wrong format kisType = %s", fileName, kisType))
	}

	if _, ok := all.Flows[cc.CName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat function_id:%s", fileName, cc.CName))
	}

	all.Conns[cc.CName] = cc

	return nil

}

func parseConfigWalkYaml(loadPath string) (*allConfig, error) {
	all := new(allConfig)

	all.Flows = make(map[string]*config.KisFlowConfig)
	all.Funcs = make(map[string]*config.KisFuncConfig)
	all.Conns = make(map[string]*config.KisConnConfig)

	err := filepath.Walk(loadPath, func(filePath string, info fs.FileInfo, err error) error {
		if suffix := path.Ext(filePath); suffix != ".yml" && suffix != ".yaml" {
			return nil
		}

		confData, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		confMap := make(map[string]interface{})

		if err := yaml.Unmarshal(confData, confMap); err != nil {
			return err
		}

		if kisType, ok := confMap["kistype"]; !ok {
			return errors.New(fmt.Sprintf("yaml file %s has no file [kistype]!", filePath))
		} else {
			switch kisType {
			case common.KisIdTypeFlow:
				return kisTypeFlowConfigure(all, confData, filePath, kisType)
			case common.KisIdTypeFunction:
				return kisTypeFuncConfigure(all, confData, filePath, kisType)
			case common.KisIdTypeConnnector:
				return kisTypeConnConfigure(all, confData, filePath, kisType)
			default:
				return errors.New(fmt.Sprintf("%s set wrong kistype %s", filePath, kisType))
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return all, nil
}

// ConfigImportYaml 全盘解析配置文件，yaml格式
func ConfigImportYaml(loadPath string) error {
	all, err := parseConfigWalkYaml(loadPath)
	if err != nil {
		return err
	}
	for flowName, flowConfig := range all.Flows {
		newFlow := flow.NewKisFlow(flowConfig)
		for _, fp := range flowConfig.Flows {
			if err := buildFlow(all, fp, newFlow, flowName); err != nil {
				return err
			}
		}

		kis.Pool().AddFlow(flowName, newFlow)
	}
	return nil
}

func buildFlow(all *allConfig, fp config.KisFlowFunctionParam, newFlow kis.Flow, flowName string) error {
	if funcConfig, ok := all.Funcs[fp.FuncName]; !ok {
		return errors.New(fmt.Sprintf("FlowName [%s] need FuncName [%s], But has No This FuncName Config", flowName, fp.FuncName))
	} else {
		if funcConfig.Option.CName != "" {
			if connConf, ok := all.Conns[funcConfig.Option.CName]; !ok {
				return errors.New(fmt.Sprintf("FuncName [%s] need ConnName [%s], But has No This ConnName Config", fp.FuncName, funcConfig.Option.CName))
			} else {
				_ = funcConfig.AddConnConfig(connConf)
			}
		}
		if err := newFlow.Link(funcConfig, fp.Params); err != nil {
			return err
		}
	}
	return nil
}
