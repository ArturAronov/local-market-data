package utils

import (
	"fmt"
)

/*
 * 		Turns object, such as
 * 		nestedObject = {
 * 			key1-1: value1-1
 * 			key1-2: {
 * 				key2-1: value2-1
 * 				key2-2: {
 * 					key3-1: value3-1
 * 					key3-2: {
 * 						key4-1: value4-1,
 * 						key4-2: value4-2,
 *					},
 *				}
 *			}
 *		}
 *
 * 		Into:
 * 		nestedObject = {
 * 			key1-1: value1-1
 * 			key1-2.key2-1: value2-1
 * 			key1-2.key2-2.key3-1: value3-1
 * 			key1-2.key2-2.key3-2: {
 * 				key4-1: value4-1,
 * 				key4-2: value4-2,
 *			},
 *		}
 */

// func flattenObj(obj json.RawMessage, parentKey *string, keyTrace []string) {
func FlattenObj(obj map[string]any, parentKey *string, keyTrace *string) {
	for key, value := range obj {
		// fmt.Printf("USING REFLETCT: key: %v, value: %v \n", reflect.TypeOf(key), reflect.TypeOf(value))

		switch value.(type) {
		case string:
			fmt.Printf("%v: string \n", key)
		case map[string]any:
			fmt.Printf("%v: map[string]any \n", key)
		case []map[string]any:
			fmt.Printf("%v: []map[string]any \n", key)
		case []string:
			fmt.Printf("%v: []string \n", key)

		}

		// fmt.Printf("key: %v, value: %v \n", key, value)

		// switch(value) {
		// 	case
		// }
	}
}
