package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm-advanced-filter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	Age       uint
	BirthDate time.Time
}

var filter *gorm_advanced_filter.Filter

func TestCreateQuery(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("./gorm_advanced_data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	filter = gorm_advanced_filter.MakeGormAdvancedFilter(db)
}

func TestGetAllUserCount(t *testing.T) {
	var users []User
	filter.ToSql().Find(&users)
	assert.Equal(t, 3, len(users))
}

func TestGetUserNameContain(t *testing.T) {
	var users []User
	filter.
		Contains("Farcha", "name", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 2, len(users))
	filter.Clear()
}

func TestGetUserNameDoesNotContain(t *testing.T) {
	var users []User
	filter.
		DoesNotContains("Farcha", "name", "AND").
		ToSql().
		Find(&users)

	filter.Clear()
	assert.Equal(t, 1, len(users))
}

func TestGetUserNameIn(t *testing.T) {
	var user User
	filter.Is("Miqdad Farcha", "name", "AND").
		ToSql().
		First(&user)

	assert.Equal(t, "Miqdad Farcha", user.Name)
	assert.Equal(t, 22, int(user.Age))
	assert.Equal(t, "1999-01-20", user.BirthDate.Format("2006-01-02"))
	filter.Clear()
}

func TestGetUserNameNotIn(t *testing.T) {
	var users []User
	filter.IsNot("Miqdad Farcha", "name", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 2, len(users))
	filter.Clear()
}

func TestGetUserNameStartWith(t *testing.T) {
	var users []User
	filter.StartWith("Bing", "name", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 1, len(users))
	filter.Clear()
}

func TestGetUserNameEndWith(t *testing.T) {
	var users []User
	filter.EndWith("Farcha", "name", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 2, len(users))
	filter.Clear()
}

func TestAgeEqual(t *testing.T) {
	var users []User
	filter.Equal("22", "age", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 1, len(users))
	filter.Clear()
}

func TestAgeLessThan(t *testing.T) {
	var users []User
	filter.LessThan("22", "age", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 2, len(users))
	filter.Clear()
}

func TestAgeMoreThan(t *testing.T) {
	var users []User
	filter.MoreThan("21", "age", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 1, len(users))
	filter.Clear()
}

func TestAgeLessThanEqual(t *testing.T) {
	var users []User
	filter.LessThanEqual("22", "age", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 3, len(users))
	filter.Clear()
}

func TestAgeMoreThanEqual(t *testing.T) {
	var users []User
	filter.MoreThanEqual("18", "age", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 3, len(users))
	filter.Clear()
}

func TestAgeNotEqual(t *testing.T) {
	var users []User
	filter.NotEqual("21", "age", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 2, len(users))
	filter.Clear()
}

func TestBirthDateAt(t *testing.T) {
	var users []User
	filter.At("1999-01-20", "birth_date", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 1, len(users))
	filter.Clear()
}

func TestBirthDateBefore(t *testing.T) {
	var users []User
	filter.Before("2000-01-01", "birth_date", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 2, len(users))
	filter.Clear()
}

func TestBirthDateAfter(t *testing.T) {
	var users []User
	filter.After("2000-01-01", "birth_date", "AND").
		ToSql().
		Find(&users)

	assert.Equal(t, 1, len(users))
	filter.Clear()
}

func TestMultipleQuery(t *testing.T) {
	var users []User

	filter.Contains("Farcha", "name", "AND").
		After("2000-01-01", "birth_date", "AND").
		Is("Angger Pangestu", "name", "OR").
		ToSql().
		Find(&users)

	assert.Equal(t, 2, len(users))
	filter.Clear()
}

func TestBuildAndParse(t *testing.T) {
	var users []User
	// Encode some query
	filterEncoded := filter.Contains("Farcha", "name", "AND").
		After("2000-01-01", "birth_date", "AND").
		Is("Angger Pangestu", "name", "OR").Build()

	fmt.Println(filterEncoded)

	filter, err := gorm_advanced_filter.Parse(filter.ToSql(), filterEncoded)
	if err != nil {
		assert.Error(t, err)
	}

	filter.ToSql().Find(&users)

	assert.Equal(t, 2, len(users))
}
