package file

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"kis-flow/common"
	"kis-flow/kis"
	"os"
)

func ConfigExportYaml(flow kis.Flow, savePath string) error {
	if data, err := yaml.Marshal(flow.GetConfig()); err != nil {
		return err
	} else {
		err := os.WriteFile(savePath+common.KisIdTypeFlow+"-"+flow.GetName()+".yaml", data, 0644)
		if err != nil {
			return err
		}

		for _, fp := range flow.GetConfig().Flows {
			fConf := flow.GetFuncConfigByName(fp.FuncName)
			if fConf == nil {
				return errors.New(fmt.Sprintf("function name = %s config is nil ", fp.FuncName))
			}

			if fdata, err := yaml.Marshal(fConf); err != nil {
				return err
			} else {
				if err := os.WriteFile(savePath+common.KisIdTypeFunction+"-"+fp.FuncName+".yaml", fdata, 0644); err != nil {
					return err
				}
			}

			if fConf.Option.CName != "" {
				cConf, err := fConf.GetConnConfig()
				if err != nil {
					return err
				}
				if cdata, err := yaml.Marshal(cConf); err != nil {
					return err
				} else {
					if err := os.WriteFile(savePath+common.KisIdTypeConnnector+"-"+cConf.CName+".yaml", cdata, 0644); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
