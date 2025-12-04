package company_facts

var naughtyObj = map[string]any{
	"cik":  "123",
	"name": "pear",
	"float": map[string]any{
		"float_key1": "float_value1",
		"float_key2": "float_value2",
		"float_key3": map[string]any{
			"float_nested_key1": "float_nested_value1",
			"float_nested_key2": "float_nested_value2",
			"float_nested_key3": []string{
				"float_nested_value3_1",
				"float_nested_value3_2",
				"float_nested_value3_3",
			},
			"float_nested_key4": []map[string]any{
				{
					"float_super_nested_1_key1": "float_super_nested_1_value1",
					"float_super_nested_1_key2": "float_super_nested_1_value2",
					"float_super_nested_1_key3": "float_super_nested_1_value3",
					"float_super_nested_1_key4": "float_super_nested_1_value4",
				},
				{
					"float_super_nested_2_key1": "float_super_nested_2_value1",
					"float_super_nested_2_key2": "float_super_nested_2_value2",
					"float_super_nested_2_key3": "float_super_nested_2_value3",
					"float_super_nested_2_key4": "float_super_nested_2_value4",
				},
			},
		},
	},
	"data": map[string]any{
		"data_key1": "data_value1",
		"data_key2": "data_value2",
		"data_key3": map[string]any{
			"data_nested_key1": "data_nested_value1",
			"data_nested_key2": "data_nested_value2",
			"data_nested_key3": []string{
				"data_nested_value3_1",
				"data_nested_value3_2",
				"data_nested_value3_3",
			},
			"data_nested_key4": []map[string]any{
				{
					"data_super_nested_1_key1": "data_super_nested_1_value1",
					"data_super_nested_1_key2": "data_super_nested_1_value2",
					"data_super_nested_1_key3": "data_super_nested_1_value3",
					"data_super_nested_1_key4": "data_super_nested_1_value4",
				},
				{
					"data_super_nested_2_key1": "data_super_nested_2_value1",
					"data_super_nested_2_key2": "data_super_nested_2_value2",
					"data_super_nested_2_key3": "data_super_nested_2_value3",
					"data_super_nested_2_key4": []map[string]any{
						{
							"data_giga_nested_1_key1": "data_giga_nested_1_value1",
							"data_giga_nested_1_key2": "data_giga_nested_1_value2",
							"data_giga_nested_1_key3": "data_giga_nested_1_value3",
							"data_giga_nested_1_key4": "data_giga_nested_1_value4",
						},
						{
							"data_giga_nested_2_key1": "data_giga_nested_2_value1",
							"data_giga_nested_2_key2": "data_giga_nested_2_value2",
							"data_giga_nested_2_key3": "data_giga_nested_2_value3",
							"data_giga_nested_2_key4": "data_giga_nested_2_value4",
						},
						{
							"data_giga_nested_3_key1": "data_giga_nested_3_value1",
							"data_giga_nested_3_key2": "data_giga_nested_3_value2",
							"data_giga_nested_3_key3": "data_giga_nested_3_value3",
							"data_giga_nested_3_key4": "data_giga_nested_3_value4",
						},
					},
				},
			},
		},
	},
	"funkey": []map[string]any{
		{
			"funkey_1_key1": "funkey_1_value1",
			"funkey_1_key2": "funkey_1_value2",
			"funkey_1_key3": "funkey_1_value3",
			"funkey_1_key4": "funkey_1_value4",
		},
		{
			"funkey_2_key1": "funkey_2_value1",
			"funkey_2_key2": "funkey_2_value2",
			"funkey_2_key3": "funkey_2_value3",
			"funkey_2_key4": "funkey_2_value4",
		},
		{
			"funkey_3_key1": "funkey_3_value1",
			"funkey_3_key2": "funkey_3_value2",
			"funkey_3_key3": "funkey_3_value3",
			"funkey_3_key4": "funkey_3_value4",
		},
	},
}

// func flattenObj(obj json.RawMessage, parentKey *string, keyTrace []string) {
func flattenObj(obj map[string]any, parentKey *string, keyTrace []string) {

}

func GetCompanyFacts() {
	flattenObj(naughtyObj, nil, nil)
}
