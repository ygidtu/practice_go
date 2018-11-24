package models

import (
	"os"
	"hash/fnv"
	"io/ioutil"
	"log"
	"math/big"
	"path/filepath"

	"github.com/astaxie/beego/orm"
)

// generate hash code based on full path of files
func Hash(s string) *big.Int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return big.NewInt(int64(h.Sum32()))
}

// check nil, cause panic
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// Update database when first init this database
func UpdateDatabase() {
	home := "./Downloads"

	if _, err := os.Stat(home); os.IsNotExist(err) {
		// path/to/whatever does not exist
		os.MkdirAll(home, os.ModePerm)
	}

	o := orm.NewOrm()
	o.Using("default")

	targetDirs := []string{home}

	for {
		if len(targetDirs) <= 0 {
			break
		}

		targetDir := targetDirs[0]
		targetDirs = targetDirs[1:]

		files, err := ioutil.ReadDir(targetDir)
		CheckErr(err)

		for _, f := range files {
			if f.IsDir() {
				targetDirs = append(targetDirs, filepath.Join(targetDir, f.Name()))
			}

			record := new(File)

			record.FileName = f.Name()
			record.IsDir = f.IsDir()
			record.Parent = targetDir
			record.Id = *Hash(filepath.Join(targetDir, f.Name()))

			o.Insert(record)
		}
	}

}


// check username and password
func CheckUser(name string, password string) string {

	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable("user")
	exist := qs.Filter("Name", name).Exist()

	if !exist {
		return "User not exists"
	}

	exist = qs.Filter("Name", name).Filter("Password", password).Exist()
	
	if !exist {
		return "Password invalid"
	}

	return "Success"
}


// create user
func CreateUser(name string, password string) {
	o := orm.NewOrm()
	o.Using("default")

	if CheckUser(name, password) != "User not exist" {
		panic("User exists")
	}

	user := new(User)

	user.Name = name
	user.Password = password
	user.Admin = false

	o.Insert(user)
}


// delete files
func DeleteFile(file big.Int) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable("file")
	
	if qs.Filter("Id", file).Exist() {
		var file_object File
		qs.Filter("Id", file).One(&file)

		if file_object.IsDir {
			var files []File
			qs.Filter("Parent", filepath.Join(file_object.Parent, file_object.FileName)).All(&files)

			for _, i := range files {
				o.Delete(&i)
			}
		}
		o.Delete(&File{Id: file})
	}
}


// get file path
func GetFile(id int) string {
	o := orm.NewOrm()
	o.Using("default")	
	qs := o.QueryTable("file").Filter("Id", id).Filter("Isdir", false)

	if !qs.Exist() {
		return ""
	}

	var file File
	qs.One(file)

	return filepath.Join(file.Parent, file.FileName)
}

// get list of directory
func GetDir(id int) []File {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable("file")

	var files []File

	directoryQuery := qs.Filter("Id", id).Filter("Isdir", true)
	if directoryQuery.Exist() {
		var directory File
		directoryQuery.One(&directory)
	
		qs.Filter("Parent", filepath.Join(directory.Parent, directory.FileName)).All(&files)
	}

	return files
}
