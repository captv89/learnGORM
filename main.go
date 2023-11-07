package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Person struct {
	gorm.Model
	Name  string `gorm:"unique"`
	Age   int
	Sex   string
	Cars  []Car
	Bikes []Bike
}

type Car struct {
	gorm.Model
	Name     string
	Brand    string
	PersonID uint
}

type Bike struct {
	gorm.Model
	Name     string
	Brand    string
	PersonID uint
}

func main() {
	dsn := "host=localhost user=postgres password=Pass1234 dbname=learn port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Person{}, &Car{}, &Bike{})
	if err != nil {
		panic(err)
	}

	// Sample data set 1
	Person1 := Person{
		Name: "John",
		Age:  20,
		Sex:  "Male",
		Cars: []Car{
			{
				Name:  "Car1",
				Brand: "Brand1",
			},
			{
				Name:  "Car2",
				Brand: "Brand2",
			},
		},
		Bikes: []Bike{
			{
				Name:  "Bike1",
				Brand: "Brand1",
			},
			{
				Name:  "Bike2",
				Brand: "Brand2",
			},
		},
	}

	// Sample data set 2
	Person2 := Person{
		Name: "Jane",
		Age:  20,
		Sex:  "Female",
		Cars: []Car{
			{
				Name:  "Car3",
				Brand: "Brand3",
			},
			{
				Name:  "Car4",
				Brand: "Brand4",
			},
		},
		Bikes: []Bike{
			{
				Name:  "Bike3",
				Brand: "Brand3",
			},
			{
				Name:  "Bike4",
				Brand: "Brand4",
			},
		},
	}

	// Insert data
	db.Save(&Person1)
	db.Save(&Person2)

	// Query data
	var p1, p2 Person
	db.Preload("Cars").Preload("Bikes").Where("name = ?", Person1.Name).First(&p1)

	// Print data
	fmt.Println("Person 1: ", p1.Name)

	db.Preload("Cars").Preload("Bikes").Where("name = ?", Person2.Name).First(&p2)

	// Print data
	fmt.Println("Person 2: ", p2.Name)

	// Update data
	Person1update := Person{
		Name: "John",
		Age:  21,
		Sex:  "Male",
		Cars: []Car{
			{
				Name:  "Car1",
				Brand: "Brand1",
			},
			{
				Name:  "Car2",
				Brand: "Brand1",
			},
		},
		Bikes: []Bike{
			{
				Name:  "Bike1",
				Brand: "Brand1",
			},
			{
				Name:  "Bike2",
				Brand: "Brand1",
			},
		},
	}

	// Person 2 update
	Person2update := Person{
		Name: "Jane",
		Age:  23,
		Sex:  "Female",
		Cars: []Car{
			{
				Name:  "Car3",
				Brand: "Brand3",
			},
			{
				Name:  "Car4",
				Brand: "Brand3",
			},
		},
		Bikes: []Bike{
			{
				Name:  "Bike3",
				Brand: "Brand3",
			},
			{
				Name:  "Bike4",
				Brand: "Brand3",
			},
		},
	}

	// METHOD 1: Upsert data
	// Upsert data: Can be used for update and insert
	// Problem with this is that in our case of has many the Associated tables are inserted with new records and not updated
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(&Person2update)

	// Print Upsert data
	fmt.Println("After Upsert", Person2)

	// Association Mode
	err = db.Model(&Person{}).Association("Cars").Error
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Association to Cars is OK")
	}

	// Find p1 Association
	var c1 []Car
	err = db.Model(&p1).Association("Cars").Find(&c1)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Found Cars: ", c1)
	}

	//METHOD 2: Update data
	// Update P1 data without updating associated data
	p1UpdateCopy := Person1update
	p1UpdateCopy.Cars = nil
	p1UpdateCopy.Bikes = nil
	p1WithoutAssociateData := p1UpdateCopy
	db.Session(&gorm.Session{FullSaveAssociations: true}).Where("name = ?", Person1update.Name).Updates(&p1WithoutAssociateData)

	// Print updated data
	fmt.Println("After Session Update", Person1)

	// Option 1: Replace p1 Association.
	// Problem here is that the previous association is not deleted but only the fKey is removed
	var c2 = []Car{
		{
			Name:  "Car5",
			Brand: "Brand5",
		},
		{
			Name:  "Car6",
			Brand: "Brand6",
		},
	}
	err = db.Model(&p1).Association("Cars").Replace(&c2)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Replace Cars: ", c2)
	}

	// Option 2 (Most Suitable): Replace p1 Association by clearing and appending
	// Delete p1 Association permanently
	err = db.Unscoped().Model(&p1).Association("Cars").Unscoped().Delete(&c2)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Cleared Cars for p1")
	}

	// Append p1 Association
	var c3 = []Car{
		{
			Name:  "Car7",
			Brand: "Brand7",
		},
		{
			Name:  "Car8",
			Brand: "Brand8",
		},
	}

	err = db.Model(&p1).Association("Cars").Append(&c3)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Append Cars: ", c3)
	}

	// Delete p2 with Association
	// delete user's has one/many/many2many relations when deleting user
	db.Unscoped().Select(clause.Associations).Where("name = ?", p2.Name).Delete(&p2)

	// Close DB connection
	DB, err := db.DB()
	if err != nil {
		panic(err)
	} else {
		err := DB.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("DB Connection Closed")
	}

}
