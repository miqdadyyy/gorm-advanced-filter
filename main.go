package gorm_advanced_filter

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
)

type Clause struct {
	Column       string    `json:"column"`
	FunctionName string    `json:"function_name"`
	Value        string    `json:"value"`
	Operator     string    `json:"operator"`
	Group        []*Clause `json:"group"`
}

type Filter struct {
	query   *gorm.DB
	encoded []*Clause
}

func (filter *Filter) Contains(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "Contains",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse("%%"+value+"%%", column, "LIKE", operator)
}

func (filter *Filter) DoesNotContains(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "DoesNotContains",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse("%%"+value+"%%", column, "NOT LIKE", operator)
}

func (filter *Filter) Is(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "Is",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(value, column, "=", operator)
}

func (filter *Filter) IsNot(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "IsNot",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(value, column, "!=", operator)
}

func (filter *Filter) StartWith(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "StartWith",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(value+"%%", column, "LIKE", operator)
}

func (filter *Filter) EndWith(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "EndWith",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse("%%"+value, column, "LIKE", operator)
}

func (filter *Filter) Equal(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "Equal",
		Value:        value,
		Operator:     operator,
	})

	return filter.Is(value, column, operator)
}

func (filter *Filter) NotEqual(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "NotEqual",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(value, column, "<>", operator)
}

func (filter *Filter) LessThan(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "LessThan",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(value, column, "<", operator)
}

func (filter *Filter) MoreThan(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "MoreThan",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(value, column, ">", operator)
}

func (filter *Filter) LessThanEqual(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "LessThanEqual",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(value, column, "<=", operator)
}

func (filter *Filter) MoreThanEqual(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "MoreThanEqual",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(value, column, ">=", operator)
}

func (filter *Filter) At(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "At",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(fmt.Sprintf("%s", value), fmt.Sprintf("%s", column), "=", operator)
}

func (filter *Filter) Before(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "Before",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(fmt.Sprintf("%s", value), fmt.Sprintf("%s", column), "<", operator)
}

func (filter *Filter) After(value, column, operator string) *Filter {
	filter.encoded = append(filter.encoded, &Clause{
		Column:       column,
		FunctionName: "After",
		Value:        value,
		Operator:     operator,
	})

	return filter.parse(fmt.Sprintf("%s", value), fmt.Sprintf("%s", column), ">", operator)
}

func (filter *Filter) parse(value, column, expression, operator string) *Filter {
	if strings.ToLower(operator) == "or" {
		filter.query = filter.query.Or(fmt.Sprintf("%s %s ?", column, expression), value)
	} else {
		filter.query = filter.query.Where(fmt.Sprintf("%s %s ?", column, expression), value)
	}

	return filter
}

func (filter *Filter) ToSql() *gorm.DB {
	return filter.query
}

func (filter *Filter) Clear() *Filter {
	filter.query.Statement.Clauses = map[string]clause.Clause{}
	filter.encoded = nil
	return filter
}

func (filter *Filter) Build() string {
	jsonFilter, _ := json.Marshal(filter.encoded)
	return string(jsonFilter)
}

func Parse(db *gorm.DB, input string) (*Filter, error) {
	filter := MakeGormAdvancedFilter(db)
	var encodedData []*Clause
	// Total 15 functions
	functions := map[string]interface{}{
		"Contains":        filter.Contains,
		"DoesNotContains": filter.DoesNotContains,
		"Is":              filter.Is,
		"IsNot":           filter.IsNot,
		"StartWith":       filter.StartWith,
		"EndWith":         filter.EndWith,
		"Equal":           filter.Equal,
		"NotEqual":        filter.NotEqual,
		"LessThan":        filter.LessThan,
		"MoreThan":        filter.MoreThan,
		"LessThanEqual":   filter.LessThanEqual,
		"MoreThanEqual":   filter.MoreThanEqual,
		"At":              filter.At,
		"Before":          filter.Before,
		"After":           filter.After,
	}

	if err := json.Unmarshal([]byte(input), &encodedData); err != nil {
		return nil, err
	}

	for _, filterClause := range encodedData {
		function := reflect.ValueOf(functions[filterClause.FunctionName])
		params := []string{
			filterClause.Value,
			filterClause.Column,
			filterClause.Operator,
		}

		input := make([]reflect.Value, len(params))
		for key, param := range params {
			input[key] = reflect.ValueOf(param)
		}

		function.Call(input)
	}

	return filter, nil
}

func MakeGormAdvancedFilter(db *gorm.DB) *Filter {
	return &Filter{
		query: db,
	}
}
