package ormdb

//CommonConnectionDef contains common connection parameters
type CommonConnectionDef struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
}

//ORMDBInstance represents a ORMDB instance
type ORMDBInstance struct {
	confs map[string]interface{}
}

func NewORMDBInstance() {

}
