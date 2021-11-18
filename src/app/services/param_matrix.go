package services

func GenerateMatrixCombinations(matrixMap MatrixMap) ([]ParamMap, error) {
	if len(matrixMap) == 0 {
		return []ParamMap{{}}, nil
	}

	params := []ParamMap{{}}
	for paramKey, paramList := range matrixMap {
		tmpParams := []ParamMap{}
		for _, param := range params {
			for _, paramElement := range paramList {
				tmpParam := ParamMap{}
				for k, v := range param {
					tmpParam[k] = v
				}
				tmpParam[paramKey] = paramElement
				tmpParams = append(tmpParams, tmpParam)
			}
		}
		params = tmpParams
	}
	return params, nil
}
