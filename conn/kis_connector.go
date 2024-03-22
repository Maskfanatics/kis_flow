package conn

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/id"
	"kis-flow/kis"
	"sync"
)

type KisConnector struct {
	CId   string
	CName string
	Conf  *config.KisConnConfig

	onceInit sync.Once
}

func NetKisConnector(conf *config.KisConnConfig) *KisConnector {
	conn := new(KisConnector)
	conn.CName = conf.CName
	conn.CId = id.KisId(common.KisIdTypeConnnector)
	conn.Conf = conf

	return conn
}

func (conn *KisConnector) Init() error {
	var err error

	conn.onceInit.Do(func() {
		err = kis.Pool().CallConnInit(conn)
	})

	return err
}

func (conn *KisConnector) Call(ctx context.Context, flow kis.Flow, args interface{}) error {
	if err := kis.Pool().CallConnector(ctx, flow, conn, args); err != nil {
		return err
	}
	return nil
}

func (conn *KisConnector) GetName() string {
	return conn.CName
}

func (conn *KisConnector) GetConfig() *config.KisConnConfig {
	return conn.Conf
}

func (conn *KisConnector) GetId() string {
	return conn.CId
}
