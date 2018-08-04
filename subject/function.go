package subject

import (
	"os"
	"encoding/json"
	"io/ioutil"
)

func saveSubjects() error{
	data,err:=json.Marshal(mSubjects)
	if err!=nil{
		return err
	}
	if err := ioutil.WriteFile(dir+"/subject_conf/subject.json",data,os.ModeAppend);err != nil{
		return err
	}
	return nil
}
