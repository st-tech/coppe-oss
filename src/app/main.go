package main

import (
	"context"
	"log"

	"zozo.com/coppe/src/app/services"
)

func main() {
	ctx := context.Background()

	fileNameToDataMap, err := services.ReadFiles("./yaml")
	if err != nil {
		log.Fatalf("file retrieval failed. %v", err)
	}

	var validations []services.Validation
	for fileName, data := range fileNameToDataMap {
		validationsInFile, err := services.UnmarshalToValidations(data)
		if err != nil {
			log.Fatalf("yaml unmarshal failed. %v", err)
		}
		for i, v := range validationsInFile {
			v.Identifier.FoundAt.FileName = fileName
			v.Identifier.FoundAt.Index = i
			validations = append(validations, v)
		}
	}

	for _, v := range validations {
		v.Sql, err = services.GetSql(v.Sql, v.SqlFile)
		if err != nil {
			log.Fatal(err)
		}

		matrixCombinations, err := services.GenerateMatrixCombinations(v.Matrix)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, matrixCombination := range matrixCombinations {
			if len(v.Params) == 0 {
				v.Params = services.ParamMap{}
			}
			for key, value := range matrixCombination {
				if len(v.Identifier.MatrixValue) == 0 {
					v.Identifier.MatrixValue = services.ParamMap{}
				}
				v.Identifier.MatrixValue[key] = value
			}
			err := services.Validate(ctx, v)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
