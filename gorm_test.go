package nrcontext

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"testing"
)

type Model struct {
	Id    int
	Value string
}

func TestGorm(t *testing.T) {
	var first Model
	m1 := Model{Value: "test1"}
	m2 := Model{Value: "test2"}
	m3 := Model{Value: "test3"}

	os.Remove("./foo.db")
	db, err := gorm.Open("sqlite3", "./foo.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.CreateTable(&Model{}).Error; err != nil {
		t.Error(err)
	}
	// shouldn't be in log
	if err := db.Create(&m1).Error; err != nil {
		t.Error(err)
	}
	if err := db.Create(&m3).Error; err != nil {
		t.Error(err)
	}

	app := &NewrelicAppMock{}
	txn := app.StartTransaction("txn-name", nil, nil)
	ctx := context.WithValue(context.Background(), txnKey, txn)

	AddGormCallbacks(db)
	ctxDB := SetTxnToGorm(ctx, db)

	// should be in log
	if err := ctxDB.Create(&m2).Error; err != nil {
		t.Error(err)
	}
	if err := ctxDB.First(&first).Error; err != nil {
		t.Error(err)
	}
	m2.Value = "new value"
	if err := ctxDB.Save(&m2).Error; err != nil {
		t.Error(err)
	}
	if err := ctxDB.Delete(&m3).Error; err != nil {
		t.Error(err)
	}
	if m1.Id != first.Id {
		t.Error("You just broke gorm")
	}

	// shouldn't be in log
	var reloadM Model
	db.First(&reloadM, m2)
}