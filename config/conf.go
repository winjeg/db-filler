package config

import (
	"github.com/winjeg/db-filler/cmd"
	config "github.com/winjeg/go-commons/conf"
	logger "github.com/winjeg/go-commons/log"

	"fmt"
	"log"
	"sync"
)

const (
	configFile      = "conf.yaml"
	defaultDbType   = "mysql"
	maxConn         = 30
	idleConn        = 5
	logFormat       = "colored"
	logLevel        = "info"
	logOutput       = "std"
	logReportCaller = true
	connStrPattern  = "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"
)

// database related config
type DbConf struct {
	DbName   string `yaml:"dbName"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbType   string `yaml:"dbType"`
	MaxConn  uint32 `yaml:"maxConn"`
	IdleConn uint32 `yaml:"idleConn"`
}

// get connection str from db settings
func (dc DbConf) GetConnectionStr() string {
	if dc.Port < 1 {
		panic("illegal port")
	}
	if len(dc.Host) < 1 {
		panic("illegal host info")
	}
	if len(dc.Username) < 1 {
		panic("illegal username")
	}
	if len(dc.DbName) < 1 {
		panic("illegal database name")
	}
	return fmt.Sprintf(connStrPattern, dc.Username, dc.Password, dc.Host, dc.Port, dc.DbName)
}

// go routines related config
type WorkerConf struct {
	TableName         string `yaml:"tableName"`
	SqlNum            int    `yaml:"sqlNum"`
	RowNum            int    `yaml:"rowNum"`
	GenerateWorkerNum int    `yaml:"generateWorkerNum"`
	InsertWorkerNum   int    `yaml:"insertWorkerNum"`
	DDLWorkerNum      int    `yaml:"ddlWorkerNum"`
}

// other configs
type ExtraConf struct {
	Sql          string `yaml:"sql"`
	ErrorFile    string `yaml:"errorFile"`
	SqlFile      string `yaml:"sqlFile"`
	OnlyGenerate bool   `yaml:"onlyGenerate"`
}

// total config
type Config struct {
	DbSettings     DbConf             `yaml:"database"`
	WorkerSettings WorkerConf         `yaml:"worker"`
	ExtraSettings  ExtraConf          `yaml:"extra"`
	LogSettings    logger.LogSettings `yaml:"log"`
}

var (
	once sync.Once
	conf *Config
)

// get config, singleton
func GetConf() *Config {
	if conf != nil {
		return conf
	} else {
		once.Do(getConf)
	}
	return conf
}

// get config firstly read from the config file
// if there is none config file that's ok
// some of the config will be given the default values
// and others will be loaded from the command line arguments
func getConf() {
	conf = new(Config)
	err := config.Yaml2Object(configFile, &conf)
	if err != nil {
		log.Printf("no config file found")
	}

	// override all config, whenever the command line exits

	// log settings
	if len(conf.LogSettings.Format) < 1 {
		conf.LogSettings.Format = logFormat
	}
	if len(conf.LogSettings.Level) < 1 {
		conf.LogSettings.Level = logLevel
	}
	if len(conf.LogSettings.Output) < 1 {
		conf.LogSettings.Output = logOutput
	}
	if err != nil {
		conf.LogSettings.ReportCaller = logReportCaller
	}

	// db settings
	if len(conf.DbSettings.DbType) < 1 {
		conf.DbSettings.DbType = defaultDbType
	}
	if conf.DbSettings.MaxConn < 1 {
		conf.DbSettings.MaxConn = maxConn
	}
	if conf.DbSettings.IdleConn < 1 {
		conf.DbSettings.IdleConn = idleConn
	}
	if len(cmd.GetDbName()) > 0 {
		conf.DbSettings.DbName = cmd.GetDbName()
	}
	if len(cmd.GetHost()) > 0 {
		conf.DbSettings.Host = cmd.GetHost()
	}
	if cmd.GetPort() > 0 {
		conf.DbSettings.Port = cmd.GetPort()
	}
	if len(cmd.GetUserName()) > 0 {
		conf.DbSettings.Username = cmd.GetUserName()
	}
	if len(cmd.GetPassword()) > 0 {
		conf.DbSettings.Password = cmd.GetPassword()
	}

	// go routine settings
	if len(cmd.GetTableName()) > 0 {
		conf.WorkerSettings.TableName = cmd.GetTableName()
	}
	if cmd.GetGenSqlNum() > 0 {
		conf.WorkerSettings.SqlNum = cmd.GetGenSqlNum()
	}
	if cmd.GetRowNumPerSql() > 0 {
		conf.WorkerSettings.RowNum = cmd.GetRowNumPerSql()
	}
	if cmd.GetDDLWorkerNum() > 0 {
		conf.WorkerSettings.DDLWorkerNum = cmd.GetDDLWorkerNum()
	}
	if cmd.GetInsertionWorkerNum() > 0 {
		conf.WorkerSettings.InsertWorkerNum = cmd.GetInsertionWorkerNum()
	}
	if cmd.GetGenerateWorkerNum() > 0 {
		conf.WorkerSettings.GenerateWorkerNum = cmd.GetGenerateWorkerNum()
	}

	// extra settings
	if len(cmd.GetSql()) > 0 {
		conf.ExtraSettings.Sql = cmd.GetSql()
	}
	if len(cmd.GetErrorFile()) > 0 {
		conf.ExtraSettings.ErrorFile = cmd.GetErrorFile()
	}
	if len(cmd.GetSqlFile()) > 0 {
		conf.ExtraSettings.SqlFile = cmd.GetSqlFile()
	}
	if cmd.OnlyGenerate() {
		conf.ExtraSettings.OnlyGenerate = cmd.OnlyGenerate()
	}
}
